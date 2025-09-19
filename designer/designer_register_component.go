package designer

// 组件设计注册

// 创建设计组件回调函数
type TNewComponent func(designerForm *FormTab, x, y int32) *DesigningComponent

// 注册设计组件
var registerComponents = make(map[string]TNewComponent)

func init() {
	registerComponents["TButton"] = NewButtonDesigner
	registerComponents["TEdit"] = NewEditDesigner
}

// 获取注册的设计组件
func GetRegisterComponent(name string) TNewComponent {
	if cb, ok := registerComponents[name]; ok {
		return cb
	}
	return nil
}
