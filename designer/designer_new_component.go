package designer

import (
	"github.com/energye/lcl/lcl"
)

// 组件设计创建管理

// 按钮 Button
func NewButtonDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewButton(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("Button"))
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 编辑框 Edit
func NewEditDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewEdit(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("Edit"))
	comp.SetText(comp.Name())
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 多选框 CheckBox
func NewCheckBoxDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewCheckBox(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("CheckBox"))
	comp.SetChecked(false)
	setBaseProp(comp, x, y)
	comp.SetOnChange(func(sender lcl.IObject) {
		comp.SetChecked(false)
	})
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 面板 Panel
func NewPanelDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewPanel(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("Panel"))
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 主菜单 MainMenu
func NewMainMenuDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newNonVisualComponent(designerForm, x, y)
	comp := lcl.NewMainMenu(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("MainMenu"))
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 弹出菜单 PopupMenu
func NewPopupMenuDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newNonVisualComponent(designerForm, x, y)
	comp := lcl.NewPopupMenu(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("PopupMenu"))
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 标签 Label
func NewLabelDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewLabel(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("Label"))
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 多行文本框 Memo
func NewMemoDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewMemo(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("Memo"))
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 状态标记 ToggleBox
func NewToggleBoxDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewToggleBox(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("ToggleBox"))
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}

// 虚拟树 LazVirtualStringTree
func NewLazVirtualStringTreeDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
	m := newVisualComponent(designerForm)
	comp := lcl.NewLazVirtualStringTree(designerForm.FormRoot.object)
	comp.SetName(designerForm.GetComponentCaptionName("LazVirtualStringTree"))
	setBaseProp(comp, x, y)
	m.drag = newDrag(designerForm.scroll, DsAll)
	m.drag.SetRelation(m)
	m.SetObject(comp)
	return m
}
