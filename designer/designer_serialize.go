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
// See the License for the specific language governing permissions and
// limitations under the License.

package designer

import (
	"strings"

	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	uiBean "github.com/energye/designer/uigen/bean"
)

// 设计组件子树序列化(.ui 同源数据结构):
// .ui 布局文件生成(uigen)、复制/粘贴剪贴板、撤销/恢复快照共用此实现。

// 默认生成属性
var defaultGenerateProperty = map[string]struct{}{
	"Left": {}, "Top": {}, "Width": {}, "Height": {}, "Caption": {}, "Name": {},
}

// BuildUIComponentTree 构建组件子树的 UI 数据结构:
// 属性列表非空的组件序列化默认属性(Left/Top/Width/Height/Caption/Name)
// 与已修改属性(IsModify), 递归子组件。
func BuildUIComponentTree(component *TDesigningComponent) uiBean.TUIComponent {
	uiComp := uiBean.TUIComponent{
		Name:       component.Name(),
		ClassName:  component.ClassName(),
		Mod:        component.Mod(),
		Properties: make([]uiBean.TProperty, 0),
		Child:      make([]uiBean.TUIComponent, 0),
		Type:       component.ComponentType,
	}

	// 获取变更的属性
	if component.PropertyList != nil {
		for _, prop := range component.PropertyList {
			// 默认生成的属性 Left Top Width Height
			if _, defGenOk := defaultGenerateProperty[prop.Name()]; defGenOk {
				// 忽略非可视化控件 Left Top, 不生成代码
				noCode := component.ComponentType == consts.CtNonVisual && tool.Equal(prop.Name(), "Left", "Top")
				propName := prop.Name()
				propValue := prop.EditValue()
				uiComp.Properties = append(uiComp.Properties, uiBean.TProperty{Name: propName, Value: propValue,
					Type: prop.Type(), NoCode: noCode})
			} else {
				// 只保存修改过的属性
				switch prop.Type() {
				case consts.PdtClass:
					var iterator func(node *vtedit.TEditNodeData)
					iterator = func(node *vtedit.TEditNodeData) {
						for _, data := range node.Child {
							if data.IsModify() {
								if data.Type() == consts.PdtClass {
									iterator(data)
								} else {
									paths := data.Paths()
									if paths != nil {
										tool.StringArrayReverse(paths)
										paths = append(paths, data.Name())
										propName := strings.Join(paths, ".")
										propValue := data.EditValue()
										uiComp.Properties = append(uiComp.Properties, uiBean.TProperty{Name: propName, Value: propValue,
											Type: data.Type()})
									} else {
										logs.Error("错误, 生成UI布局文件, 属性是 class 获取子节点路径错误 nil. 属性名: ", prop.Name(), "子节点属性名:", data.Name())
									}
								}
							}
						}
					}
					iterator(prop)
				default:
					if prop.IsModify() {
						propName := prop.Name()
						propValue := prop.EditValue()
						uiComp.Properties = append(uiComp.Properties, uiBean.TProperty{Name: propName, Value: propValue, Type: prop.Type()})
					}
				}
			}
		}
	}

	if component.EventList != nil {
		for _, event := range component.EventList {
			if event.IsModify() {
				eventName := event.Name()
				eventValue := event.EditValue()
				uiComp.Properties = append(uiComp.Properties, uiBean.TProperty{Name: eventName, Value: eventValue, Type: event.Type()})
			}
		}
	}

	// 递归处理子组件
	for _, child := range component.Child {
		uiComp.Child = append(uiComp.Child, BuildUIComponentTree(child))
	}

	return uiComp
}

// ExportUIComponentTree 序列化设计组件及其子树为 UI 数据结构
// (复制/粘贴的剪贴板内容与撤销/恢复的操作快照共用此序列化)。
// 属性列表未加载的组件先加载(GetProps), 保证 Left/Top/Width/Height 等
// 默认属性被捕获 —— BuildUIComponentTree 只序列化默认属性与已修改属性。
func ExportUIComponentTree(component *TDesigningComponent) uiBean.TUIComponent {
	var ensureProps func(comp *TDesigningComponent)
	ensureProps = func(comp *TDesigningComponent) {
		if comp == nil {
			return
		}
		if comp.PropertyList == nil {
			comp.GetProps()
		}
		for _, child := range comp.Child {
			ensureProps(child)
		}
	}
	ensureProps(component)
	return BuildUIComponentTree(component)
}
