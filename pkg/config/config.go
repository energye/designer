package config

import (
	"encoding/json"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/resources"
)

type config struct {
	Title   string `json:"title"`
	Version string `json:"version"`
	Window  window `json:"window"`
}

type window struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

var Config *config

func init() {
	Config = &config{}
	err.CheckErr(json.Unmarshal(resources.Config(), Config))
}
