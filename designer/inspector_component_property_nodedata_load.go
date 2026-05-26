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
	m.UpdatePropertyTreeWidth()
	if m.isLoadProperty {
		// 加载完 不在继续加载
		return
	}
	m.isLoadProperty = true
	configCompProp := config.ComponentProperty
	for _, nodeData := range m.PropertyList {
		// 通用属性, 排除的属性
		if configCompProp.IsExclude(nodeData.EditNodeData.Name) {
			logs.Debug("排除属性:", nodeData.EditNodeData.Metadata.ToJSON())
			continue
		}
		//logs.Debug("加载属性:", nodeData.EditNodeData.Metadata.ToJSON())
		// 属性节点数据添加到树
		vtedit.AddPropertyNodeData(m.propertyTree, 0, nodeData)
	}
}

// 加载组件事件列表
func (m *TDesigningComponent) loadEventList() {
	m.UpdateEventTreeWidth()
	if m.isLoadEvent {
		// 加载完 不在继续加载
		return
	}
	m.isLoadEvent = true
	configCompProp := config.ComponentProperty
	for _, nodeData := range m.EventList {
		// 通用属性, 排除的属性
		if configCompProp.IsExclude(nodeData.EditNodeData.Name) {
			logs.Debug("排除属性:", nodeData.EditNodeData.Metadata.ToJSON())
			continue
		}
		//logs.Debug("加载属性:", nodeData.EditNodeData.Metadata.ToJSON())
		// 属性节点数据添加到树
		vtedit.AddPropertyNodeData(m.eventTree, 0, nodeData)
	}
}

var (
	gDefaultPropertyNameColumnTreeWidth = int32(-1)
	gDefaultEventNameColumnTreeWidth    = int32(-1)
	gSwitchDefaultTreeTabPage           = int32(0)
)

func (m *TDesigningComponent) SwitchDefaultTreeTabPage() {
	m.page.SetActivePageIndex(gSwitchDefaultTreeTabPage)
}

func (m *TDesigningComponent) UpdatePropertyTreeWidth() {
	if gDefaultPropertyNameColumnTreeWidth == -1 {
		gDefaultPropertyNameColumnTreeWidth = config.Config.WindowLayout.ContentLayout.InspectorLayout.PropertyTreeWidth
	}
	columns := m.propertyTree.Header().Columns()
	nameColumn := columns.ItemsWithColumnIndexToVirtualTreeColumn(0)
	valueColumn := columns.ItemsWithColumnIndexToVirtualTreeColumn(1)
	//cn := m.propertyTree.ClientWidth() - valueColumn.Width()
	nameColumn.SetWidth(gDefaultPropertyNameColumnTreeWidth)
	cv := m.propertyTree.ClientWidth() - nameColumn.Width()
	valueColumn.SetWidth(cv)
}

func (m *TDesigningComponent) UpdateEventTreeWidth() {
	if gDefaultEventNameColumnTreeWidth == -1 {
		gDefaultEventNameColumnTreeWidth = config.Config.WindowLayout.ContentLayout.InspectorLayout.EventTreeWidth
	}
	columns := m.eventTree.Header().Columns()
	nameColumn := columns.ItemsWithColumnIndexToVirtualTreeColumn(0)
	valueColumn := columns.ItemsWithColumnIndexToVirtualTreeColumn(1)
	//cn := m.eventTree.ClientWidth() - valueColumn.Width()
	nameColumn.SetWidth(gDefaultEventNameColumnTreeWidth)
	cv := m.eventTree.ClientWidth() - nameColumn.Width()
	valueColumn.SetWidth(cv)
}
