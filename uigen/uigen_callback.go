package uigen

import (
	"github.com/energye/designer/designer"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// UI 布局文件回调函数

func InitUIGeneration() {
	designer.SetUIGenerationCallback(UIGeneration)
}

func UIGeneration(formTab *designer.FormTab, component *designer.TDesigningComponent) {
	DebouncedGenerate(formTab, component)
}

var (
	debounceTimers = make(map[string]*time.Timer)
	debounceMutex  sync.Mutex
	debounceDelay  = 500 * time.Millisecond
)

// UI生成
func DebouncedGenerate(formTab *designer.FormTab, component *designer.TDesigningComponent) {
	debounceMutex.Lock()
	defer debounceMutex.Unlock()
	formID := formTab.Name
	// 取消之前的定时器
	if timer, exists := debounceTimers[formID]; exists {
		timer.Stop()
	}

	// 创建新的定时器
	timer := time.AfterFunc(debounceDelay, func() {
		debounceMutex.Lock()
		delete(debounceTimers, formID)
		debounceMutex.Unlock()

		formId := strings.ToLower(formTab.Name)
		uiFilePath := filepath.Join(projectPath, "ui", formId+".ui")

		// 执行UI生成
		GenerateUIFile(formTab, component, uiFilePath)

		// 触发代码生成
		triggerCodeGeneration(uiFilePath)
	})

	debounceTimers[formID] = timer
}
