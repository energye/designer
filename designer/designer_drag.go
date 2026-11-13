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
	"github.com/energye/lcl/tool"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 拖拽控制

var dragBorder int32 = 6

// 最小移动距离阈值
const minDistance = 4

// 拖拽控制
type drag struct {
	relation           *TDesigningComponent  // 关联设计的组件
	ds                 consts.DragShowStatus // 显示方向
	dx, dy             int32                 // 拖拽控制
	dcl, dct           int32                 // 拖拽控制
	dcx, dcy, dcw, dch int32                 // 拖拽控制
	isDown             bool                  // 拖拽控制
	owner              lcl.IWinControl       // 拖拽控制所属组件, 非空表示未创建 ds(拖拽控制)空表示创建过
	left               lcl.IPanel
	top                lcl.IPanel
	right              lcl.IPanel
	bottom             lcl.IPanel
	leftTop            lcl.IPanel
	rightTop           lcl.IPanel
	leftBottom         lcl.IPanel
	rightBottom        lcl.IPanel
}

func (m *drag) Free() {
	m.relation = nil
	m.owner = nil
	m.left = nil
	m.top = nil
	m.right = nil
	m.bottom = nil
	m.leftTop = nil
	m.rightTop = nil
	m.leftBottom = nil
	m.rightBottom = nil
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
	//_, _ = dx, dy
	//_, _, _, _ = dcx, dcy, dcw, dch
	pnl.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, X int32, Y int32) {
		if m.isDown {
			dx := X - m.dx
			dy := Y - m.dy
			left, top := m.dcx, m.dcy
			right, bottom := m.dcx+m.dcw, m.dcy+m.dch
			moveLeft, moveTop, moveRight, moveBottom := false, false, false, false
			switch d {
			case consts.DLeft:
				left += dx
				moveLeft = true
			case consts.DTop:
				top += dy
				moveTop = true
			case consts.DRight:
				right += dx
				moveRight = true
			case consts.DBottom:
				bottom += dy
				moveBottom = true
			case consts.DLeftTop:
				left += dx
				top += dy
				moveLeft, moveTop = true, true
			case consts.DRightTop:
				right += dx
				top += dy
				moveRight, moveTop = true, true
			case consts.DLeftBottom:
				left += dx
				bottom += dy
				moveLeft, moveBottom = true, true
			case consts.DRightBottom:
				right += dx
				bottom += dy
				moveRight, moveBottom = true, true
			}
			if m.snapGridEnabled(shift) {
				if moveLeft {
					left = snapDesignerGrid(left)
				}
				if moveTop {
					top = snapDesignerGrid(top)
				}
				if moveRight {
					right = snapDesignerGrid(right)
				}
				if moveBottom {
					bottom = snapDesignerGrid(bottom)
				}
			}
			minSize := int32(designerGridSize)
			if moveLeft && right-left < minSize {
				left = right - minSize
			}
			if moveRight && right-left < minSize {
				right = left + minSize
			}
			if moveTop && bottom-top < minSize {
				top = bottom - minSize
			}
			if moveBottom && bottom-top < minSize {
				bottom = top + minSize
			}
			m.relation.SetBounds(left, top, right-left, bottom-top)
			br := m.relation.BoundsRect()
			m.relation.UpdateNodeDataBounds(br)
			msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
			message.Follow(msgContent)
			m.relation.SetStatusText(br.Left, br.Top, br.Width(), br.Height())
		}
	})
	pnl.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		logs.Debug("DRAG OnMouseDown direction:", d)
		if tool.IsWindows() {
			m.Hide()
		}
		m.dx, m.dy = X, Y
		br := m.relation.BoundsRect()
		m.dcx, m.dcy, m.dcw, m.dch = br.Left, br.Top, br.Width(), br.Height()
		m.isDown = true
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", m.dcx, m.dcy, m.dcw, m.dch)
		message.Follow(msgContent)
		m.relation.DragBegin()
	})
	pnl.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
		logs.Debug("DRAG OnMouseUP direction:", d)
		if tool.IsWindows() {
			m.Show()
		} else {
			m.Follow()
		}
		m.isDown = false
		message.FollowHide()
		br := m.relation.BoundsRect()
		m.relation.UpdateNodeDataBounds(br)
		m.relation.DragEnd()
	})
	return pnl
}

// snapDesignerGrid 四舍五入到设计网格。
func snapDesignerGrid(v int32) int32 {
	grid := int32(designerGridSize)
	if v >= 0 {
		return (v + grid/2) / grid * grid
	}
	return -((-v + grid/2) / grid * grid)
}

// snapGridEnabled 网格显示时启用吸附，按住 Alt 临时关闭。
func (m *drag) snapGridEnabled(shift types.TShiftState) bool {
	if m.relation == nil || m.relation.FormTab == nil {
		return false
	}
	return m.relation.FormTab.gridVisible && !shift.In(types.SsAlt) && !altHeld()
}

func newDrag(owner lcl.IWinControl, ds consts.DragShowStatus) *drag {
	m := new(drag)
	m.ds = ds
	m.owner = owner
	return m
}

// mustDS 初始化拖拽面板控件
// 根据拖拽敏感度设置为 drag 结构创建拖拽面板，用于不同方向的调整大小操作
func (m *drag) mustDS() {
	if m.owner == nil {
		return
	}
	if m.ds == consts.DsAll {
		m.left = m.newDragPanel(m.owner, types.CrSizeWE, consts.DLeft)
		m.top = m.newDragPanel(m.owner, types.CrSizeNS, consts.DTop)
		m.right = m.newDragPanel(m.owner, types.CrSizeWE, consts.DRight)
		m.bottom = m.newDragPanel(m.owner, types.CrSizeNS, consts.DBottom)
		m.leftTop = m.newDragPanel(m.owner, types.CrSizeNWSE, consts.DLeftTop)
		m.rightTop = m.newDragPanel(m.owner, types.CrSizeNESW, consts.DRightTop)
		m.leftBottom = m.newDragPanel(m.owner, types.CrSizeNESW, consts.DLeftBottom)
		m.rightBottom = m.newDragPanel(m.owner, types.CrSizeNWSE, consts.DRightBottom)
	} else {
		m.right = m.newDragPanel(m.owner, types.CrSizeWE, consts.DRight)
		m.bottom = m.newDragPanel(m.owner, types.CrSizeNS, consts.DBottom)
		m.rightBottom = m.newDragPanel(m.owner, types.CrSizeNWSE, consts.DRightBottom)
	}
	m.owner = nil
}

// 设置关联组件
func (m *drag) SetRelation(relation *TDesigningComponent) {
	m.relation = relation
}

// 隐藏所有
func (m *drag) Hide() {
	if m.owner != nil {
		return
	}
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
	if m.owner != nil {
		return
	}
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
	if m.owner != nil {
		return
	}
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
	if m.owner != nil {
		return
	}
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
	if m.owner != nil {
		return
	}
	if m.relation != nil {
		br := m.relation.BoundsRect()
		// 转换为 form tab 的坐标
		point := m.GetTopLeft(m.relation)
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
func (m *drag) GetTopLeft(component *TDesigningComponent) (result types.TPoint) {
	form := m.relation.FormTab.FormRoot
	parent := component.Parent()
	if form == nil || parent == nil {
		return types.NewPoint(0, 0)
	}
	result = parent.Control().ClientOrigin()
	formOrigin := form.Control().ClientOrigin()
	result.X = result.X - formOrigin.X + component.Control().Left()
	result.Y = result.Y - formOrigin.Y + component.Control().Top()
	return
}

// 设计组件鼠标移动
func (m *drag) OnMouseMove(sender *TDesigningComponent, shift types.TShiftState, X int32, Y int32) {
	br := sender.BoundsRect()
	//hint := fmt.Sprintf(`%v
	//Left: %v Top: %v
	//Width: %v Height: %v`, sender.TreeName(), br.Left, br.Top, br.Width(), br.Height())
	//message.Follow(hint)
	if m.isDown && sender.ComponentType != consts.CtForm {
		m.Hide()
		//point := sender.ClientToParent(types.TPoint{X: X, Y: Y}, sender.FormTab.FormRoot.object)
		point := lcl.Mouse.CursorPos()
		newLeft := m.dcl + point.X - m.dx
		newTop := m.dct + point.Y - m.dy
		if m.snapGridEnabled(shift) {
			newLeft = snapDesignerGrid(newLeft)
			newTop = snapDesignerGrid(newTop)
		}
		sender.SetBounds(newLeft, newTop, br.Width(), br.Height())
		// SetBounds 可能被父容器夹紧，属性以控件最终边界为准。
		actual := sender.BoundsRect()
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", actual.Left, actual.Top, actual.Width(), actual.Height())
		message.Follow(msgContent)
		sender.SetStatusText(actual.Left, actual.Top, actual.Width(), actual.Height())
		sender.UpdateNodeDataPoint(actual.Left, actual.Top)
	}
}

// 设计组件鼠标按下事件
func (m *drag) OnMouseDown(sender *TDesigningComponent, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", sender.ClassName())
	m.relation.mustComponentPropertyPage()
	m.mustDS()
	if !sender.FormTab.placeComponent(sender, X, Y) {
		m.isDown = true
		//point := sender.ClientToParent(types.TPoint{X: X, Y: Y}, sender.FormTab.FormRoot.object)
		point := lcl.Mouse.CursorPos()
		m.dx, m.dy = point.X, point.Y
		br := sender.BoundsRect()
		m.dcl = br.Left
		m.dct = br.Top
		// 更新设计查看器的属性信息
		sender.FormTab.SwitchComponentEditing(sender)
		// 更新设计查看器的组件树信息
		go lcl.RunOnMainThreadAsync(func(id uint32) {
			// 设置选中状态
			sender.SetSelected()
		})
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
		message.Follow(msgContent)
		sender.DragBegin()
	}
}

// 设计组件鼠标抬起事件
func (m *drag) OnMouseUp(sender *TDesigningComponent, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.isDown {
		m.Show()
		// 只有真正发生过拖动才补写最终坐标，普通点击不触发 UI 生成。
		br := sender.BoundsRect()
		sender.UpdateNodeDataPoint(br.Left, br.Top)
		m.isDown = false
	}
	message.FollowHide()
	sender.DragEnd()
}
