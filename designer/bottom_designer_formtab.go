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
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/dast"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strings"
)

// 设计器面板

type FormTabState int32

const (
	FtsNone FormTabState = iota
	FtsClose
	FtsHide
	FtsShow
)

// 设计窗体的 tab
type FormTab struct {
	Id           int                  // 唯一索引, 关联 forms key: index
	IsDesigner   bool                 // 当前设计窗体 Form 是否正在设计, 当显示和隐藏时设置值
	State        FormTabState         //
	sheet        *wg.TPage            // tab sheet
	scroll       lcl.IScrollBox       // 外 滚动条
	formDesigner *TEngFormDesigner    // 设计器处理器
	FormRoot     *TDesigningComponent // 设计器, 窗体 Form, 组件树的根节点
	recover      *TRecoverForm        // 恢复模式
	recvMethods  []*dast.TFuncInfo    // 属于设计窗体的自引用方法列表, 动态更新
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
	m.sheet.Close()
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

// 隐藏之前设计的组件 drag, 对象查看器属性和事件
func (m *FormTab) HideAllDesignHelpers(target *TDesigningComponent) {
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
	m.HideAllDesignHelpers(targetComp)
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
	if m.sheet.IsEnterClose() {
		// 关闭状态不处理任何逻辑
		return
	}
	logs.Debug("Designer PageControl FormTab Hide:", m.Id, "name:", m.FormRoot.Name())
	m.IsDesigner = false

	designComp := m.FindDesignComponent(m.FormRoot)
	if designComp == nil {
		designComp = m.FormRoot
	}
	// 隐藏所有设计组件
	m.HideAllDesignHelpers(designComp)
	// 隐藏掉对象查看器 tab page, 属性列表和事件列表
	designComp.page.SetVisible(false)
	m.sheet.Button().Font().SetColor(colors.ClBlack)
}

// 当前tab显示事件
func (m *FormTab) tabSheetOnShow(sender lcl.IObject) {
	if m.sheet.IsEnterClose() {
		// 关闭状态不处理任何逻辑
		return
	}
	logs.Debug("Designer PageControl FormTab Show id:", m.Id, "name:", m.FormRoot.Name())
	m.IsDesigner = true

	designComp := m.FindDesignComponent(m.FormRoot)
	if designComp == nil {
		designComp = m.FormRoot
	}
	// 显示掉对象查看器 tab page, 属性列表和事件列表
	designComp.page.SetVisible(true)
	// 恢复模式, 恢复所有设计的子组件
	lcl.RunOnMainThreadAsync(func(id uint32) {
		m.sheet.Button().Font().SetColor(0xD47800)
		//m.Recover()
		m.RecoverComponentPropertyValue()
		// 确保节点被选中
		ProjectTreeSetSelected(designComp.node)
		// 确保组件 helper 能正确显示, 因为选中已选中的节点不会再触发选中事件
		m.SwitchComponentEditing(designComp)
		m.formDesigner.Form.SetFocus()
	})
}

func (m *FormTab) HideTabPage() {
	m.sheet.Button().Hide()
	m.sheet.SetVisible(false)
	m.sheet.SetActive(false)
	designer.tab.RecalculatePosition()
	m.State = FtsHide
}

func (m *FormTab) ShowTabPage() {
	m.sheet.Button().Show()
	m.sheet.SetVisible(true)
	m.sheet.SetActive(true)
	designer.tab.RecalculatePosition()
	m.State = FtsShow
}

func (m *FormTab) tabSheetOnClose(page *wg.TPage, canClose *bool) {
	logs.Debug("Designer PageControl FormTab Close id:", m.Id, "name:", m.FormRoot.Name())
	if m.State == FtsClose {
		*canClose = true
		m.FormRoot.Free() // 关闭从根节点释放
		m.recover = nil
		// 在设计器列表删除当前窗体
		delete(designer.designerForms, m.Id)
		if len(designer.tab.Pages()) == 0 {
			designer.tab.EnableScrollButton(false)
		}
	} else {
		*canClose = false
		var (
			activeId        = -1
			activeOtherForm *FormTab
		)
		
		for id, formTab := range designer.designerForms {
			if formTab.sheet.Active() && formTab.sheet == page {
				activeId = id
			} else if activeOtherForm == nil && formTab.sheet.Button().Visible() {
				activeOtherForm = formTab
			}
			if activeId != -1 && activeOtherForm != nil {
				break
			}
		}

		// 不是 -1 时删除的是自己
		// 此时要选择一个设计表单激活设计
		if activeId != -1 && activeOtherForm != nil {
			designer.tab.HideAllActivated()
			designer.ActiveFormTab(activeOtherForm)
		}
		// 隐藏设计器 tab page
		m.HideTabPage()
	}
}

// 删除当前表单
func (m *FormTab) Remove() {
	m.State = FtsClose
	m.sheet.Close()
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
		// 绘制网格
		m.drawGrid(control)
	})
}

// 绘制风格线
func (m *FormTab) drawGrid(control lcl.ICustomControl) {
	//logs.Debug("drawGrid")
	gridSize := 9 // 小刻度
	formRoot := control
	canvas := formRoot.Canvas()
	canvas.PenToPen().SetColor(colors.ClBlack)
	width, height := formRoot.Width(), formRoot.Height()
	for i := 1; i < int(width)/gridSize; i++ {
		x := int32(i * gridSize)
		for j := 1; j < int(height)/gridSize; j++ {
			y := int32(j * gridSize)
			canvas.SetPixels(x, y, colors.ClBlack)
		}
	}
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
