package docs

import (
	_ "embed"
	"encoding/json"
	"strings"
)

// auto gen documentation.json file, use: gen/gdocs/gdocs.go
//
//go:embed documentation.json
var documentation []byte

type TDocumentation struct {
	Name    string
	Content string
}

type propertyInfo struct {
	Desc string
	Ref  string
}

func (pi *propertyInfo) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		pi.Desc = s
		return nil
	}

	type Alias propertyInfo
	var obj Alias
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	pi.Desc = obj.Desc
	pi.Ref = obj.Ref
	return nil
}

type classDoc struct {
	Desc       string                   `json:"desc"`
	Properties map[string]*propertyInfo `json:"properties"`
}

var docMap map[string]*classDoc

func init() {
	go func() {
		if err := json.Unmarshal(documentation, &docMap); err != nil {
			println("failed to parse documentation.json: " + err.Error())
		}
	}()
}

// GetClassDesc 获取类的描述信息
func GetClassDesc(className string) string {
	if docMap == nil {
		return ""
	}
	if classDoc, ok := docMap[className]; ok {
		return classDoc.Desc
	}
	return ""
}

// GetPropertyDesc 获取类中属性的描述信息
func GetPropertyDesc(className, propertyName string) string {
	classDoc, ok := docMap[className]
	if !ok {
		return ""
	}
	propInfo, ok := classDoc.Properties[propertyName]
	if !ok {
		return ""
	}
	if propInfo.Desc != "" {
		return propInfo.Desc
	}
	if propInfo.Ref != "" {
		parts := strings.SplitN(propInfo.Ref, ".", 2)
		if len(parts) == 2 {
			desc := GetPropertyDesc(parts[0], parts[1])
			return desc
		}
	}
	return ""
}

// HasClass 检查类是否存在于文档中
func HasClass(className string) bool {
	if docMap == nil {
		return false
	}
	_, ok := docMap[className]
	return ok
}

// HasProperty 检查类中是否存在指定的属性
func HasProperty(className, propertyName string) bool {
	if docMap == nil {
		return false
	}
	classDoc, ok := docMap[className]
	if !ok {
		return false
	}
	_, ok = classDoc.Properties[propertyName]
	return ok
}

// GetClassProperties 获取类的所有属性名列表
func GetClassProperties(className string) []string {
	if docMap == nil {
		return nil
	}
	classDoc, ok := docMap[className]
	if !ok {
		return nil
	}

	properties := make([]string, 0, len(classDoc.Properties))
	for propName := range classDoc.Properties {
		properties = append(properties, propName)
	}
	return properties
}
