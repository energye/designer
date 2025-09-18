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

const (
	DLeft = iota
	DTop
	DRight
	DBottom
	DLeftTop
	DRightTop
	DLeftBottom
	DRightBottom
)

type drag struct {
	relation    lcl.IControl   // 关联的控件
	ds          DragShowStatus // 显示方向
	left        lcl.IPanel
	top         lcl.IPanel
	right       lcl.IPanel
	bottom      lcl.IPanel
	leftTop     lcl.IPanel
	rightTop    lcl.IPanel
	leftBottom  lcl.IPanel
	rightBottom lcl.IPanel
}

func (m *drag) newDragPanel(owner lcl.IWinControl, cursor types.TCursor, d int) lcl.IPanel {
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

	var (
		isDown             bool
		dx, dy             int32
		dcx, dcy, dcw, dch int32
	)
	_, _ = dx, dy
	_, _, _, _ = dcx, dcy, dcw, dch
	pnl.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
		if isDown {
			//br := m.relation.BoundsRect()
			//width, _ := br.Width(), br.Height()
			switch d {
			case DLeft:
				x := X - dx
				w := dcw - x
				m.relation.SetBounds(dcx+x, dcy, w, dch)
			case DTop:
			case DRight:
			case DBottom:
			case DLeftTop:
			case DRightTop:
			case DLeftBottom:
			case DRightBottom:
			}
		}
	})
	pnl.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		m.Hide()
		dx, dy = X, Y
		br := m.relation.BoundsRect()
		dcx, dcy, dcw, dch = br.Left, br.Top, br.Width(), br.Height()
		isDown = true
		println("SetOnMouseDown")
	})
	pnl.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		m.Show()
		isDown = false
		println("SetOnMouseUp")
	})
	return pnl
}

func newDrag(owner lcl.IWinControl, ds DragShowStatus) *drag {
	m := new(drag)
	m.ds = ds
	if m.ds == DsAll {
		m.left = m.newDragPanel(owner, types.CrSizeWE, DLeft)
		m.top = m.newDragPanel(owner, types.CrSizeNS, DTop)
		m.right = m.newDragPanel(owner, types.CrSizeWE, DRight)
		m.bottom = m.newDragPanel(owner, types.CrSizeNS, DBottom)
		m.leftTop = m.newDragPanel(owner, types.CrSizeNWSE, DLeftTop)
		m.rightTop = m.newDragPanel(owner, types.CrSizeNESW, DRightTop)
		m.leftBottom = m.newDragPanel(owner, types.CrSizeNESW, DLeftBottom)
		m.rightBottom = m.newDragPanel(owner, types.CrSizeNWSE, DRightBottom)
	} else {
		m.right = m.newDragPanel(owner, types.CrSizeWE, DRight)
		m.bottom = m.newDragPanel(owner, types.CrSizeNS, DBottom)
		m.rightBottom = m.newDragPanel(owner, types.CrSizeNWSE, DRightBottom)
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
			m.left.SetBounds(x-db, y+(height/2)-db, dragBorder, dragBorder)
			m.top.SetBounds(x+(width/2)-db, y-db, dragBorder, dragBorder)
			m.right.SetBounds(x+width-db, y+(height/2)-db, dragBorder, dragBorder)
			m.bottom.SetBounds(x+(width/2)-db, y+height-db, dragBorder, dragBorder)
			m.leftTop.SetBounds(x-db, y-db, dragBorder, dragBorder)
			m.rightTop.SetBounds(x+width-db, y-db, dragBorder, dragBorder)
			m.rightBottom.SetBounds(x+width-db, y+height-db, dragBorder, dragBorder)
			m.leftBottom.SetBounds(x-db, y+height-db, dragBorder, dragBorder)
		} else {
			m.right.SetBounds(x+width-db, y+(height/2)-db, dragBorder, dragBorder)
			m.bottom.SetBounds(x+(width/2)-db, y+height-db, dragBorder, dragBorder)
			m.rightBottom.SetBounds(x+width-db, y+height-db, dragBorder, dragBorder)
		}
	}
}
