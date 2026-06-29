//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

package cef

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/energye/designer/pkg/config"
	"github.com/energye/lcl/api/libname"
)

type runtimeSourceConfig struct {
	Version string              `json:"version"`
	Source  string              `json:"source"`
	URL     string              `json:"url"`
	URLs    []string            `json:"urls"`
	Sources map[string][]string `json:"sources"`
}

type runtimeURLTemplate struct {
	value   string
	version string
}

func ensureRuntimeForCEF(ctx context.Context, oav string, onProgress func(Progress)) error {
	osName, arch, version, ok := ParseOSArchVersion(oav)
	if !ok {
		return fmt.Errorf("invalid CEF version: %s", oav)
	}
	if osName != runtime.GOOS || arch != runtime.GOARCH {
		return nil
	}
	major := MajorVersion(version)
	if major == "" {
		return fmt.Errorf("invalid CEF version major: %s", version)
	}
	runtimeDir := config.Config.FrameworkRuntimePath()
	if err := os.MkdirAll(runtimeDir, os.ModePerm); err != nil {
		return err
	}
	libName := runtimeLibName(osName, arch)
	activePath := filepath.Join(runtimeDir, libName)
	targetPath := filepath.Join(runtimeDir, runtimeVersionedLibName(libName, major))

	if fileExists(activePath) {
		activeVersion, err := CEFLibVersion(activePath)
		if err != nil {
			return fmt.Errorf("read active libenergy CEF version failed: %w", err)
		}
		if fmt.Sprint(activeVersion.Major) == major {
			return nil
		}
	}

	if !fileExists(targetPath) {
		if err := downloadRuntimeLib(ctx, version, major, osName, arch, runtimeWS(osName), targetPath, onProgress); err != nil {
			return err
		}
	}
	targetVersion, err := CEFLibVersion(targetPath)
	if err != nil {
		_ = os.Remove(targetPath)
		return fmt.Errorf("read target libenergy CEF version failed: %w", err)
	}
	if fmt.Sprint(targetVersion.Major) != major {
		_ = os.Remove(targetPath)
		return fmt.Errorf("libenergy CEF major mismatch: want %s, got %d", major, targetVersion.Major)
	}
	return switchRuntimeLib(activePath, targetPath, libName)
}

func switchRuntimeLib(activePath, targetPath, libName string) error {
	var backupPath, archivePath string
	if fileExists(activePath) {
		activeVersion, err := CEFLibVersion(activePath)
		if err != nil {
			return err
		}
		archivePath = filepath.Join(filepath.Dir(activePath), runtimeVersionedLibName(libName, fmt.Sprint(activeVersion.Major)))
		if archivePath != activePath {
			backupPath = activePath + ".switch-backup"
			_ = os.Remove(backupPath)
			if err = os.Rename(activePath, backupPath); err != nil {
				return err
			}
		}
	}
	if err := os.Rename(targetPath, activePath); err != nil {
		if backupPath != "" {
			_ = os.Rename(backupPath, activePath)
		}
		return err
	}
	if backupPath != "" {
		_ = os.Remove(archivePath)
		if err := os.Rename(backupPath, archivePath); err != nil {
			_ = os.Rename(backupPath, archivePath+".switch-backup")
		}
	}
	return nil
}

func downloadRuntimeLib(ctx context.Context, version, major, osName, arch, ws, targetPath string, onProgress func(Progress)) error {
	urls := RuntimeDownloadURLs(version, osName, arch, ws)
	if len(urls) == 0 {
		return fmt.Errorf("download URL not found for libenergy CEF major: %s", major)
	}
	var errs []error
	for i, rawURL := range urls {
		archivePath := filepath.Join(filepath.Dir(targetPath), runtimeArchiveFileName(rawURL, major, i))
		tmpTargetPath := runtimeTempLibPath(targetPath)
		_ = os.Remove(archivePath)
		_ = os.Remove(tmpTargetPath)
		progress(onProgress, Progress{Kind: ProgressInfo, Message: "Start downloading libenergy: " + rawURL})
		if err := download(ctx, rawURL, archivePath, onProgress); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", rawURL, err))
			continue
		}
		progress(onProgress, Progress{Kind: ProgressInfo, Message: "Extracting libenergy: " + archivePath})
		if err := extractRuntimeZip(archivePath, tmpTargetPath); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", rawURL, err))
			continue
		}
		targetVersion, err := CEFLibVersion(tmpTargetPath)
		if err != nil {
			_ = os.Remove(tmpTargetPath)
			errs = append(errs, fmt.Errorf("%s: read libenergy CEF version failed: %w", rawURL, err))
			continue
		}
		if fmt.Sprint(targetVersion.Major) != major {
			_ = os.Remove(tmpTargetPath)
			errs = append(errs, fmt.Errorf("%s: libenergy CEF major mismatch: want %s, got %d", rawURL, major, targetVersion.Major))
			continue
		}
		if err := os.Rename(tmpTargetPath, targetPath); err != nil {
			_ = os.Remove(tmpTargetPath)
			errs = append(errs, fmt.Errorf("%s: install libenergy failed: %w", rawURL, err))
			continue
		}
		return nil
	}
	return errors.Join(errs...)
}

func runtimeTempLibPath(targetPath string) string {
	ext := filepath.Ext(targetPath)
	base := strings.TrimSuffix(targetPath, ext)
	return base + ".download" + ext
}

func RuntimeDownloadURLs(version, osName, arch, ws string) []string {
	major := MajorVersion(version)
	configs := runtimeSourceConfigs()
	defaultReleaseVersion := runtimeReleaseVersion(configs)
	selected := runtimeSelectedSource(configs)
	var templates []runtimeURLTemplate
	addTemplate := func(value, releaseVersion string) {
		templates = append(templates, runtimeURLTemplate{value: value, version: releaseVersion})
	}
	for _, cfg := range configs {
		releaseVersion := strings.TrimSpace(cfg.Version)
		if releaseVersion == "" {
			releaseVersion = defaultReleaseVersion
		}
		addTemplate(cfg.URL, releaseVersion)
		for _, value := range cfg.URLs {
			addTemplate(value, releaseVersion)
		}
		if len(cfg.Sources) > 0 {
			names := make([]string, 0, len(cfg.Sources))
			for name := range cfg.Sources {
				names = append(names, name)
			}
			sort.Strings(names)
			if selected != "" {
				if values := cfg.Sources[selected]; len(values) > 0 {
					for _, value := range values {
						addTemplate(value, releaseVersion)
					}
				}
			}
			for _, name := range names {
				if name == selected {
					continue
				}
				for _, value := range cfg.Sources[name] {
					addTemplate(value, releaseVersion)
				}
			}
		}
	}
	if envURL := strings.TrimSpace(os.Getenv("ENERGY_CEF_RUNTIME_URL")); envURL != "" {
		envTemplates := strings.FieldsFunc(envURL, func(r rune) bool {
			return r == ';' || r == ','
		})
		prefix := make([]runtimeURLTemplate, 0, len(envTemplates))
		for _, value := range envTemplates {
			prefix = append(prefix, runtimeURLTemplate{value: value, version: defaultReleaseVersion})
		}
		templates = append(prefix, templates...)
	}
	var result []string
	seen := map[string]bool{}
	for _, tmpl := range templates {
		value := strings.TrimSpace(tmpl.value)
		if value == "" {
			continue
		}
		u := expandRuntimeURL(value, version, major, osName, arch, ws, tmpl.version)
		if !seen[u] {
			seen[u] = true
			result = append(result, u)
		}
	}
	return result
}

func runtimeSourceConfigs() []runtimeSourceConfig {
	return decodeRuntimeSourceConfigs(config.DesignerConfig.CEFRuntime)
}

func runtimeSelectedSource(configs []runtimeSourceConfig) string {
	if selected := strings.TrimSpace(os.Getenv("ENERGY_CEF_RUNTIME_SOURCE")); selected != "" {
		return selected
	}
	if selected := selectedRuntimeSource(config.Config.CEFRuntime); selected != "" {
		return selected
	}
	return firstRuntimeSource(configs)
}

func selectedRuntimeSource(raw json.RawMessage) string {
	var cfg struct {
		Source string `json:"source"`
	}
	if json.Unmarshal(raw, &cfg) != nil {
		return ""
	}
	return strings.TrimSpace(cfg.Source)
}

func firstRuntimeSource(configs []runtimeSourceConfig) string {
	for _, cfg := range configs {
		if source := strings.TrimSpace(cfg.Source); source != "" {
			return source
		}
	}
	return ""
}

func runtimeReleaseVersion(configs []runtimeSourceConfig) string {
	for _, cfg := range configs {
		if version := strings.TrimSpace(cfg.Version); version != "" {
			return version
		}
	}
	return ""
}

func decodeRuntimeSourceConfigs(raw json.RawMessage) []runtimeSourceConfig {
	var cfg runtimeSourceConfig
	if json.Unmarshal(raw, &cfg) == nil && (cfg.Version != "" || cfg.Source != "" || cfg.URL != "" || len(cfg.URLs) > 0 || len(cfg.Sources) > 0) {
		return []runtimeSourceConfig{cfg}
	}
	var list []runtimeSourceConfig
	if json.Unmarshal(raw, &list) == nil {
		return list
	}
	return nil
}

func expandRuntimeURL(tmpl, version, major, osName, arch, ws, releaseVersion string) string {
	releaseVersion = strings.TrimSpace(releaseVersion)
	if releaseVersion == "" {
		releaseVersion = version
	}
	if ws == "" {
		tmpl = removeOptionalWSPlaceholder(tmpl)
	}
	replacer := strings.NewReplacer(
		"{version}", releaseVersion,
		"{major}", major,
		"{os}", osName,
		"{arch}", arch,
		"{ws}", ws,
	)
	return replacer.Replace(tmpl)
}

func removeOptionalWSPlaceholder(tmpl string) string {
	replacer := strings.NewReplacer(
		"-{ws}-", "-",
		"_{ws}_", "_",
		"/{ws}/", "/",
		"-{ws}", "",
		"_{ws}", "",
		"/{ws}", "",
		"{ws}-", "",
		"{ws}_", "",
		"{ws}/", "",
		"{ws}", "",
	)
	return replacer.Replace(tmpl)
}

func extractRuntimeZip(archivePath, targetPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		if err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			_ = src.Close()
			return err
		}
		dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			_ = src.Close()
			return err
		}
		_, copyErr := io.Copy(dst, src)
		closeSrcErr := src.Close()
		closeDstErr := dst.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeSrcErr != nil {
			return closeSrcErr
		}
		return closeDstErr
	}
	return errors.New("runtime zip has no file")
}

func runtimeArchiveFileName(rawURL, major string, index int) string {
	u, err := url.Parse(rawURL)
	if err == nil {
		name := filepath.Base(u.Path)
		if name == "download" {
			name = filepath.Base(filepath.Dir(u.Path))
		}
		if name != "." && name != "/" && name != "" {
			if unescaped, unescapeErr := url.PathUnescape(name); unescapeErr == nil {
				return fmt.Sprintf("%d_%s", index, unescaped)
			}
			return fmt.Sprintf("%d_%s", index, name)
		}
	}
	return fmt.Sprintf("%d_libenergy_%s.zip", index, major)
}

func runtimeLibName(osName, arch string) string {
	ws, ext := "", ""
	switch osName {
	case "darwin":
		ext = "dylib"
	case "linux":
		ext = "so"
		ws = runtimeWS(osName)
		if ws != "" {
			ws = "-" + ws
		}
	case "windows":
		ext = "dll"
	}
	return fmt.Sprintf("libenergy-%s%s.%s", arch, ws, ext)
}

func runtimeVersionedLibName(libName, major string) string {
	ext := filepath.Ext(libName)
	base := strings.TrimSuffix(libName, ext)
	return base + "-" + major + ext
}

func runtimeWS(osName string) string {
	if osName != "linux" {
		return ""
	}
	if os.Getenv("ENERGY_WS") != "" {
		return os.Getenv("ENERGY_WS")
	}
	if libname.UseWS != "" {
		return libname.UseWS
	}
	return "gtk3"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type RuntimeVersion struct {
	Major   int32
	Minor   int32
	Release int32
	Build   int32
}

func MajorVersion(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return ""
	}
	if idx := strings.Index(version, "."); idx >= 0 {
		return version[:idx]
	}
	return version
}

func ResolveConfiguredVersion(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if _, _, version, ok := ParseOSArchVersion(value); ok {
		value = version
	}
	for _, ver := range Versions() {
		if ver == value {
			return ver
		}
	}
	var matches []string
	for _, ver := range Versions() {
		if strings.HasPrefix(ver, value) || MajorVersion(ver) == value {
			matches = append(matches, ver)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	return ""
}
