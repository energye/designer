package designer

import (
	"github.com/energye/designer/pkg/vtedit"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
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
		m.propertyTree.SetIndent(0)
		propTreeOptions := m.propertyTree.TreeOptions()
		propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Exclude(types.ToShowTreeLines))
		propTreeOptions.SetPaintOptions(propTreeOptions.PaintOptions().Include(types.ToShowVertGridLines, types.ToShowHorzGridLines))
		propTreeOptions.SetSelectionOptions(propTreeOptions.SelectionOptions().Include(types.ToFullRowSelect, types.ToAlwaysSelectNode))
		propTreeOptions.SetMiscOptions(propTreeOptions.MiscOptions().Include(types.ToEditable, types.ToEditOnClick, types.ToEditOnDblClick))

		// 组件事件树列表
		m.eventTree = lcl.NewLazVirtualStringTree(m.eventSheet)
		m.eventTree.SetParent(m.eventSheet)
		m.eventTree.SetBorderStyleToBorderStyle(types.BsNone)
		m.eventTree.SetAlign(types.AlClient)
		m.eventTree.SetNodeDataSize(int32(unsafe.Sizeof(TTreePropertyNodeData{})))
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
		//m.propertyTree.SetRootNodeCount(m.propertyTree.RootNodeCount() + uint32(5))
		for i := 0; i < 5; i++ {
			node := m.propertyTree.AddChild(0, 0)
			treePropertyNodeDatas[node] = &TTreePropertyNodeData{name: "Name" + strconv.Itoa(i), value: "Value" + strconv.Itoa(i)}
		}
	}
}

// 初始化组件属性树
func (m *InspectorComponentProperty) initComponentPropertyTree() {
	header := m.propertyTree.Header()
	header.SetOptions(header.Options().Include(types.HoVisible, types.HoAutoSpring)) //types.HoAutoResize
	columns := header.Columns()
	columns.Clear()
	propNameCol := columns.AddToVirtualTreeColumn()
	propNameCol.SetText("Name")
	propNameCol.SetWidth(100)
	propNameCol.SetAlignment(types.TaLeftJustify)
	//propNameCol.SetOptions(propNameCol.Options().Include(types.CoDisableAnimatedResize))

	propValueCol := columns.AddToVirtualTreeColumn()
	propValueCol.SetText("Value")
	propValueCol.SetWidth(leftBoxWidth - 100)
	propValueCol.SetAlignment(types.TaLeftJustify)
	propValueCol.SetOptions(propValueCol.Options().Include(types.CoAutoSpring))
	m.propertyTree.SetOnColumnClick(func(sender lcl.IBaseVirtualTree, column int32, shift types.TShiftState) {
		// edit: 1. 触发编辑
		log.Println("propertyTree OnColumnClick column:", column)
		m.propertyTree.EditNode(sender.FocusedNode(), column)
	})
	m.propertyTree.SetOnEditing(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, allowed *bool) {
		// edit: 2. 第二列可以编辑
		log.Println("propertyTree OnEditing column:", column)
		*allowed = column == 1
	})
	m.propertyTree.SetOnEditCancelled(func(sender lcl.IBaseVirtualTree, column int32) {
		log.Println("propertyTree OnEditCancelled column:", column)
	})
	m.propertyTree.SetOnEdited(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) {
		// edit: 4. 编辑结束
		log.Println("propertyTree OnEdited column:", column)
	})
	m.propertyTree.SetOnCreateEditor(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, outEditLink *lcl.IVTEditLink) {
		// edit: 3. 创建编辑或组件
		log.Println("propertyTree OnCreateEditor column:", column)
		if column == 1 {
			ceNode := node
			strEditLink := vtedit.NewStringEditLink()
			strEditLink.SetOnNewData(func(node types.PVirtualNode, column int32, value string) {
				log.Println("StringEditLink NewData:", value, node == ceNode)
				if data, ok := treePropertyNodeDatas[node]; ok {
					data.value = value
				}
			})
			*outEditLink = strEditLink.AsIVTEditLink()
		}
	})
	m.propertyTree.SetOnGetText(func(sender lcl.IBaseVirtualTree, node types.PVirtualNode,
		column int32, textType types.TVSTTextType, cellText *string) {
		log.Println("propertyTree OnGetText column:", column)
		if data, ok := treePropertyNodeDatas[node]; ok {
			if column == 0 {
				*cellText = data.name
			} else if column == 1 {
				*cellText = data.value
			}
		}
	})
	m.propertyTree.SetNodeDataSize(int32(unsafe.Sizeof(uintptr(0))))
}

var treePropertyNodeDatas = make(map[types.PVirtualNode]*TTreePropertyNodeData)

type TTreePropertyNodeData struct {
	name  string
	value string
}
