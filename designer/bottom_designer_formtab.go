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
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"strings"
	"widget/wg"
)

// 设计器面板

// 设计表单的 tab
type FormTab struct {
	Id         int    // 唯一索引, 关联 forms key: index
	name       string // 窗体名称, 实际是一个临时名称
	isDesigner bool   // 当前窗体Form是否正在设计
	//sheet         lcl.ITabSheet        // tab sheet
	sheet         *wg.TPage            // tab sheet
	scroll        lcl.IScrollBox       // 外 滚动条
	tree          lcl.ITreeView        // 组件树
	componentName map[string]int       // 组件分类名, 同类组件ID序号
	formDesigner  *TEngFormDesigner    // 设计器处理器
	FormRoot      *TDesigningComponent // 设计器, 窗体 Form, 组件树的根节点
	componentMenu *TComponentMenu      // 组件菜单
	recover       *TRecoverForm        // 恢复模式
}

func (m *FormTab) UIFile() string {
	return strings.ToLower(m.FormRoot.Name()) + consts.UIExt
}

func (m *FormTab) GOFile() string {
	return strings.ToLower(m.FormRoot.Name()) + consts.UIGoExt
}

func (m *FormTab) GOUserFile() string {
	return strings.ToLower(m.FormRoot.Name()) + consts.UIGoUserExt
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
		for _, comp := range comp.Child {
			if iterable(comp) {
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

// 切换组件设计
func (m *FormTab) switchComponentEditing(targetComp *TDesigningComponent) {
	// 隐藏之前设计的组件
	// 隐藏之前设计组件的属性和事件列表
	var iterable func(comp *TDesigningComponent)
	iterable = func(comp *TDesigningComponent) {
		if comp.isDesigner {
			comp.drag.Hide()
			comp.page.Hide()
		}
		for _, comp := range comp.Child {
			iterable(comp)
		}
	}
	iterable(m.FormRoot)

	// 显示当前设计组件 drag
	targetComp.drag.Show()
	// 显示当前设计组件的属性和事件列表
	targetComp.page.Show()
	// 加载属性到属性列表
	targetComp.LoadPropertyToInspector()
}

// 放置设计组件到设计面板或父组件容器
func (m *FormTab) placeComponent(owner *TDesigningComponent, x, y int32) bool {
	// 放置设计组件
	isAcceptsControl := false
	if owner.object != nil {
		isAcceptsControl = owner.object.ControlStyle().In(types.CsAcceptsControls)
	}
	if toolbar.selectComponent != nil && isAcceptsControl {
		logs.Debug("选中设计组件:", toolbar.selectComponent.index, toolbar.selectComponent.name)
		m.switchComponentEditing(m.FormRoot)
		// 获取注册的组件创建函数
		if create := GetRegisterComponent(toolbar.selectComponent.name); create != nil {
			// 创建设计组件
			newComp := create(m, x, y)
			newComp.SetParent(owner)
			newComp.formTab.switchComponentEditing(newComp)
			newComp.DragEnd()
			// 1. 加载属性到设计器
			// 此步骤会初始化并填充设计组件实例
			newComp.LoadPropertyToInspector()
			// 2. 添加到组件树
			go lcl.RunOnMainThreadAsync(func(id uint32) {
				owner.AddChild(newComp)
			})
			// 放置对象
			triggerUIGeneration(newComp)
		} else {
			logs.Warn("选中设计组件", toolbar.selectComponent.name, "未实现或未注册")
		}
		// 重置工具栏选项卡上的组件工具按钮按下
		toolbar.ResetTabComponentDown()
		return true
	}
	return false
}

// 窗体设计界面 鼠标移动
func (m *FormTab) designerOnMouseMove(sender lcl.IObject, shift types.TShiftState, x, y int32) {
	br := m.FormRoot.BoundsRect()
	hint := fmt.Sprintf(`%v: TForm
	Left: %v Top: %v
	Width: %v Height: %v`, m.FormRoot.Name(), br.Left, br.Top, br.Width(), br.Height())
	m.FormRoot.SetHint(hint)
}

// 窗体设计界面 鼠标按下, 放置设计组件, 加载组件属性
func (m *FormTab) designerOnMouseDown(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
	// 创建组件
	logs.Debug("鼠标点击设计器")
	if !m.placeComponent(m.FormRoot, x, y) {
		m.switchComponentEditing(m.FormRoot)
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

func (m *FormTab) tabSheetOnHide(sender lcl.IObject) {
	logs.Debug("Designer PageControl FormTab Hide")
	// 设计状态 关闭
	m.isDesigner = false
	m.tree.SetVisible(false)
	// 隐藏属性列表 page
	var (
		iterable    func(comp *TDesigningComponent)
		defaultComp = m.FormRoot
	)
	iterable = func(comp *TDesigningComponent) {
		if comp == nil {
			return
		}
		if comp.isDesigner {
			defaultComp = comp
		}
		for _, comp := range comp.Child {
			iterable(comp)
		}
	}
	iterable(m.FormRoot)
	defaultComp.page.SetVisible(false)
}

// 当前tab显示事件
func (m *FormTab) tabSheetOnShow(sender lcl.IObject) {
	logs.Debug("Designer PageControl FormTab Show")
	triggerUIGeneration(m.FormRoot)
	// 设计状态 开启
	m.isDesigner = true
	// 显示组件树
	m.tree.SetVisible(true)

	var (
		iterable    func(comp *TDesigningComponent) // 遍历当前正在设计的组件
		defaultComp = m.FormRoot                    // 默认选中的组件
	)
	iterable = func(comp *TDesigningComponent) {
		if comp == nil {
			return
		}
		if comp.isDesigner {
			defaultComp = comp
		}
		for _, comp := range comp.Child {
			iterable(comp)
		}
	}
	iterable(m.FormRoot)
	// 显示选中设计的组件属性列表
	defaultComp.page.SetVisible(true)
	if m.recover != nil {
		// 恢复模式, 恢复所有设计的子组件
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.Recover()
		})
	}
}

func (m *FormTab) tabSheetOnClose(sender lcl.IObject) {
	delete(designer.designerForms, m.Id)
	if len(designer.tab.Pages()) == 0 {
		designer.tab.EnableScrollButton(false)
	}

}

// 获取组件名 Caption
func (m *FormTab) GetComponentCaptionName(component string) string {
	if c, ok := m.componentName[component]; ok {
		m.componentName[component] = c + 1
	} else {
		m.componentName[component] = 1
	}
	component = fmt.Sprintf("%v%d", component, m.componentName[component])
	return component
}

func (m *FormTab) designerOnPaint(control lcl.ICustomControl) {
	//control.SetOnPaint(func(sender lcl.IObject) {
	//	// 绘制网格
	//	m.drawGrid(control)
	//})
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
func (m *FormTab) AddFormNode() {
	// 窗体 根节点
	m.tree.BeginUpdate()
	defer m.tree.EndUpdate()
	items := m.tree.Items()
	m.FormRoot.id = nextTreeDataId()
	newNode := items.AddChild(nil, m.FormRoot.TreeName())
	newNode.SetImageIndex(m.FormRoot.IconIndex())    // 显示图标索引
	newNode.SetSelectedIndex(m.FormRoot.IconIndex()) // 选中图标索引
	newNode.SetSelected(true)
	newNode.SetData(m.FormRoot.instance())
	m.FormRoot.node = newNode
	// 添加到设计组件列表
	m.AddComponentToList(m.FormRoot)
}

// 添加组件节点
func (m *FormTab) AddComponentNode(parent, child *TDesigningComponent) {
	if parent == nil {
		logs.Error("添加组件节点失败, 父节点为空")
		return
	} else if child == nil {
		logs.Error("添加组件节点失败, 子节点为空")
		return
	}
	if child.ComponentType == consts.CtVisual || child.ComponentType == consts.CtNonVisual {
		m.tree.BeginUpdate()
		defer m.tree.EndUpdate()
		items := m.tree.Items()
		// 组件 子节点
		child.id = nextTreeDataId()
		node := items.AddChild(parent.node, child.TreeName())
		child.node = node
		node.SetImageIndex(child.IconIndex())    // 显示图标索引
		node.SetSelectedIndex(child.IconIndex()) // 选中图标索引
		node.SetSelected(true)                   // 选中
		node.SetData(child.instance())           // 设置数据为当前实例
		// 添加到设计组件列表
		m.AddComponentToList(child)
	} else {
		logs.Error("添加组件节点失败, 子节点非组件节点")
	}
}
