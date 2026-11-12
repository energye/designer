package webview

import (
	"github.com/energye/designer/designer/editor"
	"github.com/energye/energy/v3/ipc"
)

func init() {
	// 前端页面加载完成后主动同步一次 gopls 状态。
	ipc.On("gopls-status-request", func(context ipc.IContext) {
		ok, reason := editor.PLSStatus()
		context.Result(map[string]interface{}{
			"ok":       ok,
			"reason":   reason,
			"indexing": editor.PLSIndexing(),
		})
	})
}
