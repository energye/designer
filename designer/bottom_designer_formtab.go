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
	"strings"

	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

// 设计器面板

type FormTabState int32

const (
	FtsNone FormTabState = iota
	FtsClose
	FtsHide
	FtsShow
)

const (
	designerGridSize     = 8
	designerGridTileSize = 512
)

var designerGridColor = colors.RGBToColor(205, 205, 205)

// 设计窗体的 tab
type FormTab struct {
	Id         int          // 唯一索引, 关联 forms key: index
	IsDesigner bool         // 当前设计窗体 Form 是否正在设计, 当显示和隐藏时设置值
	State      FormTabState //
	mainPage   *wg.TPage    // tab sheet

	designerBox   lcl.ICustomControl // 设计画布背景，网格绘制目标
	gridVisible   bool               // 是否显示网格
	gridTile      lcl.IBitmap        // 固定大小网格瓦片，按需平铺绘制
	gridTileColor types.TColor
	//scroll         lcl.IScrollBox       // 外 滚动条
	formDesigner   *TEngFormDesigner    // 设计器处理器
	FormRoot       *TDesigningComponent // 设计器, 窗体 Form, 组件树的根节点
	recover        *TRecoverForm        // 恢复模式
	recvMethods    []*dast.TFuncInfo    // 属于设计窗体的自引用方法列表, 动态更新
	formDesignPage *TFormDesignPage     // 窗体设计页, 只针对设计窗体和对应的代码 tab page
}

// UIFile 返回UI文件名 xxx.ui
func (m *FormTab) UIFile() string {
	return strings.ToLower(m.FormRoot.Name()) + consts.UIExt
}

// GOFile 返回 Go UI 文件名 xxx.ui.go
func (m *FormTab) GOFile() string {
	return strings.ToLower(m.FormRoot.Name()) + consts.UIGoExt
}

// GOUserFile 返回 Go 用户文件名 xxx.go
func (m *FormTab) GOUserFile() string {
	return strings.ToLower(m.FormRoot.Name()) + consts.UIGoUserExt
}

func (m *FormTab) SetRecvMethods(methods []*dast.TFuncInfo) {
	m.recvMethods = methods
}

// 强制关闭当前tab
func (m *FormTab) Close() {
	m.mainPage.Close()
}

func (m *FormTab) IsDuplicateName(currComp *TDesigningComponent, name string) bool {
	if m.FormRoot != currComp && m.FormRoot.Name() == name {
		return true
	}
	var iterable func(comp *TDesigningComponent) bool
	iterable = func(comp *TDesigningComponent) bool {
		if comp != currComp && comp.Name() == name {
			return true
		}
		for _, child := range comp.Child {
			if iterable(child) {
				return true
			}
		}
		return false
	}
	return iterable(m.FormRoot)
}

// 添加设计组件到组件列表
func (m *FormTab) AddComponentToList(component *TDesigningComponent) {
	m.formDesigner.AddComponentToList(component)
}

// 返回设计组件
func (m *FormTab) GetComponentFormList(instance uintptr) *TDesigningComponent {
	return m.formDesigner.GetComponentFormList(instance)
}

// 删除一个设计组件
func (m *FormTab) RemoveComponentFormList(instance uintptr) {
	m.formDesigner.RemoveComponentFormList(instance)
}

// 隐藏其它设计的组件 drag, 对象查看器属性和事件
func (m *FormTab) HideOtherDesignHelpers(target *TDesigningComponent) {
	var iterable func(component *TDesigningComponent)
	iterable = func(component *TDesigningComponent) {
		if component.IsDesign && target != component {
			component.IsDesign = false // 标记为非设计状态
			component.HideDesignHelpers()
		}
		for _, child := range component.Child {
			iterable(child)
		}
	}
	iterable(m.FormRoot)
}

// 隐藏所有下在设计的组件 drag, 对象查看器属性和事件
func (m *FormTab) HideAllDesignHelpers() {
	var iterable func(component *TDesigningComponent)
	iterable = func(component *TDesigningComponent) {
		if component.IsDesign {
			component.IsDesign = false // 标记为非设计状态
			component.HideDesignHelpers()
		}
		for _, child := range component.Child {
			iterable(child)
		}
	}
	iterable(m.FormRoot)
}

// 查找当前设计窗体正在设计的组件
func (m *FormTab) FindDesignComponent(component *TDesigningComponent) *TDesigningComponent {
	if component == nil {
		return nil
	}
	if component.IsDesign {
		return component
	}
	for _, child := range component.Child {
		designComp := m.FindDesignComponent(child)
		if designComp != nil {
			return designComp
		}
	}
	return nil
}

// 切换组件到设计状态
func (m *FormTab) SwitchComponentEditing(targetComp *TDesigningComponent) {
	targetComp.mustComponentPropertyPage()
	targetComp.drag.mustDS()
	m.HideOtherDesignHelpers(targetComp)
	targetComp.IsDesign = true
	if m.State == FtsHide {
		designer.tab.HideAllActivated()
		m.ShowTabPage()
	}
	if !m.IsDesigner {
		// 切换到设计窗体
		designer.tab.HideAllActivated()
		designer.ActiveFormTab(m)
	}
	// helper 显示位置
	targetComp.ShowDesignHelpers()
	// 加载属性到属性列表
	targetComp.LoadPropertyToInspector()
	m.formDesigner.Form.SetFocus()
}

// 放置设计组件到设计面板或父组件容器
func (m *FormTab) placeComponent(owner *TDesigningComponent, x, y int32) bool {
	// 放置设计组件
	isAcceptsControl := false
	if owner.object != nil {
		isAcceptsControl = owner.object.ControlStyle().In(types.CsAcceptsControls)
	}
	selectComponent := SelectedComponent()
	if selectComponent != nil && isAcceptsControl {
		SetStatusCenterText("-")
		logs.Debug("选中设计组件:", selectComponent.name)
		m.SwitchComponentEditing(m.FormRoot)

		if newDesComp := GetDesignerComponent(m, x, y, selectComponent.name); newDesComp != nil {
			// 创建设计组件
			newDesComp.SetParent(owner)
			newDesComp.FormTab.SwitchComponentEditing(newDesComp)
			newDesComp.DragEnd()
			// 1. 加载属性到设计器
			// 此步骤会初始化并填充设计组件实例
			newDesComp.LoadPropertyToInspector()
			// 2. 添加到组件树
			go lcl.RunOnMainThreadAsync(func(id uint32) {
				owner.AddChild(newDesComp)
				newDesComp.node.SetSelected(true) // 选中
			})
			// 放置对象创建全量UI
			triggerUIGeneration(newDesComp, nil, event.CodeGenUI)
		} else {
			logs.Warn("选中设计组件", selectComponent.name, "未实现或未注册")
		}
		// 重置工具栏选项卡上的组件工具按钮按下
		//MainWindow.toolLayout.ResetTabComponentDown()
		ResetSelectedComponent()
		return true
	}
	return false
}

// 窗体设计界面 鼠标按下, 放置设计组件, 加载组件属性
func (m *FormTab) designerOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	// 创建组件
	logs.Debug("鼠标点击设计器")
	if !m.placeComponent(m.FormRoot, x, y) {
		m.SwitchComponentEditing(m.FormRoot)
		logs.Debug("加载窗体属性")
		// 设置选中状态到设计器组件树
		m.FormRoot.SetSelected()
		//lcl.Mouse.SetCapture(m.FormRoot.object.Handle())
	}
}

// 窗体设计界面 鼠标抬起
func (m *FormTab) designerOnMouseUp(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	//lcl.Mouse.SetCapture(0)
}

// 当前tab隐藏事件
func (m *FormTab) tabSheetOnHide(sender lcl.IObject) {
	if m.mainPage.IsEnterClose() {
		// 关闭状态不处理任何逻辑
		return
	}
	logs.Debug("Designer PageControl FormTab Hide:", m.Id, "name:", m.FormRoot.Name())
	if m.formDesignPage != nil {
		m.formDesignPage.stopRulerSync()
	}
	m.IsDesigner = false
	designComp := m.FindDesignComponent(m.FormRoot)
	if designComp == nil {
		designComp = m.FormRoot
	}
	// 隐藏所有设计组件
	m.HideOtherDesignHelpers(designComp)
	// 隐藏掉对象查看器 tab page, 属性列表和事件列表
	designComp.page.SetVisible(false)
	m.mainPage.Button().Font().SetColor(colors.ClBlack)
}

// 当前tab显示事件
func (m *FormTab) tabSheetOnShow(sender lcl.IObject) {
	if m.mainPage.IsEnterClose() {
		// 关闭状态不处理任何逻辑
		return
	}
	logs.Debug("Designer PageControl FormTab Show id:", m.Id, "name:", m.FormRoot.Name())

	if m.formDesignPage != nil {
		m.formDesignPage.startRulerSync()
	}
	m.IsDesigner = true
	designComp := m.FindDesignComponent(m.FormRoot)
	if designComp == nil {
		designComp = m.FormRoot
	}
	// 显示掉对象查看器 tab page, 属性列表和事件列表
	designComp.page.SetVisible(true)
	// 切换设计窗体代码编辑器
	m.formDesignPage.ActiveCodeEditorTab()

	// 恢复模式, 恢复所有设计的子组件
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.mainPage.Button().Font().SetColor(0xD47800)
		m.RecoverComponentPropertyValue()
		// 在当前选中的是组件节点时才切换选中, 保持源码文件节点选中状态
		selectedNode := ProjectTreeSelectNode()
		if selectedNode == nil || !selectedNode.IsValid() {
			selectedNode = designComp.node
		}
		isComponentNode := false
		if selectedNode != nil && selectedNode.IsValid() {
			data := selectedNode.Data()
			comp := TreeNodeDataToDesigningComponent(data)
			if comp != nil {
				isComponentNode = true
			}
		}
		if isComponentNode {
			// 确保节点被选中
			ProjectTreeSetSelected(designComp.node)
			// 确保组件 helper 能正确显示, 因为选中已选中的节点不会再触发选中事件
			m.SwitchComponentEditing(designComp)
			// 聚焦设计窗体
			m.formDesigner.Form.SetFocus()
		} else {
			// 源码文件节点选中时, 只显示组件 helper 不切换树节点选中
			m.HideOtherDesignHelpers(designComp)
			designComp.IsDesign = true
			designComp.ShowDesignHelpers()
		}
	})
}

func (m *FormTab) HideTabPage() {
	m.mainPage.Button().Hide()
	m.mainPage.SetVisible(false)
	m.mainPage.SetActive(false)
	designer.tab.RecalculatePosition()
	m.State = FtsHide
}

func (m *FormTab) ShowTabPage() {
	// 先置状态再激活：SetActive 会同步触发 OnShow，回调里据此避免重复恢复。
	m.State = FtsShow
	m.mainPage.Button().Show()
	m.mainPage.SetVisible(true)
	m.mainPage.SetActive(true)
	refreshTabScrollButtons()
	designer.tab.RecalculatePosition()
}

// formTabNeedsRestore 隐藏状态的窗体页必须走 ShowTabPage 才能恢复顶部页签按钮；
// 只调 SetActive(true) 会出现「内容可见但没有 tab 标题」的半隐藏状态。
func formTabNeedsRestore(state FormTabState) bool {
	return state == FtsHide
}

func (m *FormTab) tabSheetOnClose(page *wg.TPage, canClose *bool) {
	logs.Debug("Designer PageControl FormTab Close id:", m.Id, "name:", m.FormRoot.Name())
	if m.State == FtsClose {
		*canClose = true
		m.FormRoot.Free() // 关闭从根节点释放
		m.recover = nil
		// 在设计器列表删除当前窗体
		delete(designer.designerForms, m.Id)
		// 当前页在 onClose 之后才真正移除，延后一拍再刷新滚动按钮。
		lcl.RunOnMainThreadAsync(func(id uint32) {
			refreshTabScrollButtons()
		})
	} else {
		*canClose = false
		var (
			activeId        = -1
			activeOtherForm *FormTab
		)
		// 查找当前激活的窗体, 和要激活其它的窗体
		for id, formTab := range designer.designerForms {
			if formTab.mainPage.Active() && formTab.mainPage == page {
				activeId = id
			} else if activeOtherForm == nil && formTab.mainPage.Button().Visible() {
				activeOtherForm = formTab
			}
			if activeId != -1 && activeOtherForm != nil {
				break
			}
		}
		// 隐藏设计器 tab page
		m.HideAllDesignHelpers()
		m.HideTabPage()

		// 此时要选择一个设计窗体激活设计
		if activeId != -1 && activeOtherForm != nil {
			designer.tab.HideAllActivated()
			designer.ActiveFormTab(activeOtherForm)
		}
		// 代码页或其它窗体页仍可能可见，按所有可见页签刷新滚动按钮。
		refreshTabScrollButtons()
	}
}

// 删除当前表单
func (m *FormTab) Remove() {
	m.freeGridTile()
	m.designerBox = nil
	if m.formDesignPage != nil {
		m.formDesignPage.closeRuler()
	}
	m.State = FtsClose
	m.mainPage.Close()
}

// 获取组件名 Caption
func (m *FormTab) GetComponentCaptionName(component string) string {
	var matchIterateComponent func(component *TDesigningComponent, newName string) bool
	matchIterateComponent = func(component *TDesigningComponent, newName string) bool {
		if component.Name() == newName {
			return true
		}
		for _, child := range component.Child {
			ok := matchIterateComponent(child, newName)
			if ok {
				return true
			}
		}
		return false
	}

	nextId := 1
	for {
		newTmpName := fmt.Sprintf("%v%v", component, nextId)
		if !matchIterateComponent(m.FormRoot, newTmpName) {
			return newTmpName
		}
		nextId++
	}
}

func (m *FormTab) designerOnPaint(control lcl.ICustomControl) {
	control.SetOnPaint(func(sender lcl.IObject) {
		m.drawGrid(control)
	})
}

// drawGrid 每帧把固定大小的网格瓦片平铺到画布。
// 瓦片与画布尺寸无关，拖动/缩放控件时无需重建，三平台都走 LCL Canvas。
func (m *FormTab) drawGrid(control lcl.ICustomControl) {
	if !m.gridVisible || control == nil {
		return
	}
	canvas := control.Canvas()
	if canvas == nil || !canvas.IsValid() {
		return
	}
	m.ensureGridTile(control.Color())
	if m.gridTile == nil {
		return
	}
	w, h := control.ClientWidth(), control.ClientHeight()
	tileSize := int32(designerGridTileSize)

	// 只平铺当前重绘区域；ClipRect 无效时回退整张画布。
	left, top, right, bottom := int32(0), int32(0), w, h
	clip := canvas.ClipRect()
	if clip.Right > clip.Left && clip.Bottom > clip.Top {
		if clip.Left > left {
			left = clip.Left
		}
		if clip.Top > top {
			top = clip.Top
		}
		if clip.Right < right {
			right = clip.Right
		}
		if clip.Bottom < bottom {
			bottom = clip.Bottom
		}
	}
	if right <= left || bottom <= top {
		return
	}

	// 起点按瓦片对齐，避免局部重绘时网格错位。
	startX := left / tileSize * tileSize
	startY := top / tileSize * tileSize
	for y := startY; y < bottom; y += tileSize {
		for x := startX; x < right; x += tileSize {
			canvas.DrawWithIntX2Graphic(x, y, m.gridTile)
		}
	}
}

func (m *FormTab) ensureGridTile(color types.TColor) {
	if m.gridTile != nil && m.gridTileColor == color {
		return
	}
	m.freeGridTile()
	bmp := lcl.NewBitmap()
	bmp.SetSize(designerGridTileSize, designerGridTileSize)
	canvas := bmp.Canvas()
	canvas.BrushToBrush().SetColor(color)
	canvas.FillRectWithRect(types.Rect(0, 0, designerGridTileSize, designerGridTileSize))
	pen := canvas.PenToPen()
	pen.SetColor(designerGridColor)
	pen.SetWidth(1)
	pen.SetStyle(types.PsSolid)
	// 不画右/下边界线，相邻瓦片从 0 开始绘制，保证接缝只有一条线。
	for x := int32(0); x < designerGridTileSize; x += designerGridSize {
		canvas.MoveToWithIntX2(x, 0)
		canvas.LineToWithIntX2(x, designerGridTileSize-1)
	}
	for y := int32(0); y < designerGridTileSize; y += designerGridSize {
		canvas.MoveToWithIntX2(0, y)
		canvas.LineToWithIntX2(designerGridTileSize-1, y)
	}
	m.gridTile = bmp
	m.gridTileColor = color
}

func (m *FormTab) freeGridTile() {
	if m.gridTile != nil {
		m.gridTile.Free()
		m.gridTile = nil
	}
	m.gridTileColor = 0
}

// 添加窗体表单根节点
func (m *FormTab) AddFormNode() lcl.ITreeNode {
	// 窗体 根节点
	newNode := ProjectTreeAddComponentNode(nil, m.FormRoot)
	m.FormRoot.node = newNode
	// 添加到设计组件列表
	m.AddComponentToList(m.FormRoot)
	return newNode
}

// 添加组件节点
func (m *FormTab) AddComponentNode(parent, child *TDesigningComponent) {
	if parent == nil {
		logs.Error("添加组件节点失败, 父节点为空")
		return
	}
	if child == nil {
		logs.Error("添加组件节点失败, 子节点为空")
		return
	}
	if child.ComponentType == consts.CtVisual || child.ComponentType == consts.CtNonVisual {
		newNode := ProjectTreeAddComponentNode(parent, child)
		child.node = newNode
		// 添加到设计组件列表
		m.AddComponentToList(child)
	} else {
		logs.Error("添加组件节点失败, 子节点非组件节点")
	}
}
