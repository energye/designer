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
	"bytes"
	"errors"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"text/template"
)

const (
	IsWindows = runtime.GOOS == "windows"
	IsLinux   = runtime.GOOS == "linux"
	IsDarwin  = runtime.GOOS == "darwin"
)

// Go 关键字列表（完整官方版）
var goKeywords = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
	"any":         true,
}

// 匹配 Go 变量名规则
var varNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// IsValidVariableName 判断是否为合法的 Go 变量名/字段名
func IsValidVariableName(s string) bool {
	if s == "" {
		return false
	}
	if goKeywords[s] {
		return false
	}
	return varNameRegex.MatchString(s)
}

// 加载图像到列表
func ImageListAddPng(imageList lcl.IImageList, filePath string) {
	data := resources.Images(filePath)
	if data != nil {
		pic := lcl.NewPicture()
		defer pic.Free()

		mem := lcl.NewMemoryStream()
		defer mem.Free()
		lcl.StreamHelper.WriteBuffer(mem, data)
		mem.SetPosition(0)
		pic.LoadFromStream(mem)
		imageList.Add(pic.Bitmap(), nil)
	} else {
		data = resources.Images("components/default_150.png")
		pic := lcl.NewPicture()
		defer pic.Free()

		mem := lcl.NewMemoryStream()
		defer mem.Free()
		lcl.StreamHelper.WriteBuffer(mem, data)
		mem.SetPosition(0)
		pic.LoadFromStream(mem)
		imageList.Add(pic.Bitmap(), nil)
	}
}

// 加载图片列表
func LoadImageList(owner lcl.IComponent, imageList []string, width, height int32) lcl.IImageList {
	images := lcl.NewImageList(owner)
	if width > 0 {
		images.SetWidth(width)
	}
	if height > 0 {
		images.SetHeight(height)
	}
	for _, image := range imageList {
		ImageListAddPng(images, image)
	}
	return images
}

// 判断字符串相等, 忽略大小写
func Equal(s1 string, s2 ...string) bool {
	s1 = strings.ToLower(s1)
	for _, s := range s2 {
		s = strings.ToLower(s)
		if s1 == s {
			return true
		}
	}
	return false
}

// 删除道字母 T
func RemoveT(name string) string {
	if name == "" {
		return ""
	}
	if name[0] == 'T' {
		return name[1:]
	}
	return name
}

// 第一个字母转为大写
func FirstToUpper(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// 判断文件是否目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	if s.IsDir() {
		return true
	}
	return false
}

// 字符串数组元素反转
func StringArrayReverse(array []string) {
	n := len(array)
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		array[i], array[j] = array[j], array[i]
	}
}

// 分割字符串并去除空字符串 aa,bb,,cc > [aa,bb,cc]
func Split(s, sep string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var result []string
	for _, v := range strings.Split(s, sep) {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		result = append(result, v)
	}
	return result[:]
}

// 字符串替换
func Replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

// 判断当前是否为主线程
func IsMainThread() bool {
	return api.MainThreadId() == api.CurrentThreadId()
}

// ExtractFile 从zip文件中提取指定文件到目标路径
//
//	zipFile: 要提取的zip文件对象
//	targetFile: 目标文件路径
//	error: 提取过程中发生的错误，如果成功则返回nil
func ExtractFile(zipFile *zip.File, outputPath, outputFilename string) (string, error) {
	srcFile, err := zipFile.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()
	if outputFilename == "" {
		outputPath = filepath.Join(outputPath, zipFile.Name)
	} else {
		outputPath = filepath.Join(outputPath, outputFilename)
	}
	if zipFile.Mode().IsDir() {
		return outputPath, os.MkdirAll(outputPath, 0755)
	}
	dstFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		if runtime.GOOS == "windows" {
			if pathErr, ok := err.(*os.PathError); ok {
				if errno, ok := pathErr.Err.(syscall.Errno); ok && errno == 32 {
					// 文件正在使用中, 路过提取
					logs.Error("File is busy, skipping extraction:", outputPath)
					return outputPath, nil
				}
			}
		}
		return outputPath, err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return outputPath, err
}

func RenderTemplate(templateText string, data any) ([]byte, error) {
	tmpl, err := template.New("RenderTemplate").Parse(templateText)
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	if err = tmpl.Execute(&out, data); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func CopyFile(srcFilePath, dstFilePath string) error {
	srcFileInfo, err := os.Stat(srcFilePath)
	if err != nil {
		return err
	}
	if srcFileInfo.IsDir() {
		return errors.New("source path is a directory, not a file")
	}
	dstDir := filepath.Dir(dstFilePath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	if err := dstFile.Sync(); err != nil {
		return err
	}
	if err := os.Chmod(dstFilePath, srcFileInfo.Mode()); err != nil {
		return err
	}
	return nil
}

func ContainsStrings(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

//// DirPermissions 结构体用于存储目录权限检查结果
//type DirPermissions struct {
//	Exists     bool
//	Readable   bool
//	Writable   bool
//	Executable bool // 在Windows上始终为false
//}
//
//// CheckDirPermissions 检查指定目录的存在性及读写执行权限
//func CheckDirPermissions(dirPath string) (DirPermissions, error) {
//	var perm DirPermissions
//	// 1. 检查目录是否存在
//	info, err := os.Stat(dirPath)
//	if err != nil {
//		if os.IsNotExist(err) {
//			perm.Exists = false
//			// 即使目录不存在，也要检查父目录是否可写（用于创建目录）
//			parentDir := filepath.Dir(dirPath)
//			_, parentErr := os.Stat(parentDir)
//			if parentErr != nil {
//				return perm, nil // 父目录也不存在或无法访问
//			}
//			// 检查父目录写权限
//			tempFile, createErr := os.CreateTemp(parentDir, ".perm-check-")
//			if createErr != nil {
//				perm.Writable = false
//			} else {
//				perm.Writable = true
//				defer tempFile.Close()
//				defer os.Remove(tempFile.Name())
//			}
//			return perm, nil
//		}
//		// 如果是因为权限不足导致无法stat，Exists仍为true，但其他权限会是false
//		perm.Exists = true
//		return perm, fmt.Errorf("failed to stat directory: %w", err)
//	}
//	perm.Exists = true
//	// 确认路径是一个目录
//	if !info.IsDir() {
//		return perm, errors.New("path is not a directory")
//	}
//	// 2. 检查读权限 (尝试打开目录)
//	dir, err := os.Open(dirPath)
//	if err != nil {
//		perm.Readable = false
//	} else {
//		perm.Readable = true
//		defer dir.Close()
//	}
//	// 3. 检查写权限 (尝试创建一个临时文件)
//	tempFile, err := os.CreateTemp(dirPath, ".perm-check-")
//	if err != nil {
//		perm.Writable = false
//	} else {
//		perm.Writable = true
//		defer tempFile.Close()
//		defer os.Remove(tempFile.Name()) // 清理临时文件
//	}
//	// 4. 检查执行权限 (仅在非Windows系统上有效)
//	if runtime.GOOS != "windows" {
//		// 目录的执行权限位
//		// 对于所有者、组、其他用户，只要有一个有执行权限，就认为可执行
//		perm.Executable = info.Mode()&0111 != 0
//	} else {
//		perm.Executable = false
//	}
//	return perm, nil
//}
