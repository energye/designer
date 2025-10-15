package resources

import (
	"embed"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/logs"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	// 测试
	//go:embed lib/liblcl.dll
	lib embed.FS
	// 配置
	//go:embed config.json
	config embed.FS
	// 组件属性配置
	//go:embed component-property.json
	componentProperty embed.FS
	// 图标资源
	//go:embed images
	images embed.FS
)

var (
	LibPath string
)

// 配置文件
func Config() []byte {
	if d, err := config.ReadFile("config.json"); err == nil {
		return d
	}
	return nil
}

// 图标资源
func Images(filePath string) []byte {
	if d, err := images.ReadFile("images/" + filePath); err == nil {
		return d
	}
	return nil
}

// 组件属性配置
// 用于通用属性和定制属性
// 通用属性 1. 排除 2. 包含
// 定制属性 组件特有的属性
func ComponentProperty() []byte {
	if d, err := componentProperty.ReadFile("component-property.json"); err == nil {
		return d
	}
	return nil
}

func init() {
	tempDir := os.TempDir()
	outPath := filepath.Join(tempDir, "lib-energy.dll")
	libByte, e := lib.ReadFile("lib/liblcl.dll")
	err.CheckErr(e)
	os.WriteFile(outPath, libByte, fs.ModePerm)
	LibPath = outPath
	logs.Info("Lib Path:", outPath)
}
