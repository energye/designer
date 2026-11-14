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
// See the License for the specific language governing permissions and
// limitations under the License.

package designer

import (
	"sync/atomic"

	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/tool"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 设计 - 覆盖层(玻璃层, 参考 Lazarus TDesigner.DoPaintDesignerItems):
//
// 每个页签一个置顶面板覆盖整个设计画布, 绘制主选缩放控制点与次选选中边框。
// 窗口区域 = 所有锚点/边框矩形并集(LCL SetShapeWithRegion, 跨平台):
//   - 区域外: 不绘制也不拦截鼠标, 视觉完全透明, 鼠标穿透到下层组件;
//   - 区域内: 显示锚点/边框, 且始终位于所有设计组件之上(Z 序遮挡问题消失)。
//
// 鼠标交互: 区域内点击由覆盖层处理(缩放/选中), 区域外穿透到下层组件,
// 由设计钩子的模型级命中测试(dragHandleAt 优先于组件命中)兜底。
// 资源: 每页签一个窗口, 任意组件数量/任意多组锚点/任意重叠层级都不增加窗口。
//
// 光标: Screen 全局光标 + OnUserInput 兜底复位(见 setOverlayCursor 注释)。
//
// 平台: Win32(SetWindowRgn)/GTK2(gdk shape)/Qt5(setMask) 走覆盖层;
// Cocoa 无区域实现, 回退容器画布绘制(paintContainerOverlays)。

// overlayGlyph 覆盖层绘制单元
type overlayGlyph struct {
	rect     types.TRect
	isHandle bool // true: 主选缩放控制点(黑); false: 次选选中边框(设计蓝)
}

// overlayGlyphs 收集当前覆盖层绘制单元(设计画布客户区坐标):
// 主选组件缩放控制点(拖拽控制点面板) + 次选组件选中边框。
func (m *FormTab) overlayGlyphs() []overlayGlyph {
	var glyphs []overlayGlyph
	primary := m.primary
	if primary != nil && primary.drag != nil {
		off := m.overlayOffset(primary, m.FormRoot)
		for _, item := range primary.drag.handleRects() {
			glyphs = append(glyphs, overlayGlyph{
				rect:     types.Rect(item.rect.Left+off.X, item.rect.Top+off.Y, item.rect.Right+off.X, item.rect.Bottom+off.Y),
				isHandle: true,
			})
		}
	}
	for _, x := range m.overlaySecondary() {
		ctrl := x.Control()
		if ctrl == nil || !ctrl.IsValid() {
			continue
		}
		off := m.overlayOffset(x, m.FormRoot)
		br := ctrl.BoundsRect()
		l, t := br.Left+off.X, br.Top+off.Y
		w, h := br.Width(), br.Height()
		// 边框: 组件外侧 1px, 4 条矩形
		glyphs = append(glyphs,
			overlayGlyph{rect: types.Rect(l-1, t-1, l+w+1, t)},
			overlayGlyph{rect: types.Rect(l-1, t+h, l+w+1, t+h+1)},
			overlayGlyph{rect: types.Rect(l-1, t, l, t+h)},
			overlayGlyph{rect: types.Rect(l+w, t, l+w+1, t+h)},
		)
	}
	return glyphs
}

// mustOverlayPane 创建页签级覆盖层窗口(启用 + AlClient + 置顶)。
// 父容器为滚动容器(而非设计画布): 覆盖层略大于设计画布,
// 窗体边缘的锚点(跨在窗体边缘)不会被画布边界裁剪。
// 窗口区域=锚点矩形并集: 区域内点击由覆盖层处理(缩放/选中),
// 区域外鼠标穿透到下层组件(选中/移动等设计交互不受影响)。
func (m *FormTab) mustOverlayPane() {
	if m.overlayPane != nil && m.overlayPane.IsValid() {
		return
	}
	if m.formDesignPage == nil || m.formDesignPage.formDesignScroll == nil || !m.formDesignPage.formDesignScroll.IsValid() {
		return
	}
	scroll := m.formDesignPage.formDesignScroll
	p := lcl.NewPanel(scroll)
	p.SetBevelOuter(types.BvNone)
	p.SetParent(scroll)
	p.SetAlign(types.AlClient)
	p.SetOnPaint(func(sender lcl.IObject) {
		m.paintOverlayPane(p)
	})
	p.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X, Y int32) {
		m.overlayPaneMouseDown(button, shift, X, Y)
	})
	p.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, X, Y int32) {
		m.overlayPaneMouseMove(shift, X, Y)
	})
	p.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, X, Y int32) {
		m.overlayPaneMouseUp(button, shift)
	})
	// 鼠标离开覆盖层: 复位缩放光标(兜底; 全局兜底见 overlayCursorRestore)
	p.SetOnMouseLeave(func(sender lcl.IObject) {
		resetOverlayCursor()
	})
	p.BringToFront()
	m.overlayPane = p
}

// panePosToCanvas 覆盖层客户区坐标 → 设计画布客户区坐标
func (m *FormTab) panePosToCanvas(X, Y int32) types.TPoint {
	off := m.canvasToPaneOffset()
	return types.TPoint{X: X - off.X, Y: Y - off.Y}
}

// overlayPaneMouseDown 覆盖层按下: 命中主选控制点→开始缩放;
// 未命中(次选边框上)→按修饰键转发组件选中。
func (m *FormTab) overlayPaneMouseDown(button types.TMouseButton, shift types.TShiftState, X, Y int32) {
	pos := m.panePosToCanvas(X, Y)
	if button != types.MbLeft {
		return
	}
	if comp, dir, ok := m.dragHandleAt(pos); ok && comp.drag != nil {
		comp.drag.beginResize(dir, pos)
		return
	}
	// 次选边框上的点击: 转发组件选中 (Ctrl 切换 / Shift 追加 / 无修饰单选)
	if comp := m.formDesigner.GetComponentAtPos(pos.X, pos.Y); comp != nil && comp.ComponentType != consts.CtForm {
		if shift.In(types.SsCtrl) {
			m.ToggleComponentSelection(comp)
		} else if shift.In(types.SsShift) {
			m.AddComponentSelection(comp)
		} else if !m.IsComponentSelected(comp) {
			m.SwitchComponentEditing(comp)
		}
	}
}

// overlayPaneMouseMove 覆盖层移动: 缩放中→按方向调整;
// 非缩放→按控制点方向设置/复位全局缩放光标(移出覆盖层区域后
// 由 overlayCursorRestore 兜底复位, 见 setOverlayCursor 注释)。
func (m *FormTab) overlayPaneMouseMove(shift types.TShiftState, X, Y int32) {
	pos := m.panePosToCanvas(X, Y)
	if primary := m.primary; primary != nil && primary.drag != nil && primary.drag.resizeDir >= 0 {
		primary.drag.applyResize(shift, pos.X, pos.Y)
		setOverlayCursor(primary.drag.resizeDir)
		return
	}
	if _, dir, ok := m.dragHandleAt(pos); ok {
		setOverlayCursor(dir)
	} else {
		resetOverlayCursor()
	}
}

// overlayPaneMouseUp 覆盖层抬起: 结束缩放并复位光标
func (m *FormTab) overlayPaneMouseUp(button types.TMouseButton, shift types.TShiftState) {
	if primary := m.primary; primary != nil && primary.drag != nil && primary.drag.resizeDir >= 0 {
		primary.drag.finishResize()
	}
	resetOverlayCursor()
}

// canvasToPaneOffset 画布客户区坐标 → 覆盖层窗口客户区坐标 的偏移
func (m *FormTab) canvasToPaneOffset() types.TPoint {
	var off types.TPoint
	if m.overlayPane == nil || !m.overlayPane.IsValid() || m.FormRoot == nil {
		return off
	}
	if root := m.FormRoot.Control(); root != nil && root.IsValid() {
		co := root.ClientOrigin()
		po := m.overlayPane.ClientOrigin()
		off = types.TPoint{X: co.X - po.X, Y: co.Y - po.Y}
	}
	return off
}

// paintOverlayPane 绘制覆盖层内容(仅窗口区域内可见)
func (m *FormTab) paintOverlayPane(p lcl.IPanel) {
	glyphs := m.overlayGlyphs()
	if len(glyphs) == 0 {
		return
	}
	canvas := p.Canvas()
	if canvas == nil || !canvas.IsValid() {
		return
	}
	paneOff := m.canvasToPaneOffset()
	brush := canvas.BrushToBrush()
	for _, g := range glyphs {
		if g.isHandle {
			brush.SetColor(colors.ClBlack)
		} else {
			brush.SetColor(designerSelectionColor)
		}
		canvas.FillRectWithRect(types.Rect(g.rect.Left+paneOff.X, g.rect.Top+paneOff.Y,
			g.rect.Right+paneOff.X, g.rect.Bottom+paneOff.Y))
	}
}

// updateOverlayPane 覆盖层跟随当前选中状态: 重建窗口区域并重绘。
// 几何变化(选中变化/移动/缩放/对齐)后调用; 拖动期间由 selectionDragActive 隐藏。
func (m *FormTab) updateOverlayPane() {
	if !overlayPaneSupported() {
		return
	}
	logs.Debug("PERF updateOverlayPane enter primary:", func() any {
		if m.primary != nil {
			return m.primary.Name()
		}
		return "nil"
	}(), "dragActive:", m.selectionDragActive)
	m.mustOverlayPane()
	if m.overlayPane == nil || !m.overlayPane.IsValid() {
		return
	}
	if m.selectionDragActive || m.primary == nil {
		if m.overlayPane.Visible() {
			m.overlayPane.SetVisible(false)
		}
		return
	}
	glyphs := m.overlayGlyphs()
	if len(glyphs) == 0 {
		if m.overlayPane.Visible() {
			m.overlayPane.SetVisible(false)
		}
		return
	}
	if !m.overlayPane.Visible() {
		m.overlayPane.SetVisible(true)
	}
	// 重新显示前复位光标: 上一次悬停/缩放可能残留方向光标,
	// 鼠标位置未变时不会触发 mousemove 设置, 先清一次避免错样。
	resetOverlayCursor()
	m.overlayPane.BringToFront()
	paneOff := m.canvasToPaneOffset()
	paneRects := make([]types.TRect, 0, len(glyphs))
	for _, g := range glyphs {
		paneRects = append(paneRects, types.Rect(g.rect.Left+paneOff.X, g.rect.Top+paneOff.Y,
			g.rect.Right+paneOff.X, g.rect.Bottom+paneOff.Y))
	}
	applyOverlayRegion(m.overlayPane, paneRects)
	m.overlayPane.Invalidate()
	logs.Debug("PERF updateOverlayPane glyphs:", len(glyphs))
}

// overlayPaneSupported 当前平台是否支持覆盖层窗口区域。
// LCL SetShapeWithRegion → LCLIntf.SetWindowRgn 按组件集实现:
// Win32(SetWindowRgn)/GTK2(gdk_window_shape_combine_region)/Qt5(QWidget::setMask);
// Cocoa 无实现, 回退容器画布绘制(paintContainerOverlays)。
func overlayPaneSupported() bool {
	return !tool.IsDarwin()
}

// applyOverlayRegion 将覆盖层窗口区域设置为绘制单元矩形并集(跨平台)。
// 每次重建新 Region 后交给 SetShapeWithRegion, 区域句柄由组件集/系统
// 接管(Win32 系统持有, GTK2 内部 DeleteObject), 不复用不手动释放。
func applyOverlayRegion(p lcl.IPanel, rects []types.TRect) {
	if p == nil || !p.IsValid() || len(rects) == 0 {
		return
	}
	rgn := lcl.NewRegion()
	for _, r := range rects {
		rgn.AddRectangle(r.Left, r.Top, r.Right, r.Bottom)
	}
	p.SetShapeWithRegion(rgn)
}

// 锚点缩放光标(Screen 全局光标 + OnUserInput/定时器 双通道兜底复位):
//
// LCL 的窗口级光标机制在这里全部失效, 必须用 Screen 全局光标:
//   - 设计模式窗口: LCL WindowProcHelper.HandleSetCursor 显式跳过
//     (csDesigning), 且 LCL 不会为 WM_SETCURSOR 调 DefWindowProc(类光标),
//     无人调用 Windows.SetCursor 时光标永远停留 —— 控件 Cursor 属性
//     只在 WM_SETCURSOR 时生效, 对设计窗口同样无效;
//   - 普通窗口: 光标为 CrDefault 时 LCL 也不调用 SetCursor, 同样停留。
//
// Screen.SetCursor 内部立即调用 Windows.SetCursor(可见即生效), 但它是
// 全局态: 残留时所有窗口都会被强制显示该光标(LCL HandleSetCursor 优先
// 采用 Screen.RealCursor)。因此复位必须有兜底(initGlobalCursorGuard
// 双通道):
//   - Application.OnUserInput: 每次用户输入, 且先于控件事件处理,
//     即时复位; 部分应用面板的鼠标移动不触发该事件;
//   - 120ms 定时器: 覆盖上述不触发 OnUserInput 的窗口路径。
// 兜底判定 overlayCursorRestore: 非拖拽态且 Screen 光标为缩放光标且
// 鼠标不在缩放控制点上 → 恢复箭头; 仍在控制点上则保持(悬停设置由
// 覆盖层/设计钩子的 mousemove 负责, 兜底只复位)。
// 缩放拖拽期间(overlayCursorResize >= 0)跳过兜底, 光标保持方向样式。

// overlayCursorResize 进行中的缩放拖拽方向(仅拖拽生命周期:
// beginResize 置位, finishResize/cancelStaleDrag 复位; -1=无拖拽)。
// 悬停(非拖拽)设置的光标不置位 —— 兜底复位判定只看该标志,
// 鼠标移出锚点后的下一次用户输入(先于控件事件)即恢复箭头,
// 而仍悬停在锚点上时覆盖层 mousemove 会随后重新设置方向光标。
var overlayCursorResize int32 = -1

// setOverlayCursor 设置锚点缩放光标(全局, 立即生效)
func setOverlayCursor(dir int) {
	if cur := resizeCursor(dir); cur != types.CrDefault {
		lcl.Screen.SetCursor(cur)
	} else {
		resetOverlayCursor()
	}
}

// resetOverlayCursor 恢复默认光标(全局, 立即生效)
func resetOverlayCursor() {
	if lcl.Screen.Cursor() != types.CrDefault {
		lcl.Screen.SetCursor(types.CrDefault)
	}
}

// overlayCursorRestore 光标兜底复位: 缩放拖拽中保持方向光标;
// 否则 Screen 光标为缩放样式且鼠标不在缩放控制点上 → 恢复默认箭头。
// 悬停在控制点上时不动作(悬停设置由覆盖层/设计钩子的 mousemove 负责)。
// 设计画布箭头(CrArrow, 见 TEngFormDesigner.setCursor)在鼠标离开
// 设计画布客户区后同样复位为默认, 避免全局箭头影响设计器之外
// 的光标(如代码编辑器的 I-beam)。
func overlayCursorRestore() {
	if atomic.LoadInt32(&overlayCursorResize) >= 0 {
		return // 缩放拖拽进行中: 光标固定方向样式
	}
	cur := lcl.Screen.Cursor()
	if cur == types.CrDefault {
		return
	}
	tab := activeFormTab()
	pos, ok := overlayGuardMousePos(tab)
	if isOverlayResizeCursor(cur) {
		if !ok {
			lcl.Screen.SetCursor(types.CrDefault)
			return
		}
		if _, _, hit := tab.dragHandleAt(pos); !hit {
			lcl.Screen.SetCursor(types.CrDefault)
		}
		return
	}
	// 设计画布箭头(CrArrow): 鼠标不在设计画布客户区内时复位为默认
	if cur == types.CrArrow {
		if !ok || !posInDesignCanvas(tab, pos) {
			lcl.Screen.SetCursor(types.CrDefault)
		}
	}
}

// posInDesignCanvas pos(设计画布客户区坐标)是否在设计画布客户区内
func posInDesignCanvas(tab *FormTab, pos types.TPoint) bool {
	if tab == nil || tab.FormRoot == nil {
		return false
	}
	root := tab.FormRoot.Control()
	if root == nil || !root.IsValid() {
		return false
	}
	return pos.X >= 0 && pos.Y >= 0 && pos.X < root.ClientWidth() && pos.Y < root.ClientHeight()
}

// isOverlayResizeCursor 是否为本设计器设置的缩放方向光标
func isOverlayResizeCursor(c types.TCursor) bool {
	switch c {
	case types.CrSizeWE, types.CrSizeNS, types.CrSizeNWSE, types.CrSizeNESW:
		return true
	}
	return false
}

// activeFormTab 当前激活(正在设计)的窗体页签
func activeFormTab() *FormTab {
	if designer == nil {
		return nil
	}
	for _, form := range designer.designerForms {
		if form != nil && form.IsDesigner {
			return form
		}
	}
	return nil
}

// overlayGuardMousePos 鼠标屏幕坐标 → 设计画布客户区坐标
func overlayGuardMousePos(tab *FormTab) (types.TPoint, bool) {
	if tab == nil || tab.FormRoot == nil {
		return types.TPoint{}, false
	}
	root := tab.FormRoot.Control()
	if root == nil || !root.IsValid() {
		return types.TPoint{}, false
	}
	mp := lcl.Mouse.CursorPos()
	co := root.ClientOrigin()
	return types.TPoint{X: mp.X - co.X, Y: mp.Y - co.Y}, true
}

// dragHandleAt 命中测试: pos 为设计画布客户区坐标,
// 命中主选组件缩放控制点时返回组件与方向(优先于组件命中)。
func (m *FormTab) dragHandleAt(pos types.TPoint) (*TDesigningComponent, int, bool) {
	primary := m.primary
	if primary == nil || primary.drag == nil || primary.drag.relation == nil {
		logs.Debug("PERF dragHandleAt skip: primary nil")
		return nil, 0, false
	}
	off := m.overlayOffset(primary, m.FormRoot)
	for _, item := range primary.drag.handleRects() {
		r := types.Rect(item.rect.Left+off.X, item.rect.Top+off.Y, item.rect.Right+off.X, item.rect.Bottom+off.Y)
		if pos.X >= r.Left && pos.X < r.Right && pos.Y >= r.Top && pos.Y < r.Bottom {
			return primary, item.dir, true
		}
	}
	return nil, 0, false
}

// overlayOffset 覆盖层坐标系偏移:
// 组件 x 的覆盖层矩形位于 x 的父容器客户区(根节点为设计画布客户区),
// 返回该坐标系相对 comp 客户区坐标系的偏移。
func (m *FormTab) overlayOffset(x, comp *TDesigningComponent) types.TPoint {
	var base, so types.TPoint
	if x == nil {
		return base
	}
	if x == m.FormRoot || x.Parent() == nil {
		if c := m.FormRoot.Control(); c != nil && c.IsValid() {
			base = c.ClientOrigin()
		}
	} else if c := x.Parent().Control(); c != nil && c.IsValid() {
		base = c.ClientOrigin()
	}
	if comp != nil && comp != m.FormRoot && comp.Control() != nil && comp.Control().IsValid() {
		so = comp.Control().ClientOrigin()
	} else if c := m.FormRoot.Control(); c != nil && c.IsValid() {
		so = c.ClientOrigin()
	}
	return types.TPoint{X: base.X - so.X, Y: base.Y - so.Y}
}

// overlaySecondary 返回次选组件列表(排除主选/窗体根节点/非可视组件)
func (m *FormTab) overlaySecondary() []*TDesigningComponent {
	var result []*TDesigningComponent
	for _, comp := range m.SelectedComponents() {
		if comp == m.primary || comp == m.FormRoot || comp.ComponentType == consts.CtNonVisual {
			continue
		}
		result = append(result, comp)
	}
	return result
}
