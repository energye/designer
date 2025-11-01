package codegen

import (
	"github.com/energye/designer/uigen"
)

// 构建模板数据

const packageName = "forms"

// 自动代码模板数据
type TTemplateData struct {
	PackageName string            // 包名
	FormName    string            // 窗体名称
	ClassName   string            // 窗体类名
	Components  []*TComponentData // 组件列表
	Properties  []uigen.TProperty // 窗体属性列表
	BaseInfo    TBaseInfo         // 基础信息
}

// 基础信息
type TBaseInfo struct {
	DesignerVersion string // 生成工具版本
	DateTime        string // 生成时间
	UIFile          string // UI 文件
	UserFile        string // 用户文件
}

// 组件数据
type TComponentData struct {
	Name       string            // 组件名称
	ClassName  string            // 组件类名
	Properties []uigen.TProperty // 组件属性
	Parent     *TComponentData   // 组件所属父类
	Children   []*TComponentData // 子组件列表
}

// 构建自动代码模板数据
func buildAutoTemplateData(component *uigen.TUIComponent) TTemplateData {
	data := TTemplateData{
		PackageName: packageName,
		FormName:    component.Name,
		ClassName:   component.ClassName,
		Properties:  component.Properties,
	}
	data.Components = data.buildComponents(component)
	return data
}

// 构建用户代码模板数据
func buildUserTemplateData(component *uigen.TUIComponent) TTemplateData {
	data := TTemplateData{
		PackageName: packageName,
		FormName:    component.Name,
	}
	data.Components = data.buildComponents(component)
	return data
}

// 构建组件列表
func (m *TTemplateData) buildComponents(component *uigen.TUIComponent) []*TComponentData {
	var components []*TComponentData
	// 创建根组件（窗体）的模板数据
	rootComponent := &TComponentData{
		Name:       component.Name,
		ClassName:  component.ClassName,
		Properties: component.Properties,
		Children:   make([]*TComponentData, 0),
	}
	// 递归构建所有子组件
	m.buildChildComponents(component, rootComponent, &components)
	return components
}

// 递归构建子组件
func (m *TTemplateData) buildChildComponents(uiParent *uigen.TUIComponent, templateParent *TComponentData, result *[]*TComponentData) {
	for _, child := range uiParent.Child {
		// 为每个子组件创建模板数据
		childTemplate := &TComponentData{
			Name:       child.Name,
			ClassName:  child.ClassName,
			Properties: child.Properties,
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

// 获取属性类型
func getPropertyType(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "float64"
	case bool:
		return "bool"
	default:
		return "any"
	}
}
