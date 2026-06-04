// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package cef

import (
	"archive/tar"
	"compress/bzip2"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
)

type ProgressKind string

const (
	ProgressDownload ProgressKind = "download"
	ProgressExtract  ProgressKind = "extract"
	ProgressInfo     ProgressKind = "info"
)

type Progress struct {
	Kind    ProgressKind
	Message string
	Current int64
	Total   int64
}

type InstallOptions struct {
	Dir        string
	Version    string
	OS         string
	Arch       string
	Project    *bean.TProject
	OnProgress func(Progress)
}

type InstallResult struct {
	OSArchVersion string
	Dir           string
	ArchivePath   string
	Installed     bool
}

var SupportedOSList = []string{"windows", "linux", "darwin"}

var OSArchMap = map[string][]string{
	"windows": {"amd64", "386"},
	"linux":   {"amd64", "386", "arm64", "arm"},
	"darwin":  {"amd64", "arm64"},
}

var cefOSArchMap = map[string]map[string]string{
	"windows": {"amd64": "windows64", "386": "windows32"},
	"linux":   {"amd64": "linux64", "386": "linux32", "arm64": "linuxarm64", "arm": "linuxarm"},
	"darwin":  {"amd64": "macosx64", "arm64": "macosarm64"},
}

func Versions() []string {
	versions := make([]string, 0, len(config.DesignerConfig.Chromium))
	for ver := range config.DesignerConfig.Chromium {
		versions = append(versions, ver)
	}
	sort.Slice(versions, func(i, j int) bool {
		return CompareVersion(versions[i], versions[j]) < 0
	})
	return versions
}

func LatestVersion() string {
	versions := Versions()
	if len(versions) == 0 {
		return ""
	}
	return versions[len(versions)-1]
}

func InstalledVersions(osName, arch string) []string {
	if osName == "" {
		osName = runtime.GOOS
	}
	if arch == "" {
		arch = runtime.GOARCH
	}
	manifest := config.Config.Chromium.LoadCEFManifest()
	prefix := fmt.Sprintf("%s_%s_", osName, arch)
	var versions []string
	for oav := range manifest {
		if strings.HasPrefix(oav, prefix) && config.Config.Chromium.IsCEFInstalled(oav) {
			versions = append(versions, oav)
		}
	}
	sort.Slice(versions, func(i, j int) bool {
		return CompareVersion(ExtractVersionFromOAV(versions[i]), ExtractVersionFromOAV(versions[j])) > 0
	})
	return versions
}

func OSArchVersion(osName, arch, version string) string {
	return fmt.Sprintf("%s_%s_%s", osName, arch, version)
}

func ParseOSArchVersion(oav string) (osName, arch, version string, ok bool) {
	parts := strings.SplitN(oav, "_", 3)
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func ExtractVersionFromOAV(oav string) string {
	_, _, version, ok := ParseOSArchVersion(oav)
	if !ok {
		return ""
	}
	return version
}

func IsSupportedOSArch(osName, arch string) bool {
	archs, ok := OSArchMap[osName]
	if !ok {
		return false
	}
	for _, item := range archs {
		if item == arch {
			return true
		}
	}
	return false
}

func BuildDownloadURL(version, osName, arch string) string {
	urlTemplate := config.DesignerConfig.Chromium.Get(version)
	if urlTemplate == "" {
		return ""
	}
	result := strings.ReplaceAll(urlTemplate, "{version}", version)
	result = strings.ReplaceAll(result, "{osarch}", cefOSArch(osName, arch))
	return result
}

func ArchiveFileName(version, osName, arch string) string {
	dlURL := BuildDownloadURL(version, osName, arch)
	if dlURL == "" {
		return fmt.Sprintf("cef_%s_%s_%s.tar.bz2", osName, arch, version)
	}
	u, err := url.Parse(dlURL)
	if err != nil {
		return filepath.Base(dlURL)
	}
	name, err := url.PathUnescape(filepath.Base(u.Path))
	if err != nil {
		return filepath.Base(u.Path)
	}
	return name
}

func EnsureInstalled(ctx context.Context, options InstallOptions) (*InstallResult, error) {
	options, oav, err := normalizeInstallOptions(options, true, true, true)
	if err != nil {
		return nil, err
	}
	config.Config.Chromium.Dir = options.Dir
	if config.Config.Chromium.IsCEFInstalled(oav) {
		if err = useInstalled(oav, options.Project); err != nil {
			return nil, err
		}
		return &InstallResult{OSArchVersion: oav, Dir: config.Config.Chromium.CEFVersionDir(oav), Installed: true}, nil
	}
	return Install(ctx, options)
}

func Install(ctx context.Context, options InstallOptions) (*InstallResult, error) {
	options, oav, err := normalizeInstallOptions(options, false, false, false)
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(options.Dir, os.ModePerm); err != nil {
		return nil, err
	}
	downloadURL := BuildDownloadURL(options.Version, options.OS, options.Arch)
	if downloadURL == "" {
		return nil, fmt.Errorf("download URL not found for CEF version: %s", options.Version)
	}
	config.Config.Chromium.Dir = options.Dir
	archivePath := filepath.Join(options.Dir, ArchiveFileName(options.Version, options.OS, options.Arch))
	progress(options.OnProgress, Progress{Kind: ProgressInfo, Message: "Start downloading CEF: " + downloadURL})
	if err = download(ctx, downloadURL, archivePath, options.OnProgress); err != nil {
		return nil, err
	}
	progress(options.OnProgress, Progress{Kind: ProgressInfo, Message: "Extracting CEF: " + archivePath})
	destDir := filepath.Join(options.Dir, oav)
	files, err := ExtractTarBz2(ctx, archivePath, destDir, options.OnProgress)
	if err != nil {
		return nil, err
	}
	if err = config.Config.Chromium.SaveCEFManifest(oav, files); err != nil {
		return nil, err
	}
	if !config.Config.Chromium.IsCEFInstalled(oav) {
		return nil, errors.New("CEF installation verification failed")
	}
	if err = useInstalled(oav, options.Project); err != nil {
		return nil, err
	}
	progress(options.OnProgress, Progress{Kind: ProgressInfo, Message: "CEF installed: " + destDir})
	return &InstallResult{OSArchVersion: oav, Dir: destDir, ArchivePath: archivePath, Installed: true}, nil
}

func normalizeInstallOptions(options InstallOptions, useLatestVersion, useConfiguredDir, validateTarget bool) (InstallOptions, string, error) {
	if options.Version == "" && useLatestVersion {
		options.Version = LatestVersion()
	}
	if options.Version == "" {
		return options, "", errors.New("CEF version is empty")
	}
	if options.OS == "" {
		options.OS = runtime.GOOS
	}
	if options.Arch == "" {
		options.Arch = runtime.GOARCH
	}
	if validateTarget && !IsSupportedOSArch(options.OS, options.Arch) {
		return options, "", fmt.Errorf("unsupported CEF target: %s/%s", options.OS, options.Arch)
	}
	if options.Dir == "" {
		if useConfiguredDir && config.Config.Chromium.Dir != "" {
			options.Dir = config.Config.Chromium.Dir
		} else {
			options.Dir = config.Config.Chromium.DefaultDir()
		}
	}
	absDir, err := filepath.Abs(options.Dir)
	if err != nil {
		return options, "", err
	}
	options.Dir = absDir
	return options, OSArchVersion(options.OS, options.Arch, options.Version), nil
}

func UseInstalled(oav string, project *bean.TProject) error {
	if !config.Config.Chromium.IsCEFInstalled(oav) {
		return fmt.Errorf("CEF is not installed or incomplete: %s", oav)
	}
	return useInstalled(oav, project)
}

func useInstalled(oav string, project *bean.TProject) error {
	config.Config.Chromium.Version = oav
	if project != nil && project.GUIRenderFramework == bean.GUIRenderFramework_CEF {
		project.FrameworkVersion = oav
		if bean.GProject == project {
			if err := project.Save(); err != nil {
				return err
			}
		}
	}
	config.UpdateConfig()
	return nil
}

func ExtractTarBz2(ctx context.Context, archivePath, destDir string, onProgress func(Progress)) ([]config.CEFFileInfo, error) {
	rootDir, totalFiles, err := scanArchive(ctx, archivePath)
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(destDir, os.ModePerm); err != nil {
		return nil, err
	}
	f, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tarReader := tar.NewReader(bzip2.NewReader(&contextReader{ctx: ctx, r: f}))
	var files []config.CEFFileInfo
	var current int64
	copyBuf := make([]byte, 32*1024)
	for {
		if err = ctx.Err(); err != nil {
			return nil, err
		}
		header, nextErr := tarReader.Next()
		if nextErr != nil {
			if nextErr == io.EOF {
				break
			}
			return nil, nextErr
		}
		rel := extractRelPath(header.Name, rootDir)
		if rel == "" {
			continue
		}
		target := filepath.Join(destDir, rel)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(target, os.ModePerm); err != nil {
				return nil, err
			}
		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				return nil, err
			}
			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return nil, err
			}
			copyErr := copyWithContext(ctx, outFile, tarReader, copyBuf)
			closeErr := outFile.Close()
			if copyErr != nil {
				return nil, copyErr
			}
			if closeErr != nil {
				return nil, closeErr
			}
			if !config.IsCEFExcludeFile(filepath.Base(rel)) {
				files = append(files, config.CEFFileInfo{Name: rel, Size: header.Size})
				current++
				progress(onProgress, Progress{Kind: ProgressExtract, Current: current, Total: totalFiles, Message: fmt.Sprintf("Extracting %d / %d", current, totalFiles)})
			}
		case tar.TypeSymlink:
			if err = os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				return nil, err
			}
			_ = os.Remove(target)
			_ = os.Symlink(header.Linkname, target)
		}
	}
	return files, nil
}

func download(ctx context.Context, rawURL, targetPath string, onProgress func(Progress)) error {
	var existingSize int64
	if info, err := os.Stat(targetPath); err == nil {
		existingSize = info.Size()
	}
	remoteSize := checkRemoteFileSize(ctx, rawURL)
	if existingSize > 0 && remoteSize > 0 && existingSize >= remoteSize {
		progress(onProgress, Progress{Kind: ProgressDownload, Current: existingSize, Total: remoteSize, Message: "CEF archive already downloaded"})
		return nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	if existingSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var startSize int64
	total := resp.ContentLength
	if resp.StatusCode == http.StatusPartialContent {
		startSize = existingSize
		total = remoteSize
	} else if resp.StatusCode == http.StatusOK {
		startSize = 0
		existingSize = 0
	} else {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}
	var out *os.File
	if startSize > 0 {
		out, err = os.OpenFile(targetPath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		out, err = os.Create(targetPath)
	}
	if err != nil {
		return err
	}
	defer out.Close()
	buf := make([]byte, 32*1024)
	current := existingSize
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err = out.Write(buf[:n]); err != nil {
				return err
			}
			current += int64(n)
			progress(onProgress, Progress{Kind: ProgressDownload, Current: current, Total: total, Message: fmt.Sprintf("Downloading %s / %s", FormatSize(current), FormatSize(total))})
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}
	}
	return nil
}

func scanArchive(ctx context.Context, archivePath string) (string, int64, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	tarReader := tar.NewReader(bzip2.NewReader(&contextReader{ctx: ctx, r: f}))
	var rootDir string
	var totalFiles int64
	for {
		if err = ctx.Err(); err != nil {
			return "", 0, err
		}
		header, nextErr := tarReader.Next()
		if nextErr != nil {
			if nextErr == io.EOF {
				break
			}
			return "", 0, nextErr
		}
		name := header.Name
		if rootDir == "" {
			if idx := strings.Index(name, "/Release/"); idx >= 0 {
				rootDir = name[:idx+1]
			}
		}
		if rootDir != "" {
			rel := extractRelPath(name, rootDir)
			if rel != "" && header.Typeflag == tar.TypeReg && !config.IsCEFExcludeFile(filepath.Base(rel)) {
				totalFiles++
			}
		}
	}
	if rootDir == "" {
		return "", 0, errors.New("release directory not found in archive")
	}
	return rootDir, totalFiles, nil
}

func checkRemoteFileSize(ctx context.Context, rawURL string) int64 {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
	if err != nil {
		return -1
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return resp.ContentLength
	}
	return -1
}

func extractRelPath(name, rootDir string) string {
	if strings.HasPrefix(name, rootDir+"Release/") {
		return strings.TrimPrefix(name, rootDir+"Release/")
	}
	if strings.HasPrefix(name, rootDir+"Resources/") {
		return strings.TrimPrefix(name, rootDir+"Resources/")
	}
	return ""
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := r.r.Read(p)
	if ctxErr := r.ctx.Err(); ctxErr != nil {
		return n, ctxErr
	}
	return n, err
}

func copyWithContext(ctx context.Context, dst io.Writer, src io.Reader, buf []byte) error {
	if len(buf) == 0 {
		buf = make([]byte, 32*1024)
	}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		nr, readErr := src.Read(buf)
		if nr > 0 {
			if err := ctx.Err(); err != nil {
				return err
			}
			nw, writeErr := dst.Write(buf[:nr])
			if writeErr != nil {
				return writeErr
			}
			if nw != nr {
				return io.ErrShortWrite
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				return nil
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			return readErr
		}
	}
}

func cefOSArch(osName, arch string) string {
	if osMap, ok := cefOSArchMap[osName]; ok {
		if cefArch, ok := osMap[arch]; ok {
			return cefArch
		}
	}
	return "windows64"
}

func CompareVersion(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}
	for i := 0; i < maxLen; i++ {
		var aNum, bNum int
		if i < len(aParts) {
			aNum, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bNum, _ = strconv.Atoi(bParts[i])
		}
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
	}
	return 0
}

func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func progress(onProgress func(Progress), p Progress) {
	if onProgress != nil {
		onProgress(p)
	}
}
