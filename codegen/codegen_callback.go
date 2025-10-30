package codegen

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/uigen"
)

func InitCodeGeneration() {
	uigen.SetCodeGenerationCallback(CodeGeneration)
}

func CodeGeneration(uiFilePath string) {
	err := GenerateCode(uiFilePath)
	if err != nil {
		logs.Error("代码生成错误:", err.Error())
	}
}
