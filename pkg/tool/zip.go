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

package tool

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ZipFolder 将指定的文件夹打包成一个 ZIP 文件。
// 参数：
//   - srcDir: 需要压缩的源文件夹路径，支持相对路径或绝对路径。
//   - destZip: 输出的 ZIP 文件路径。
//   - skip: 可选的 HashSet，用于指定需要跳过的顶层目录或文件名。如果为 nil，则不跳过任何内容。
//
// 返回值：
//   - error: 如果在执行过程中发生错误则返回相应的错误信息，否则返回 nil。
func ZipFolder(srcDir string, destZip string, skip *HashSet[string]) error {
	// 清理路径（处理 ./、../ 等）
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return fmt.Errorf("获取源目录绝对路径失败: %w", err)
	}

	// 检查源目录是否存在且是一个目录
	info, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("检查源目录失败: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("源路径不是文件夹: %s", srcDir)
	}

	// 创建目标 ZIP 文件
	zipFile, err := os.Create(destZip)
	if err != nil {
		return fmt.Errorf("创建 zip 文件失败: %w", err)
	}
	defer zipFile.Close()

	// 初始化 ZIP 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录中的所有文件和子目录，并逐个写入 ZIP 文件中
	err = filepath.Walk(srcDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("遍历文件失败 %s: %w", filePath, err)
		}
		filePath = filepath.ToSlash(filePath)
		relPath, err := filepath.Rel(srcDir, filePath)
		if err != nil {
			return fmt.Errorf("计算相对路径失败 %s: %w", filePath, err)
		}
		relPath = filepath.ToSlash(relPath)
		if relPath == "." || relPath == "" {
			// 跳过根目录本身
			return nil
		}
		if skip != nil {
			// 分割路径并判断是否应该跳过当前项
			parts := strings.FieldsFunc(relPath, func(c rune) bool { return c == '/' })
			if len(parts) == 0 {
				return nil
			}
			tempBasePath := parts[0]
			if skip.Contains(tempBasePath) || skip.Contains(relPath) {
				return nil
			}
		}
		// 如果是目录，则在 ZIP 中创建对应的空目录条目
		if fileInfo.IsDir() {
			_, err := zipWriter.CreateHeader(&zip.FileHeader{
				Name:     relPath + "/", // 目录必须以 / 结尾
				Method:   zip.Store,     // 不进行压缩
				Modified: fileInfo.ModTime(),
			})
			return err
		}

		// 打开源文件准备读取内容
		srcFile, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("打开文件失败 %s: %w", filePath, err)
		}
		defer srcFile.Close()

		// 在 ZIP 中创建对应文件条目
		zipEntry, err := zipWriter.CreateHeader(&zip.FileHeader{
			Name:     relPath,
			Method:   zip.Deflate, // 使用 deflate 压缩算法
			Modified: fileInfo.ModTime(),
		})
		if err != nil {
			return fmt.Errorf("创建 zip 条目失败 %s: %w", relPath, err)
		}

		// 将源文件的内容复制到 ZIP 条目中
		_, err = io.Copy(zipEntry, srcFile)
		if err != nil {
			return fmt.Errorf("写入文件到 zip 失败 %s: %w", relPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("遍历目录失败: %w", err)
	}

	return nil
}
