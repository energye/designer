package codegen

import (
	"github.com/energye/designer/uigen"
)

const packageName = "forms"

// TAutoTemplateData 自动代码模板数据
type TAutoTemplateData struct {
	PackageName string
	FormName    string
	ClassName   string
	Components  []TComponentData
	Properties  []TPropertyData
}

// TUserTemplateData 用户代码模板数据
type TUserTemplateData struct {
	PackageName string
	FormName    string
	Components  []TComponentData
}

// TComponentData 组件数据
type TComponentData struct {
	Name      string
	ClassName string
	Parent    string
}

// TPropertyData 属性数据
type TPropertyData struct {
	Name  string
	Value any
	Type  string
}

// buildAutoTemplateData 构建自动代码模板数据
func buildAutoTemplateData(component uigen.UIComponent) TAutoTemplateData {
	data := TAutoTemplateData{
		PackageName: packageName,
		FormName:    component.Name,
		ClassName:   component.ClassName,
	}

	// 处理属性
	for name, value := range component.Properties {
		prop := TPropertyData{
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
func buildUserTemplateData(component uigen.UIComponent) TUserTemplateData {
	return TUserTemplateData{
		PackageName: packageName,
		FormName:    component.Name,
		Components:  buildComponents(component.Child, component.Name),
	}
}

// buildComponents 构建组件列表
func buildComponents(children []uigen.UIComponent, parentName string) []TComponentData {
	var components []TComponentData
	for _, child := range children {
		comp := TComponentData{
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
