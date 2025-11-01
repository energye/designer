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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
)

// 组件属性配置
type componentProperty struct {
	Common             common       `json:"common"` // 通用属性配置
	CustomPropertyList propertyList `json:"custom"` // 自定义属性列表配置
}

// 通用配置
type common struct {
	Exclude []string                  `json:"exclude"` // 排除的属性
	Include []lcl.ComponentProperties `json:"include"` // 包含的属性
}

// 自定义属性列表
type propertyList map[string][]lcl.ComponentProperties

// 获取指定组件的自定义属性配置
func (m propertyList) Get(componentName string) []lcl.ComponentProperties {
	if info, ok := m[componentName]; ok {
		return info
	}
	return nil
}

var ComponentProperty *componentProperty

func init() {
	ComponentProperty = &componentProperty{}
	err.CheckErr(json.Unmarshal(resources.ComponentProperty(), ComponentProperty))
}

// 通用属性 是否排除的属性
func (m *componentProperty) IsExclude(propertyName string) bool {
	for _, name := range m.Common.Exclude {
		if tool.Equal(name, propertyName) {
			return true
		}
	}
	return false
}

// 通用属性 获取包含属性
func (m *componentProperty) Include() (propertyList []lcl.ComponentProperties) {
	return m.Common.Include
}

// 获取自定义组件属性
func (m *componentProperty) GetCustomPropertyList(componentName string) []lcl.ComponentProperties {
	return m.CustomPropertyList.Get(componentName)
}

// 删除自定义组件属性
func (m *componentProperty) DeleteCustomPropertyList(componentName string) {
	delete(m.CustomPropertyList, componentName)
}
