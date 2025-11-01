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
