package uigen

import (
	"encoding/json"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/vtedit"
	"os"
	"path/filepath"
	"strings"
)

// UI 布局文件生成, JSON 格式, 自动时时生成
// 依赖 designer 生成 JSON UI文件(form[n].ui)
// 只存放被修改的组件属性
// xxx.ui 文件内容是 tree JSON 结构, 数据格式为组件[变更的属性列表]
// 生成触发条件: 即时触发 防抖

var (
	// TODO 这是一个测试变量 --test
	projectPath string
)

func init() {
	projectPath = "C:\\app\\workspace\\test"
}

// 表示UI组件的结构
type TUIComponent struct {
	Name       string                 `json:"name"`            // 组件名称
	ClassName  string                 `json:"class_name"`      // 组件类名
	Type       designer.ComponentType `json:"type"`            // 组件类型
	Properties []TProperty            `json:"properties"`      // 组件属性列表
	Child      []TUIComponent         `json:"child,omitempty"` // 组件子组件列表
}

// 组件属性
type TProperty struct {
	Name  string                  // 属性名称
	Value any                     // 属性值
	Type  vtedit.PropertyDataType // 属性类型
	//ClassInstance uintptr                 // 属性值 class 实例, 当属性类型是 PdtClass 时有效
}

// 生成UI文件
func GenerateUIFile(formTab *designer.FormTab, component *designer.TDesigningComponent, filePath string) error {
	// 构建UI树结构
	uiTree := buildUITree(formTab.FormRoot)

	// 序列化为JSON
	data, err := json.MarshalIndent(uiTree, "", "  ")
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, data, 0644)
}

// buildUITree 构建UI树结构
func buildUITree(component *designer.TDesigningComponent) TUIComponent {
	uiComp := TUIComponent{
		Name:       component.Name(),
		ClassName:  component.ClassName(),
		Properties: make([]TProperty, 0),
		Child:      make([]TUIComponent, 0),
		Type:       component.ComponentType,
	}

	// 获取变更的属性
	if component.PropertyList != nil {
		for _, prop := range component.PropertyList {
			// 默认生成的属性 Left Top Width Height
			if tool.Equal(prop.Name(), "Left", "Top", "Width", "Height", "Caption") {
				propName := prop.Name()
				propValue := prop.EditValue()
				uiComp.Properties = append(uiComp.Properties, TProperty{Name: propName, Value: propValue,
					Type: prop.Type()})
			} else {
				// 只保存修改过的属性
				switch prop.Type() {
				case vtedit.PdtClass:
					var iterator func(node *vtedit.TEditNodeData)
					iterator = func(node *vtedit.TEditNodeData) {
						for _, data := range node.Child {
							if data.IsModify() {
								if data.Type() == vtedit.PdtClass {
									iterator(data)
								} else {
									paths := data.Paths()
									if paths != nil {
										tool.StringArrayReverse(paths)
										paths = append(paths, data.Name())
										propName := strings.Join(paths, ".")
										propValue := data.EditValue()
										uiComp.Properties = append(uiComp.Properties, TProperty{Name: propName, Value: propValue,
											Type: data.Type()})
									} else {
										logs.Error("错误, 生成UI布局文件, 属性是 class 获取子节点路径错误 nil. 属性名: ", prop.Name(), "子节点属性名:", data.Name())
									}
								}
							}
						}
					}
					iterator(prop)
				default:
					if prop.IsModify() {
						propName := prop.Name()
						propValue := prop.EditValue()
						uiComp.Properties = append(uiComp.Properties, TProperty{Name: propName, Value: propValue,
							Type: prop.Type()})
					}
				}
			}
		}
	}

	// 递归处理子组件
	for _, child := range component.Child {
		uiComp.Child = append(uiComp.Child, buildUITree(child))
	}

	return uiComp
}
