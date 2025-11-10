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
	"bytes"
	"encoding/json"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/resources"
)

// 弹窗过滤
var DialogFilter = &dialogFilter{}

type dialogFilter struct {
	Image []string `json:"image"`
	File  []string `json:"file"`
	UI    []string `json:"ui"`
}

func init() {
	data := resources.DialogFilter()
	if data == nil {
		logs.Error("加载弹窗过滤配置为nil")
		return
	}
	err := json.Unmarshal(data, DialogFilter)
	if data == nil {
		logs.Error("加载弹窗过滤配置错误:", err.Error())
		return
	}
}

func (m *dialogFilter) ImageFilter() string {
	buf := bytes.Buffer{}
	for i, item := range m.Image {
		if i > 0 {
			buf.WriteString("|")
		}
		buf.WriteString(item)
	}
	return buf.String()
}

func (m *dialogFilter) FileFilter() string {
	buf := bytes.Buffer{}
	for i, item := range m.File {
		if i > 0 {
			buf.WriteString("|")
		}
		buf.WriteString(item)
	}
	return buf.String()
}
