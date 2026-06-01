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

package lib

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"github.com/energye/designer/cmd/env"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/rtl/version"
	"github.com/energye/lcl/tool/command"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const (
	PathAMD64Cocoa    = "darwin/libenergy-amd64.zip"
	PathARM64Cocoa    = "darwin/libenergy-arm64.zip"
	PathAMD64Gtk2     = "linux/libenergy-amd64-gtk2.zip"
	PathAMD64Gtk3     = "linux/libenergy-amd64-gtk3.zip"
	PathARMGtk2       = "linux/libenergy-arm-gtk2.zip"
	PathARMGtk3       = "linux/libenergy-arm-gtk3.zip"
	PathARM64Gtk2     = "linux/libenergy-arm64-gtk2.zip"
	PathARM64Gtk3     = "linux/libenergy-arm64-gtk3.zip"
	Path386Gtk2       = "linux/libenergy-386-gtk2.zip"
	Path386Gtk3       = "linux/libenergy-386-gtk3.zip"
	PathLoong64Gtk2   = "linux/libenergy-loong64-gtk2.zip"
	PathLoong64Gtk3   = "linux/libenergy-loong64-gtk3.zip"
	PathAMD64Win32    = "windows/libenergy-amd64.zip"
	Path386Win32      = "windows/libenergy-386.zip"
	PathWV2AMD64Win32 = "windows/WebView2Loader-amd64.zip"
	PathWV2386Win32   = "windows/WebView2Loader-386.zip"
	PathWV2Setup      = "windows/MicrosoftEdgeWebview2Setup.zip"
)

// 存放内嵌资源
var libs = tool.NewHashMap[string, *EmbedFS]()

func Add(fs *EmbedFS) {
	if fs == nil || fs.Path == "" {
		return
	}
	libs.Add(fs.Path, fs)
}

type EmbedFS struct {
	Path           string
	Lib            *embed.FS
	OutputFilename string
	NotReleased    bool
}

func (m *EmbedFS) Release(outputPath string, files ...string) error {
	libByte, e := m.Lib.ReadFile(m.Path)
	if e != nil {
		return e
	}
	zipReader, e := zip.NewReader(bytes.NewReader(libByte), int64(len(libByte)))
	if e != nil {
		return e
	}
	for _, file := range zipReader.File {
		if len(files) > 0 {
			// 按指定文件提取, TODO 未实现
		} else {
			// 提取默认只有一个文件
			_, e := tool.ExtractFile(file, outputPath, m.OutputFilename)
			if e != nil {
				return e
			}
			break
		}
	}
	return nil
}

// ExtractLibrary 从内置资源中提取库文件到指定输出路径
//
//   - outputPath: 库文件的输出目录路径
//   - libPath: 提取后的库文件完整路径
func ExtractLibrary(outputPath string) (libPath string) {
	libs.Iterate(func(path string, lib *EmbedFS) bool {
		if lib.NotReleased {
			return false
		}
		tempPath := filepath.Join(outputPath, lib.OutputFilename)
		if DefaultLibName(lib.OutputFilename) {
			libPath = tempPath
		}
		// 读取内嵌 zip 数据
		libByte, e := lib.Lib.ReadFile(path)
		err.CheckErr(e)
		// 在内存中解压，拿到原始文件字节
		fileBytes, e := extractFromZIP(libByte)
		err.CheckErr(e)
		// 计算内嵌文件 hash
		embedHash := sha256.Sum256(fileBytes)
		// 检查磁盘文件是否需要更新
		needWrite := true
		if tool.IsExist(tempPath) {
			// 先比大小，大小不同直接判定需要更新
			diskInfo, e := os.Stat(tempPath)
			if e == nil && diskInfo.Size() == int64(len(fileBytes)) {
				// 大小相同，计算磁盘文件 hash 对比
				diskHash := hashFile(tempPath)
				if diskHash != "" && diskHash == fmt.Sprintf("%x", embedHash) {
					needWrite = false
				}
			}
		}
		if needWrite {
			_ = os.WriteFile(tempPath, fileBytes, 0755)
		}
		return false
	})
	if tool.IsDarwin {
		go macOSUniversalBinary(outputPath)
	}
	return
}

// hashFile 计算文件的 SHA256 hash，返回十六进制字符串
func hashFile(filePath string) string {
	f, e := os.Open(filePath)
	if e != nil {
		return ""
	}
	defer f.Close()
	h := sha256.New()
	if _, e := io.Copy(h, f); e != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// extractFromZIP 从 zip 数据中解压第一个非目录文件，返回其原始字节
func extractFromZIP(zipData []byte) ([]byte, error) {
	zipReader, e := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if e != nil {
		return nil, e
	}
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		srcFile, e := file.Open()
		if e != nil {
			return nil, e
		}
		data, e := io.ReadAll(srcFile)
		_ = srcFile.Close()
		if e != nil {
			return nil, e
		}
		return data, nil
	}
	return nil, fmt.Errorf("file not found in ZIP")
}

func DefaultLibName(filename string) bool {
	name := GetDLLName()
	return tool.Equal(filename, name)
}

func Libs() *tool.HashMap[string, *EmbedFS] {
	return libs
}

func macOSUniversalBinary(outputPath string) {
	version.Init()
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2
		return
	}
	universalLibFilePath := filepath.Join(outputPath, libname.DarwinUniversalBinaryName)
	if tool.IsExist(universalLibFilePath) {
		return
	}
	libArm64 := libs.Get(PathARM64Cocoa)
	if libArm64 == nil {
		panic("libArm64 is nil")
	}
	libAmd64 := libs.Get(PathAMD64Cocoa)
	if libAmd64 == nil {
		panic("libAmd64 is nil")
	}
	arm64LibFilePath := filepath.Join(outputPath, libArm64.OutputFilename)
	amd64LibFilePath := filepath.Join(outputPath, libAmd64.OutputFilename)
	cmd := command.NewCMD()
	cmd.HideWindow = true
	cmd.Command("lipo", "-create", amd64LibFilePath, arm64LibFilePath, "-output", universalLibFilePath)
}

const DarwinUniversalBinaryName = "libenergy-universal.dylib"

// GetDLLName 用于获取当前系统架构的 lib 库
func GetDLLName() string {
	goos := GOOS()
	goarch := GOARCH()
	ws, ext := "", ""
	switch goos {
	case "darwin":
		ext = "dylib"
	case "linux":
		ext = "so"
		envws := env.Get("ENERGY_WS")
		if envws == "" {
			envws = libname.UseWS
		}
		if envws == "gtk3" {
			ws = "gtk3"
		} else {
			ws = "gtk2"
		}
		if len(ws) > 0 && ws[0] != '-' {
			ws = "-" + ws // add first str "-"
		}
	case "windows":
		ext = "dll"
	}
	// windows, macOS: libenergy-[arch].xx
	// linux:  libenergy-[arch]-[ws].xx
	name := fmt.Sprintf("libenergy-%s%s.%s", goarch, ws, ext)
	return name
}

func GOOS() (goos string) {
	goos = env.Get("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	return
}

func GOARCH() (goarch string) {
	goarch = env.Get("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return
}
