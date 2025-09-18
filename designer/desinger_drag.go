package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strconv"
)

type DragShowStatus int32

var dragBorder int32 = 4

const (
	DsAll         DragShowStatus = iota // 显示所有
	DsRightBottom                       // 显示 右 右下 下
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

func (m *drag) newDragPanel(owner lcl.IWinControl, cursor types.TCursor) lcl.IPanel {
	pnl := lcl.NewPanel(owner)
	pnl.SetParent(owner)
	pnl.SetWidth(dragBorder)
	pnl.SetHeight(dragBorder)
	pnl.SetBevelOuter(types.BvNone)
	pnl.SetDoubleBuffered(true)
	pnl.SetColor(colors.ClBlack)
	pnl.SetVisible(false)
	//pnl.SetTag(uintptr(tag))
	pnl.SetShowHint(true)
	pnl.SetHint(strconv.Itoa(int(cursor)))
	pnl.SetCursor(cursor)
	pnl.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {

	})
	pnl.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

	})
	pnl.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {

	})
	return pnl
}

func newDrag(owner lcl.IWinControl, ds DragShowStatus) *drag {
	m := new(drag)
	m.ds = ds
	if m.ds == DsAll {
		m.left = m.newDragPanel(owner, types.CrSizeWE)
		m.top = m.newDragPanel(owner, types.CrSizeNS)
		m.right = m.newDragPanel(owner, types.CrSizeWE)
		m.bottom = m.newDragPanel(owner, types.CrSizeNS)
		m.leftTop = m.newDragPanel(owner, types.CrSizeNWSE)
		m.rightTop = m.newDragPanel(owner, types.CrSizeNESW)
		m.leftBottom = m.newDragPanel(owner, types.CrSizeNESW)
		m.rightBottom = m.newDragPanel(owner, types.CrSizeNWSE)
	} else {
		m.right = m.newDragPanel(owner, types.CrSizeWE)
		m.bottom = m.newDragPanel(owner, types.CrSizeNS)
		m.rightBottom = m.newDragPanel(owner, types.CrSizeNWSE)
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
		m.rightBottom.SetVisible(false)
	}
}

// 显示
func (m *drag) Show() {
	if m.ds == DsAll {
		m.left.SetVisible(true)
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
		_ = width
		db := dragBorder / 2
		if m.ds == DsAll {
			m.left.SetLeft(x - db)
			m.left.SetTop(y + (height / 2) - db)
			m.top.SetLeft(x + (width / 2) - db)
			m.top.SetTop(y - db)
			m.right.SetLeft(x + width - db)
			m.right.SetTop(y + (height / 2) - db)
			m.bottom.SetLeft(x + (width / 2) - db)
			m.bottom.SetTop(y + height - db)
			m.leftTop.SetLeft(x - db)
			m.leftTop.SetTop(y - db)
			m.rightTop.SetLeft(x + width - db)
			m.rightTop.SetTop(y - db)
			m.rightBottom.SetLeft(x + width - db)
			m.rightBottom.SetTop(y + height - db)
			m.leftBottom.SetLeft(x - db)
			m.leftBottom.SetTop(y + height - db)
		} else {
			m.right.SetLeft(x + width - db)
			m.right.SetTop(y + (height / 2) - db)
			m.bottom.SetLeft(x + (width / 2) - db)
			m.bottom.SetTop(y + height - db)
			m.rightBottom.SetLeft(x + width - db)
			m.rightBottom.SetTop(y + height - db)
		}
	}
}
