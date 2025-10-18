package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"unsafe"
)

// 初始化组件属性树
func (m *InspectorComponentProperty) initComponentPropertyTreeEvent() {
	tree := m.propertyTree
	tree.SetOnScroll(func(sender lcl.IBaseVirtualTree, deltaX int32, deltaY int32) {
		tree.EndEditNode()
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
				font.SetColor(colors.ClBlack)
			case 1:
				font.SetColor(colors.ClBlue)
			default:
				font.SetColor(colors.ClGreen)
			}
		} else if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				//logs.Debug("object inspector-property OnPaintText column:", column, "IsModify:", data.IsModify())
				font := targetCanvas.FontToFont()
				// 编辑列 需要动态控制时
				switch data.EditNodeData.Type {
				case vtedit.PdtColorSelect:
					font.SetStyle(font.Style().Include(types.FsBold))
					font.SetColor(colors.TColor(data.EditNodeData.IntValue))
				case vtedit.PdtClass:
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
		node := sender.FocusedNode()
		data := vtedit.GetPropertyNodeData(node)
		if data != nil {
			if data.EditNodeData.Type == vtedit.PdtClass {
				switch data.EditNodeData.Name {
				case "Icon":
					m.currentComponent.compPropTreeState.selectPropName = data.EditNodeData.Name
					tree.EditNode(node, 1)
				default:

				}
			} else if data.EditNodeData.Type == vtedit.PdtCheckBoxList {
			} else {
				m.currentComponent.compPropTreeState.selectPropName = data.EditNodeData.Name
				tree.EditNode(node, 1)
			}
		}
	})
	tree.SetOnEditing(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32, allowed *bool) {
		logs.Debug("[object inspector-property] OnEditing column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil && data.EditNodeData.Type == vtedit.PdtText {
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
				switch data.EditNodeData.Type {
				case vtedit.PdtText:
					link := vtedit.NewStringEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtInt, vtedit.PdtInt64:
					link := vtedit.NewIntEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtFloat:
					link := vtedit.NewFloatEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtCheckBox:
					link := vtedit.NewCheckBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtCheckBoxList:
					//link := vtedit.NewStringEditLink(data)
					//link.SetReadOnly(true)
					//*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtComboBox:
					link := vtedit.NewComboBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtColorSelect:
					link := vtedit.NewColorSelectEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtClass:
					propName := data.EditNodeData.Name
					switch propName {
					case "Icon":
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
				*cellText = data.EditValue()
			}
			if data.EditNodeData.Name == m.currentComponent.compPropTreeState.selectPropName {
				m.currentComponent.compPropTreeState.selectNode = node
			}
		}
	})
	tree.SetNodeDataSize(int32(unsafe.Sizeof(uintptr(0))))
}

func (m *InspectorComponentProperty) PropertyEndEdit() {
	m.propertyTree.EndEditNode()
}

func (m *InspectorComponentProperty) EventEndEdit() {
	m.eventTree.EndEditNode()
}
