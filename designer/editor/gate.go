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
// See the License for the specific language governing permissions and limitations under the License.

package editor

import (
	"sync"
	"time"

	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/ipc"
	"github.com/energye/lcl/lcl"
)

// 代码生成写盘门禁。
//
// 不变量：codegen 覆盖 .go 文件之前，磁盘内容必须等于编辑器内容。
// 因此有未保存修改时绝不写盘，也不在 Go 侧缓存文本副本（副本会滞后）；
// 而是让 Monaco 按当前 model 走 save-file 保存，保存完成后重放生成。
// 测试时可缩短超时；生产环境固定 10s。
var generateSaveTimeout = 10 * time.Second

// emitSaveAll 抽成变量，便于在不启动 LCL/WebView 的单元测试里替换。
var emitSaveAll = func(requestID uint64) {
	lcl.RunOnMainThreadAsync(func(id uint32) {
		ipc.Emit("save-all-request", requestID)
	})
}

var (
	gateMu           sync.Mutex
	pendingGenerates []func()
	gateTimer        *time.Timer
	gateSeq          uint64
)

// HasDirtyFiles 是否有打开的文件存在未保存修改。
func HasDirtyFiles() bool {
	ed := GetCurrentEditor()
	if ed == nil {
		return false
	}
	fm := ed.FileManager()
	if fm == nil {
		return false
	}
	for p := range fm.GetAllFiles() {
		if fm.IsDirty(p) {
			return true
		}
	}
	return false
}

// IsFileDirty 指定文件是否存在未保存修改。
func IsFileDirty(filePath string) bool {
	ed := GetCurrentEditor()
	if ed == nil {
		return false
	}
	fm := ed.FileManager()
	if fm == nil {
		return false
	}
	return fm.IsDirty(filePath)
}

// AllowGenerate 写盘门禁：无脏文件返回 true；有脏文件时把 resume 排队、
// 请求前端保存全部脏文件，并返回 false —— 调用方必须跳过本次生成。
// resume 会在所有脏文件保存完成后被重放，不在 UI 线程里同步执行。
func AllowGenerate(resume func()) bool {
	gateMu.Lock()
	if !HasDirtyFiles() {
		gateMu.Unlock()
		return true
	}
	pendingGenerates = append(pendingGenerates, resume)
	var requestID uint64
	if gateTimer == nil {
		gateSeq++
		requestID = gateSeq
		gateTimer = time.AfterFunc(generateSaveTimeout, func() {
			gateMu.Lock()
			defer gateMu.Unlock()
			if requestID != gateSeq {
				return
			}
			pendingGenerates = nil
			gateTimer = nil
			logs.Error("等待编辑器保存超时，已取消排队的代码生成；请手动保存后重试")
		})
	}
	gateMu.Unlock()

	if requestID != 0 {
		emitSaveAll(requestID)
	}
	return false
}

// CancelPendingGenerates 丢弃排队的代码生成重放。
// 构建前会基于最新 UI 统一 flush，不再需要旧请求重放。
func CancelPendingGenerates() {
	gateMu.Lock()
	pendingGenerates = nil
	if gateTimer != nil {
		gateTimer.Stop()
		gateTimer = nil
	}
	gateSeq++
	gateMu.Unlock()
}

// SaveAllAndWait 请求前端保存全部脏文件并等待完成。
// 构建/运行前使用，避免 build 读到旧内容；必须在非 UI 线程调用。
func SaveAllAndWait(timeout time.Duration) bool {
	if !HasDirtyFiles() {
		return true
	}
	emitSaveAll(0)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		time.Sleep(50 * time.Millisecond)
		if !HasDirtyFiles() {
			return true
		}
	}
	return !HasDirtyFiles()
}

// NotifyFileSaved 由 save-file 在保存成功后调用。
// 只有当 Go 侧确认已无脏文件时，才重放排队的代码生成。
func NotifyFileSaved() {
	gateMu.Lock()
	if len(pendingGenerates) == 0 || HasDirtyFiles() {
		gateMu.Unlock()
		return
	}
	list := pendingGenerates
	pendingGenerates = nil
	if gateTimer != nil {
		gateTimer.Stop()
		gateTimer = nil
	}
	gateSeq++
	gateMu.Unlock()

	go func() {
		for _, resume := range list {
			resume()
		}
	}()
}
