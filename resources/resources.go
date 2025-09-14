package resources

import "embed"

var (
	//go:embed config.json
	config embed.FS
	//go:embed assets
	assets embed.FS
)

// 配置文件
func Config() []byte {
	if d, err := config.ReadFile("config.json"); err == nil {
		return d
	}
	return nil
}

// 静态资源
func Assets(fileName string) []byte {
	if d, err := assets.ReadFile("assets/" + fileName); err == nil {
		return d
	}
	return nil
}
