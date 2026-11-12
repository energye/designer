package editor

import (
	"sync"

	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
)

type plsStatus struct {
	ok       bool
	reason   string
	indexing bool
}

var (
	plsStatusMu sync.RWMutex
	gPLSStatus  = plsStatus{ok: false, reason: "gopls 尚未启动"}

	// 抽成变量，便于在不启动 LCL/WebView 的单元测试里替换。
	plsStatusEmitter = func(ok bool, reason string) {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			ipc.Emit("gopls-status", ok, reason)
		})
	}
	plsIndexingEmitter = func(indexing bool) {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			ipc.Emit("gopls-indexing", indexing)
		})
	}
)

// setPLSStatus 记录状态并推给代码编辑页。前端也会主动查询一次，
// 避免状态上报早于 WebView 的 IPC 注册而丢失。
func setPLSStatus(ok bool, reason string) {
	plsStatusMu.Lock()
	gPLSStatus.ok = ok
	gPLSStatus.reason = reason
	plsStatusMu.Unlock()
	plsStatusEmitter(ok, reason)
}

// PLSStatus 返回当前 gopls 状态。
func PLSStatus() (bool, string) {
	plsStatusMu.RLock()
	defer plsStatusMu.RUnlock()
	return gPLSStatus.ok, gPLSStatus.reason
}

// PLSIndexing 是否处于首次索引窗口（请求可能需要几十秒）。
func PLSIndexing() bool {
	plsStatusMu.RLock()
	defer plsStatusMu.RUnlock()
	return gPLSStatus.indexing
}

// setPLSIndexing 索引开始/结束时通知前端延长或恢复 pending 超时。
func setPLSIndexing(indexing bool) {
	plsStatusMu.Lock()
	if gPLSStatus.indexing == indexing {
		plsStatusMu.Unlock()
		return
	}
	gPLSStatus.indexing = indexing
	plsStatusMu.Unlock()
	plsIndexingEmitter(indexing)
}
