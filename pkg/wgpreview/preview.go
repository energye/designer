// Package wgpreview 让 designer 的一个页面用新渲染引擎（github.com/energye/lcl/wg）绘制。
//
// 它只做一件事：给定一个 lcl.IEngForm，创建一个 wg.Container，
// 挂载指定的场景组件，把窗体的绘制/输入/尺寸消息接进引擎循环。
// designer 原有的 LCL 原生页面不受影响——这是**增量接入**，
// 不替换设计器主画布（那会牵动 12k 行既有代码）。
//
// 用法（在某个菜单动作里）：
//
//	f := lcl.NewEngForm(nil)
//	f.SetCaption("wg 预览")
//	p := wgpreview.New(f, wgpreview.Showcase)
//	f.Show()
//	defer p.Close()
package wgpreview

import (
	"sync"

	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/wg"
	"github.com/energye/lcl/wg/migration"
	"github.com/energye/lcl/wg/ui"
)

// Scene 是可选场景。
type Scene uint8

const (
	// Showcase 是展示台（9 卡片 + 5 类图表）。
	Showcase Scene = iota
	// IME 是输入法演示页。
	IME
	// DBModeler 是数据库建模页。
	DBModeler
)

// Preview 是一个已接入引擎的预览页控制器。
type Preview struct {
	mu   sync.Mutex
	form lcl.IEngForm
	c    *wg.Container
}

// New 在 form 上创建预览页并启动引擎循环。
func New(form lcl.IEngForm, scene Scene) (*Preview, error) {
	if form == nil {
		return nil, wg.ErrNilForm
	}
	c := wg.NewContainer(form, wg.DefaultOptions())
	c.MountWidget(widgetFor(scene))
	if err := c.Run(); err != nil {
		return nil, err
	}
	return &Preview{form: form, c: c}, nil
}

func widgetFor(s Scene) ui.Widget {
	switch s {
	case IME:
		return migration.NewCanvasIME()
	case DBModeler:
		return migration.NewDBModeler()
	default:
		return migration.NewShowcase()
	}
}

// Container 返回底层容器（诊断/自动化用）。
func (p *Preview) Container() *wg.Container { return p.c }

// Close 停止引擎循环并释放资源（幂等）。
func (p *Preview) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.c != nil {
		p.c.Close()
		p.c = nil
	}
}
