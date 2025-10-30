package uigen

// 生成回调事件

type CodeGenerationCallback func(uiFilePath string)

var (
	// Go代码生成回调
	onCodeGeneration CodeGenerationCallback
)

// 设置Go代码生成回调
func SetCodeGenerationCallback(callback CodeGenerationCallback) {
	onCodeGeneration = callback
}

// 触发Go代码生成事件
func triggerCodeGeneration(uiFilePath string) {
	if onCodeGeneration != nil {
		onCodeGeneration(uiFilePath)
	}
}
