package resources

import "embed"

var (
	//go:embed config.json
	config embed.FS
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
