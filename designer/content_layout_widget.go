package designer

import (
	"fmt"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

type ContentLayoutWidget struct {
	searchEdit        lcl.ITreeFilterEdit // 组件搜索框
	topBox            lcl.IPanel
	title             lcl.ILabel
	tree              lcl.ITreeView
	components        map[uintptr]*TWidgetTreeItem // 组件选项卡： 标准，附加，通用等等
	selectedComponent *TWidgetTreeItem             // 选中的组件
}

type TWidgetTreeItem struct {
	parent   *TWidgetTreeItem
	name     string
	node     lcl.ITreeNode
	child    map[uintptr]*TWidgetTreeItem
	selected bool
}

func (m *TWidgetTreeItem) IsSelectTool() bool {
	return m.name == "SelectTool"
}

func SelectedComponent() *TWidgetTreeItem {
	if MainWindow.contentLayout != nil && MainWindow.contentLayout.layoutWidget != nil {
		return MainWindow.contentLayout.layoutWidget.selectedComponent
	}
	return nil
}

func ResetSelectedComponent() {
	if MainWindow.contentLayout != nil && MainWindow.contentLayout.layoutWidget != nil {
		MainWindow.contentLayout.layoutWidget.resetAllNoSelected()
		MainWindow.contentLayout.layoutWidget.selectedComponent = nil
		MainWindow.contentLayout.layoutWidget.tree.Invalidate()
	}
}

func initContentLayoutWidget(owner *ContentLayout) *ContentLayoutWidget {
	logs.Debug("创建组件选项卡面板")
	m := &ContentLayoutWidget{components: make(map[uintptr]*TWidgetTreeItem)}
	m.searchEdit = lcl.NewTreeFilterEdit(owner.widgetPanel)
	m.searchEdit.SetTextHint("搜索组件")
	m.searchEdit.SetAlign(types.AlTop)
	m.searchEdit.SetAutoSelect(false)
	borderSpacing := m.searchEdit.BorderSpacing()
	borderSpacing.SetLeft(3)
	borderSpacing.SetRight(3)
	borderSpacing.SetTop(3)
	borderSpacing.SetBottom(3)
	m.searchEdit.SetParent(owner.widgetPanel)

	m.topBox = lcl.NewPanel(owner.widgetPanel)
	m.topBox.SetBorderStyleToBorderStyle(types.BsNone)
	m.topBox.SetBevelOuter(types.BvNone)
	m.topBox.SetAlign(types.AlTop)
	m.topBox.SetHeight(30)
	borderSpacing = m.topBox.BorderSpacing()
	m.topBox.SetParent(owner.widgetPanel)

	title := lcl.NewLabel(m.topBox)
	title.SetCaption("组件库")
	title.SetLeft(5)
	title.SetTop(5)
	font := title.Font()
	font.SetSize(10)
	font.SetStyle(types.NewSet(types.FsBold))
	title.SetParent(m.topBox)

	m.tree = lcl.NewTreeView(owner.widgetPanel)
	m.tree.SetAutoExpand(true)
	m.tree.SetReadOnly(true)
	m.tree.SetDoubleBuffered(true)
	m.tree.SetTreeLinePenStyle(types.PsClear)
	m.tree.SetAlign(types.AlClient)
	m.tree.SetBorderStyleToBorderStyle(types.BsNone)
	m.tree.SetImages(imageComponents.ImageList100())
	m.tree.SetIndent(18)
	m.tree.SetDefaultItemHeight(26)
	m.tree.SetRowSelect(true)
	m.tree.SetHideSelection(false)
	m.tree.Font().SetSize(8)
	m.tree.SetScrollBars(types.SsVertical)

	var (
		hoverNode lcl.ITreeNode
		pressNode lcl.ITreeNode
		isDown    bool
	)

	m.tree.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, x int32, y int32) {
		if isDown {
			return
		}
		n := m.tree.GetNodeAt(x, y)
		if n != nil && hoverNode != n {
			hoverNode = n
			m.tree.Invalidate()
			owner.contentStatus.SetSimpleText(fmt.Sprintf(n.Text()))
		} else {
			hoverNode = nil
		}
	})
	m.tree.SetOnMouseLeave(func(sender lcl.IObject) {
		hoverNode = nil
		m.tree.Invalidate()
	})
	m.tree.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x int32, y int32) {
		isDown = true
		if button != types.MbLeft {
			return
		}
		pressNode = m.tree.GetNodeAt(x, y)
		m.tree.Invalidate()
	})
	m.tree.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x int32, y int32) {
		defer m.tree.Invalidate()
		isDown = false
		pressNode = nil
		if hoverNode != nil {
			if hoverNode.Level() == 0 {
				hoverNode.SetExpanded(!hoverNode.Expanded())
			} else {
				m.updateAllNoSelected(hoverNode)
				m.selectedComponent = m.findComponentTreeItem(hoverNode)
				m.selectedComponent.selected = true // 如果是 nil 错误, 说明逻辑有问题
				if m.selectedComponent.IsSelectTool() {
					m.selectedComponent = nil
				}
				fmt.Println("click:", hoverNode.Level(), hoverNode.Text())
				return
			}
		}
		m.selectedComponent = nil
		if tool.IsDarwin {
			// 强制刷新, 滚动条出现再隐藏后渲染有问题, 先这样解决
			lcl.RunOnMainThreadAsync(func(id uint32) {
				br := owner.widgetPanel.BoundsRect()
				br.SetWidth(br.Width() - 1)
				owner.widgetPanel.SetWidth(br.Width())
				br.SetWidth(br.Width() + 1)
				owner.widgetPanel.SetWidth(br.Width())
			})
		}
	})
	m.tree.SetOnAdvancedCustomDrawItem(func(sender lcl.ICustomTreeView, node lcl.ITreeNode, state types.TCustomDrawState,
		stage types.TCustomDrawStage, paintImages *bool, defaultDraw *bool) {
		*paintImages = false
		*defaultDraw = false
		if stage != types.CdPrePaint {
			return
		}
		canvas := sender.Canvas()
		brush := canvas.BrushToBrush()
		font := canvas.FontToFont()
		pen := canvas.PenToPen()
		//fmt.Println("OnAdvancedCustomDrawItem Width:", m.tree.Width(), canvas.Width())

		r := node.DisplayRect(false)
		rowH := r.Bottom - r.Top

		isHover := hoverNode != nil && hoverNode.Equals(node)
		isPress := pressNode != nil && pressNode.Equals(node)
		//isSelected := node.Selected()

		if node.Level() == 0 {
			bg := types.TColor(0xF4F6F9)
			if isHover {
				bg = 0xF1E9E3
			}
			if isPress {
				bg = 0xE1D8D0
			}
			// 背景
			brush.SetStyle(types.BsSolid)
			brush.SetColor(bg)
			canvas.FillRectWithRect(types.Rect(0, r.Top, m.tree.ClientWidth(), r.Bottom))

			// 上下边框线
			pen.SetStyle(types.PsSolid)
			pen.SetColor(0xD9D9D9)

			// 下边线
			canvas.MoveToWithIntX2(0, r.Bottom-1)
			canvas.LineToWithIntX2(m.tree.ClientWidth(), r.Bottom-1)

			// 字体
			font.SetColor(colors.ClBlack)

			arrow := "▶"
			if node.Expanded() {
				arrow = "▼"
			}

			titleY := r.Top + (rowH-14)/2
			canvas.TextOutWithIntX2Str(6, titleY, arrow+" "+node.Text())
			return
		} else {
			iconSize := int32(24)
			iconX := int32(20)
			iconY := r.Top + (rowH-iconSize)/2
			gap := int32(10)

			textX := iconX + iconSize + gap
			textY := r.Top + (rowH-14)/2

			// 背景
			bg := colors.ClWhite
			if isHover {
				bg = 0xF6F6F6
			}
			if isPress {
				bg = 0xECECEC
			}
			brush.SetStyle(types.BsSolid)
			brush.SetColor(bg)
			canvas.FillRectWithRect(types.Rect(0, r.Top, m.tree.ClientWidth(), r.Bottom))
			isSelected := false
			if item := m.findComponentTreeItem(node); item != nil && item.selected {
				isSelected = true
			}

			// Selected 图标背景
			if isSelected {
				offset := int32(1)
				left := iconX - offset
				top := iconY - offset
				right := iconX + iconSize + offset
				bottom := iconY + iconSize + offset
				brush.SetColor(0xE6A23C)
				pen.SetStyle(types.PsClear)
				canvas.RoundRectWithIntX6(left, top, right, bottom, 6, 6)
			}

			// 图标
			imgs := m.tree.Images()
			if imgs != nil {
				imgIndex := node.ImageIndex()
				if isSelected && node.SelectedIndex() >= 0 {
					imgIndex = node.SelectedIndex()
				}
				if imgIndex >= 0 {
					imgs.DrawWithCanvasIntX3Bool(canvas, iconX, iconY, imgIndex, true)
				}
			}

			// 文本
			brush.SetStyle(types.BsClear)
			font.SetColor(colors.ClBlack)
			//font.SetSize(8)
			if isPress {
				font.SetColor(0x00333333)
			}
			canvas.TextOutWithIntX2Str(textX, textY, node.Text())
			brush.SetStyle(types.BsSolid)
		}
	})

	// ======================================================
	// 创建数据
	// ======================================================
	m.tree.SetParent(owner.widgetPanel)

	m.initComponentTreeData()

	return m
}

func (m *ContentLayoutWidget) findComponentTreeItem(node lcl.ITreeNode) *TWidgetTreeItem {
	if item, ok := m.components[node.Instance()]; ok {
		return item
	}
	return nil
}

func (m *ContentLayoutWidget) resetAllNoSelected() {
	// 先重置所有
	for _, item := range m.components {
		item.selected = false
		if item.IsSelectTool() {
			item.selected = true
		}
	}
}

// 更新所有节点状态未选中
// 只工具选项默认选中
func (m *ContentLayoutWidget) updateAllNoSelected(node lcl.ITreeNode) {
	// 先重置所有
	m.resetAllNoSelected()
	// 重置当前节点所属组的所有节点为未选中
	if item := m.findComponentTreeItem(node); item != nil {
		if item.parent != nil {
			tagGroup := m.findComponentTreeItem(item.parent.node)
			for _, child := range tagGroup.child {
				child.selected = false
			}
		}
	}
}

// 初始化组件选项面板树数据
func (m *ContentLayoutWidget) initComponentTreeData() {
	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()
	items := m.tree.Items()
	items.Clear()
	// 创建组件选项卡
	newComponentData := func(tab config.Tab) *TWidgetTreeItem {
		logs.Debug("创建组件选项卡:", tab.Cn)
		// 一级
		root := items.AddChild(nil, tab.Cn)
		rootTab := &TWidgetTreeItem{child: make(map[uintptr]*TWidgetTreeItem), node: root, name: tab.En}
		m.components[root.Instance()] = rootTab
		//root.SetImageIndex(6)
		//root.SetSelectedIndex(6)

		var (
			child      lcl.ITreeNode
			imageIndex int32
			item       *TWidgetTreeItem
		)

		// 二级
		// 选择工具 鼠标
		child = items.AddChild(root, "选择指针")
		imageIndex = imageComponents.ImageIndex("cursortool.png")
		child.SetImageIndex(imageIndex)
		child.SetSelectedIndex(imageIndex)
		item = &TWidgetTreeItem{name: "SelectTool", node: child, selected: true, parent: rootTab}
		rootTab.child[child.Instance()] = item
		m.components[child.Instance()] = item

		// 创建组件按钮
		for _, name := range tab.Component {
			child = items.AddChild(root, name)
			imageIndex = imageComponents.ImageIndex(name + ".png")
			child.SetImageIndex(imageIndex)
			child.SetSelectedIndex(imageIndex)
			item = &TWidgetTreeItem{name: name, node: child, parent: rootTab}
			rootTab.child[child.Instance()] = item
			m.components[child.Instance()] = item
		}
		root.SetExpanded(false)
		return rootTab
	}
	// 创建组件选项卡
	newComponentData(config.DesignerConfig.ComponentTabs.Standard)
	newComponentData(config.DesignerConfig.ComponentTabs.Additional)
	newComponentData(config.DesignerConfig.ComponentTabs.Common)
	newComponentData(config.DesignerConfig.ComponentTabs.Dialogs)
	newComponentData(config.DesignerConfig.ComponentTabs.Misc)
	newComponentData(config.DesignerConfig.ComponentTabs.System)
	newComponentData(config.DesignerConfig.ComponentTabs.LazControl)
	newComponentData(config.DesignerConfig.ComponentTabs.SynEdit)
	webview := newComponentData(config.DesignerConfig.ComponentTabs.WebView)
	webview.node.SetExpanded(true)
}
