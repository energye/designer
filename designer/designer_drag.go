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
	"sync/atomic"
	"time"

	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/message"
	"github.com/energye/lcl/types"
)

// 拖拽控制

var dragBorder int32 = 6

// 拖拽控制
// 缩放控制点绘制在页签级覆盖层(玻璃层, 见 designer_overlay.go):
// 普通组件(DsAll)与窗体根节点(DsRightBottom)的控制点/次选边框都在覆盖层上,
// 置顶于所有设计组件之上 —— 被覆盖组件的锚点可见可操作; 鼠标事件里按
// dragHandleAt 命中测试, 不为任何组件创建锚点子窗口。
// 无窗口区域平台(Cocoa)回退容器画布绘制(paintContainerOverlays)。
type drag struct {
	relation           *TDesigningComponent  // 关联设计的组件
	ds                 consts.DragShowStatus // 显示方向
	resizeDir          int                   // 缩放方向(consts.DLeft..DRightBottom), -1 表示非缩放
	dx, dy             int32                 // 缩放起始鼠标位置(相对设计窗体)
	dcl, dct           int32                 // 移动拖拽起始组件位置
	dcx, dcy, dcw, dch int32                 // 缩放起始组件边界
	isDown             bool                  // 移动拖拽按下中
	moved              bool                  // 本次按下是否真的改变了组件边界
	// 对齐目标缓存，拖动过程中复用避免频繁分配
	alignXs []designerAlignTarget
	alignYs []designerAlignTarget
	// 移动节流: 上次应用的位置(吸附后), 网格步进内跳过整组 SetBounds
	lastMoveX, lastMoveY int32
	// 缩放节流: 上次应用的边界
	lastResizeL, lastResizeT, lastResizeR, lastResizeB int32
	// 多选拖动时其余选中组件的起始边界
	groupItems []dragGroupItem
}

type dragGroupItem struct {
	comp           *TDesigningComponent
	left, top      int32
	width, height  int32
	designStripped bool // 拖动期间被临时摘除设计模式, 结束时恢复
}

// groupClassNames 返回多选拖动组内组件类名(诊断用)
func (m *drag) groupClassNames() []string {
	names := make([]string, 0, len(m.groupItems))
	for i := range m.groupItems {
		names = append(names, m.groupItems[i].comp.ClassName())
	}
	return names
}

// flushParentRepaints 同步冲刷拖动组件父窗口的待重绘脏区。
// 窗口位移只会把旧位置的暴露区域排入父窗口重绘队列(异步),
// 拖动过程中队列持续繁忙(设计模式通知/覆盖层/属性刷新), 重绘被推迟,
// 高对比背景的控件(如白底输入框)旧图像滞留可见(残影);
// Update 只处理已排队的脏区, 不做全量重绘, 把补绘提前到本帧完成。
func (m *drag) flushParentRepaints() {
	for i := range m.groupItems {
		p := m.groupItems[i].comp.parent
		if p != nil && p.object != nil && p.object.IsValid() {
			p.object.Update()
		}
	}
	r := m.relation
	if r != nil && r.parent != nil && r.parent.object != nil && r.parent.object.IsValid() {
		r.parent.object.Update()
	}
}

// restoreGroupDesignModes 恢复拖动期间被摘除设计模式的组件。
// mouseup 可能被缩放控制点面板截获而走不到 OnMouseUp, 必须在所有结束路径上恢复。
func (m *drag) restoreGroupDesignModes() {
	restored := false
	for i := range m.groupItems {
		if m.groupItems[i].designStripped {
			RestoreComponentDesignMode(m.groupItems[i].comp.Instance())
			m.groupItems[i].designStripped = false
			restored = true
		}
	}
	if restored && m.relation != nil && m.relation.FormTab != nil {
		// 设计模式切换改变组件的渲染路径, 恢复后必须强制重绘,
		// 否则组件保持拖动期间的旧渲染状态(参见 EnableDeferredDesignMode),
		// 表现为旧位置残影/视觉上未跟随移动
		if box := m.relation.FormTab.designerBox; box != nil && box.IsValid() {
			box.Invalidate()
		}
	}
}

// cancelStaleDrag 结束可能残留的移动拖拽状态。
// MouseUp 被缩放控制点截获时 OnMouseUp 不会执行, isDown/groupItems/设计模式摘除
// 会一直残留, 后续 mousemove 会驱动陈旧拖拽移动整组组件。
func (m *drag) cancelStaleDrag() {
	if !m.isDown && m.resizeDir < 0 {
		return
	}
	m.restoreGroupDesignModes()
	m.groupItems = m.groupItems[:0]
	m.moved = false
	m.isDown = false
	m.resizeDir = -1
	atomic.StoreInt32(&overlayCursorResize, -1)
	resetOverlayCursor()
	// 拖拽中可能带着坐标信息气泡(Mouse 移动/缩放展示), 随残留状态一并隐藏
	message.FollowHide()
	if m.relation != nil && m.relation.FormTab != nil {
		m.relation.FormTab.endSelectionDrag()
		m.relation.FormTab.hideGuides()
	}
}

func (m *drag) Free() {
	m.relation = nil
	m.resizeDir = -1
	m.groupItems = nil
}

// dragHandleRect 单个缩放控制点: 方向 + 矩形(父容器客户区坐标)
type dragHandleRect struct {
	dir  int
	rect types.TRect
}

// resizeCursor 返回缩放方向对应的光标
func resizeCursor(dir int) types.TCursor {
	switch dir {
	case consts.DLeft, consts.DRight:
		return types.CrSizeWE
	case consts.DTop, consts.DBottom:
		return types.CrSizeNS
	case consts.DLeftTop, consts.DRightBottom:
		return types.CrSizeNWSE
	case consts.DRightTop, consts.DLeftBottom:
		return types.CrSizeNESW
	}
	return types.CrDefault
}

// handleRects 计算关联组件的缩放控制点矩形(父容器客户区坐标)。
// 控制点绘制在组件外侧 1px 起, 完全不与组件窗口重叠。
func (m *drag) handleRects() []dragHandleRect {
	if m.relation == nil {
		return nil
	}
	control := m.relation.Control()
	if control == nil || !control.IsValid() {
		return nil
	}
	br := control.BoundsRect()
	x, y, w, h := br.Left, br.Top, br.Width(), br.Height()
	db := dragBorder
	var dirs []int
	if m.ds == consts.DsAll {
		// 角控制点优先命中(小尺寸组件时角与边可能重叠)
		dirs = []int{consts.DLeftTop, consts.DRightTop, consts.DLeftBottom, consts.DRightBottom,
			consts.DLeft, consts.DTop, consts.DRight, consts.DBottom}
	} else {
		dirs = []int{consts.DRightBottom, consts.DRight, consts.DBottom}
	}
	rects := make([]dragHandleRect, 0, len(dirs))
	for _, dir := range dirs {
		var r types.TRect
		switch dir {
		case consts.DLeft:
			r = types.Rect(x-db, y+h/2-db, x, y+h/2+db)
		case consts.DTop:
			r = types.Rect(x+w/2-db, y-db, x+w/2+db, y)
		case consts.DRight:
			r = types.Rect(x+w, y+h/2-db, x+w+db, y+h/2+db)
		case consts.DBottom:
			r = types.Rect(x+w/2-db, y+h, x+w/2+db, y+h+db)
		case consts.DLeftTop:
			r = types.Rect(x-db, y-db, x, y)
		case consts.DRightTop:
			r = types.Rect(x+w, y-db, x+w+db, y)
		case consts.DLeftBottom:
			r = types.Rect(x-db, y+h, x, y+h+db)
		case consts.DRightBottom:
			r = types.Rect(x+w, y+h, x+w+db, y+h+db)
		}
		rects = append(rects, dragHandleRect{dir: dir, rect: r})
	}
	return rects
}

// beginResize 从画布控制点命中开始缩放, pos 为设计窗体客户区坐标
func (m *drag) beginResize(dir int, pos types.TPoint) {
	if m.relation == nil {
		return
	}
	logs.Debug("PERF beginResize dir:", dir, "pos:", pos.X, pos.Y)
	m.resizeDir = dir
	atomic.StoreInt32(&overlayCursorResize, int32(dir))
	m.dx, m.dy = pos.X, pos.Y
	br := m.relation.BoundsRect()
	m.dcx, m.dcy, m.dcw, m.dch = br.Left, br.Top, br.Width(), br.Height()
	m.lastResizeL, m.lastResizeT, m.lastResizeR, m.lastResizeB = br.Left, br.Top, br.Right, br.Bottom
	m.isDown = false
	m.moved = false
	if m.relation.FormTab != nil {
		m.relation.FormTab.hideGuides()
	}
	msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", m.dcx, m.dcy, m.dcw, m.dch)
	message.Follow(msgContent)
	m.relation.DragBegin()
}

// applyResize 缩放中: 按方向调整组件边界(含对齐吸附/网格吸附/最小尺寸)。
// X/Y 为设计画布客户区坐标(消息坐标换算, 与组件边界同一坐标系)。
func (m *drag) applyResize(shift types.TShiftState, X, Y int32) {
	if m.relation == nil || m.relation.FormTab == nil {
		return
	}
	dx := X - m.dx
	dy := Y - m.dy
	left, top := m.dcx, m.dcy
	right, bottom := m.dcx+m.dcw, m.dcy+m.dch
	moveLeft, moveTop, moveRight, moveBottom := false, false, false, false
	switch m.resizeDir {
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
	targets := m.alignmentTargets(m.relation)
	alignEnabled := !altHeld() && !shift.In(types.SsAlt)
	gridEnabled := m.snapGridEnabled(shift)
	var matchV, matchH designerAlignTarget
	hitV, hitH := false, false
	if alignEnabled {
		if moveLeft {
			if edge, target, hit := snapEdgeToTargets(left, targets.xs, designerAlignThreshold); hit {
				left = edge
				matchV, hitV = target, true
			} else if gridEnabled {
				left = snapDesignerGrid(left)
			}
		} else if moveRight {
			if edge, target, hit := snapEdgeToTargets(right, targets.xs, designerAlignThreshold); hit {
				right = edge
				matchV, hitV = target, true
			} else if gridEnabled {
				right = snapDesignerGrid(right)
			}
		}
		if moveTop {
			if edge, target, hit := snapEdgeToTargets(top, targets.ys, designerAlignThreshold); hit {
				top = edge
				matchH, hitH = target, true
			} else if gridEnabled {
				top = snapDesignerGrid(top)
			}
		} else if moveBottom {
			if edge, target, hit := snapEdgeToTargets(bottom, targets.ys, designerAlignThreshold); hit {
				bottom = edge
				matchH, hitH = target, true
			} else if gridEnabled {
				bottom = snapDesignerGrid(bottom)
			}
		}
	} else if gridEnabled {
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
	clampedX, clampedY := false, false
	if moveLeft && right-left < minSize {
		left = right - minSize
		clampedX = true
	}
	if moveRight && right-left < minSize {
		right = left + minSize
		clampedX = true
	}
	if moveTop && bottom-top < minSize {
		top = bottom - minSize
		clampedY = true
	}
	if moveBottom && bottom-top < minSize {
		bottom = top + minSize
		clampedY = true
	}
	if clampedX {
		hitV = false
	}
	if clampedY {
		hitH = false
	}
	// 步进内: 边界未变化时跳过本次事件(整组/单组件 SetBounds 代价高)
	if m.lastResizeL == left && m.lastResizeT == top && m.lastResizeR == right && m.lastResizeB == bottom {
		return
	}
	m.lastResizeL, m.lastResizeT, m.lastResizeR, m.lastResizeB = left, top, right, bottom
	m.relation.SetBounds(left, top, right-left, bottom-top)
	// 缩放位移同样冲刷父窗口暴露区域, 避免拖尾
	if p := m.relation.parent; p != nil && p.object != nil && p.object.IsValid() {
		p.object.Update()
	}
	if m.relation.FormTab != nil {
		var guidesV, guidesH []designerGuide
		if hitV {
			guidesV = append(guidesV, guideFor(top, bottom-top, matchV))
		}
		if hitH {
			guidesH = append(guidesH, guideFor(left, right-left, matchH))
		}
		if len(guidesV) > 0 || len(guidesH) > 0 {
			m.relation.FormTab.showGuides(guidesV, guidesH, targets.origin)
		} else {
			m.relation.FormTab.hideGuides()
		}
	}
	br := m.relation.BoundsRect()
	if !m.moved && (br.Left != m.dcx || br.Top != m.dcy || br.Width() != m.dcw || br.Height() != m.dch) {
		m.moved = true
	}
	// 缩放中只刷新属性面板，不触发 uigen，避免 .ui 在鼠标未松开时落盘。
	m.relation.SyncBoundsToProperties(br)
	// 覆盖层窗口(锚点绘制于其上)跟随主选组件边界更新区域与重绘
	m.relation.FormTab.updateOverlayPane()
	msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
	message.Follow(msgContent)
	m.relation.SetStatusText(br.Left, br.Top, br.Width(), br.Height())
}

const designerAlignThreshold = 5

// designerAlignTarget 单个对齐目标。
// pos 为吸附坐标；from/to 为对齐对象沿另一轴的显示范围，用于绘制局部对齐线。
type designerAlignTarget struct {
	pos      int32
	from, to int32
}

// designerAlignment 某父容器内的对齐目标（父边/中心 + 兄弟控件边/中心）。
type designerAlignment struct {
	xs     []designerAlignTarget
	ys     []designerAlignTarget
	origin types.TPoint // 父容器客户区 -> 设计画布客户区偏移
}

// snapEdgeToTargets 找到最近的对齐目标；命中时返回目标坐标。
func snapEdgeToTargets(edge int32, targets []designerAlignTarget, threshold int32) (newEdge int32, matched designerAlignTarget, hit bool) {
	newEdge = edge
	bestDelta := threshold + 1
	for _, target := range targets {
		delta := target.pos - edge
		if delta < 0 {
			delta = -delta
		}
		if delta <= threshold && delta < bestDelta {
			bestDelta = delta
			newEdge = target.pos
			matched = target
			hit = true
		}
	}
	return
}

// snapMoveAxis 用矩形的起点/中心/终点匹配对齐目标，返回吸附后的起点。
func snapMoveAxis(start, size int32, targets []designerAlignTarget, threshold int32) (newStart int32, matched designerAlignTarget, hit bool) {
	newStart = start
	anchors := [3]int32{start, start + size/2, start + size}
	bestDelta := threshold + 1
	for _, anchor := range anchors {
		for _, target := range targets {
			delta := target.pos - anchor
			if delta < 0 {
				delta = -delta
			}
			if delta <= threshold && delta < bestDelta {
				bestDelta = delta
				newStart = start + (target.pos - anchor)
				matched = target
				hit = true
			}
		}
	}
	return
}

// guideFor 生成对齐线：绘制范围为拖动控件与对齐目标的联合区间。
func guideFor(start, size int32, target designerAlignTarget) designerGuide {
	guide := designerGuide{pos: target.pos, from: target.from, to: target.to}
	if start < guide.from {
		guide.from = start
	}
	if start+size > guide.to {
		guide.to = start + size
	}
	return guide
}

func (m *drag) alignmentTargets(comp *TDesigningComponent) designerAlignment {
	var result designerAlignment
	if comp == nil || comp.parent == nil {
		return result
	}
	parent := comp.parent
	parentObj := parent.object
	if parentObj == nil || !parentObj.IsValid() {
		return result
	}
	width, height := parentObj.ClientWidth(), parentObj.ClientHeight()
	result.xs = m.alignXs[:0]
	result.ys = m.alignYs[:0]
	// 父容器参考线贯穿整个父容器。
	result.xs = append(result.xs,
		designerAlignTarget{pos: 0, from: 0, to: height},
		designerAlignTarget{pos: width / 2, from: 0, to: height},
		designerAlignTarget{pos: width, from: 0, to: height})
	result.ys = append(result.ys,
		designerAlignTarget{pos: 0, from: 0, to: width},
		designerAlignTarget{pos: height / 2, from: 0, to: width},
		designerAlignTarget{pos: height, from: 0, to: width})
	for _, sibling := range parent.Child {
		if sibling == nil || sibling == comp || sibling.ComponentType == consts.CtNonVisual {
			continue
		}
		control := sibling.Control()
		if control == nil || !control.IsValid() || !control.Visible() {
			continue
		}
		br := sibling.BoundsRect()
		if br.Width() <= 0 || br.Height() <= 0 {
			continue
		}
		result.xs = append(result.xs,
			designerAlignTarget{pos: br.Left, from: br.Top, to: br.Bottom},
			designerAlignTarget{pos: br.Left + br.Width()/2, from: br.Top, to: br.Bottom},
			designerAlignTarget{pos: br.Left + br.Width(), from: br.Top, to: br.Bottom})
		result.ys = append(result.ys,
			designerAlignTarget{pos: br.Top, from: br.Left, to: br.Right},
			designerAlignTarget{pos: br.Top + br.Height()/2, from: br.Left, to: br.Right},
			designerAlignTarget{pos: br.Top + br.Height(), from: br.Left, to: br.Right})
	}
	if comp.FormTab != nil {
		box := comp.FormTab.designerBox
		if box != nil && box.IsValid() && box.Instance() != parentObj.Instance() {
			po, bo := parentObj.ClientOrigin(), box.ClientOrigin()
			result.origin = types.TPoint{X: po.X - bo.X, Y: po.Y - bo.Y}
		}
	}
	m.alignXs, m.alignYs = result.xs, result.ys
	return result
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

func newDrag(ds consts.DragShowStatus) *drag {
	m := new(drag)
	m.ds = ds
	m.resizeDir = -1
	return m
}

// 设置关联组件
func (m *drag) SetRelation(relation *TDesigningComponent) {
	m.relation = relation
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
	// 缩放中: 控制点命中后由 applyResize 处理 (X/Y 为设计画布坐标)
	if m.resizeDir >= 0 {
		m.applyResize(shift, X, Y)
		return
	}
	br := sender.BoundsRect()
	//hint := fmt.Sprintf(`%v
	//Left: %v Top: %v
	//Width: %v Height: %v`, sender.TreeName(), br.Left, br.Top, br.Width(), br.Height())
	//message.Follow(hint)
	if m.isDown && sender.ComponentType != consts.CtForm {
		perfStart := time.Now()
		defer func() {
			logs.Debug("PERF OnMouseMove group:", len(m.groupItems), "cost:", time.Since(perfStart))
		}()
		// X/Y 已是设计画布坐标(设计钩子换算), 直接用于移动计算
		newLeft := m.dcl + X - m.dx
		newTop := m.dct + Y - m.dy
		targets := m.alignmentTargets(sender)
		alignEnabled := !altHeld() && !shift.In(types.SsAlt)
		gridEnabled := m.snapGridEnabled(shift)
		var matchV, matchH designerAlignTarget
		hitV, hitH := false, false
		if alignEnabled {
			if v, target, hit := snapMoveAxis(newLeft, br.Width(), targets.xs, designerAlignThreshold); hit {
				newLeft = v
				matchV, hitV = target, true
			} else if gridEnabled {
				newLeft = snapDesignerGrid(newLeft)
			}
			if v, target, hit := snapMoveAxis(newTop, br.Height(), targets.ys, designerAlignThreshold); hit {
				newTop = v
				matchH, hitH = target, true
			} else if gridEnabled {
				newTop = snapDesignerGrid(newTop)
			}
		} else if gridEnabled {
			newLeft = snapDesignerGrid(newLeft)
			newTop = snapDesignerGrid(newTop)
		}
		// 网格/吸附步进内: 目标位置未变化时跳过本次事件,
		// 避免鼠标每移动 1px 都触发整组组件 SetBounds(多选拖动卡顿根因)。
		if m.lastMoveX == newLeft && m.lastMoveY == newTop {
			return
		}
		m.lastMoveX, m.lastMoveY = newLeft, newTop
		// 若当前组件本身位于某个被选中祖先内，直接移动它会让祖先和它各移动一次。
		// 这种情况改为移动组内祖先，子组件跟随祖先移动。
		senderHasSelectedAncestor := sender.FormTab != nil && sender.FormTab.HasSelectedAncestor(sender)
		var actual types.TRect
		if !senderHasSelectedAncestor {
			sender.SetBounds(newLeft, newTop, br.Width(), br.Height())
			// SetBounds 可能被父容器夹紧，属性以控件最终边界为准。
			actual = sender.BoundsRect()
		}
		if len(m.groupItems) > 0 {
			dx, dy := newLeft-m.dcl, newTop-m.dct
			if !senderHasSelectedAncestor {
				dx, dy = actual.Left-m.dcl, actual.Top-m.dct
			}
			var anchor *dragGroupItem
			if senderHasSelectedAncestor {
				for i := range m.groupItems {
					if sender.FormTab.IsComponentAncestor(m.groupItems[i].comp, sender) {
						anchor = &m.groupItems[i]
						break
					}
				}
			}
			if anchor != nil {
				anchor.comp.SetBounds(anchor.left+dx, anchor.top+dy, anchor.width, anchor.height)
				anchorActual := anchor.comp.BoundsRect()
				dx, dy = anchorActual.Left-anchor.left, anchorActual.Top-anchor.top
			}
			for _, item := range m.groupItems {
				if anchor != nil && item.comp == anchor.comp {
					continue
				}
				item.comp.SetBounds(item.left+dx, item.top+dy, item.width, item.height)
			}
		}
		// 位移后同步冲刷父窗口暴露区域, 消除拖动残影(单选/多选共用)
		m.flushParentRepaints()
		if senderHasSelectedAncestor {
			actual = sender.BoundsRect()
		}
		// 记录是否真的发生位移：普通点击（含网格吸附后回到原位）不应触发写盘。
		if !m.moved {
			if actual.Left != m.dcl || actual.Top != m.dct {
				m.moved = true
			} else {
				for i := range m.groupItems {
					cur := m.groupItems[i].comp.BoundsRect()
					if cur.Left != m.groupItems[i].left || cur.Top != m.groupItems[i].top {
						m.moved = true
						break
					}
				}
			}
		}
		if sender.FormTab != nil {
			var guidesV, guidesH []designerGuide
			if hitV {
				guidesV = append(guidesV, guideFor(actual.Top, actual.Height(), matchV))
			}
			if hitH {
				guidesH = append(guidesH, guideFor(actual.Left, actual.Width(), matchH))
			}
			if len(guidesV) > 0 || len(guidesH) > 0 {
				sender.FormTab.showGuides(guidesV, guidesH, targets.origin)
			} else {
				sender.FormTab.hideGuides()
			}
			// 大量选中时每步刷新数百个边框面板会拖垮 UI 线程：
			// 拖动期间隐藏次选边框，MouseUp 时统一恢复一次。
			sender.FormTab.beginSelectionDrag()
		}
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", actual.Left, actual.Top, actual.Width(), actual.Height())
		message.Follow(msgContent)
		sender.SetStatusText(actual.Left, actual.Top, actual.Width(), actual.Height())
		// 拖动中只刷新属性面板，不触发 uigen，避免 .ui 在鼠标未松开时落盘。
		sender.SyncPointToProperties(actual.Left, actual.Top)
	}
}

// 设计组件鼠标按下事件
func (m *drag) OnMouseDown(sender *TDesigningComponent, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	logs.Debug("OnMouseDown 设计组件", sender.ClassName())
	m.relation.mustComponentPropertyPage()
	// placeComponent 需要点击点相对组件父容器客户区的坐标: 画布坐标减去组件画布偏移
	px, py := X, Y
	if sctrl := sender.Control(); sctrl != nil && sctrl.IsValid() {
		off := sender.FormTab.overlayOffset(sender, sender.FormTab.FormRoot)
		px -= off.X
		py -= off.Y
	}
	if !sender.FormTab.placeComponent(sender, px, py) {
		formTab := sender.FormTab
		if sender.ComponentType != consts.CtForm {
			// Ctrl 切换选中, Shift 追加选中; 无修饰键时保持已有多选否则单选
			if shift.In(types.SsCtrl) {
				if !formTab.ToggleComponentSelection(sender) {
					return
				}
			} else if shift.In(types.SsShift) {
				formTab.AddComponentSelection(sender)
			} else if !formTab.IsComponentSelected(sender) || len(formTab.SelectedComponents()) <= 1 {
				formTab.SwitchComponentEditing(sender)
			} else {
				formTab.SelectComponents(formTab.SelectedComponents(), sender)
			}
		}
		m.isDown = true
		m.moved = false
		// 恢复上一次未正常结束拖动遗留的摘除状态, 再重建多选拖动组
		m.restoreGroupDesignModes()
		// 记录多选拖动组，拖动时整组同 delta 位移
		m.groupItems = m.groupItems[:0]
		if selected := formTab.SelectedComponents(); len(selected) > 1 && formTab.IsComponentSelected(sender) && sender.ComponentType != consts.CtForm {
			for _, comp := range selected {
				if comp == sender || formTab.HasSelectedAncestor(comp) {
					continue
				}
				br := comp.BoundsRect()
				m.groupItems = append(m.groupItems, dragGroupItem{comp: comp, left: br.Left, top: br.Top, width: br.Width(), height: br.Height()})
			}
			// 设计模式下组件真实 bounds 变更每次要走消息通知路径(20ms+),
			// 整组拖动时摘除次选组件的设计模式走快速路径。
			// 主选组件保留设计模式: 拖动中的鼠标消息依赖它路由。
			for i := range m.groupItems {
				if instance := m.groupItems[i].comp.Instance(); instance != 0 {
					StripComponentDesignMode(instance)
					m.groupItems[i].designStripped = true
				}
				// 摘除切换了渲染路径, 强制重绘一次,
				// 避免拖动首个帧保留设计模式的旧渲染状态(残影来源之一)
				if ctrl := m.groupItems[i].comp.Control(); ctrl != nil && ctrl.IsValid() {
					ctrl.Invalidate()
				}
			}
			if perfDiag {
				logs.Info("PERF groupDrag primary:", sender.ClassName(),
					"items:", fmt.Sprintf("%v", m.groupClassNames()))
			}
		}
		formTab.hideGuides()
		// X/Y 为设计画布坐标(设计钩子换算), 记录按下位置用于拖动计算
		m.dx, m.dy = X, Y
		br := sender.BoundsRect()
		m.dcl = br.Left
		m.dct = br.Top
		m.lastMoveX, m.lastMoveY = br.Left, br.Top
		msgContent := fmt.Sprintf("X: %v Y: %v\nW: %v H: %v", br.Left, br.Top, br.Width(), br.Height())
		message.Follow(msgContent)
		sender.DragBegin()
	}
}

// finishResize 缩放结束(画布控制点与窗体控制点面板共用):
// 只有真正改变过边界才写回属性并触发生成。
func (m *drag) finishResize() {
	if m.resizeDir < 0 {
		logs.Debug("PERF finishResize skipped: resizeDir < 0")
		return
	}
	logs.Debug("PERF finishResize dir:", m.resizeDir, "moved:", m.moved)
	if m.moved {
		br := m.relation.BoundsRect()
		m.relation.UpdateNodeDataBounds(br)
	}
	m.resizeDir = -1
	m.moved = false
	atomic.StoreInt32(&overlayCursorResize, -1)
	resetOverlayCursor()
	if m.relation.FormTab != nil {
		m.relation.FormTab.hideGuides()
		m.relation.FormTab.endSelectionDrag()
	}
	message.FollowHide()
	m.relation.DragEnd()
}

// 设计组件鼠标抬起事件
func (m *drag) OnMouseUp(sender *TDesigningComponent, button types.TMouseButton, shift types.TShiftState, X int32, Y int32) {
	if m.resizeDir >= 0 {
		m.finishResize()
		return
	}
	if m.isDown {
		perfStart := time.Now()
		defer func() {
			logs.Debug("PERF OnMouseUp cost:", time.Since(perfStart))
		}()
		m.isDown = false
		// 恢复拖动期间被摘除设计模式的次选组件
		m.restoreGroupDesignModes()
		// 只有真正发生过拖动才补写最终坐标，普通点击不触发 UI 生成。
		if m.moved {
			br := sender.BoundsRect()
			sender.UpdateNodeDataPoint(br.Left, br.Top)
			for _, item := range m.groupItems {
				actual := item.comp.BoundsRect()
				item.comp.UpdateNodeDataPoint(actual.Left, actual.Top)
			}
		}
		m.groupItems = m.groupItems[:0]
		m.moved = false
		if sender.FormTab != nil {
			// 拖动期间隐藏的次选边框在结束时统一恢复一次。
			sender.FormTab.endSelectionDrag()
			sender.FormTab.hideGuides()
		}
	}
	message.FollowHide()
	sender.DragEnd()
}
