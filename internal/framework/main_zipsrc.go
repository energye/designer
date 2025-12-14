package main

import (
	"fmt"
	"github.com/energye/designer/internal/framework/genzip"
	"os"
	"path/filepath"
)

func main() {
	// 压缩源文件夹所在的根目录
	// 其中包含 lcl cef wv
	zipSrcRoot := "E:\\app\\workspace"
	// 压缩后复制的框架目标目录
	// 其中包含 lcl cef wv
	descFrameworksDir := "E:\\app\\workspace\\designer\\resources\\frameworks"
	// 开始压缩文件
	lcl, cef, wv := genzip.ZipSRC(zipSrcRoot)
	fmt.Println(lcl, cef, wv)
	// 开始移动文件到框架目录
	moveZipToFrameworks(lcl, filepath.Join(descFrameworksDir, "lcl"))
	moveZipToFrameworks(cef, filepath.Join(descFrameworksDir, "cef"))
	moveZipToFrameworks(wv, filepath.Join(descFrameworksDir, "wv"))
}

// 压缩包复制到框架目录
func moveZipToFrameworks(srcZip, descFrameworksDir string) {
	fmt.Println("开始移动文件 ", srcZip, "=>", descFrameworksDir)
	// 源文件为空时直接返回
	if srcZip == "" {
		fmt.Println("源压缩包路径为空，跳过复制")
		return
	}
	// 检查源文件是否存在
	if _, err := os.Stat(srcZip); os.IsNotExist(err) {
		fmt.Printf("源压缩包不存在: %s, 错误: %v\n", srcZip, err)
		return
	} else if err != nil {
		fmt.Printf("检查源压缩包失败: %s, 错误: %v\n", srcZip, err)
		return
	}
	// 获取源文件名
	fileName := filepath.Base(srcZip)
	// 构建目标文件路径
	dstFile := filepath.Join(descFrameworksDir, fileName)

	// 如果目标文件已存在，先删除
	if _, err := os.Stat(dstFile); err == nil {
		if err := os.Remove(dstFile); err != nil {
			fmt.Printf("删除已存在的目标文件失败: %s, 错误: %v\n", dstFile, err)
			return
		}
	}
	// 使用 os.Rename 移动文件
	if err := os.Rename(srcZip, dstFile); err == nil {
		fmt.Printf("成功移动文件: %s => %s\n", srcZip, dstFile)
		return
	} else {
		fmt.Println("移动文件失败", err)

	}
}
