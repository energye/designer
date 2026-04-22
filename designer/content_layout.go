package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/widget/wg"
)

var (
	defaultSplitterWidth   = int32(1)
	defaultSplitterMinSize = int32(1)
)

type ContentLayout struct {
	box      lcl.IPanel // 主内容布局盒子
	rightBox lcl.IPanel // 右侧盒子

	widgetSplitter lcl.ISplitter // 组件面板分隔器
	widgetPanel    *TViewPanel   // 组件面板
	widgetSearch   lcl.IEdit     // 组件搜索框

	projectSplitter lcl.ISplitter // 项目面板分隔器
	projectPanel    *TViewPanel   // 项目面板

	designerPanel lcl.IPanel // 设计器面板

	inspectorSplitter lcl.ISplitter // 属性检查器面板分隔器
	inspectorPanel    *TViewPanel   // 属性检查器面板

	consoleLogSplitter lcl.ISplitter // 日志面板分隔器
	consoleLogPanel    lcl.IPanel    // 控制台输出
}

// 底部布局 - 左: 组件库, 左中: 项目查看, 中: 中间画布(自适应), 右: 属性, 下: 日志控制
func initBottomBox(owner lcl.IWinControl) *ContentLayout {
	m := &ContentLayout{}

	// 主内容布局盒子
	m.box = lcl.NewPanel(owner)
	m.box.SetColor(colors.Cl3DDkShadow)
	m.box.SetBevelOuter(types.BvNone)
	m.box.SetDoubleBuffered(true)
	m.box.SetAlign(types.AlClient)
	m.box.SetCaption("主内容布局")
	m.box.SetParent(owner)

	// 组件面板分隔器
	m.widgetSplitter = lcl.NewSplitter(owner)
	m.widgetSplitter.SetAlign(types.AlLeft)
	m.widgetSplitter.SetWidth(defaultSplitterWidth)
	m.widgetSplitter.SetMinSize(defaultSplitterMinSize)
	m.widgetSplitter.SetParent(m.box)
	// 组件面板
	m.widgetPanel = NewViewPanel(owner)
	m.widgetPanel.SetColor(wg.LightenColor(colors.ClAqua, 0.3))
	m.widgetPanel.SetAlign(types.AlLeft)
	m.widgetPanel.SetCaption("组件库")
	m.widgetPanel.SetWidth(170)    //动态控制
	m.widgetPanel.SetVisible(true) //动态控制
	m.widgetPanel.SetTitle("组件库")
	m.widgetPanel.Constraints().SetMinWidth(30)
	m.widgetPanel.SetParent(m.box)
	//m.widgetPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("widgetPanel.SetOnResize BoundsRect:", m.widgetPanel.BoundsRect())
	//})

	//m.widgetSearch = lcl.NewEdit(owner)
	//m.widgetSearch.SetTextHint("搜索组件")
	//m.widgetSearch.SetAlign(types.AlTop)
	//m.widgetSearch.SetParent(m.widgetPanel)

	// 右侧盒子
	m.rightBox = lcl.NewPanel(owner)
	m.rightBox.SetColor(wg.LightenColor(colors.ClAqua, 1.5))
	m.rightBox.SetBevelOuter(types.BvNone)
	m.rightBox.SetDoubleBuffered(true)
	m.rightBox.SetAlign(types.AlClient)
	m.rightBox.SetCaption("右内容")
	m.rightBox.SetHeight(150)
	m.rightBox.SetParent(m.box)

	// 项目面板分隔器
	m.projectSplitter = lcl.NewSplitter(owner)
	m.projectSplitter.SetAlign(types.AlLeft)
	m.projectSplitter.SetWidth(defaultSplitterWidth)
	m.projectSplitter.SetMinSize(defaultSplitterMinSize)
	m.projectSplitter.SetParent(m.rightBox)
	// 项目面板
	m.projectPanel = NewViewPanel(owner)
	m.projectPanel.SetColor(wg.LightenColor(colors.ClAqua, 0.6))
	m.projectPanel.SetAlign(types.AlLeft)
	m.projectPanel.SetWidth(150)    //动态控制
	m.projectPanel.SetVisible(true) //动态控制
	m.projectPanel.SetCaption("项目管理器")
	m.projectPanel.SetTitle("项目管理器")
	m.projectPanel.Constraints().SetMinWidth(30)
	m.projectPanel.SetParent(m.rightBox)
	//m.projectPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("projectPanel.SetOnResize BoundsRect:", m.projectPanel.BoundsRect())
	//})

	// 设计器
	m.designerPanel = lcl.NewPanel(owner)
	m.designerPanel.SetBevelOuter(types.BvNone)
	m.designerPanel.SetDoubleBuffered(true)
	m.designerPanel.SetCaption("设计器画布")
	m.designerPanel.SetColor(wg.LightenColor(colors.ClAqua, 0.9))
	m.designerPanel.SetAlign(types.AlClient)
	m.designerPanel.SetParent(m.rightBox)

	// 查看器面板分隔器
	m.inspectorSplitter = lcl.NewSplitter(owner)
	m.inspectorSplitter.SetAlign(types.AlRight)
	m.inspectorSplitter.SetWidth(defaultSplitterWidth)
	m.inspectorSplitter.SetMinSize(defaultSplitterMinSize)
	m.inspectorSplitter.SetParent(m.rightBox)
	// 查看器
	m.inspectorPanel = NewViewPanel(owner)
	m.inspectorPanel.SetColor(wg.LightenColor(colors.ClAqua, 1.2))
	m.inspectorPanel.SetAlign(types.AlRight)
	m.inspectorPanel.SetCaption("属性查看器")
	m.inspectorPanel.SetTitle("属性查看器")
	m.inspectorPanel.SetWidth(225)    //动态控制
	m.inspectorPanel.SetVisible(true) //动态控制
	m.inspectorPanel.Constraints().SetMinWidth(30)
	m.inspectorPanel.SetParent(m.rightBox)
	//m.inspectorPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("inspectorPanel.SetOnResize BoundsRect:", m.inspectorPanel.BoundsRect())
	//})

	// 日志面板分隔器
	m.consoleLogSplitter = lcl.NewSplitter(owner)
	m.consoleLogSplitter.SetAlign(types.AlBottom)
	m.consoleLogSplitter.SetWidth(defaultSplitterWidth)
	m.consoleLogSplitter.SetMinSize(defaultSplitterMinSize)
	m.consoleLogSplitter.SetParent(m.rightBox)
	// 日志
	m.consoleLogPanel = lcl.NewPanel(owner)
	m.consoleLogPanel.SetColor(wg.LightenColor(colors.ClAqua, 1.5))
	m.consoleLogPanel.SetBevelOuter(types.BvNone)
	m.consoleLogPanel.SetDoubleBuffered(true)
	m.consoleLogPanel.SetAlign(types.AlBottom)
	m.consoleLogPanel.SetCaption("控制台输出")
	m.consoleLogPanel.SetHeight(150)   //动态控制
	m.consoleLogPanel.SetVisible(true) //动态控制
	m.consoleLogPanel.Constraints().SetMinHeight(30)
	m.consoleLogPanel.SetParent(m.rightBox)
	//m.consoleLogPanel.SetOnResize(func(sender lcl.IObject) {
	//	fmt.Println("consoleLogPanel.SetOnResize BoundsRect:", m.consoleLogPanel.BoundsRect())
	//})

	m.initWidgetPanel()
	return m
}

// 初始化组件面板布局
func (m *ContentLayout) initWidgetPanel() {
	tree := lcl.NewTreeView(m.widgetPanel)
	tree.SetAutoExpand(true)
	tree.SetReadOnly(true)
	tree.SetDoubleBuffered(true)
	tree.SetTreeLineColor(colors.RGBToColor(128, 128, 128))
	tree.SetTreeLinePenStyle(types.PsClear)
	tree.SetAlign(types.AlClient)
	tree.SetBorderStyleToBorderStyle(types.BsNone)
	tree.SetImages(imageComponents.ImageList100())
	tree.SetIndent(18)
	tree.SetDefaultItemHeight(26)
	tree.SetRowSelect(true)
	tree.SetHideSelection(false)

	// ======================================================
	// 状态变量（放在外层）
	// ======================================================
	var hoverNode lcl.ITreeNode
	var pressNode lcl.ITreeNode

	tree.SetRowSelect(true)
	tree.SetDoubleBuffered(true)
	tree.SetDefaultItemHeight(28)

	// ======================================================
	// Hover
	// ======================================================
	tree.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, x int32, y int32) {
		n := tree.GetNodeAt(x, y)
		if hoverNode != n {
			hoverNode = n
			tree.Invalidate()
		}
	})

	tree.SetOnMouseLeave(func(sender lcl.IObject) {
		hoverNode = nil
		tree.Invalidate()
	})

	// ======================================================
	// Pressed
	// ======================================================
	tree.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x int32, y int32) {
		if button != types.MbLeft {
			return
		}

		pressNode = tree.GetNodeAt(x, y)
		tree.Invalidate()
	})

	tree.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x int32, y int32) {
		pressNode = nil
		tree.Invalidate()
	})

	// ======================================================
	// 自绘
	// ======================================================
	tree.SetOnAdvancedCustomDrawItem(func(sender lcl.ICustomTreeView, node lcl.ITreeNode, state types.TCustomDrawState, stage types.TCustomDrawStage, paintImages *bool, defaultDraw *bool) {
		*paintImages = false
		*defaultDraw = false
		if stage != types.CdPrePaint {
			return
		}
		canvas := tree.Canvas()
		brush := canvas.BrushToBrush()
		font := canvas.FontToFont()
		pen := canvas.PenToPen()

		r := node.DisplayRect(false)
		rowH := r.Bottom - r.Top

		isHover := hoverNode != nil && hoverNode.Equals(node)
		isPress := pressNode != nil && pressNode.Equals(node)
		isSelected := node.Selected()

		// ==================================================
		// 一级节点（分类标题）
		// ==================================================
		if node.Level() == 0 {
			bg := types.TColor(0x00F4F6F9)
			if isHover {
				bg = 0x00E3E9F1
			}
			if isPress {
				bg = 0x00D0D8E1
			}
			// 背景
			brush.SetStyle(types.BsSolid)
			brush.SetColor(bg)
			canvas.FillRectWithRect(types.Rect(0, r.Top, tree.ClientWidth(), r.Bottom))

			// 上下边框线
			pen.SetStyle(types.PsSolid)
			pen.SetColor(0x00D9D9D9)

			// 下边线
			canvas.MoveToWithIntX2(0, r.Bottom-1)
			canvas.LineToWithIntX2(tree.ClientWidth(), r.Bottom-1)

			// 字体
			//font.SetStyle(font.Style().Include(types.FsBold))
			font.SetColor(colors.ClBlack)

			arrow := "▶"
			if node.Expanded() {
				arrow = "▼"
			}

			titleY := r.Top + (rowH-14)/2
			canvas.TextOutWithIntX2Str(6, titleY, arrow+" "+node.Text())
			return
		}

		// ==================================================
		// 二级节点（组件项）
		// ==================================================
		iconSize := int32(24)
		iconX := int32(20)
		iconY := r.Top + (rowH-iconSize)/2
		gap := int32(10)

		textX := iconX + iconSize + gap
		textY := r.Top + (rowH-14)/2

		// 背景
		bg := colors.ClWhite
		if isHover {
			bg = 0x00F6F6F6
		}
		if isPress {
			bg = 0x00ECECEC
		}
		brush.SetStyle(types.BsSolid)
		brush.SetColor(bg)
		canvas.FillRectWithRect(types.Rect(0, r.Top, tree.ClientWidth(), r.Bottom))

		// Selected 图标背景
		if isSelected {
			offset := int32(1)
			left := iconX - offset
			top := iconY - offset
			right := iconX + iconSize + offset
			bottom := iconY + iconSize + offset
			brush.SetColor(0x00E6A23C)
			pen.SetStyle(types.PsClear)
			canvas.RoundRectWithIntX6(left, top, right, bottom, 6, 6)
		}

		// 图标
		imgs := tree.Images()
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
		font.SetSize(8)
		if isPress {
			font.SetColor(0x00333333)
		}
		canvas.TextOutWithIntX2Str(textX, textY, node.Text())
		brush.SetStyle(types.BsSolid)
	})

	// ======================================================
	// 点击一级节点展开 / 收起
	// ======================================================
	tree.SetOnClick(func(sender lcl.IObject) {
		n := tree.Selected()
		if n == nil {
			return
		}
		if n.Level() == 0 {
			n.SetExpanded(!n.Expanded())
		}
	})

	// ======================================================
	// 创建数据
	// ======================================================
	tree.SetParent(m.widgetPanel)

	tree.BeginUpdate()
	defer tree.EndUpdate()

	items := tree.Items()
	items.Clear()

	// ======================
	// 常用
	// ======================
	cat := items.Add(nil, "常用")
	cat.SetExpanded(true)

	node := items.AddChild(cat, "按钮 Button")
	node.SetImageIndex(0)
	node.SetSelectedIndex(0)

	node = items.AddChild(cat, "输入框 Edit")
	node.SetImageIndex(1)
	node.SetSelectedIndex(1)

	node = items.AddChild(cat, "标签 Label")
	node.SetImageIndex(2)
	node.SetSelectedIndex(2)

	// ======================
	// 容器
	// ======================
	cat2 := items.Add(nil, "容器")
	cat2.SetExpanded(false)

	node = items.AddChild(cat2, "面板 Panel")
	node.SetImageIndex(3)
	node.SetSelectedIndex(3)

	node = items.AddChild(cat2, "分页 Tabs")
	node.SetImageIndex(4)
	node.SetSelectedIndex(4)

	node = items.AddChild(cat2, "分割条 Splitter")
	node.SetImageIndex(5)
	node.SetSelectedIndex(5)

	// ======================
	// Web
	// ======================
	cat3 := items.Add(nil, "Web")
	cat3.SetExpanded(true)

	node = items.AddChild(cat3, "WebView")
	node.SetImageIndex(6)
	node.SetSelectedIndex(6)

	node = items.AddChild(cat3, "CEF")
	node.SetImageIndex(7)
	node.SetSelectedIndex(7)
}
