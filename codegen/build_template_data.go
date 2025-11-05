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

package codegen

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/uigen"
	"strings"
)

// 构建模板数据

// 窗体数据
type TFormData struct {
	BaseInfo    *TBaseInfo      // 基础信息
	PackageName string          // 包名
	Imports     *tool.HashSet   // 导入包
	Form        *TComponentData // 窗体根组件节点
}

// 组件数据
type TComponentData struct {
	Name            string                 // 组件名称
	ClassName       string                 // 组件类名
	ComponentModule consts.ComponentModule // 组件所属模块
	Type            consts.ComponentType   // 组件类型
	Properties      []TPropertyData        // 组件属性
	Parent          *TComponentData        // 组件所属父类
	Children        []*TComponentData      // 子组件列表
}

// 基础信息
type TBaseInfo struct {
	DesignerVersion string // 生成工具版本
	DateTime        string // 生成时间
	UIFile          string // UI 文件
	UserFile        string // 用户文件
}

// 属性数组
type TPropertyData struct {
	Name  string                  // 属性名称
	Value any                     // 属性值
	Type  consts.PropertyDataType // 属性类型
}

// 模板调用函数 - 返回组件在Go定义的接口名
func (m *TComponentData) GoIntfName() string {
	intfName := "lcl.I" + tool.RemoveT(m.ClassName)
	return intfName
}

// 模板调用函数 - 返回组件在Go创建方法名
func (m *TComponentData) GoNewObject() string {
	// m.{{$comp.GoFieldName}} =
	newObject := tool.Buffer{}
	newObject.WriteString("m.", m.GoFieldName(), " = ")
	newObject.WriteString("lcl.New", tool.RemoveT(m.ClassName), "(m)", "\n")
	return newObject.String()
}

// 模板调用函数 - 返回组件在Go创建方法名
func (m *TComponentData) GoSetObjectParent() string {
	newObject := tool.Buffer{}
	switch m.Type {
	case consts.CtVisual:
		// 组件所属父类
		parentName := ""
		if m.Parent != nil {
			switch m.Parent.Type {
			case consts.CtForm:
				parentName = "m"
			default:
				parentName = "m." + m.Parent.GoFieldName()
			}
		} else {
			logs.Warn("设置对象父类时父类对象为 nil, 父类为空, 字段名:", m.GoFieldName(), "类名:", m.ClassName)
		}
		newObject.WriteString("m.", m.GoFieldName(), ".SetParent(", parentName, ")", "\n")
	}
	return newObject.String()
}

// 模板调用函数 - 返回组件字段名
func (m *TComponentData) GoFieldName() string {
	return m.Name
}

// 模板调用函数 - 返回当前使用的包
func (m *TFormData) GoGetImports() string {
	imports := strings.Join(m.Imports.Values(), "\n")
	return imports
}

// 包含的包
// 该函数在渲染之前执行, 获取当前所依赖的模块包
// 根据不同的属性和组件所属模板区分包
func (m *TFormData) IncludePackage() {
	var findPackages func(comp *TComponentData)
	findPackages = func(comp *TComponentData) {
		//comp.ComponentModule // TODO 后面要根据组件所属模块区分 types
		for _, prop := range comp.Properties {
			switch prop.Type {
			case consts.PdtCheckBoxList, consts.PdtComboBox:
				m.Imports.Add(lclTypes)
				// 其它 ...
			}
		}
		for _, childComp := range comp.Children {
			findPackages(childComp)
		}
	}
	findPackages(m.Form)
}

// 模板调用函数 - 设置对象属性
func (m *TPropertyData) GoPropertySet(comp *TComponentData, form *TFormData) string {
	prop := tool.Buffer{}
	object := ""
	switch comp.Type {
	case consts.CtForm:
		object = "m"
	default:
		object = "m." + comp.GoFieldName()
	}
	prop.WriteString(object)
	// 属性路径 Font.Style
	namePaths := tool.Split(m.Name, ".")
	for i := 0; i < len(namePaths)-1; i++ {
		prop.WriteString(".", namePaths[i], "()")
	}
	// 参数, 多种类型, 每种类型传入的方式不一样
	value := ""
	switch m.Type {
	case consts.PdtText: // string
		value = `"` + m.Value.(string) + `"`
	case consts.PdtInt, consts.PdtInt64: // int32 or int64
		value = tool.IntToString(m.Value)
	case consts.PdtUint16: // uint16
		if size := len(m.Value.(string)); size == 1 {
			// 单字符使用单引号 "'" 直接转换
			value = `'` + m.Value.(string) + `'`
		} else if size > 1 {
			// 多字符串串转换成 uint16
			if uv, err := tool.StrToUint16(m.Value.(string)); err == nil {
				value = tool.IntToString(uv)
			}
		}
	case consts.PdtFloat: // float32 or float64
		value = tool.FloatToString(m.Value)
	case consts.PdtCheckBox: // bool
		value = tool.BoolToString(m.Value)
	case consts.PdtCheckBoxList: // set: types.NewSet
		setStr := tool.SetToString(m.Value)
		items := tool.Split(setStr, ",")
		sets := tool.Buffer{}
		for i, item := range items {
			if i > 0 {
				sets.WriteString(",")
			}
			sets.WriteString("lclTypes.", item)
		}
		value = "lclTypes.NewSet(" + sets.String() + ")"
		form.Imports.Add(lclTypes)
	case consts.PdtComboBox: // package: mapper.GetLCL([name])
		value = "lclTypes." + m.Value.(string)
		form.Imports.Add(lclTypes)
	case consts.PdtColorSelect: // uint32: types.Color([value])
		value = tool.IntToString(m.Value)
	case consts.PdtClass: // Class instance
		logs.Debug("属性类对象实例未设置:", m.Value)
	}
	// 属性名
	name := namePaths[len(namePaths)-1]
	prop.WriteString(".Set", name, "(", value, ")")
	return prop.String()
}

// 构建自动代码模板数据
func buildAutoTemplateData(component *uigen.TUIComponent) *TFormData {
	formData := &TFormData{BaseInfo: &TBaseInfo{}, PackageName: packageName, Imports: tool.NewHashSet()}
	formData.Form = &TComponentData{
		Name:       component.Name,
		ClassName:  component.ClassName,
		Type:       component.Type,
		Properties: uiPropertiesToTemplateProperties(component.Properties),
		Children:   make([]*TComponentData, 0),
	}
	formData.Form.Children = formData.Form.buildComponents(component)
	return formData
}

// 构建用户代码模板数据
func buildUserTemplateData(component *uigen.TUIComponent) *TFormData {
	formData := &TFormData{BaseInfo: &TBaseInfo{}, PackageName: packageName, Imports: tool.NewHashSet()}
	formData.Form = &TComponentData{
		Name:       component.Name,
		ClassName:  component.ClassName,
		Type:       component.Type,
		Properties: uiPropertiesToTemplateProperties(component.Properties),
		Children:   make([]*TComponentData, 0),
	}
	formData.Form.Children = formData.Form.buildComponents(component)
	return formData
}

// 构建组件列表
func (m *TComponentData) buildComponents(component *uigen.TUIComponent) []*TComponentData {
	var components []*TComponentData
	// 创建根组件（窗体）的模板数据
	rootComponent := &TComponentData{
		Name:       component.Name,
		ClassName:  component.ClassName,
		Type:       component.Type,
		Properties: uiPropertiesToTemplateProperties(component.Properties),
		Children:   make([]*TComponentData, 0),
	}
	// 递归构建所有子组件
	m.buildChildComponents(component, rootComponent, &components)
	return components
}

// 递归构建子组件
func (m *TComponentData) buildChildComponents(uiParent *uigen.TUIComponent, templateParent *TComponentData, result *[]*TComponentData) {
	for _, child := range uiParent.Child {
		// 为每个子组件创建模板数据
		childTemplate := &TComponentData{
			Name:       child.Name,
			ClassName:  child.ClassName,
			Type:       child.Type,
			Properties: uiPropertiesToTemplateProperties(child.Properties),
			Parent:     templateParent,
			Children:   make([]*TComponentData, 0),
		}
		// 建立父子关系
		templateParent.Children = append(templateParent.Children, childTemplate)
		*result = append(*result, childTemplate)

		// 递归处理子组件的子组件
		if len(child.Child) > 0 {
			m.buildChildComponents(&child, childTemplate, result)
		}
	}
}

// UI 而已属性列表转模板数据列表
func uiPropertiesToTemplateProperties(uiProperties []uigen.TProperty) []TPropertyData {
	var templateProperties []TPropertyData
	for _, uiProperty := range uiProperties {
		templateProperty := uiPropertyToTemplateProperty(uiProperty)
		templateProperties = append(templateProperties, templateProperty)
	}
	return templateProperties
}

// UI 而已属性转列表模板数据列表
func uiPropertyToTemplateProperty(uiProperty uigen.TProperty) TPropertyData {
	return TPropertyData{
		Name:  uiProperty.Name,
		Value: uiProperty.Value,
		Type:  uiProperty.Type,
	}
}
