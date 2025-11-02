// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"fmt"
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/message"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 拖拽控制

var dragBorder int32 = 4

// 最小移动距离阈值
const minDistance = 4

// 拖拽控制
type drag struct {
	relation    *TDesigningComponent  // 关联设计的组件
	ds          consts.DragShowStatus // 显示方向
	isShow      bool                  // 是否显示
	dx, dy      int32                 // 拖拽控制
	dcl, dct    int32                 // 拖拽控制
	isDown      bool                  // 拖拽控制
	left        lcl.IPanel
	top         lcl.IPanel
	right       lcl.IPanel
	bottom      lcl.IPanel
	leftTop     lcl.IPanel
	rightTop    lcl.IPanel
	leftBottom  lcl.IPanel
	rightBottom lcl.IPanel
}

func (m *drag) Free() {
	m.relation = nil
	if m.ds == consts.DsAll {
		m.left.Free()
		m.top.Free()
		m.right.Free()
		m.bottom.Free()
		m.leftTop.Free()
		m.rightTop.Free()
		m.leftBottom.Free()
		m.rightBottom.Free()
	} else {
		m.right.Free()
		m.bottom.Free()
		m.rightBottom.Free()
	}
}
func (m *drag) newDragPanel(owner lcl.IWinControl, cursor types.TCursor, d int) lcl.IPanel {
	pnl := lcl.NewPanel(owner)
	pnl.SetWidth(dragBorder)
	pnl.SetHeight(dragBorder)
	pnl.SetBevelOuter(types.BvNone)
	pnl.SetDoubleBuffered(true)
	pnl.SetColor(colors.ClBlack)
	pnl.SetVisible(false)
	//pnl.SetTag(uintptr(tag))
	pnl.SetShowHint(true)
	pnl.SetCursor(cursor)
	hint := ""
	switch cursor {
	case types.CrSizeWE:
		hint = "WE"
	case types.CrSizeNS:
		hint = "NS"
	case types.CrSizeNWSE:
		hint = "NWSE"
	case types.CrSizeNESW:
		hint = "NESW"
	}
	pnl.SetHint(hint)
	pnl.SetParent(owner)
	var (
		isDown             bool
		dx, dy             int32
		dcx, dcy, dcw, dch int32
	)
	_, _ = dx, dy
	_, _, _, _ = dcx, dcy, dcw, dch
	pnl.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
		if isDown {
			switch d {
			case consts.DLeft:
				x := X - dx
				w := dcw - x
				m.relation.SetBounds(dcx+x, dcy, w, dch)
			case consts.DTop:
				y := Y - dy
				h := dch - y
				m.relation.SetBounds(dcx, dcy+y, dcw, h)
			case consts.DRight:
				x := X - dx
				m.relation.SetBounds(dcx, dcy, dcw+x, dch)
			case consts.DBottom:
				y := Y - dy
				m.relation.SetBounds(dcx, dcy, dcw, dch+y)
			case consts.DLeftTop:
				x := X - dx
				w := dcw - x
				y := Y - dy
				h := dch - y
				m.relation.SetBounds(dcx+x, dcy+y, w, h)
			case consts.DRightTop:
				y := Y - dy
				h := dch - y
				x := X - dx
				m.relation.SetBounds(dcx, dcy+y, dcw+x, h)
			case consts.DLeftBottom:
				x := X - dx
				w := dcw - x
				y := Y - dy
				m.relation.SetBounds(dcx+x, dcy, w, dch+y)
			case consts.DRightBottom:
				x := X - dx
				y := Y - dy
				m.relation.SetBounds(dcx, dcy, dcw+x, dch+y)
			}
			br := m.relation.BoundsRect()
			go m.relation.UpdateNodeDataSize(br.Width(), br.Height())
			msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
			message.Follow(msgContent)
		}
	})
	pnl.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		logs.Debug("DRAG OnMouseDown direction:", d)
		m.Hide()
		dx, dy = X, Y
		br := m.relation.BoundsRect()
		dcx, dcy, dcw, dch = br.Left, br.Top, br.Width(), br.Height()
		isDown = true

		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", dcx, dcy, dcw, dch)
		message.Follow(msgContent)
		m.relation.DragBegin()
	})
	pnl.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		logs.Debug("DRAG OnMouseUP direction:", d)
		m.Show()
		isDown = false
		message.FollowHide()
		br := m.relation.BoundsRect()
		go m.relation.UpdateNodeDataSize(br.Width(), br.Height())
		m.relation.DragEnd()
	})
	return pnl
}

func newDrag(owner lcl.IWinControl, ds consts.DragShowStatus) *drag {
	m := new(drag)
	m.ds = ds
	if m.ds == consts.DsAll {
		m.left = m.newDragPanel(owner, types.CrSizeWE, consts.DLeft)
		m.top = m.newDragPanel(owner, types.CrSizeNS, consts.DTop)
		m.right = m.newDragPanel(owner, types.CrSizeWE, consts.DRight)
		m.bottom = m.newDragPanel(owner, types.CrSizeNS, consts.DBottom)
		m.leftTop = m.newDragPanel(owner, types.CrSizeNWSE, consts.DLeftTop)
		m.rightTop = m.newDragPanel(owner, types.CrSizeNESW, consts.DRightTop)
		m.leftBottom = m.newDragPanel(owner, types.CrSizeNESW, consts.DLeftBottom)
		m.rightBottom = m.newDragPanel(owner, types.CrSizeNWSE, consts.DRightBottom)
	} else {
		m.right = m.newDragPanel(owner, types.CrSizeWE, consts.DRight)
		m.bottom = m.newDragPanel(owner, types.CrSizeNS, consts.DBottom)
		m.rightBottom = m.newDragPanel(owner, types.CrSizeNWSE, consts.DRightBottom)
	}
	return m
}

// 设置关联组件
func (m *drag) SetRelation(relation *TDesigningComponent) {
	m.relation = relation
}

// 隐藏所有
func (m *drag) Hide() {
	if !m.isShow {
		return
	}
	m.relation.isDesigner = false
	m.isShow = false
	if m.ds == consts.DsAll {
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
	if m.isShow {
		return
	}
	m.relation.isDesigner = true
	m.isShow = true
	if m.ds == consts.DsAll {
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
	if m.ds == consts.DsAll {
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

func (m *drag) SetParent(value lcl.IWinControl) {
	if m.ds == consts.DsAll {
		m.left.SetParent(value)
		m.top.SetParent(value)
		m.right.SetParent(value)
		m.bottom.SetParent(value)
		m.leftTop.SetParent(value)
		m.rightTop.SetParent(value)
		m.leftBottom.SetParent(value)
		m.rightBottom.SetParent(value)
	} else {
		m.right.SetParent(value)
		m.bottom.SetParent(value)
		m.rightBottom.SetParent(value)
	}
}

// 跟随关联组件
func (m *drag) Follow() {
	if m.relation != nil {
		br := m.relation.BoundsRect()
		// 转换为 form tab 的坐标
		point := m.relation.ClientToParent(types.TPoint{X: 0, Y: 0}, m.relation.formTab.scroll)
		x, y := point.X, point.Y
		width, height := br.Width(), br.Height()
		db := dragBorder / 2
		if m.ds == consts.DsAll {
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

// 设计组件鼠标移动
func (m *drag) OnMouseMove(sender *TDesigningComponent, shift types.TShiftState, X int32, Y int32) {
	br := sender.BoundsRect()
	hint := fmt.Sprintf(`%v
	Left: %v Top: %v
	Width: %v Height: %v`, sender.TreeName(), br.Left, br.Top, br.Width(), br.Height())
	sender.SetHint(hint)
	if m.isDown {
		m.Hide()
		point := sender.ClientToParent(types.TPoint{X: X, Y: Y}, sender.formTab.FormRoot.object)
		x := point.X - m.dx
		y := point.Y - m.dy
		sender.SetBounds(m.dcl+x, m.dct+y, br.Width(), br.Height())
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", m.dcl+x, m.dct+y, br.Width(), br.Height())
		message.Follow(msgContent)
		go sender.UpdateNodeDataPoint(br.Left, br.Top)
	}
}

// 设计组件鼠标按下事件
func (m *drag) OnMouseDown(sender *TDesigningComponent, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", sender.ClassName())
	if !sender.formTab.placeComponent(sender, X, Y) {
		m.isDown = true
		point := sender.ClientToParent(types.TPoint{X: X, Y: Y}, sender.formTab.FormRoot.object)
		m.dx, m.dy = point.X, point.Y
		br := sender.BoundsRect()
		m.dcl = br.Left
		m.dct = br.Top
		// 更新设计查看器的属性信息
		sender.formTab.switchComponentEditing(sender)
		// 更新设计查看器的组件树信息
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			// 设置选中状态
			sender.SetSelected()
		})
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
		message.Follow(msgContent)
		if sender.object != nil {
			lcl.Mouse.SetCapture(sender.object.Handle())
		}
		sender.DragBegin()
	}
}

// 设计组件鼠标抬起事件
func (m *drag) OnMouseUp(sender *TDesigningComponent, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDown {
		m.Show()
	}
	br := sender.BoundsRect()
	go sender.UpdateNodeDataPoint(br.Left, br.Top)
	m.isDown = false
	message.FollowHide()
	lcl.Mouse.SetCapture(0)
	sender.DragEnd()
}
