package codegen

import (
	"github.com/energye/designer/uigen"
)

// 构建模板数据

const packageName = "forms"

// 自动代码模板数据
type TAutoTemplateData struct {
	PackageName string
	FormName    string
	ClassName   string
	Components  []TComponentData
	Properties  []uigen.TProperty
	BaseInfo    TBaseInfo
}

type TBaseInfo struct {
	DesignerVersion string
	DateTime        string
	UIFile          string
	UserFile        string
}

// 用户代码模板数据
type TUserTemplateData struct {
	PackageName string
	FormName    string
	Components  []TComponentData
}

// 组件数据
type TComponentData struct {
	Name       string            // 组件名称
	ClassName  string            // 组件类名
	Parent     string            // 组件所属父类
	Properties []uigen.TProperty // 组件属性
}

// 构建自动代码模板数据
func buildAutoTemplateData(component uigen.TUIComponent) TAutoTemplateData {
	data := TAutoTemplateData{
		PackageName: packageName,
		FormName:    component.Name,
		ClassName:   component.ClassName,
		Properties:  component.Properties,
	}

	// 处理子组件
	data.Components = buildComponents(component.Child, component.Name)

	return data
}

// 构建用户代码模板数据
func buildUserTemplateData(component uigen.TUIComponent) TUserTemplateData {
	return TUserTemplateData{
		PackageName: packageName,
		FormName:    component.Name,
		Components:  buildComponents(component.Child, component.Name),
	}
}

// 构建组件列表
func buildComponents(children []uigen.TUIComponent, parentName string) []TComponentData {
	var components []TComponentData
	for _, child := range children {
		comp := TComponentData{
			Name:       child.Name,
			ClassName:  child.ClassName,
			Parent:     parentName,
			Properties: child.Properties,
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
