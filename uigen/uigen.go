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
	projectPath string
)

func init() {
	projectPath = "C:\\app\\workspace\\test"
}

// UIComponent 表示UI组件的结构
type UIComponent struct {
	Name       string                 `json:"name"`
	ClassName  string                 `json:"class_name"`
	Properties map[string]interface{} `json:"properties"`
	Children   []UIComponent          `json:"children,omitempty"`
}

// GenerateUIFile 生成UI文件
func GenerateUIFile(formTab *designer.FormTab, component *designer.TDesigningComponent) error {
	formId := strings.ToLower(formTab.Name)
	filePath := filepath.Join(projectPath, "ui", formId+".ui")
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
func buildUITree(component *designer.TDesigningComponent) UIComponent {
	uiComp := UIComponent{
		Name:       component.Name(),
		ClassName:  component.ClassName(),
		Properties: make(map[string]interface{}),
		Children:   make([]UIComponent, 0),
	}

	// 获取变更的属性
	if component.PropertyList != nil {
		for _, prop := range component.PropertyList {
			if tool.Equal(prop.Name(), "ChildSizing") {
				println()
			}
			// 只保存修改过的属性
			switch prop.Type() {
			case vtedit.PdtClass:
				// 获得类里被修改的节点
				for _, child := range prop.Child {
					if child.IsModify() {
						paths := child.Paths()
						if paths != nil {
							tool.StringArrayReverse(paths)
							paths = append(paths, child.Name())
							propName := strings.Join(paths, ".")
							propValue := child.EditValue()
							uiComp.Properties[propName] = propValue
						} else {
							logs.Error("错误, 生成UI布局文件, 属性是 class 获取子节点路径错误 nil. 属性名: ", prop.Name(), "子节点属性名:", child.Name())
						}
					}
				}
			case vtedit.PdtCheckBox:

			default:
				if prop.IsModify() {
					propName := prop.Name()
					propValue := prop.EditValue()
					uiComp.Properties[propName] = propValue
				}
			}
		}
	}

	// 递归处理子组件
	for _, child := range component.Child {
		uiComp.Children = append(uiComp.Children, buildUITree(child))
	}

	return uiComp
}
