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

//go:build darwin

package tool

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/energye/lcl/rtl/version"
	"os"
	"path/filepath"

	"github.com/energye/lcl/tool/command"
)

// UniversalBinary 创建 macOS 通用二进制文件（Universal Binary）
// 该函数将 amd64 和 arm64 两个架构的动态库合并为一个通用二进制文件，
// 适用于 macOS ≥ 11.0 且 Xcode ≥ 12.2 的环境
//
//	amd64ZipPath: x86_64 架构的 zip 压缩包路径
//	arm64ZipPath: ARM64 架构的 zip 压缩包路径
//	outputPath: 输出文件的目标目录路径
//
//	error: 错误信息，如果系统版本不满足要求、文件不存在或处理失败则返回相应错误
func UniversalBinary(amd64ZipPath, arm64ZipPath, outputPath string) (string, error) {
	if version.OSVersion.Major <= 10 {
		// 非 macOS ≥ 11.0 Xcode ≥ 12.2 禁用通用二进制生成
		return "", errors.New("must macOS ≥ 11.0 Xcode ≥ 12.2")
	}
	if !IsExist(amd64ZipPath) {
		return "", errors.New("amd64 zip binary not found")
	}
	if !IsExist(arm64ZipPath) {
		return "", errors.New("arm64 zip binary not found")
	}

	tempArm64LibName := "temp-libenergy-darwin-arm64-cocoa.dylib"
	tempAmd64LibName := "temp-libenergy-darwin-amd64-cocoa.dylib"
	arm64LibFilePath := filepath.Join(outputPath, tempArm64LibName)
	amd64LibFilePath := filepath.Join(outputPath, tempAmd64LibName)
	universalLibFilePath := filepath.Join(outputPath, "libenergy-darwin-universal-cocoa.dylib")
	defer func() {
		_ = os.Remove(arm64LibFilePath)
		_ = os.Remove(amd64LibFilePath)
	}()

	if err := readLib(amd64ZipPath, outputPath, tempArm64LibName); err != nil {
		return "", err
	}
	if err := readLib(arm64ZipPath, outputPath, tempAmd64LibName); err != nil {
		return "", err
	}
	// lipo -create [x86_64架构文件路径] [arm64架构文件路径] -output [输出通用二进制文件路径]
	cmd := command.NewCMD()
	cmd.Command("lipo", "-create", amd64LibFilePath, arm64LibFilePath, "-output", universalLibFilePath)
	if !IsExist(universalLibFilePath) {
		return "", errors.New("universal dylib not exist")
	}
	return universalLibFilePath, nil
}

func readLib(inZipPath, outputPath, outputFileName string) error {
	zipData, err := os.ReadFile(inZipPath)
	if err != nil {
		return err
	}
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}
	for _, file := range zipReader.File {
		_, err := ExtractFile(file, outputPath, outputFileName)
		return err
	}
	return nil
}
