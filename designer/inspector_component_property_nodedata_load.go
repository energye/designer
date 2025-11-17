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

package designer

import (
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
)

// 设计 - 组件的设计属性和设计事件数据加载

// 加载组件属性列表
func (m *TDesigningComponent) loadPropertyList() {
	if m.isLoadProperty {
		// 加载完 不在继续加载
		return
	}
	m.isLoadProperty = true
	tempPropertyMap := make(map[string]struct{}) // 用于下面判断
	configCompProp := config.ComponentProperty
	for i, nodeData := range m.PropertyList {
		// 通用属性, 排除的属性
		if configCompProp.IsExclude(nodeData.EditNodeData.Name) {
			logs.Debug("排除属性:", nodeData.EditNodeData.Metadata.ToJSON())
			continue
		}
		tempPropertyMap[nodeData.Name()] = struct{}{}
		logs.Debug("加载属性:", nodeData.EditNodeData.Metadata.ToJSON())
		if !nodeData.IsFinal {
			// 自定义属性, 使用会覆蓋掉
			// 返回数组
			if customProps := configCompProp.GetCustomPropertyList(nodeData.EditNodeData.Name); customProps != nil {
				if len(customProps) == 1 {
					// 数组只有一个元素，规则为直接作用在当前属性上
					customProp := &customProps[0]
					// 默认值设置
					if customProp.Options == "" {
						// 未配置options时,使用默认值
						customProp.Options = nodeData.EditNodeData.Metadata.Options
					} else {
						// TODO 其它元数据的字段默认值??
					}
					customProperty := vtedit.NewEditLinkNodeData(customProp)
					newEditNodeData := &vtedit.TEditNodeData{IsFinal: true, EditNodeData: customProperty,
						OriginNodeData: customProperty.Clone(), AffiliatedComponent: m}
					m.PropertyList[i] = newEditNodeData // 更新到组件属性
					nodeData = m.PropertyList[i]
					newEditNodeData.Build()
				} else {
					// 自定义属性添加？？
				}
			}
			nodeData.IsFinal = true
		}
		// 属性节点数据添加到树
		vtedit.AddPropertyNodeData(m.propertyTree, 0, nodeData)
	}
	// 通用属性 包含属性, 如果重复忽略配置的
	for _, prop := range configCompProp.Include() {
		if _, ok := tempPropertyMap[prop.Name]; ok {
			// 忽略配置属性
			continue
		}
		customProperty := vtedit.NewEditLinkNodeData(&prop)
		newEditNodeData := &vtedit.TEditNodeData{IsFinal: true, EditNodeData: customProperty,
			OriginNodeData: customProperty.Clone(), AffiliatedComponent: m}
		m.PropertyList = append(m.PropertyList, newEditNodeData) // 添加到组件属性
		newEditNodeData.Build()
	}
}

// 加载组件事件列表
func (m *TDesigningComponent) loadEventList() {

}
