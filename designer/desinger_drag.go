package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types/colors"
)

type DragShowStatus int32

var dragBorder int32 = 8

const (
	DsAll         DragShowStatus = iota // 显示所有
	DsRightBottom                       // 显示 右 右下 下
)
const (
	DLeft        = 1
	DTop         = 2
	DRight       = 3
	DBottom      = 4
	DLeftTop     = 5
	DRightTop    = 6
	DLeftBottom  = 7
	DRightBottom = 8
)

type drag struct {
	relation    lcl.IControl   // 关联的控件
	ds          DragShowStatus // 显示状态
	left        lcl.IPanel
	top         lcl.IPanel
	right       lcl.IPanel
	bottom      lcl.IPanel
	leftTop     lcl.IPanel
	rightTop    lcl.IPanel
	leftBottom  lcl.IPanel
	rightBottom lcl.IPanel
}

func newDragPanel(owner lcl.IWinControl, tag int) lcl.IPanel {
	pnl := lcl.NewPanel(owner)
	pnl.SetParent(owner)
	pnl.SetWidth(dragBorder)
	pnl.SetHeight(dragBorder)
	pnl.SetDoubleBuffered(true)
	pnl.SetColor(colors.ClGradientActiveCaption)
	pnl.SetBevelColor(colors.ClHighlight)
	pnl.SetVisible(false)
	pnl.SetTag(uintptr(tag))
	return pnl
}

func newDrag(owner lcl.IWinControl, ds DragShowStatus) *drag {
	m := new(drag)
	m.ds = ds
	if m.ds == DsAll {
		m.left = newDragPanel(owner, DLeft)
		m.top = newDragPanel(owner, DTop)
		m.right = newDragPanel(owner, DRight)
		m.bottom = newDragPanel(owner, DBottom)
		m.leftTop = newDragPanel(owner, DLeftTop)
		m.rightTop = newDragPanel(owner, DRightTop)
		m.leftBottom = newDragPanel(owner, DLeftBottom)
		m.rightBottom = newDragPanel(owner, DRightBottom)
	} else {
		m.right = newDragPanel(owner, DRight)
		m.bottom = newDragPanel(owner, DBottom)
		m.rightBottom = newDragPanel(owner, DRightBottom)
	}
	return m
}

// 设置关联控件
func (m *drag) SetRelation(relation lcl.IControl) {
	m.relation = relation
}

// 隐藏所有
func (m *drag) Hide() {
	if m.ds == DsAll {
		m.left.SetVisible(false)
		m.top.SetVisible(false)
		m.right.SetVisible(false)
		m.bottom.SetVisible(false)
		m.leftTop.SetVisible(false)
		m.rightTop.SetVisible(false)
		m.leftBottom.SetVisible(false)
		m.rightBottom.SetVisible(false)
	} else {
		m.right.SetVisible(false)
		m.bottom.SetVisible(false)
		m.leftBottom.SetVisible(false)
	}
}

// 显示
func (m *drag) Show() {
	if m.ds == DsAll {
		m.left.SetVisible(true)
		m.left.BringToFront()
		m.top.SetVisible(true)
		m.right.SetVisible(true)
		m.bottom.SetVisible(true)
		m.leftTop.SetVisible(true)
		m.rightTop.SetVisible(true)
		m.leftBottom.SetVisible(true)
		m.rightBottom.SetVisible(true)
	} else {
		m.right.SetVisible(true)
		m.bottom.SetVisible(true)
		m.rightBottom.SetVisible(true)
	}
	m.BringToFront()
	m.Follow()
}
func (m *drag) BringToFront() {
	if m.ds == DsAll {
		m.left.BringToFront()
		m.top.BringToFront()
		m.right.BringToFront()
		m.bottom.BringToFront()
		m.leftTop.BringToFront()
		m.rightTop.BringToFront()
		m.leftBottom.BringToFront()
		m.rightBottom.BringToFront()
	} else {
		m.right.BringToFront()
		m.bottom.BringToFront()
		m.rightBottom.BringToFront()
	}
}

// 跟随关联控件
func (m *drag) Follow() {
	if m.relation != nil {
		br := m.relation.BoundsRect()
		x, y := br.Left, br.Top
		width, height := br.Width(), br.Height()
		if m.ds == DsAll {
			m.left.SetLeft(x - (dragBorder / 2))
			m.left.SetTop(y - (dragBorder / 2))
			m.top.SetLeft(x + (width / 2))
			m.top.SetTop(y + height - (dragBorder / 2))
		} else {
			m.right.SetLeft(x + width - (dragBorder / 2))
			m.right.SetTop(y + (height / 2))
			m.bottom.SetLeft(x + (width / 2))
			m.bottom.SetTop(y + height - (dragBorder / 2))
			m.rightBottom.SetLeft(x + width - (dragBorder / 2))
			m.rightBottom.SetTop(y + height - (dragBorder / 2))
		}
	}
}
