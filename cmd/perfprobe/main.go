// LCL 调用性能探针: 定位设计器大量组件时每次 API 调用 ~30ms 的根因。
// 仅用于本地诊断, 不参与设计器构建。
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/api"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/messages"
)

type ProbeForm struct {
	lcl.TEngForm
	designBox  lcl.IWinControl
	buttons    []lcl.IWinControl
	lb         lcl.ILabel
	designer   lcl.IDesigner
	withHook   bool
	paintCount int32
}

// designerOnIsDesignMsg 与设计器一致的 Go 回调: 所有设计模式控件消息都路由到这里。
func (m *ProbeForm) designerOnIsDesignMsg(sender lcl.IControl, message *types.TLMessage) bool {
	if !sender.ComponentState().In(types.CsDesigning) {
		return false
	}
	switch message.Msg {
	case messages.LM_PAINT:
		m.paintCount++
	case messages.LM_MOUSEMOVE:
		m.paintCount++
	case messages.LM_SETCURSOR:
		m.paintCount++
	}
	return false
}

func elapsed(start time.Time) string {
	d := time.Since(start)
	return fmt.Sprintf("%8.2fms", float64(d.Microseconds())/1000.0)
}

// 按设计器 GetDesignerComponent 相同方式创建一个按钮
// ownerInstance: 组件 Owner (应用中是设计窗体, 决定设计消息路由)
func (m *ProbeForm) createButton(ownerInstance uintptr, designMode bool, name string) lcl.IWinControl {
	instance := api.NewInstanceByComponentClass(lcl.TButtonClass())
	if designMode {
		api.SetComponentDesignMode(instance, true)
		api.SetComponentInlineMode(instance, true)
		api.SetComponentDesignInstanceMode(instance, true)
		api.SetWidgetSetDesigning(instance)
	}
	api.CreateObjectByComponent(instance, ownerInstance)
	comp := lcl.AsWinControl(instance)
	comp.SetControlStyle(comp.ControlStyle().Include(types.CsOwnedChildrenNotSelectable, types.CsNoFocus))
	comp.SetName(name)
	return comp
}

func (m *ProbeForm) attachDesignerTo(form lcl.IForm) {
	nd := lcl.NewDesigner()
	nd.SetOnIsDesignMsg(m.designerOnIsDesignMsg)
	m.designer = nd
	form.SetDesigner(m.designer)
	m.withHook = true
}

func (m *ProbeForm) bench(name string, fn func()) {
	fn() // 预热
	var total time.Duration
	for i := 0; i < 10; i++ {
		t := time.Now()
		fn()
		total += time.Since(t)
	}
	fmt.Printf("%-34s avg=%v\n", name, time.Duration(total/10))
}

func (m *ProbeForm) benchOne(name string, fn func()) {
	t := time.Now()
	fn()
	fmt.Printf("%-34s once=%v\n", name, time.Since(t))
}

func (m *ProbeForm) benchButton(prefix string, comp lcl.IWinControl) {
	m.bench(prefix+" SetLeft", func() { comp.SetLeft(comp.Left() + 1) })
	m.bench(prefix+" SetTop", func() { comp.SetTop(comp.Top() + 1) })
	m.bench(prefix+" SetHeight", func() { comp.SetHeight(comp.Height() + 1) })
	m.bench(prefix+" SetWidth", func() { comp.SetWidth(comp.Width() + 1) })
	m.bench(prefix+" SetBounds", func() { comp.SetBounds(20, 20, 60, 30) })
	m.bench(prefix+" SetCaption", func() { comp.SetCaption("X") })
	m.bench(prefix+" SetName", func() { comp.SetName("X") })
	m.bench(prefix+" SetVisible(true)", func() { comp.SetVisible(true) })
	m.bench(prefix+" BringToFront", func() { comp.BringToFront() })
	m.bench(prefix+" BoundsRect读", func() { _ = comp.BoundsRect() })
}

// 设计画布结构体: 复刻设计器 NewFormDesigner 的嵌套结构
type ProbeCanvas struct {
	*lcl.TEngForm
}

func (m *ProbeForm) FormCreate(sender lcl.IObject) {
	m.SetWidth(2004)
	m.SetHeight(1303)

	// 复刻 formDesignScroll + 嵌套设计窗体 (同 NewFormDesigner: 包装 TEngForm)
	scroll := lcl.NewScrollBox(m)
	scroll.SetParent(m)
	scroll.SetAlign(types.AlClient)

	canvas := api.NewInstanceByComponentClass(lcl.TEngFormClass())
	api.CreateObjectByComponent(canvas, 0)
	desForm := &ProbeCanvas{TEngForm: lcl.AsEngForm(canvas).(*lcl.TEngForm)}
	desForm.SetWidth(2004)
	desForm.SetHeight(1303)
	desForm.SetAlign(types.AlCustom)
	desForm.SetBorderStyleToFormBorderStyle(types.BsNone)
	desForm.SetParent(scroll)
	desForm.Show()
	// 设计器挂到嵌套设计窗体上 (同应用: designerForm.SetDesigner(...))
	m.attachDesignerTo(desForm)

	// 设计画布面板: 设计模式 + 双缓冲 (同 NewFormDesigner)
	boxInstance := api.NewInstanceByComponentClass(lcl.TCustomPanelClass())
	api.SetComponentDesignMode(boxInstance, true)
	api.SetComponentInlineMode(boxInstance, true)
	api.SetComponentDesignInstanceMode(boxInstance, true)
	api.SetWidgetSetDesigning(boxInstance)
	api.CreateObjectByComponent(boxInstance, canvas)
	box := lcl.AsCustomPanel(boxInstance)
	box.SetAlign(types.AlClient)
	box.SetBevelOuter(types.BvNone)
	box.SetDoubleBuffered(true)
	box.SetParent(desForm)
	m.designBox = box

	m.lb = lcl.NewLabel(m)
	m.lb.SetParent(m)
	m.lb.SetAlign(types.AlBottom)
	m.lb.SetHeight(60)

	// 先创建全部(与恢复流程一致: 创建时未挂父), 再统一 SetParent
	// Owner 用设计窗体 (与应用一致), 而非画布面板
	m.buttons = make([]lcl.IWinControl, 0, 100)
	t0 := time.Now()
	for i := 0; i < 100; i++ {
		tCreate := time.Now()
		comp := m.createButton(desForm.Instance(), true, fmt.Sprintf("Button%d", i+1))
		tSetParent := time.Now()
		comp.SetParent(box)
		fmt.Printf("创建[%d] newInstance+create=%v setParent=%v\n",
			i, elapsed(tCreate), elapsed(tSetParent))
		m.buttons = append(m.buttons, comp)
	}
	fmt.Printf("== 100 组件创建总耗时: %v\n", time.Since(t0))

	// 位置类属性设置 (模拟恢复流程: 每 Button 设置 Caption/Height/Left/Top/Width)
	t0 = time.Now()
	for i, comp := range m.buttons {
		comp.SetCaption(fmt.Sprintf("Btn%d", i))
		comp.SetHeight(30)
		comp.SetLeft(20 + int32(i%10)*70)
		comp.SetTop(20 + int32(i/10)*45)
		comp.SetWidth(60)
	}
	fmt.Printf("== 100 组件属性设置(5项/个) 总耗时: %v\n", time.Since(t0))

	// 基准: 设计模式按钮
	m.benchButton("[设计模式]", m.buttons[0])

	// 基准: 无设计模式按钮
	norm := m.createButton(desForm.Instance(), false, "NormalButton")
	norm.SetParent(box)
	norm.SetBounds(760, 20, 60, 30)
	m.benchButton("[普通模式]", norm)

	// 隐藏窗体时 100 次批量 SetBounds
	m.benchOne("[隐藏窗体] 100组件SetBounds", func() {
		m.SetVisible(false)
		t := time.Now()
		for _, comp := range m.buttons {
			comp.SetBounds(20, 20, 60, 30)
		}
		fmt.Printf("    (内部) 100次 SetBounds 总耗时=%v\n", time.Since(t))
		m.SetVisible(true)
	})

	// 可见窗体 100 次批量 SetBounds (模拟拖动一步)
	m.benchOne("[可见窗体] 100组件SetBounds", func() {
		t := time.Now()
		for _, comp := range m.buttons {
			comp.SetBounds(20, 20, 60, 30)
		}
		fmt.Printf("    (内部) 100次 SetBounds 总耗时=%v\n", time.Since(t))
	})

	m.lb.SetCaption("PROBE DONE")
	fmt.Printf("design hook 收到消息数: %d\n", m.paintCount)
	fmt.Println("PROBE DONE")
}

func main() {
	libname.LibName = lib.ExtractLibrary(os.TempDir())
	lcl.Init()
	lcl.Application.Initialize()
	lcl.Application.SetTitle("perfprobe")
	var probe ProbeForm
	lcl.Application.NewForms(&probe)
	lcl.Application.Run()
}
