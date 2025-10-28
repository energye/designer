package designer

// 组件设计注册
// 所有要实现设计的组件都在此处注册

// 创建设计组件回调函数
type TNewComponent func(designerForm *FormTab, x, y int32) *TDesigningComponent

// 注册设计组件
// key: 组件类名, value: 组件创建函数
var registerComponents = make(map[string]TNewComponent)

func init() {
	registerComponents["TButton"] = NewButtonDesigner
	registerComponents["TEdit"] = NewEditDesigner
	registerComponents["TCheckBox"] = NewCheckBoxDesigner
	registerComponents["TPanel"] = NewPanelDesigner
	registerComponents["TMainMenu"] = NewMainMenuDesigner
	registerComponents["TPopupMenu"] = NewPopupMenuDesigner
	registerComponents["TLabel"] = NewLabelDesigner
	registerComponents["TMemo"] = NewMemoDesigner
	registerComponents["TToggleBox"] = NewToggleBoxDesigner
	registerComponents["TLazVirtualStringTree"] = NewLazVirtualStringTreeDesigner
}

// 获取注册的设计组件
func GetRegisterComponent(name string) TNewComponent {
	if cb, ok := registerComponents[name]; ok {
		return cb
	}
	return nil
}
