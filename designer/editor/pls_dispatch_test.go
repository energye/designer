package editor

import (
	"sync"
	"testing"
	"time"

	"github.com/energye/designer/designer/editor/gopls"
)

// newRecorderDispatcher 构造一个不启动 LCL/WebView、同步可排空的派发器。
func newRecorderDispatcher(t *testing.T) (*plsDispatcher, *[]PLSNotification, *sync.Mutex) {
	t.Helper()
	var mu sync.Mutex
	sent := make([]PLSNotification, 0)
	fake := &gopls.PLSClient{}
	d := newPLSDispatcher()
	d.clientFn = func() *gopls.PLSClient { return fake }
	d.readyFn = func() bool { return true }
	d.sendFn = func(_ *gopls.PLSClient, n PLSNotification) error {
		mu.Lock()
		sent = append(sent, n)
		mu.Unlock()
		return nil
	}
	return d, &sent, &mu
}

func snapshot(mu *sync.Mutex, sent *[]PLSNotification) []PLSNotification {
	mu.Lock()
	defer mu.Unlock()
	return append([]PLSNotification(nil), (*sent)...)
}

func TestPLSDispatcher_FIFOOrder(t *testing.T) {
	d, sent, mu := newRecorderDispatcher(t)
	const uri = "file:///e:/app/a.go"
	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: uri, LanguageID: "go", Content: "v1", Version: 1})
	d.Enqueue(PLSNotification{Kind: PLSNotifyChange, URI: uri, Content: "v2", Version: 2})
	d.Enqueue(PLSNotification{Kind: PLSNotifyClose, URI: uri})
	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: uri, LanguageID: "go", Content: "disk", Version: 1})
	d.Enqueue(PLSNotification{Kind: PLSNotifyChange, URI: uri, Content: "v3", Version: 2})
	d.Enqueue(PLSNotification{Kind: PLSNotifySave, URI: uri, Content: "v3"})
	d.drain()

	got := snapshot(mu, sent)
	wantKinds := []PLSNotifyKind{
		PLSNotifyOpen, PLSNotifyChange, PLSNotifyClose,
		PLSNotifyOpen, PLSNotifyChange, PLSNotifySave,
	}
	if len(got) != len(wantKinds) {
		t.Fatalf("通知数量不对: %d", len(got))
	}
	for i, k := range wantKinds {
		if got[i].Kind != k {
			t.Fatalf("第 %d 条通知顺序错误: got %v want %v", i, got[i].Kind, k)
		}
	}
	// version 由 Go 侧单调分配：didClose/didOpen 重载不重置为 1。
	wantVersions := []int{1, 2, 0, 3, 4, 0}
	for i, v := range wantVersions {
		if got[i].Version != v {
			t.Fatalf("第 %d 条 version 错误: got %d want %d", i, got[i].Version, v)
		}
	}
}

func TestPLSDispatcher_VersionNeverDecreases(t *testing.T) {
	d, sent, mu := newRecorderDispatcher(t)
	const uri = "file:///e:/app/a.go"
	d.Enqueue(PLSNotification{Kind: PLSNotifyChange, URI: uri, Version: 100})
	d.Enqueue(PLSNotification{Kind: PLSNotifyChange, URI: uri, Version: 1})
	d.Enqueue(PLSNotification{Kind: PLSNotifyClose, URI: uri})
	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: uri, Version: 1})
	d.drain()

	got := snapshot(mu, sent)
	// Close/Save 不发 version，Open/Change 必须严格递增。
	versions := []int{got[0].Version, got[1].Version, got[3].Version}
	for i, v := range versions {
		want := []int{100, 101, 102}[i]
		if v != want {
			t.Fatalf("version 应单调递增: [%d] got %d want %d", i, v, want)
		}
	}
}

func TestPLSDispatcher_QueuesUntilGoplsReady(t *testing.T) {
	d, sent, mu := newRecorderDispatcher(t)
	var readyMu sync.Mutex
	client := (*gopls.PLSClient)(nil)
	ready := false
	d.clientFn = func() *gopls.PLSClient {
		readyMu.Lock()
		defer readyMu.Unlock()
		return client
	}
	d.readyFn = func() bool {
		readyMu.Lock()
		defer readyMu.Unlock()
		return ready
	}
	d.pollInterval = time.Millisecond
	sentCh := make(chan struct{}, 1)
	d.sendFn = func(_ *gopls.PLSClient, n PLSNotification) error {
		mu.Lock()
		*sent = append(*sent, n)
		mu.Unlock()
		select {
		case sentCh <- struct{}{}:
		default:
		}
		return nil
	}
	go d.loop()

	// gopls 尚未 initialize 完成时，didOpen 必须保留在队列里而不是被丢弃。
	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: "file:///a.go", Content: "x", Version: 1})
	time.Sleep(20 * time.Millisecond)
	if len(snapshot(mu, sent)) != 0 {
		t.Fatal("gopls 未就绪时不应发送通知")
	}

	readyMu.Lock()
	client = &gopls.PLSClient{}
	ready = true
	readyMu.Unlock()
	select {
	case <-sentCh:
	case <-time.After(2 * time.Second):
		t.Fatal("gopls 就绪后应重放排队的 didOpen")
	}
	got := snapshot(mu, sent)
	if len(got) != 1 || got[0].Kind != PLSNotifyOpen || got[0].Content != "x" {
		t.Fatalf("排队的 didOpen 内容丢失: %+v", got)
	}
}

func TestPLSDispatcher_DropsNotificationWhenNeverReady(t *testing.T) {
	d, sent, mu := newRecorderDispatcher(t)
	d.clientFn = func() *gopls.PLSClient { return nil }
	d.readyFn = func() bool { return false }
	d.waitTimeout = 5 * time.Millisecond
	d.pollInterval = time.Millisecond
	go d.loop()

	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: "file:///a.go", Version: 1})
	deadline := time.Now().Add(150 * time.Millisecond)
	for time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if len(snapshot(mu, sent)) != 0 {
		t.Fatal("gopls 永久不可用时不应向 nil 客户端发通知")
	}
}

func TestPLSDispatcher_ResetClearsPendingAndVersions(t *testing.T) {
	d, sent, mu := newRecorderDispatcher(t)
	const uri = "file:///a.go"
	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: uri, Version: 7})
	d.drain()
	if got := snapshot(mu, sent); len(got) != 1 || got[0].Version != 7 {
		t.Fatalf("首个 didOpen version 不正确: %+v", got)
	}

	d.Enqueue(PLSNotification{Kind: PLSNotifyChange, URI: uri, Version: 8})
	d.Reset() // 客户端被替换/关闭，旧通知与旧版本全部作废
	d.Enqueue(PLSNotification{Kind: PLSNotifyOpen, URI: uri, Version: 1})
	d.drain()

	got := snapshot(mu, sent)
	if len(got) != 2 {
		t.Fatalf("Reset 应丢弃未发送的通知，实际 %d 条", len(got))
	}
	if got[1].Version != 1 {
		t.Fatalf("Reset 后版本应重新从 1 开始，实际 %d", got[1].Version)
	}
}

func TestPLSDispatcher_ConcurrentEnqueue(t *testing.T) {
	d, sent, mu := newRecorderDispatcher(t)
	const n = 100
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.Enqueue(PLSNotification{Kind: PLSNotifyChange, URI: "file:///a.go"})
		}()
	}
	wg.Wait()
	d.drain()

	got := snapshot(mu, sent)
	if len(got) != n {
		t.Fatalf("应有 %d 条通知，实际 %d", n, len(got))
	}
	for i, item := range got {
		if item.Version != i+1 {
			t.Fatalf("并发下 version 必须严格递增: [%d]=%d", i, item.Version)
		}
	}
}

func TestSetPLSIndexing_EmitsOnlyOnChange(t *testing.T) {
	orig := plsIndexingEmitter
	var mu sync.Mutex
	var events []bool
	plsIndexingEmitter = func(indexing bool) {
		mu.Lock()
		events = append(events, indexing)
		mu.Unlock()
	}
	t.Cleanup(func() {
		plsIndexingEmitter = orig
		setPLSIndexing(false)
	})

	setPLSIndexing(true)
	setPLSIndexing(true)
	if !PLSIndexing() {
		t.Fatal("索引状态应为 true")
	}
	setPLSIndexing(false)
	if PLSIndexing() {
		t.Fatal("索引状态应为 false")
	}
	mu.Lock()
	defer mu.Unlock()
	if len(events) != 2 || !events[0] || events[1] {
		t.Fatalf("状态变化应只推两次: %v", events)
	}
}
