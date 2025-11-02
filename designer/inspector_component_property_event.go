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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"unsafe"
)

// 设计 - 组件属性 - 事件

// 初始化组件属性树事件
func (m *TDesigningComponent) initComponentPropertyTreeEvent() {
	tree := m.propertyTree
	tree.SetOnScroll(func(sender lcl.IBaseVirtualTree, deltaX int32, deltaY int32) {
		tree.EndEditNode()
	})
	tree.SetOnExpanding(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, allowed *bool) {
		tree.EndEditNode()
		//data := vtedit.GetPropertyNodeData(node)
		//if data != nil {
		//	m.compPropTreeState.selectPropName = data.EditNodeData.Name
		//}
		tree.SetFocusedNode(node)
		sender.SetSelected(node, true)
		sender.ScrollIntoViewWithPVirtualNodeBoolX2(node, true, true)
	})
	tree.SetOnPaintText(func(sender lcl.IBaseVirtualTree, targetCanvas lcl.ICanvas, node types.PVirtualNode,
		column int32, textType types.TVSTTextType) {
		//logs.Debug("object inspector-property OnPaintText column:", column)
		if column == 0 {
			font := targetCanvas.FontToFont()
			font.SetStyle(font.Style().Include(types.FsBold))
			level := sender.GetNodeLevel(node)
			//logs.Info("  OnPaintText level:", level)
			switch level {
			case 0:
				font.SetColor(colors.RGBToColor(0, 32, 96))
			case 1:
				font.SetColor(colors.RGBToColor(0, 80, 239))
			case 2:
				font.SetColor(colors.RGBToColor(61, 133, 224))
			default:
				font.SetColor(colors.RGBToColor(100, 195, 255))
			}
		} else if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				//logs.Debug("object inspector-property OnPaintText column:", column, "IsModify:", data.IsModify())
				font := targetCanvas.FontToFont()
				// 编辑列 需要动态控制时
				switch data.EditNodeData.Type {
				case consts.PdtColorSelect:
					font.SetStyle(font.Style().Include(types.FsBold))
					font.SetColor(colors.TColor(data.EditNodeData.IntValue))
				case consts.PdtClass:
					// class 样式
					font.SetStyle(font.Style().Include(types.FsBold))
					font.SetColor(0x2D5BC4)
				default:
					if data.IsModify() {
						// 值被修改样式
						font.SetStyle(font.Style().Include(types.FsBold))
						font.SetColor(0x007DFF)
					}
				}
			}
		}
	})
	//tree.SetOnBeforeCellPaint(func(sender lcl.IBaseVirtualTree, targetCanvas lcl.ICanvas, node types.PVirtualNode,
	//	column int32, cellPaintMode types.TVTCellPaintMode, cellRect types.TRect, contentRect *types.TRect) {
	//	logs.Debug("[object inspector-property] OnBeforeCellPaint column:", column)
	//})
	//tree.SetOnAfterCellPaint(func(sender lcl.IBaseVirtualTree, targetCanvas lcl.ICanvas, node types.PVirtualNode,
	//	column int32, cellRect types.TRect) {
	//	logs.Debug("[object inspector-property] OnAfterCellPaint column:", column)
	//})
	tree.SetOnColumnClick(func(sender lcl.IBaseVirtualTree, column int32, shift types.TShiftState) {
		logs.Debug("[object inspector-property] OnColumnClick column:", column, "node:", sender.FocusedNode())
		tree.EndEditNode()
		node := sender.FocusedNode()
		data := vtedit.GetPropertyNodeData(node)
		if data != nil {
			if data.Type() == consts.PdtClass {
				switch data.EditNodeData.Name {
				case "Icon":
					//m.compPropTreeState.selectPropName = data.EditNodeData.Name
					tree.EditNode(node, 1)
				default:

				}
			} else if data.EditNodeData.Type == consts.PdtCheckBoxList {
			} else {
				//m.compPropTreeState.selectPropName = data.EditNodeData.Name
				tree.EditNode(node, 1)
			}
		}
	})
	tree.SetOnEditing(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32, allowed *bool) {
		logs.Debug("[object inspector-property] OnEditing column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil && data.EditNodeData.Type == consts.PdtText {
				*allowed = true
				return
			}
		}
	})
	//tree.SetOnEditCancelled(func(sender lcl.IBaseVirtualTree, column int32) {
	//	logs.Debug("[object inspector-property] OnEditCancelled column:", column)
	//})
	tree.SetOnCreateEditor(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32, outEditLink *lcl.IVTEditLink) {
		logs.Debug("[object inspector-property] OnCreateEditor column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				switch data.Type() {
				case consts.PdtText:
					link := vtedit.NewStringEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case consts.PdtInt, consts.PdtInt64:
					link := vtedit.NewIntEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case consts.PdtFloat:
					link := vtedit.NewFloatEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case consts.PdtCheckBox:
					link := vtedit.NewCheckBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case consts.PdtCheckBoxList:
					//link := vtedit.NewStringEditLink(data)
					//link.SetReadOnly(true)
					//*outEditLink = link.AsIVTEditLink()
				case consts.PdtComboBox:
					link := vtedit.NewComboBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case consts.PdtColorSelect:
					link := vtedit.NewColorSelectEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case consts.PdtClass:
					// class 属性实例, 根据属性名控制不同的操作
					propName := data.Name()
					switch propName {
					case "Icon":
						// 图标弹窗
						link := vtedit.NewIconEditLink(data)
						*outEditLink = link.AsIVTEditLink()
					default:

					}
					//link := vtedit.NewStringEditLink(data)
					//link.SetReadOnly(true)
					//*outEditLink = link.AsIVTEditLink()
				}
			}
		}
	})
	tree.SetOnEdited(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) {
		logs.Debug("[object inspector-property] OnEdited column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				// 从设计属性更新到组件属性
				data.FormInspectorPropertyToComponentProperty()
			}
		}
	})
	tree.SetOnExit(func(sender lcl.IObject) {
		logs.Debug("[object inspector-property] OnExit")
		tree.EndEditNode()
	})
	tree.SetOnGetText(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, textType types.TVSTTextType, cellText *string) {
		//logs.Debug("[object inspector-property] OnGetText column:", column)
		if data := vtedit.GetPropertyNodeData(node); data != nil {
			if column == 0 {
				*cellText = data.EditNodeData.Name
			} else if column == 1 {
				*cellText = data.EditStringValue()
			}
		}
	})
	tree.SetNodeDataSize(int32(unsafe.Sizeof(uintptr(0))))
}

func (m *TDesigningComponent) PropertyEndEdit() {
	m.propertyTree.EndEditNode()
}

func (m *TDesigningComponent) EventEndEdit() {
	m.eventTree.EndEditNode()
}
