package uigen

import (
	"github.com/energye/designer/designer"
	"github.com/energye/designer/event"
)

// 生成器实例
var gen = &TGenUI{trigger: make(chan event.TGeneratorTrigger, 1), cancel: make(chan bool, 1)}

// TGenUI UI生成器
type TGenUI struct {
	trigger chan event.TGeneratorTrigger // 触发UI生成事件
	cancel  chan bool                    // 取消UI生成事件
}

// Start 启动UI生成器
func (m *TGenUI) Start() {
	for {
		select {
		case trigger := <-m.trigger:
			// 处理UI生成事件
			if trigger.GenType == event.GtUI { //增强判断, 确保是UI生成事件
				if formTab, ok := trigger.Payload.(*designer.FormTab); ok {
					DebouncedGenerate(formTab)
				}
			}
		case <-m.cancel:
			// 停止UI生成器
			return
		}
	}
}

func init() {
	// 注册UI生成事件
	event.GenUI = event.NewGenerator(gen.trigger, gen.cancel)
	// 启动UI生成器
	go gen.Start()
}
