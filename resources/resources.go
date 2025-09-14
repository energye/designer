package resources

import "embed"

//go:embed config.json
var config embed.FS

func Config() []byte {
	if d, err := config.ReadFile("config.json"); err == nil {
		return d
	}
	return nil
}
