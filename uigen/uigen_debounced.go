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

package uigen

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/energye/designer/codegen"
	"github.com/energye/designer/designer"
	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/event"
	projBean "github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
)

// UI 布局文件生成

// runOnUIThread 同步在 UI 线程执行；已在 UI 线程时直接执行，避免同步等待死锁。
func runOnUIThread(fn func()) {
	// CLI/单测场景没有 LCL 运行时，直接执行，避免加载 UI 库。
	if lcl.Application == nil || !lcl.Application.Initialed() {
		fn()
		return
	}
	if tool.IsMainThread() {
		fn()
		return
	}
	lcl.RunOnMainThreadSync(fn)
}

type pendingGeneration struct {
	data  designer.TUIGenerationData
	type_ event.Type
	timer *time.Timer
	epoch uint64
}

var (
	pendingGenerations = make(map[int]*pendingGeneration)
	debounceMutex      sync.Mutex
	debounceDelay      = 500 * time.Millisecond // 最大更新间隔

	// generationEpoch 用于 FlushAll 使已触发的旧定时器失效。
	generationEpoch uint64
	// generationMu 串行化 .ui 写盘与代码生成，避免构建前 flush 与定时器并发。
	generationMu sync.Mutex
)

// UI布局文件生成
func runDebouncedGenerate(uiGenData designer.TUIGenerationData, type_ event.Type) {
	formTab := uiGenData.Component.FormTab
	if formTab == nil {
		return
	}
	formId := formTab.Id

	debounceMutex.Lock()
	if old, exists := pendingGenerations[formId]; exists {
		old.timer.Stop()
	}
	pending := &pendingGeneration{data: uiGenData, type_: type_, epoch: generationEpoch}
	pending.timer = time.AfterFunc(debounceDelay, func() {
		flushPending(formId, pending.epoch)
	})
	pendingGenerations[formId] = pending
	debounceMutex.Unlock()
}

// flushPending 执行某个窗体到期的防抖生成；epoch 变化说明已被 FlushAll 接管。
func flushPending(formId int, epoch uint64) {
	debounceMutex.Lock()
	pending, ok := pendingGenerations[formId]
	if !ok || pending.epoch != epoch || pending.epoch != generationEpoch {
		debounceMutex.Unlock()
		return
	}
	delete(pendingGenerations, formId)
	debounceMutex.Unlock()

	pending.timer.Stop()
	generateAndCodegen(pending.data, pending.type_, epoch)
}

// generateAndCodegen 同步落盘 .ui 并立即生成代码。
// .ui 的序列化放到 UI 线程，避免与设计面板的属性更新并发。
func generateAndCodegen(uiGenData designer.TUIGenerationData, type_ event.Type, epoch uint64) {
	formTab := uiGenData.Component.FormTab
	if formTab == nil {
		return
	}

	generationMu.Lock()
	defer generationMu.Unlock()

	// 等待 generationMu 期间可能已被 FlushAll 接管，此时旧定时器不能再落盘。
	debounceMutex.Lock()
	stale := epoch != generationEpoch
	debounceMutex.Unlock()
	if stale {
		return
	}

	var (
		ok           bool
		isUpdateSelf bool
	)
	runOnUIThread(func() {
		if formTab.FormRoot.Name() == "" {
			// 正在加载设计时同时关闭设计窗口，跳过生成
			return
		}
		isUpdateSelf = tryRenameFileName(formTab)
		// 以窗体实时宽高兜底，避免快速拖拽/丢失 MouseUp 时属性滞后。
		size := formTab.FormRoot.BoundsRect()
		formTab.FormRoot.SyncSizeToProperties(size.Width(), size.Height())
		uiFilePath := filepath.Join(projBean.LayoutsPath(), formTab.UIFile())
		if err := generateUIFile(formTab.FormRoot, uiFilePath); err != nil {
			logs.Error("UI布局文件生成错误:", err.Error())
			return
		}
		ok = true
	})
	if !ok {
		return
	}

	// uiGenData 必须保留 NodeData：绑定事件更新依赖它定位具体事件属性。
	if isUpdateSelf {
		// 自引用 > 用户代码
		codegen.Generate(uiGenData, event.CodeGenSelf)
	} else if type_ == event.CodeGenEvent {
		// 绑定事件 > 用户代码
		codegen.Generate(uiGenData, event.CodeGenEvent)
	}
	// UI 布局, 全量更新
	codegen.Generate(uiGenData, event.CodeGenUI)
	triggerProjectUpdate(formTab)
}

// FlushAll 同步落盘所有打开窗体的 .ui 并生成代码。
// 构建/运行前调用，保证 build 读取到的布局与生成代码都是最新状态。
// 注意：不得从 UI 线程调用。
func FlushAll() {
	if projBean.GProject == nil {
		return
	}
	debounceMutex.Lock()
	generationEpoch++
	types := make(map[int]event.Type, len(pendingGenerations))
	datas := make(map[int]designer.TUIGenerationData, len(pendingGenerations))
	for id, pending := range pendingGenerations {
		pending.timer.Stop()
		types[id] = pending.type_
		datas[id] = pending.data
		delete(pendingGenerations, id)
	}
	debounceMutex.Unlock()

	generationMu.Lock()
	defer generationMu.Unlock()

	var forms []*designer.FormTab
	runOnUIThread(func() {
		forms = designer.FormTabs()
	})
	for _, formTab := range forms {
		if formTab == nil {
			continue
		}
		var (
			ok           bool
			isUpdateSelf bool
		)
		runOnUIThread(func() {
			if formTab.FormRoot.Name() == "" {
				return
			}
			isUpdateSelf = tryRenameFileName(formTab)
			// 构建前以窗体实时宽高兜底，保证窗口大小不依赖最后一次鼠标事件。
			size := formTab.FormRoot.BoundsRect()
			formTab.FormRoot.SyncSizeToProperties(size.Width(), size.Height())
			uiFilePath := filepath.Join(projBean.LayoutsPath(), formTab.UIFile())
			if err := generateUIFile(formTab.FormRoot, uiFilePath); err != nil {
				logs.Error("UI布局文件生成错误:", err.Error())
				return
			}
			ok = true
		})
		if !ok {
			continue
		}
		uiGenData, hasData := datas[formTab.Id]
		if !hasData {
			uiGenData = designer.TUIGenerationData{Component: formTab.FormRoot}
		}
		if isUpdateSelf {
			codegen.Generate(uiGenData, event.CodeGenSelf)
		} else if types[formTab.Id] == event.CodeGenEvent {
			codegen.Generate(uiGenData, event.CodeGenEvent)
		}
		codegen.Generate(uiGenData, event.CodeGenUI)
		triggerProjectUpdate(formTab)
	}

	// flush 期间若又产生了 GenUI 事件，统一丢弃：本次已按实时状态全量落盘，
	// 若让它们留在队列里，构建期间会再次写文件。
	debounceMutex.Lock()
	generationEpoch++
	for id, pending := range pendingGenerations {
		pending.timer.Stop()
		delete(pendingGenerations, id)
	}
	debounceMutex.Unlock()
}

// PrepareForBuild 构建/运行前准备：先保存编辑器脏文件，再同步落盘并生成代码。
// 由 cmd/build 在构建前调用；必须在非 UI 线程执行。
func PrepareForBuild() bool {
	// 旧的重放请求由本次 FlushAll 统一重新生成，先丢弃避免构建期间重复写文件。
	editor.CancelPendingGenerates()
	if !editor.SaveAllAndWait(10 * time.Second) {
		logs.Error("构建前等待编辑器保存超时")
		event.ConsoleWriteError("代码编辑器保存超时，已取消构建")
		return false
	}
	FlushAll()
	return true
}

// 尝试更新文件名
// 如果窗体名称被改变, 修改文件名
//
//	修改文件:
//	xxx.ui
//	xxx.ui.go
//	xxx.go
func tryRenameFileName(tempFormTab *designer.FormTab) bool {
	// ui 布局文件名
	uiFileName := tempFormTab.UIFile()

	// 验证UI布局文件名。必须取切片元素地址，否则只会改到循环副本。
	var uiForm *projBean.TUIForm
	for i := range projBean.GProject.UIForms {
		if projBean.GProject.UIForms[i].Id == tempFormTab.Id {
			uiForm = &projBean.GProject.UIForms[i]
			break
		}
	}
	// 判断 UI 文件名是否相同, 如果不相同则更改当前窗体的文件名
	if uiForm != nil && uiForm.UIFile != uiFileName {
		// 修改 xxx.ui 布局文件名
		oldUIFilePath := filepath.Join(projBean.LayoutsPath(), uiForm.UIFile)
		newUIFilePath := filepath.Join(projBean.LayoutsPath(), uiFileName)
		if err := os.Rename(oldUIFilePath, newUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return false
		}
		uiForm.UIFile = uiFileName

		// 修改 xxx.ui.go 布局代码文件名
		oldGoUIFilePath := filepath.Join(projBean.CodePath(), uiForm.GOFile)
		newGoUIFilePath := filepath.Join(projBean.CodePath(), tempFormTab.GOFile())
		if err := os.Rename(oldGoUIFilePath, newGoUIFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return false
		}
		uiForm.GOFile = tempFormTab.GOFile()

		// 修改 xxx.go 用户代码文件名
		oldGoUIUserFilePath := filepath.Join(projBean.CodePath(), uiForm.GOUserFile)
		newGoUIUserFilePath := filepath.Join(projBean.CodePath(), tempFormTab.GOUserFile())
		if err := os.Rename(oldGoUIUserFilePath, newGoUIUserFilePath); err != nil {
			logs.Error("UI布局文件重命名错误:", err.Error())
			return false
		}
		uiForm.GOUserFile = tempFormTab.GOUserFile()
		return true
	}
	return false
}
