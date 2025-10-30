package designer

// 生成回调事件

type UIGenerationCallback func(formTab *FormTab, component *TDesigningComponent)

var (
	// UI布局回调
	onUIGeneration UIGenerationCallback
)

// 设置UI布局文件生成回调
func SetUIGenerationCallback(callback UIGenerationCallback) {
	onUIGeneration = callback
}

// 触发UI布局生成事件
func triggerUIGeneration(component *TDesigningComponent) {
	if onUIGeneration != nil {
		onUIGeneration(component.formTab, component)
	}
}
