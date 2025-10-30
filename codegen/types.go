package codegen

import (
	"github.com/energye/designer/uigen"
)

// AutoTemplateData 自动代码模板数据
type AutoTemplateData struct {
	PackageName string
	FormName    string
	ClassName   string
	Components  []ComponentData
	Properties  []PropertyData
}

// UserTemplateData 用户代码模板数据
type UserTemplateData struct {
	PackageName string
	FormName    string
	Components  []ComponentData
}

// ComponentData 组件数据
type ComponentData struct {
	Name      string
	ClassName string
	Parent    string
}

// PropertyData 属性数据
type PropertyData struct {
	Name  string
	Value interface{}
	Type  string
}

// buildAutoTemplateData 构建自动代码模板数据
func buildAutoTemplateData(component uigen.UIComponent) AutoTemplateData {
	data := AutoTemplateData{
		PackageName: "main",
		FormName:    component.Name,
		ClassName:   component.ClassName,
	}

	// 处理属性
	for name, value := range component.Properties {
		prop := PropertyData{
			Name:  name,
			Value: value,
			Type:  getPropertyType(value),
		}
		data.Properties = append(data.Properties, prop)
	}

	// 处理子组件
	data.Components = buildComponents(component.Child, component.Name)

	return data
}

// buildUserTemplateData 构建用户代码模板数据
func buildUserTemplateData(component uigen.UIComponent) UserTemplateData {
	return UserTemplateData{
		PackageName: "main",
		FormName:    component.Name,
		Components:  buildComponents(component.Child, component.Name),
	}
}

// buildComponents 构建组件列表
func buildComponents(children []uigen.UIComponent, parentName string) []ComponentData {
	var components []ComponentData

	for _, child := range children {
		comp := ComponentData{
			Name:      child.Name,
			ClassName: child.ClassName,
			Parent:    parentName,
		}
		components = append(components, comp)

		// 递归处理子组件
		if len(child.Child) > 0 {
			subComponents := buildComponents(child.Child, child.Name)
			components = append(components, subComponents...)
		}
	}

	return components
}

// getPropertyType 获取属性类型
func getPropertyType(value interface{}) string {
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
		return "interface{}"
	}
}
