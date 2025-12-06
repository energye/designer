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
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"reflect"
)

// 组件设计创建管理

//// 按钮 Button
//func NewButtonDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewButton(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("Button"))
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 编辑框 Edit
//func NewEditDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewEdit(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("Edit"))
//	comp.SetText(comp.Name())
//	comp.SetReadOnly(true)
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 多选框 CheckBox
//func NewCheckBoxDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewCheckBox(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("CheckBox"))
//	comp.SetChecked(false)
//	setBaseProp(comp, x, y)
//	comp.SetOnChange(func(sender lcl.IObject) {
//		comp.SetChecked(false)
//	})
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 面板 Panel
//func NewPanelDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewPanel(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("Panel"))
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 主菜单 MainMenu
//func NewMainMenuDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newNonVisualComponent(designerForm, x, y)
//	comp := lcl.NewMainMenu(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("MainMenu"))
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 弹出菜单 PopupMenu
//func NewPopupMenuDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newNonVisualComponent(designerForm, x, y)
//	comp := lcl.NewPopupMenu(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("PopupMenu"))
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 标签 Label
//func NewLabelDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewLabel(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("Label"))
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 多行文本框 Memo
//func NewMemoDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewMemo(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("Memo"))
//	comp.SetReadOnly(true)
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 状态标记 ToggleBox
//func NewToggleBoxDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewToggleBox(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("ToggleBox"))
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}
//
//// 虚拟树 LazVirtualStringTree
//func NewLazVirtualStringTreeDesigner(designerForm *FormTab, x, y int32) *TDesigningComponent {
//	m := newVisualComponent(designerForm)
//	comp := lcl.NewLazVirtualStringTree(designerForm.FormRoot.object)
//	comp.SetName(designerForm.GetComponentCaptionName("LazVirtualStringTree"))
//	setBaseProp(comp, x, y)
//	m.drag = newDrag(designerForm.scroll, consts.DsAll)
//	m.drag.SetRelation(m)
//	m.SetObject(comp)
//	return m
//}

// GetDesignerComponent 创建一个新的设计器组件
// designerForm: 设计器表单对象，用于承载组件
// x: 组件在表单中的横坐标位置
// y: 组件在表单中的纵坐标位置
// component: 注册组件信息，包含组件的类型和其他元数据
func GetDesignerComponent(designerForm *FormTab, x, y int32, className string) *TDesigningComponent {
	component := registerComponentsTest.Get(className)
	if designerForm == nil || component == nil || x < 0 || y < 0 {
		return nil
	}
	method := reflect.ValueOf(component.Func)
	if method.Kind() != reflect.Func {
		logs.Error("获取设计组件, 创建类函数不是函数, 类名:", className)
		return nil
	}
	asMethod := reflect.ValueOf(component.AsFunc)
	if asMethod.Kind() != reflect.Func {
		logs.Error("获取设计组件, 转换对象函数不是函数, 类名:", className)
		return nil
	}
	callAsMethod := func(instance uintptr) any {
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(instance)
		resultValues := asMethod.Call(in)
		// result
		results := make([]any, len(resultValues))
		for i, value := range resultValues {
			results[i] = value.Interface()
		}
		// object
		object := results[0]
		return object
	}
	// 创建类实例, 并标记为设计模式
	newInstance := func() uintptr {
		resultValues := method.Call(nil)
		// result
		results := make([]any, len(resultValues))
		for i, value := range resultValues {
			results[i] = value.Interface()
		}
		// object
		class := results[0].(uintptr)
		instance := api.NewInstanceByComponentClass(class)
		SetComponentDesignMode(instance)
		return instance
	}
	compName := tool.RemoveT(className)
	// 创建对象, 所属为 表单
	ownerForm := designerForm.FormRoot.originObject.(*TDesignerForm)
	switch component.Type {
	case consts.CtForm:
		// TODO
		_ = ""
	case consts.CtVisual:
		// 创建常规可视对象设计组件
		m := newVisualComponent(designerForm)
		m.mod = component.Mod
		instance := newInstance()
		comp := lcl.AsWinControl(instance)
		m.SetObject(callAsMethod(instance))
		comp.SetControlStyle(comp.ControlStyle().Include(types.CsOwnedChildrenNotSelectable, types.CsNoFocus))
		// 创建对象, 所属为 表单
		api.CreateObjectByComponent(instance, ownerForm.Instance())
		comp.SetName(designerForm.GetComponentCaptionName(compName))
		setBaseProp(comp, x, y)
		m.drag = newDrag(designerForm.scroll, consts.DsAll)
		m.drag.SetRelation(m)
		return m
	case consts.CtNonVisual:
		// 创建非可视对象设计组件
		m := newNonVisualComponent(designerForm, x, y)
		m.mod = component.Mod
		instance := newInstance()
		comp := lcl.AsComponent(instance)
		m.SetObject(callAsMethod(instance))
		// 创建对象, 所属为 表单
		api.CreateObjectByComponent(instance, ownerForm.Instance())
		comp.SetName(designerForm.GetComponentCaptionName(compName))
		SetDesignMode(m.objectNonWrap.icon)
		m.drag = newDrag(designerForm.scroll, consts.DsAll)
		m.drag.SetRelation(m)
		m.objectNonWrap.SetImage()
		return m
	case consts.CtWrapVisual:
		// 创建自定义包裹非常规可视组件
		//m := newWrapVisualComponent(designerForm, x, y)
		//m.mod = component.Mod
		//instance := newInstance()
		//comp := lcl.AsWinControl(instance)
		//m.SetObject(callAsMethod(instance))
		//comp.SetControlStyle(comp.ControlStyle().Include(types.CsOwnedChildrenNotSelectable, types.CsNoFocus))
		//api.CreateObjectByComponent(instance, ownerForm.Instance())
		//comp.SetName(designerForm.GetComponentCaptionName(compName))
		//setBaseProp(comp, x, y)
		//m.drag = newDrag(designerForm.scroll, consts.DsAll)
		//m.drag.SetRelation(m)
		//return m
	}
	return nil
}
