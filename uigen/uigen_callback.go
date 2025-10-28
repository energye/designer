package uigen

import (
	"github.com/energye/designer/designer"
	"sync"
	"time"
)

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

// 带防抖的UI生成
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

		// 执行UI生成
		GenerateUIFile(formTab, component)

		// 触发代码生成
		//codegen.GenerateCode(filePath)
	})

	debounceTimers[formID] = timer
}
