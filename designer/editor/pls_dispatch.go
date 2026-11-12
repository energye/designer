// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations
// under the License.

package editor

import (
	"sync"
	"time"

	"github.com/energye/designer/designer/editor/gopls"
	"github.com/energye/designer/pkg/logs"
)

// PLSNotifyKind 发往 gopls 的 textDocument 通知类型。
type PLSNotifyKind uint8

const (
	PLSNotifyOpen PLSNotifyKind = iota
	PLSNotifyChange
	PLSNotifyClose
	PLSNotifySave
)

// PLSNotification 一条待串行发送的通知。
// Version 只是前端给出的单调提示值；真正发给 gopls 的 version 由派发器
// 统一分配，保证 didClose/didOpen 重载后也不会回退成 1。
type PLSNotification struct {
	Kind       PLSNotifyKind
	URI        string
	LanguageID string
	Content    string
	Version    int

	enqueuedAt time.Time
	queueGen   uint64
}

const (
	defaultPLSReadyWaitTimeout = 90 * time.Second
	defaultPLSReadyPoll        = 25 * time.Millisecond
)

// plsDispatcher 把 didOpen/didChange/didClose/didSave 收敛到单个 FIFO worker。
//
// 旧实现每个 IPC 各起一个 goroutine，重载时 didClose→didOpen 可能被在途的
// didChange 插队，gopls 会拿到顺序颠倒（且 version 回退）的文档状态。
// 同时这里还负责等 gopls 完成 initialize：客户端未发布之前不丢通知，
// 否则启动/首次索引期间打开的文件永远不会被 didOpen。
type plsDispatcher struct {
	mu       sync.Mutex
	cond     *sync.Cond
	queue    []PLSNotification
	versions map[string]int
	queueGen uint64

	waitTimeout  time.Duration
	pollInterval time.Duration
	clientFn     func() *gopls.PLSClient
	readyFn      func() bool
	sendFn       func(*gopls.PLSClient, PLSNotification) error
	nowFn        func() time.Time
}

var (
	plsDispatchOnce sync.Once
	gPLSDispatcher  = newPLSDispatcher()
)

func newPLSDispatcher() *plsDispatcher {
	d := &plsDispatcher{
		versions:     make(map[string]int),
		waitTimeout:  defaultPLSReadyWaitTimeout,
		pollInterval: defaultPLSReadyPoll,
		clientFn:     func() *gopls.PLSClient { return PLSClient() },
		readyFn:      func() bool { return PLSClient() != nil },
		sendFn:       defaultPLSSend,
		nowFn:        time.Now,
	}
	d.cond = sync.NewCond(&d.mu)
	return d
}

// StartPLSDispatcher 进程级只启动一次 worker；重复启动编辑器不会重复消费队列。
func StartPLSDispatcher() {
	plsDispatchOnce.Do(func() { go gPLSDispatcher.loop() })
}

// EnqueuePLSNotification 按 IPC 到达顺序排队，由单 worker 串行写出。
func EnqueuePLSNotification(n PLSNotification) {
	gPLSDispatcher.Enqueue(n)
}

// ResetPLSQueue 丢弃尚未发送的通知和文档版本：gopls 客户端被替换/关闭后，
// 旧客户端的文档状态已失效，新页面会重新发送 didOpen。
func ResetPLSQueue() {
	gPLSDispatcher.Reset()
}

func (d *plsDispatcher) Enqueue(n PLSNotification) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if n.enqueuedAt.IsZero() {
		n.enqueuedAt = d.nowFn()
	}
	n.queueGen = d.queueGen
	d.queue = append(d.queue, n)
	d.cond.Signal()
}

func (d *plsDispatcher) Reset() {
	d.mu.Lock()
	d.queue = nil
	d.versions = make(map[string]int)
	d.queueGen++
	d.cond.Broadcast()
	d.mu.Unlock()
}

func (d *plsDispatcher) loop() {
	for {
		d.process(d.waitNext())
	}
}

func (d *plsDispatcher) waitNext() PLSNotification {
	d.mu.Lock()
	defer d.mu.Unlock()
	for len(d.queue) == 0 {
		d.cond.Wait()
	}
	n := d.queue[0]
	d.queue[0] = PLSNotification{}
	d.queue = d.queue[1:]
	return n
}

func (d *plsDispatcher) process(n PLSNotification) {
	// 未就绪时保留队首并阻塞后续通知：既保证 didOpen 不早于 initialize，
	// 也保证重载序列 didClose→didOpen 不被在途 didChange 插队。
	for {
		if d.aborted(n) {
			return
		}
		client := d.clientFn()
		if client != nil && d.readyFn() {
			d.dispatch(client, n)
			return
		}
		if d.nowFn().Sub(n.enqueuedAt) >= d.waitTimeout {
			logs.Error("gopls 未就绪，丢弃过期的文档通知:", n.URI)
			return
		}
		time.Sleep(d.pollInterval)
	}
}

func (d *plsDispatcher) aborted(n PLSNotification) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return n.queueGen != d.queueGen
}

func (d *plsDispatcher) dispatch(client *gopls.PLSClient, n PLSNotification) {
	if d.aborted(n) {
		return
	}
	if n.Kind == PLSNotifyOpen || n.Kind == PLSNotifyChange {
		n.Version = d.nextVersion(n)
	}
	if err := d.sendFn(client, n); err != nil {
		logs.Error("gopls 文档通知发送失败:", n.URI, err)
	}
}

// nextVersion 保证同一文档的 version 单调递增：didClose/didOpen 重载不重置，
// 前端 reload 传回的 1 也不会把版本拉回去（只取 max）。
func (d *plsDispatcher) nextVersion(n PLSNotification) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	key := CanonicalPathKey(n.URI)
	next := d.versions[key] + 1
	if n.Version > next {
		next = n.Version
	}
	if next < 1 {
		next = 1
	}
	d.versions[key] = next
	return next
}

func defaultPLSSend(c *gopls.PLSClient, n PLSNotification) error {
	switch n.Kind {
	case PLSNotifyOpen:
		return c.DidOpen(n.URI, n.LanguageID, n.Content, n.Version)
	case PLSNotifyChange:
		return c.DidChange(n.URI, n.Version, n.Content)
	case PLSNotifyClose:
		return c.DidClose(n.URI)
	case PLSNotifySave:
		return c.DidSave(n.URI, n.Content)
	}
	return nil
}

// drain 同步排空队列，仅用于单元测试（注入的 client/ready 均为真）。
func (d *plsDispatcher) drain() {
	for {
		d.mu.Lock()
		if len(d.queue) == 0 {
			d.mu.Unlock()
			return
		}
		n := d.queue[0]
		d.queue = d.queue[1:]
		d.mu.Unlock()
		d.process(n)
	}
}
