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
	"embed"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/rtl/version"
	"github.com/energye/lcl/tool/command"
	"path/filepath"
)

const (
	PathAMD64Cocoa    = "darwin/libenergy-macos-amd64.zip"
	PathARM64Cocoa    = "darwin/libenergy-macos-arm64.zip"
	PathAMD64Gtk2     = "linux/libenergy-linux-amd64-gtk2.zip"
	PathAMD64Gtk3     = "linux/libenergy-linux-amd64-gtk3.zip"
	PathARMGtk2       = "linux/libenergy-linux-armhf-gtk2.zip"
	PathARMGtk3       = "linux/libenergy-linux-armhf-gtk3.zip"
	PathARM64Gtk2     = "linux/libenergy-linux-arm64-gtk2.zip"
	PathARM64Gtk3     = "linux/libenergy-linux-arm64-gtk3.zip"
	PathI386Gtk2      = "linux/libenergy-linux-i386-gtk2.zip"
	PathI386Gtk3      = "linux/libenergy-linux-i386-gtk3.zip"
	PathLoong64Gtk2   = "linux/libenergy-linux-loong64-gtk2.zip"
	PathLoong64Gtk3   = "linux/libenergy-linux-loong64-gtk3.zip"
	PathAMD64Win64    = "windows/libenergy-windows-amd64.zip"
	PathWV2AMD64Win64 = "windows/WebView2Loader-amd64.zip"
	PathI386Win32     = "windows/libenergy-windows-i386.zip"
	PathWV2I386Win32  = "windows/WebView2Loader-i386.zip"
)

var libs = tool.NewHashMap[string, *EmbedFS]()

type EmbedFS struct {
	Lib            *embed.FS
	OutputFilename string
}

// ExtractLibrary 从内置资源中提取库文件到指定输出路径
//
//   - outputPath: 库文件的输出目录路径
//   - libPath: 提取后的库文件完整路径
func ExtractLibrary(outputPath string) (libPath string) {
	libs.Iterate(func(path string, lib *EmbedFS) bool {
		tempPath := filepath.Join(outputPath, lib.OutputFilename)
		if DefaultLibName(lib.OutputFilename) {
			libPath = tempPath
		}
		if tool.IsExist(tempPath) {
			return false
		}
		libByte, e := lib.Lib.ReadFile(path)
		err.CheckErr(e)
		zipReader, e := zip.NewReader(bytes.NewReader(libByte), int64(len(libByte)))
		err.CheckErr(e)
		for _, file := range zipReader.File {
			_, e := tool.ExtractFile(file, outputPath, lib.OutputFilename)
			err.CheckErr(e)
			break // 只有一个文件
		}
		return false
	})
	if tool.IsDarwin {
		go macOSUniversalBinary(outputPath)
	}
	return
}

func DefaultLibName(filename string) bool {
	name := libname.GetDLLName()
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
	cmd.Command("lipo", "-create", amd64LibFilePath, arm64LibFilePath, "-output", universalLibFilePath)
}
