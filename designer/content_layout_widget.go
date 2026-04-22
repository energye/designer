package designer

import (
	"fmt"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strconv"
)

type ContentLayoutWidget struct {
	searchEdit lcl.ITreeFilterEdit // 组件搜索框
	topBox     lcl.IPanel
	title      lcl.ILabel
	tree       lcl.ITreeView
}

func initContentLayoutWidget(owner *ContentLayout) *ContentLayoutWidget {
	m := &ContentLayoutWidget{}
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
	m.topBox.SetBorderStyleToBorderStyle(types.BsSingle)
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
	m.tree.SetTreeLineColor(colors.RGBToColor(128, 128, 128))
	m.tree.SetTreeLinePenStyle(types.PsClear)
	m.tree.SetAlign(types.AlClient)
	//tree.SetAlign(types.AlCustom)
	//tree.SetBoundsRect(m.widgetPanel.ClientRect())
	m.tree.SetBorderStyleToBorderStyle(types.BsNone)
	m.tree.SetImages(imageComponents.ImageList100())
	m.tree.SetIndent(18)
	m.tree.SetDefaultItemHeight(26)
	m.tree.SetRowSelect(true)
	m.tree.SetHideSelection(false)
	m.tree.SetScrollBars(types.SsVertical)

	// ======================================================
	// 状态变量（放在外层）
	// ======================================================
	var (
		hoverNode lcl.ITreeNode
		pressNode lcl.ITreeNode
		isDown    bool
	)

	m.tree.SetRowSelect(true)
	m.tree.SetDoubleBuffered(true)
	m.tree.SetDefaultItemHeight(28)

	m.tree.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, x int32, y int32) {
		if isDown {
			return
		}
		n := m.tree.GetNodeAt(x, y)
		if n != nil && hoverNode != n {
			hoverNode = n
			owner.contentStatus.SetSimpleText(fmt.Sprintf(n.Text()))
			m.tree.Invalidate()
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
		isDown = false
		pressNode = nil
		if hoverNode != nil {
			if hoverNode.Level() == 0 {
				hoverNode.SetExpanded(!hoverNode.Expanded())
			} else {
				fmt.Println("click:", hoverNode.Level(), hoverNode.Text())
			}
		}
		m.tree.Invalidate()
		lcl.RunOnMainThreadAsync(func(id uint32) {
			// 强制刷新
			br := owner.widgetPanel.BoundsRect()
			br.SetWidth(br.Width() - 1)
			owner.widgetPanel.SetWidth(br.Width())
			br.SetWidth(br.Width() + 1)
			owner.widgetPanel.SetWidth(br.Width())
		})
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
		isSelected := node.Selected()

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
			canvas.FillRectWithRect(types.Rect(0, r.Top, m.tree.ClientWidth(), r.Bottom))

			// 上下边框线
			pen.SetStyle(types.PsSolid)
			pen.SetColor(0x00D9D9D9)

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
				bg = 0x00F6F6F6
			}
			if isPress {
				bg = 0x00ECECEC
			}
			brush.SetStyle(types.BsSolid)
			brush.SetColor(bg)
			canvas.FillRectWithRect(types.Rect(0, r.Top, m.tree.ClientWidth(), r.Bottom))

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
			font.SetSize(8)
			if isPress {
				font.SetColor(0x00333333)
			}
			canvas.TextOutWithIntX2Str(textX, textY, node.Text())
			brush.SetStyle(types.BsSolid)
		}
	})

	//br := m.widgetPanel.BoundsRect()
	//br.SetWidth(br.Width() - 1)
	//m.widgetPanel.SetBoundsRect(br)
	//br.SetWidth(br.Width() + 1)
	//m.widgetPanel.SetBoundsRect(br)
	//m.tree.SetWidth(br.Width())
	//m.widgetPanel.Invalidate()
	// ======================================================
	// 点击一级节点展开 / 收起
	// ======================================================
	//m.tree.SetOnClick(func(sender lcl.IObject) {
	//})

	// ======================================================
	// 创建数据
	// ======================================================
	m.tree.SetParent(owner.widgetPanel)

	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()

	items := m.tree.Items()
	items.Clear()

	// ======================
	// 常用
	// ======================
	cat := items.Add(nil, "常用")
	cat.SetExpanded(true)
	var node lcl.ITreeNode
	for i := 0; i < 10; i++ {
		node = items.AddChild(cat, "按钮 Button"+strconv.Itoa(i))
		node.SetImageIndex(0)
		node.SetSelectedIndex(0)

		node = items.AddChild(cat, "输入框 Edit"+strconv.Itoa(i))
		node.SetImageIndex(1)
		node.SetSelectedIndex(1)

		node = items.AddChild(cat, "标签 Label"+strconv.Itoa(i))
		node.SetImageIndex(2)
		node.SetSelectedIndex(2)
	}

	// ======================
	// 容器
	// ======================
	cat2 := items.Add(nil, "容器")
	cat2.SetExpanded(false)
	for i := 0; i < 10; i++ {
		node = items.AddChild(cat2, "面板 Panel"+strconv.Itoa(i))
		node.SetImageIndex(3)
		node.SetSelectedIndex(3)

		node = items.AddChild(cat2, "分页 Tabs"+strconv.Itoa(i))
		node.SetImageIndex(4)
		node.SetSelectedIndex(4)

		node = items.AddChild(cat2, "分割条 Splitter"+strconv.Itoa(i))
		node.SetImageIndex(5)
		node.SetSelectedIndex(5)
	}

	// ======================
	// Web
	// ======================
	cat3 := items.Add(nil, "Web")
	cat3.SetExpanded(true)
	for i := 0; i < 10; i++ {
		node = items.AddChild(cat3, "WebView "+strconv.Itoa(i))
		node.SetImageIndex(6)
		node.SetSelectedIndex(6)

		node = items.AddChild(cat3, "CEF "+strconv.Itoa(i))
		node.SetImageIndex(7)
		node.SetSelectedIndex(7)
	}

	return m
}
