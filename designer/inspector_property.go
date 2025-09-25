package designer

import (
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"log"
	"strconv"
	"unsafe"
)

// 设计 - 组件属性

type InspectorComponentProperty struct {
	box           lcl.IPanel                // 组件属性盒子
	filter        lcl.ITreeFilterEdit       // 组件属性过滤框
	page          lcl.IPageControl          // 属性和事件页
	propertySheet lcl.ITabSheet             // 属性页
	eventSheet    lcl.ITabSheet             // 事件页
	propertyTree  lcl.ILazVirtualStringTree // 组件属性
	eventTree     lcl.ILazVirtualStringTree // 组件事件
}

func (m *InspectorComponentProperty) init(leftBoxWidth int32) {
	componentPropertyTitle := lcl.NewLabel(m.box)
	componentPropertyTitle.SetParent(m.box)
	componentPropertyTitle.SetCaption("属性")
	componentPropertyTitle.Font().SetStyle(types.NewSet(types.FsBold))
	componentPropertyTitle.SetTop(5)
	componentPropertyTitle.SetLeft(5)

	m.filter = lcl.NewTreeFilterEdit(m.box)
	m.filter.SetParent(m.box)
	m.filter.SetTop(2)
	m.filter.SetLeft(30)
	m.filter.SetWidth(leftBoxWidth - m.filter.Left())
	m.filter.SetAlign(types.AlCustom)
	m.filter.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkRight))
	// 选项卡
	{
		m.page = lcl.NewPageControl(m.box)
		m.page.SetParent(m.box)
		m.page.SetTabStop(true)
		m.page.SetTop(32)
		m.page.SetWidth(leftBoxWidth)
		m.page.SetHeight(m.box.Height() - m.page.Top())
		m.page.SetAlign(types.AlCustom)
		m.page.SetAnchors(types.NewSet(types.AkLeft, types.AkTop, types.AkBottom, types.AkRight))

		m.propertySheet = lcl.NewTabSheet(m.page)
		m.propertySheet.SetParent(m.page)
		m.propertySheet.SetCaption("  属性  ")
		m.propertySheet.SetAlign(types.AlClient)

		m.eventSheet = lcl.NewTabSheet(m.page)
		m.eventSheet.SetParent(m.page)
		m.eventSheet.SetCaption("  事件  ")
		m.eventSheet.SetAlign(types.AlClient)
	}
	// 组件的属性列表和事件列表"树"
	{
		// 组件属性树列表
		m.propertyTree = lcl.NewLazVirtualStringTree(m.propertySheet)
		m.propertyTree.SetParent(m.propertySheet)
		m.propertyTree.SetBorderStyleToBorderStyle(types.BsNone)
		m.propertyTree.SetAlign(types.AlClient)
		m.propertyTree.SetLineStyle(types.LsSolid)
		m.propertyTree.SetDefaultNodeHeight(28)
		m.propertyTree.SetIndent(8)
		propTreeOptions := m.propertyTree.TreeOptions()
		propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Exclude(types.ToShowTreeLines))
		propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Include(types.ToShowVertGridLines, types.ToShowHorzGridLines))
		propTreeOptions.SetSelectionOptions(propTreeOptions.SelectionOptions().Include(types.ToLevelSelectConstraint))
		propTreeOptions.SetMiscOptions(propTreeOptions.MiscOptions().Include(types.ToEditable, types.ToEditOnClick, types.ToEditOnDblClick))
		propColors := m.propertyTree.Colors()
		propColors.SetFocusedSelectionColor(colors.RGBToColor(43, 169, 241))
		propColors.SetUnfocusedSelectionColor(colors.RGBToColor(43, 169, 241))

		// 组件事件树列表
		m.eventTree = lcl.NewLazVirtualStringTree(m.eventSheet)
		m.eventTree.SetParent(m.eventSheet)
		m.eventTree.SetBorderStyleToBorderStyle(types.BsNone)
		m.eventTree.SetAlign(types.AlClient)
		eventTreeOptions := m.propertyTree.TreeOptions()
		eventTreeOptions.SetPaintOptions(eventTreeOptions.PaintOptions().Exclude(types.ToShowTreeLines))
		eventTreeOptions.SetPaintOptions(eventTreeOptions.PaintOptions().Include(types.ToShowVertGridLines))
		eventTreeOptions.SetPaintOptions(eventTreeOptions.PaintOptions().Include(types.ToShowHorzGridLines))
		eventTreeOptions.SetSelectionOptions(eventTreeOptions.SelectionOptions().Include(types.ToFullRowSelect))

	}
	// 初始化组件属性树
	m.initComponentPropertyTree()

	// 测试
	{
		data := &vtedit.TEditLinkNodeData{Type: vtedit.PdtText, Name: "TextEdit", StringValue: "Value"}
		vtedit.AddPropertyNodeData(m.propertyTree, 0, data)

		data = &vtedit.TEditLinkNodeData{Type: vtedit.PdtInt, Name: "IntEdit", IntValue: 1}
		vtedit.AddPropertyNodeData(m.propertyTree, 0, data)

		data = &vtedit.TEditLinkNodeData{Type: vtedit.PdtFloat, Name: "FloatEdit", FloatValue: 1.99}
		vtedit.AddPropertyNodeData(m.propertyTree, 0, data)

		data = &vtedit.TEditLinkNodeData{Type: vtedit.PdtCheckBox, Name: "CheckBox", Checked: true}
		vtedit.AddPropertyNodeData(m.propertyTree, 0, data)

		data = &vtedit.TEditLinkNodeData{Type: vtedit.PdtCheckBoxList, Name: "Anchors", BoolValue: true, StringValue: "",
			CheckBoxValue: []*vtedit.TEditLinkNodeData{{Type: vtedit.PdtCheckBox, Name: "Value1", Checked: true}, {Type: vtedit.PdtCheckBox, Name: "Value2", Checked: false}}}
		vtedit.AddPropertyNodeData(m.propertyTree, 0, data)

		data = &vtedit.TEditLinkNodeData{Type: vtedit.PdtComboBox, Name: "CombBox", StringValue: "Value1",
			ComboBoxValue: []*vtedit.TEditLinkNodeData{{StringValue: "Value1"}, {StringValue: "Value2"}}}
		vtedit.AddPropertyNodeData(m.propertyTree, 0, data)

		//node = m.propertyTree.AddChild(0, 0)
		//data = &vtedit.TEditLinkNodeData{Type: vtedit.PdtCheckBoxDraw, Name: "CheckBoxDraw", BoolValue: true,
		//	CheckBoxValue: []vtedit.TEditLinkNodeData{{Name: "Value1", BoolValue: true}, {Name: "Value2", BoolValue: false}}}
		//AddPropertyNodeData(node, data)
	}
}

// 初始化组件属性树
func (m *InspectorComponentProperty) initComponentPropertyTree() {
	header := m.propertyTree.Header()
	header.SetOptions(header.Options().Include(types.HoVisible, types.HoAutoSpring)) //types.HoAutoResize
	header.Font().SetStyle(header.Font().Style().Include(types.FsBold))
	header.Font().SetColor(colors.ClGray)
	columns := header.Columns()
	columns.Clear()
	propNameCol := columns.AddToVirtualTreeColumn()
	propNameCol.SetText("名")
	propNameCol.SetWidth(100)
	propNameCol.SetAlignment(types.TaLeftJustify)
	//propNameCol.SetOptions(propNameCol.Options().Include(types.CoDisableAnimatedResize))

	propValueCol := columns.AddToVirtualTreeColumn()
	propValueCol.SetText("值")
	//propValueCol.SetWidth(leftBoxWidth - 150)
	propValueCol.SetWidth(leftBoxWidth - 100)
	propValueCol.SetAlignment(types.TaLeftJustify)
	propValueCol.SetOptions(propValueCol.Options().Include(types.CoAutoSpring))

	m.propertyTree.SetOnPaintText(func(sender lcl.IBaseVirtualTree, targetCanvas lcl.ICanvas, node types.PVirtualNode,
		column int32, textType types.TVSTTextType) {
		//log.Println("property-inspector OnPaintText column:", column)
		if column == 0 {
			font := targetCanvas.FontToFont()
			font.SetStyle(font.Style().Include(types.FsBold))
			level := sender.GetNodeLevel(node)
			//log.Println("  OnPaintText level:", level)
			switch level {
			case 0:
				font.SetColor(colors.ClBlack)
			case 1:
				font.SetColor(colors.ClBlue)
			default:
				font.SetColor(colors.ClGreen)
			}
		}
	})
	//m.propertyTree.SetOnBeforeCellPaint(func(sender lcl.IBaseVirtualTree, targetCanvas lcl.ICanvas, node types.PVirtualNode,
	//	column int32, cellPaintMode types.TVTCellPaintMode, cellRect types.TRect, contentRect *types.TRect) {
	//	log.Println("[property-inspector] OnBeforeCellPaint column:", column)
	//})
	//m.propertyTree.SetOnAfterCellPaint(func(sender lcl.IBaseVirtualTree, targetCanvas lcl.ICanvas, node types.PVirtualNode,
	//	column int32, cellRect types.TRect) {
	//	log.Println("[property-inspector] OnAfterCellPaint column:", column)
	//})
	m.propertyTree.SetOnColumnClick(func(sender lcl.IBaseVirtualTree, column int32, shift types.TShiftState) {
		// edit: 1. 触发编辑
		log.Println("[property-inspector] OnColumnClick column:", column)
		if column == 1 {
			node := sender.FocusedNode()
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				m.propertyTree.EditNode(node, column)
			}
		}
	})
	m.propertyTree.SetOnEditing(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, allowed *bool) {
		// edit: 2. 第二列可以编辑
		log.Println("[property-inspector] OnEditing column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil && data.Type == vtedit.PdtText {
				*allowed = true
				return
			}
		}
	})
	m.propertyTree.SetOnEditCancelled(func(sender lcl.IBaseVirtualTree, column int32) {
		log.Println("[property-inspector] OnEditCancelled column:", column)
	})
	m.propertyTree.SetOnEdited(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) {
		// edit: 4. 编辑结束
		log.Println("[property-inspector] OnEdited column:", column)
	})
	m.propertyTree.SetOnCreateEditor(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, outEditLink *lcl.IVTEditLink) {
		// edit: 3. 创建编辑或组件
		log.Println("[property-inspector] OnCreateEditor column:", column)
		if column == 1 {
			if data := vtedit.GetPropertyNodeData(node); data != nil {
				switch data.Type {
				case vtedit.PdtText:
					link := vtedit.NewStringEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtInt:
					link := vtedit.NewIntEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtFloat:
					link := vtedit.NewFloatEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtCheckBox:
					link := vtedit.NewCheckBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtCheckBoxList:
					link := vtedit.NewCheckBoxListEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				case vtedit.PdtComboBox:
					link := vtedit.NewComboBoxEditLink(data)
					*outEditLink = link.AsIVTEditLink()
				}
			}
		}
	})
	m.propertyTree.SetOnGetText(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, textType types.TVSTTextType, cellText *string) {
		//log.Println("[property-inspector] OnGetText column:", column)
		if data := vtedit.GetPropertyNodeData(node); data != nil {
			if column == 0 {
				*cellText = data.Name
			} else if column == 1 {
				switch data.Type {
				case vtedit.PdtText:
					*cellText = data.StringValue
				case vtedit.PdtInt:
					*cellText = strconv.Itoa(data.IntValue)
				case vtedit.PdtFloat:
					val := strconv.FormatFloat(data.FloatValue, 'f', 2, 64)
					*cellText = val
				case vtedit.PdtCheckBox:
					*cellText = strconv.FormatBool(data.Checked)
				case vtedit.PdtCheckBoxList:
					*cellText = data.StringValue
				case vtedit.PdtComboBox:
					*cellText = data.StringValue
				default:
					*cellText = ""
				}
			}
		}
	})
	m.propertyTree.SetNodeDataSize(int32(unsafe.Sizeof(uintptr(0))))
}
