package config

import (
	"encoding/json"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/resources"
)

type config struct {
	Title         string        `json:"title"`
	Version       string        `json:"version"`
	Window        Window        `json:"window"`
	ComponentTabs ComponentTabs `json:"componentTabs"`
}

type Window struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ComponentTabs struct {
	Standard   Tab `json:"standard"`
	Additional Tab `json:"additional"`
	Common     Tab `json:"common"`
	Dialogs    Tab `json:"dialogs"`
	Misc       Tab `json:"misc"`
	System     Tab `json:"system"`
	LazControl Tab `json:"lazcontrol"`
	WebView    Tab `json:"webview"`
}

type Tab struct {
	Cn        string   `json:"cn"`
	En        string   `json:"en"`
	Component []string `json:"component"`
}

var Config *config

func init() {
	Config = &config{}
	err.CheckErr(json.Unmarshal(resources.Config(), Config))
}
