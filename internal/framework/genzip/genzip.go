package genzip

import (
	"fmt"
	"github.com/energye/designer/pkg/tool"
	"path/filepath"
)

// ZipSRC 将指定根目录下的lcl、cef和wv文件夹分别压缩为ZIP文件
//
// 参数:
//   - root: 要压缩的文件夹的根目录路径
//
// 返回值:
//   - lcl: lcl文件夹压缩后的ZIP文件路径，如果压缩失败则为空字符串
//   - cef: cef文件夹压缩后的ZIP文件路径，如果压缩失败则为空字符串
//   - wv: wv文件夹压缩后的ZIP文件路径，如果压缩失败则为空字符串
func ZipSRC(root string) (lcl, cef, wv string) {
	skip := tool.NewHashSet[string]()
	skip.Add(".git")
	{
		src := filepath.Join(root, "lcl")
		desc := filepath.Join(root, "lcl.zip")
		if err := tool.ZipFolder(src, desc, skip); err == nil {
			lcl = desc
		} else {
			fmt.Println("zip lcl:", err)
		}
	}
	{
		src := filepath.Join(root, "cef")
		desc := filepath.Join(root, "cef.zip")
		tool.ZipFolder(src, desc, skip)
		if err := tool.ZipFolder(src, desc, skip); err == nil {
			cef = desc
		} else {
			fmt.Println("zip cef:", err)
		}
	}
	{
		src := filepath.Join(root, "wv")
		desc := filepath.Join(root, "wv.zip")
		tool.ZipFolder(src, desc, skip)
		if err := tool.ZipFolder(src, desc, skip); err == nil {
			wv = desc
		} else {
			fmt.Println("zip wv:", err)
		}
	}
	return
}
