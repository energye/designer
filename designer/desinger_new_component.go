package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 组件设计创建管理

// 设置组件模式为设计模式
func SetDesignMode(component lcl.IComponent) {
	lcl.DesigningComponent().SetComponentDesignMode(component, true)
	lcl.DesigningComponent().SetComponentDesignInstanceMode(component, true)
	lcl.DesigningComponent().SetComponentInlineMode(component, true)
	lcl.DesigningComponent().SetWidgetSetDesigning(component)
}

// 创建可视组件
func newVisualComponent(designerForm *FormTab) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtVisual
	m.ownerFormTab = designerForm
	return m
}

// 创建非可视组件
func newNonVisualComponent(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := new(DesigningComponent)
	m.componentType = CtNonVisual
	m.ownerFormTab = designerForm
	objectWrap := NewNonVisualComponentWrap(designerForm.designerBox.object, m)
	objectWrap.SetLeftTop(x, y)
	m.objectNonWrap = objectWrap
	return m
}

// 按钮 Button
func NewButtonDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewButton(designerForm.designerBox.object)
	comp.SetName(designerForm.GetComponentCaptionName("Button"))
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(comp.Name())
	comp.SetShowHint(true)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 编辑框 Edit
func NewEditDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewEdit(designerForm.designerBox.object)
	comp.SetName(designerForm.GetComponentCaptionName("Edit"))
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetText(comp.Name())
	comp.SetCaption(comp.Name())
	comp.SetShowHint(true)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 多选框 CheckBox
func NewCheckBoxDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewCheckBox(designerForm.designerBox.object)
	comp.SetName(designerForm.GetComponentCaptionName("CheckBox"))
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(comp.Caption())
	comp.SetChecked(false)
	comp.SetShowHint(true)
	comp.SetOnChange(func(sender lcl.IObject) {
		comp.SetChecked(false)
	})
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 面板 Panel
func NewPanelDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewPanel(designerForm.designerBox.object)
	comp.SetName(designerForm.GetComponentCaptionName("Panel"))
	comp.SetLeft(x)
	comp.SetTop(y)
	comp.SetCursor(types.CrSize)
	comp.SetCaption(comp.Caption())
	comp.SetShowHint(true)
	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 主菜单 MainMenu
func NewMainMenuDesigner(designerForm *FormTab, x, y int32) *DesigningComponent {
	m := newNonVisualComponent(designerForm, x, y)
	comp := lcl.NewMainMenu(designerForm.designerBox.object)
	comp.SetName(designerForm.GetComponentCaptionName("MainMenu"))

	m.drag = newDrag(designerForm.designerBox.object, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}
