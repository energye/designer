package editor

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// mockEditor 只实现 IEditor 里门禁需要的方法。
type mockEditor struct {
	fm *FileManager
}

func (m *mockEditor) OpenFile(string, ...bool)  {}
func (m *mockEditor) CloseFile(string)          {}
func (m *mockEditor) SaveCurrentFile()          {}
func (m *mockEditor) FileManager() *FileManager { return m.fm }
func (m *mockEditor) Stop()                     {}

func resetGateForTest() {
	gateMu.Lock()
	defer gateMu.Unlock()
	if gateTimer != nil {
		gateTimer.Stop()
		gateTimer = nil
	}
	pendingGenerates = nil
	gateSeq++
}

func newTestEditor(t *testing.T, names ...string) (*FileManager, []string, func()) {
	t.Helper()
	dir := t.TempDir()
	fm := NewFileManager()
	paths := make([]string, 0, len(names))
	for _, name := range names {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte("package main\n"), 0644); err != nil {
			t.Fatal(err)
		}
		fm.RegisterFile(p)
		paths = append(paths, p)
	}
	SetCurrentEditor(&mockEditor{fm: fm})

	origEmit := emitSaveAll
	origTimeout := generateSaveTimeout
	emitSaveAll = func(uint64) {}
	cleanup := func() {
		SetCurrentEditor(nil)
		emitSaveAll = origEmit
		generateSaveTimeout = origTimeout
		resetGateForTest()
	}
	return fm, paths, cleanup
}

func TestCancelPendingGenerates_DropsReplay(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()
	fm.SetDirty(paths[0], true)

	replayed := 0
	if AllowGenerate(func() { replayed++ }) {
		t.Fatal("脏文件存在时不应立即生成")
	}
	CancelPendingGenerates()
	fm.SetDirty(paths[0], false)
	NotifyFileSaved()
	time.Sleep(20 * time.Millisecond)
	if replayed != 0 {
		t.Fatalf("取消后不应重放, got %d", replayed)
	}
}

func TestSaveAllAndWait_NoDirty(t *testing.T) {
	_, _, cleanup := newTestEditor(t, "a.go")
	defer cleanup()

	if !SaveAllAndWait(10 * time.Millisecond) {
		t.Fatal("无脏文件应立即返回 true")
	}
}

func TestSaveAllAndWait_EmitsSaveAllAndTimesOut(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()
	fm.SetDirty(paths[0], true)

	emitted := 0
	emitSaveAll = func(uint64) { emitted++ }
	if SaveAllAndWait(10 * time.Millisecond) {
		t.Fatal("脏文件未保存时应超时返回 false")
	}
	if emitted != 1 {
		t.Fatalf("应只请求一次 save-all, got %d", emitted)
	}
}

func TestHasDirtyFiles_NoEditor(t *testing.T) {
	SetCurrentEditor(nil)
	resetGateForTest()
	if HasDirtyFiles() {
		t.Fatal("无编辑器实例时不应报告脏文件")
	}
	if IsFileDirty("anything.go") {
		t.Fatal("无编辑器实例时不应报告脏文件")
	}
}

func TestGate_NoDirtyAllowsImmediately(t *testing.T) {
	_, _, cleanup := newTestEditor(t, "a.go")
	defer cleanup()

	called := false
	if !AllowGenerate(func() { called = true }) {
		t.Fatal("无脏文件时应直接放行生成")
	}
	if called {
		t.Fatal("放行时不应执行 resume")
	}
}

func TestGate_DirtyBlocksThenReplaysOnSave(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()

	fm.SetDirty(paths[0], true)
	if !HasDirtyFiles() || !IsFileDirty(paths[0]) {
		t.Fatal("脏状态未生效")
	}

	done := make(chan struct{})
	if AllowGenerate(func() { close(done) }) {
		t.Fatal("有脏文件时应阻止生成")
	}
	select {
	case <-done:
		t.Fatal("保存完成前不应重放生成")
	case <-time.After(50 * time.Millisecond):
	}

	// 模拟 Monaco 走 save-file 保存完毕
	fm.SetDirty(paths[0], false)
	NotifyFileSaved()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("脏状态清空后应重放生成")
	}
}

func TestGate_PartialSaveDoesNotReplay(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go", "b.go")
	defer cleanup()

	fm.SetDirty(paths[0], true)
	fm.SetDirty(paths[1], true)
	done := make(chan struct{})
	if AllowGenerate(func() { close(done) }) {
		t.Fatal("有脏文件时应阻止生成")
	}

	fm.SetDirty(paths[0], false)
	NotifyFileSaved()
	select {
	case <-done:
		t.Fatal("仍有其它脏文件时不应重放")
	case <-time.After(50 * time.Millisecond):
	}

	fm.SetDirty(paths[1], false)
	NotifyFileSaved()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("全部保存后应重放")
	}
}

func TestGate_MultipleGenerationsReplayInOrder(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()

	fm.SetDirty(paths[0], true)
	var mu sync.Mutex
	var order []int
	done := make(chan struct{})
	for i := 0; i < 3; i++ {
		i := i
		if AllowGenerate(func() {
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
			if i == 2 {
				close(done)
			}
		}) {
			t.Fatal("有脏文件时应阻止生成")
		}
	}
	gateMu.Lock()
	queued := len(pendingGenerates)
	gateMu.Unlock()
	if queued != 3 {
		t.Fatalf("排队的生成数应为 3，实际 %d", queued)
	}

	fm.SetDirty(paths[0], false)
	NotifyFileSaved()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("排队的生成应全部重放")
	}
	mu.Lock()
	defer mu.Unlock()
	for i, v := range order {
		if v != i {
			t.Fatalf("重放顺序错误: %v", order)
		}
	}
}

func TestGate_TimeoutDropsQueue(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()
	generateSaveTimeout = 30 * time.Millisecond

	fm.SetDirty(paths[0], true)
	done := make(chan struct{})
	AllowGenerate(func() { close(done) })

	time.Sleep(100 * time.Millisecond)
	gateMu.Lock()
	queued := len(pendingGenerates)
	gateMu.Unlock()
	if queued != 0 {
		t.Fatalf("超时后应丢弃队列，实际仍有 %d 个", queued)
	}

	fm.SetDirty(paths[0], false)
	NotifyFileSaved()
	select {
	case <-done:
		t.Fatal("超时丢弃的生成不应被重放")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestGate_SaveAllRequestedOnlyOncePerBatch(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()

	var mu sync.Mutex
	requestCount := 0
	emitSaveAll = func(uint64) {
		mu.Lock()
		requestCount++
		mu.Unlock()
	}

	fm.SetDirty(paths[0], true)
	AllowGenerate(func() {})
	AllowGenerate(func() {})
	AllowGenerate(func() {})

	mu.Lock()
	defer mu.Unlock()
	if requestCount != 1 {
		t.Fatalf("同一批生成应只发一次 save-all-request，实际 %d 次", requestCount)
	}
}

func TestStopPLS_ClearsClientAndIsIdempotent(t *testing.T) {
	StopPLS()
	if PLSClient() != nil {
		t.Fatal("StopPLS 后客户端应为空")
	}
	StopPLS() // 二次调用不应 panic
}

func TestIsFileReadOnly_ModuleCache(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("GOMODCACHE", dir)
	p := filepath.Join(dir, "pkg.go")
	if err := os.WriteFile(p, []byte("package pkg\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if !IsFileReadOnly(p) {
		t.Fatal("GOMODCACHE 下的文件应判为只读")
	}
}

func TestFileManager_DirtyStateConcurrent(t *testing.T) {
	fm, paths, cleanup := newTestEditor(t, "a.go")
	defer cleanup()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fm.SetDirty(paths[0], i%2 == 0)
			_ = fm.IsDirty(paths[0])
			_ = fm.GetAllFiles()
			_ = HasDirtyFiles()
		}(i)
	}
	wg.Wait()
}
