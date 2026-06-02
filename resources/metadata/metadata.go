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

package metadata

import (
	"embed"
	"strings"
)

var (
	//go:embed i18n
	i18n  embed.FS
	GI18n = &TI18n{dict: make(map[string]*value)}
)

type TI18n struct {
	Lang string
	dict map[string]*value
}

type value struct {
	value    string
	variable *string
}

func (m *TI18n) Bind(name string, variable *string) {
	name = strings.ToLower(name)
	if val, ok := m.dict[name]; ok {
		val.variable = variable
		*val.variable = val.value
	} else {
		val = &value{variable: variable}
		m.dict[name] = val
	}
}

func (m *TI18n) Get(lang string) ([]byte, error) {
	if lang == "" {
		lang = "zh-CN"
	}
	data, err := i18n.ReadFile("i18n/" + lang + ".ini")
	if err == nil && m.Lang != lang {
		m.Lang = lang
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || line[0] == '#' {
				continue
			}
			kv := strings.Split(line, "=")
			if len(kv) != 2 {
				continue
			}
			name := strings.TrimSpace(strings.ToLower(kv[0]))
			if val, ok := m.dict[name]; ok {
				val.value = kv[1]
				if val.variable != nil {
					*val.variable = val.value
				}
			} else {
				val = &value{value: kv[1]}
				m.dict[name] = val
			}
		}
	}
	return data, err
}

func (m *TI18n) Dict(name string) string {
	name = strings.TrimSpace(name)
	if m.dict == nil || name == "" {
		return ""
	}
	name = strings.ToLower(name)
	if v, ok := m.dict[name]; ok {
		return v.value
	}
	return ""
}

func (m *TI18n) Iterate(fn func(name, value string) bool) {
	if fn == nil {
		return
	}
	for name, val := range m.dict {
		if fn(name, val.value) {
			break
		}
	}
}
