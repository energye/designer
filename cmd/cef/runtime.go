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
	release string
}

type runtimeDownloadURL struct {
	value   string
	release string
}

// RuntimeOverride 单次安装可覆盖的运行时库参数, 为空时回退到 ~/.energy/config.json
type RuntimeOverride struct {
	Release string // libenergy 发行版本, 如 v3.0.1
	Source  string // 下载源名: github / sourceforge
	URL     string // 自定义下载地址模板, 多个用 , 或 ; 分隔
	Module  string // 模块名模板, 如 cef{major} / cef{major}-webview / webview; "-" 表示无模块前缀
	WS      string // Linux 窗口系统: gtk2 / gtk3
}

// RuntimeOptions 一次 libenergy 运行时库安装所需的全部参数
type RuntimeOptions struct {
	Version string   // CEF 完整版本, 如 127.3.5
	Major   string   // CEF major, 如 127
	OS      string   // windows / linux / darwin
	Arch    string   // amd64 / 386 / arm64 / arm
	WS      string   // 仅 Linux: gtk2 / gtk3, 其它平台为空
	Module  string   // 资源名中的模块段, 如 cef127, 为空表示不带模块
	Release string   // libenergy 发行版本, 如 v3.0.1
	Source  string   // 当前下载源
	URLs    []string // 用户自定义 URL 模板, 优先级最高

	releaseOverride bool // 发行版本是否由用户显式指定
}

const (
	// RuntimeModuleNone 用于把 module 段置空, 对应 libenergy-{os}-{arch}{-ws}.zip
	RuntimeModuleNone = "-"
	// DefaultRuntimeModule 默认模块名, {major} 会替换为 CEF major
	DefaultRuntimeModule = "cef{major}"
)

// NewRuntimeOptions 根据 CEF 版本、目标平台和自定义参数解析出完整安装参数
func NewRuntimeOptions(version, osName, arch string, override RuntimeOverride) RuntimeOptions {
	return resolveRuntimeOptions(version, osName, arch, override)
}

// AssetName 返回 libenergy 资源的 zip 文件名
// 规则: libenergy-{module}-{os}-{arch}{-ws}.zip, {-ws} 仅 Linux 存在
func (o RuntimeOptions) AssetName() string {
	name := "libenergy"
	if o.Module != "" {
		name += "-" + o.Module
	}
	name += "-" + o.OS + "-" + o.Arch
	if o.WS != "" {
		name += "-" + o.WS
	}
	return name + ".zip"
}

// DownloadURLs 返回按优先级排列的运行库下载地址
func (o RuntimeOptions) DownloadURLs() []string {
	return RuntimeDownloadURLsFor(o)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func splitRuntimeURLTemplates(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool {
		return r == ';' || r == ','
	})
}

func resolveRuntimeModule(value, major, version string) string {
	value = strings.TrimSpace(value)
	if value == RuntimeModuleNone {
		return ""
	}
	if value == "" {
		value = DefaultRuntimeModule
	}
	value = strings.NewReplacer(
		"{major}", major,
		"{version}", version,
		"{cef}", "cef"+major,
	).Replace(value)
	return strings.Trim(value, "-")
}

func normalizeWS(osName, value string) string {
	if osName != "linux" {
		return ""
	}
	check := func(v string) string {
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "gtk2":
			return "gtk2"
		case "gtk3":
			return "gtk3"
		}
		return ""
	}
	if ws := check(value); ws != "" {
		return ws
	}
	if ws := check(os.Getenv("ENERGY_WS")); ws != "" {
		return ws
	}
	if ws := check(libname.UseWS); ws != "" {
		return ws
	}
	return "gtk3"
}

func resolveRuntimeOptions(version, osName, arch string, override RuntimeOverride) RuntimeOptions {
	major := MajorVersion(version)
	cfg := config.Config.Chromium
	configs := runtimeSourceConfigs()
	release := firstNonEmpty(override.Release, cfg.RuntimeVersion)
	releaseOverride := release != ""
	if !releaseOverride {
		release = runtimeReleaseVersion(configs)
	}
	return RuntimeOptions{
		Version:         version,
		Major:           major,
		OS:              osName,
		Arch:            arch,
		WS:              normalizeWS(osName, firstNonEmpty(override.WS, cfg.RuntimeWS)),
		Module:          resolveRuntimeModule(firstNonEmpty(override.Module, cfg.RuntimeModule), major, version),
		Release:         release,
		Source:          firstNonEmpty(os.Getenv("ENERGY_CEF_RUNTIME_SOURCE"), override.Source, cfg.Source, firstRuntimeSource(configs)),
		URLs:            splitRuntimeURLTemplates(firstNonEmpty(override.URL, cfg.RuntimeURL)),
		releaseOverride: releaseOverride,
	}
}

func ensureRuntimeForCEF(ctx context.Context, oav string, override RuntimeOverride, onProgress func(Progress)) error {
	osName, arch, version, ok := ParseOSArchVersion(oav)
	if !ok {
		return fmt.Errorf("invalid CEF version: %s", oav)
	}
	if osName != runtime.GOOS || arch != runtime.GOARCH {
		progress(onProgress, Progress{Kind: ProgressInfo, Message: "Skip libenergy runtime: target differs from current platform " + runtime.GOOS + "/" + runtime.GOARCH})
		return nil
	}
	if MajorVersion(version) == "" {
		return fmt.Errorf("invalid CEF version major: %s", version)
	}
	options := resolveRuntimeOptions(version, osName, arch, override)
	libName := runtimeLibName(osName, arch, options.WS)
	targetPath := filepath.Join(config.Config.Chromium.CEFVersionDir(oav), libName)

	if !isCachedRuntimeValid(oav, targetPath, options) {
		if err := downloadRuntimeLib(ctx, options, targetPath, onProgress); err != nil {
			return err
		}
		if err := saveRuntimeManifest(oav, options, targetPath); err != nil {
			return err
		}
	}
	return nil
}

// runtimeFingerprint 记录影响资源选择的所有自定义参数
func runtimeFingerprint(options RuntimeOptions) string {
	return strings.Join([]string{options.Release, options.Module, options.WS, options.Source}, "|")
}

func isCachedRuntimeValid(oav, targetPath string, options RuntimeOptions) bool {
	info, ok := config.Config.Chromium.CEFRuntimeManifest(oav)
	if !ok || info.Release != options.Release || info.File.Name != filepath.Base(targetPath) {
		return false
	}
	// 新清单带指纹, 任一自定义项变化都要重新安装 (包括模块段置空的情况)
	if info.Fingerprint != "" && info.Fingerprint != runtimeFingerprint(options) {
		return false
	}
	// 旧清单没有指纹时按已记录的字段做兼容比较
	if info.Fingerprint == "" {
		if info.Module != "" && info.Module != options.Module {
			return false
		}
		if info.WS != "" && info.WS != options.WS {
			return false
		}
		if info.Source != "" && info.Source != options.Source {
			return false
		}
	}
	stat, err := os.Stat(targetPath)
	if err != nil || stat.IsDir() || stat.Size() != info.File.Size {
		return false
	}
	version, err := CEFLibVersion(targetPath)
	return err == nil && fmt.Sprint(version.Major) == options.Major
}

func saveRuntimeManifest(oav string, options RuntimeOptions, targetPath string) error {
	stat, err := os.Stat(targetPath)
	if err != nil {
		return err
	}
	return config.Config.Chromium.SaveCEFRuntimeManifestDetail(oav, config.CEFRuntimeInfo{
		Release:     options.Release,
		Module:      options.Module,
		WS:          options.WS,
		Source:      options.Source,
		Fingerprint: runtimeFingerprint(options),
		File: config.CEFFileInfo{
			Name: filepath.Base(targetPath),
			Size: stat.Size(),
		},
	})
}

func downloadRuntimeLib(ctx context.Context, options RuntimeOptions, targetPath string, onProgress func(Progress)) error {
	urls := runtimeDownloadURLsFor(options)
	if len(urls) == 0 {
		return fmt.Errorf("download URL not found for libenergy asset: %s", options.AssetName())
	}
	var errs []error
	for _, item := range urls {
		rawURL := item.value
		archivePath := runtimeArchivePath(config.Config.Chromium.Dir, rawURL, options.Major, item.release)
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
		if fmt.Sprint(targetVersion.Major) != options.Major {
			_ = os.Remove(tmpTargetPath)
			errs = append(errs, fmt.Errorf("%s: libenergy CEF major mismatch: want %s, got %d", rawURL, options.Major, targetVersion.Major))
			continue
		}
		if err := replaceFile(tmpTargetPath, targetPath); err != nil {
			_ = os.Remove(tmpTargetPath)
			errs = append(errs, fmt.Errorf("%s: install libenergy failed: %w", rawURL, err))
			continue
		}
		return nil
	}
	return errors.Join(errs...)
}

func runtimeArchivePath(chromiumDir, rawURL, major, releaseVersion string) string {
	return filepath.Join(chromiumDir, runtimeArchiveFileName(rawURL, major, releaseVersion))
}

func replaceFile(srcPath, dstPath string) error {
	backupPath := dstPath + ".bak"
	_ = os.Remove(backupPath)
	hadDst := false
	if _, err := os.Stat(dstPath); err == nil {
		hadDst = true
		if err = os.Rename(dstPath, backupPath); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(srcPath, dstPath); err != nil {
		if hadDst {
			_ = os.Rename(backupPath, dstPath)
		}
		return err
	}
	if hadDst {
		_ = os.Remove(backupPath)
	}
	return nil
}

func runtimeTempLibPath(targetPath string) string {
	ext := filepath.Ext(targetPath)
	base := strings.TrimSuffix(targetPath, ext)
	return base + ".download" + ext
}

func RuntimeDownloadURLs(version, osName, arch, ws string) []string {
	return urlsOf(runtimeDownloadURLs(version, osName, arch, ws))
}

// RuntimeDownloadURLsFor 按自定义参数返回排序后的下载地址
func RuntimeDownloadURLsFor(options RuntimeOptions) []string {
	return urlsOf(runtimeDownloadURLsFor(options))
}

func urlsOf(downloads []runtimeDownloadURL) []string {
	result := make([]string, 0, len(downloads))
	for _, item := range downloads {
		result = append(result, item.value)
	}
	return result
}

func runtimeDownloadURLs(version, osName, arch, ws string) []runtimeDownloadURL {
	return runtimeDownloadURLsFor(resolveRuntimeOptions(version, osName, arch, RuntimeOverride{WS: ws}))
}

// runtimeDownloadURLsFor 地址优先级: 自定义 URL > 环境变量 > 选中下载源 > 其它下载源
func runtimeDownloadURLsFor(options RuntimeOptions) []runtimeDownloadURL {
	configs := runtimeSourceConfigs()
	defaultReleaseVersion := firstNonEmpty(options.Release, options.Version)
	selected := runtimeSelectedSource(configs, options.Source)
	releaseFor := func(cfgVersion string) string {
		if options.releaseOverride {
			return defaultReleaseVersion
		}
		return firstNonEmpty(cfgVersion, defaultReleaseVersion)
	}
	var templates []runtimeURLTemplate
	addTemplate := func(value, releaseVersion string) {
		templates = append(templates, runtimeURLTemplate{value: value, release: releaseVersion})
	}
	for _, cfg := range configs {
		releaseVersion := releaseFor(cfg.Version)
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
				for _, value := range cfg.Sources[selected] {
					addTemplate(value, releaseVersion)
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
		envTemplates := splitRuntimeURLTemplates(envURL)
		prefix := make([]runtimeURLTemplate, 0, len(envTemplates))
		for _, value := range envTemplates {
			prefix = append(prefix, runtimeURLTemplate{value: value, release: defaultReleaseVersion})
		}
		templates = append(prefix, templates...)
	}
	// 用户显式指定的地址放在最前, 保证自定义配置优先级最高
	if len(options.URLs) > 0 {
		custom := make([]runtimeURLTemplate, 0, len(options.URLs))
		for _, value := range options.URLs {
			custom = append(custom, runtimeURLTemplate{value: value, release: defaultReleaseVersion})
		}
		templates = append(custom, templates...)
	}
	var result []runtimeDownloadURL
	seen := map[string]bool{}
	for _, tmpl := range templates {
		value := strings.TrimSpace(tmpl.value)
		if value == "" {
			continue
		}
		u := expandRuntimeURLFor(value, options, tmpl.release)
		if !seen[u] {
			seen[u] = true
			result = append(result, runtimeDownloadURL{value: u, release: tmpl.release})
		}
	}
	return result
}

func runtimeSourceConfigs() []runtimeSourceConfig {
	return decodeRuntimeSourceConfigs(config.DesignerConfig.CEFRuntime)
}

func runtimeSelectedSource(configs []runtimeSourceConfig, preferred string) string {
	if selected := strings.TrimSpace(os.Getenv("ENERGY_CEF_RUNTIME_SOURCE")); selected != "" {
		return selected
	}
	if selected := strings.TrimSpace(preferred); selected != "" {
		return selected
	}
	return firstRuntimeSource(configs)
}

func firstRuntimeSource(configs []runtimeSourceConfig) string {
	for _, cfg := range configs {
		if source := strings.TrimSpace(cfg.Source); source != "" {
			return source
		}
	}
	return ""
}

// RuntimeSources 返回内嵌配置中可用的 libenergy 下载源
func RuntimeSources() []string {
	seen := map[string]bool{}
	for _, cfg := range runtimeSourceConfigs() {
		if source := strings.TrimSpace(cfg.Source); source != "" {
			seen[source] = true
		}
		for name := range cfg.Sources {
			if name = strings.TrimSpace(name); name != "" {
				seen[name] = true
			}
		}
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// RuntimeReleaseVersion 返回内嵌配置中的 libenergy 发行版本
func RuntimeReleaseVersion() string {
	return runtimeReleaseVersion(runtimeSourceConfigs())
}

// RuntimeDefaultSource 返回内嵌配置中的默认下载源
func RuntimeDefaultSource() string {
	return firstRuntimeSource(runtimeSourceConfigs())
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

// expandRuntimeURLFor 展开 libenergy 下载地址模板
//
// 占位符分两类:
//
//   - 纯值: {version} {major} {module} {os} {arch} {ws}
//     只替换为值本身, 值不存在时为空字符串
//   - 可选段: {module-} {-ws}
//     连字符归属占位符, 值为 "值+连字符" 或 "连字符+值"; 值不存在时整个段(含连字符)消失
//
// 资源命名规则固定为 libenergy-{module-}{os}-{arch}{-ws}.zip,
// 可选段是原子的, 不需要展开后折叠 "--", 也不需要猜测该删哪一个连字符。
func expandRuntimeURLFor(tmpl string, options RuntimeOptions, releaseVersion string) string {
	releaseVersion = firstNonEmpty(releaseVersion, options.Release, options.Version)
	moduleSegment, wsSegment := "", ""
	if options.Module != "" {
		moduleSegment = options.Module + "-"
	}
	if options.WS != "" {
		wsSegment = "-" + options.WS
	}
	replacer := strings.NewReplacer(
		"{version}", releaseVersion,
		"{major}", options.Major,
		"{module-}", moduleSegment,
		"{-ws}", wsSegment,
		"{module}", options.Module,
		"{os}", options.OS,
		"{arch}", options.Arch,
		"{ws}", options.WS,
	)
	return replacer.Replace(tmpl)
}

// extractRuntimeZip 从资源包中提取动态库并按目标文件名安装
// 新版本资源包可能同时带有说明文件, 因此优先选取与目标库扩展名一致的文件
func extractRuntimeZip(archivePath, targetPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	var candidate *zip.File
	ext := strings.ToLower(filepath.Ext(targetPath))
	for _, file := range reader.File {
		if file.FileInfo().IsDir() || file.FileInfo().Size() <= 0 {
			continue
		}
		if candidate == nil {
			candidate = file
		}
		if ext != "" && strings.EqualFold(filepath.Ext(file.Name), ext) {
			candidate = file
			break
		}
	}
	if candidate == nil {
		return errors.New("runtime zip has no library file")
	}
	src, err := candidate.Open()
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

func runtimeArchiveFileName(rawURL, major, releaseVersion string) string {
	prefix := runtimeArchivePrefix(releaseVersion, major)
	u, err := url.Parse(rawURL)
	if err == nil {
		name := filepath.Base(u.Path)
		if name == "download" {
			name = filepath.Base(filepath.Dir(u.Path))
		}
		if name != "." && name != "/" && name != "" {
			if unescaped, unescapeErr := url.PathUnescape(name); unescapeErr == nil {
				return fmt.Sprintf("%s_%s", prefix, unescaped)
			}
			return fmt.Sprintf("%s_%s", prefix, name)
		}
	}
	return fmt.Sprintf("%s_libenergy_%s.zip", prefix, major)
}

func runtimeArchivePrefix(releaseVersion, major string) string {
	prefix := strings.TrimSpace(releaseVersion)
	if prefix == "" {
		prefix = major
	}
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		`"`, "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(prefix)
}

// runtimeLibName 返回运行时库安装后的文件名, 需与 lcl/api/libname 的查找规则一致
// windows, macOS: libenergy-<arch>.<ext>; linux: libenergy-<arch>-<ws>.<ext>
func runtimeLibName(osName, arch, ws string) string {
	ws, ext := normalizeWS(osName, ws), ""
	switch osName {
	case "darwin":
		ext = "dylib"
	case "linux":
		ext = "so"
		if ws != "" {
			ws = "-" + ws
		}
	case "windows":
		ext = "dll"
	}
	return fmt.Sprintf("libenergy-%s%s.%s", arch, ws, ext)
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
