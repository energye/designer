package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strconv"
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
				logs.Debug("object inspector-property OnPaintText column:", column, "IsModify:", data.IsModify())
				font := targetCanvas.FontToFont()
				// 编辑列 需要动态控制时
				switch data.EditNodeData.Type {
				case vtedit.PdtColorSelect:
					font.SetStyle(font.Style().Include(types.FsBold))
					font.SetColor(colors.TColor(data.EditNodeData.IntValue))
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
		// edit: 1. 触发编辑
		logs.Debug("[object inspector-property] OnColumnClick column:", column)
		if column == 1 {
			node := sender.FocusedNode()
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				tree.EditNode(node, column)
			}
		}
	})
	tree.SetOnEditing(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, allowed *bool) {
		// edit: 2. 第二列可以编辑
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
	tree.SetOnEdited(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) {
		// edit: 4. 编辑结束
		logs.Debug("[object inspector-property] OnEdited column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				go data.UpdateComponentProperties()
			}
		}
	})
	tree.SetOnExit(func(sender lcl.IObject) {
		logs.Debug("[object inspector-property] OnExit")
		tree.EndEditNode()
	})
	tree.SetOnCreateEditor(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, outEditLink *lcl.IVTEditLink) {
		// edit: 3. 创建编辑或组件
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
				case vtedit.PdtCheckBoxList, vtedit.PdtClass:
					link := vtedit.NewStringEditLink(data)
					link.SetReadOnly(true)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtComboBox:
					link := vtedit.NewComboBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtColorSelect:
					link := vtedit.NewColorSelectEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				}
			}
		}
	})
	tree.SetOnGetText(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, textType types.TVSTTextType, cellText *string) {
		//logs.Debug("[object inspector-property] OnGetText column:", column)
		if data := vtedit.GetPropertyNodeData(node); data != nil {
			if column == 0 {
				*cellText = data.EditNodeData.Name
			} else if column == 1 {
				dataType := data.EditNodeData.Type
				switch dataType {
				case vtedit.PdtText:
					*cellText = data.EditNodeData.StringValue
				case vtedit.PdtInt, vtedit.PdtInt64:
					*cellText = strconv.Itoa(data.EditNodeData.IntValue)
				case vtedit.PdtFloat:
					val := strconv.FormatFloat(data.EditNodeData.FloatValue, 'f', 2, 64)
					*cellText = val
				case vtedit.PdtCheckBox:
					*cellText = strconv.FormatBool(data.EditNodeData.Checked)
				case vtedit.PdtCheckBoxList:
					*cellText = data.EditNodeData.StringValue
				case vtedit.PdtComboBox:
					*cellText = data.EditNodeData.StringValue
				case vtedit.PdtColorSelect:
					*cellText = fmt.Sprintf("0x%X", data.EditNodeData.IntValue)
				default:
					*cellText = ""
				}
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
