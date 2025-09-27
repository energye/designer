package resources

import "embed"

var (
	//go:embed config.json
	config embed.FS
	//go:embed component_property.json
	componentProperty embed.FS
	//go:embed images
	images embed.FS
)

// 配置文件
func Config() []byte {
	if d, err := config.ReadFile("config.json"); err == nil {
		return d
	}
	return nil
}

// 静态资源
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
	if d, err := componentProperty.ReadFile("component_property.json"); err == nil {
		return d
	}
	return nil
}
