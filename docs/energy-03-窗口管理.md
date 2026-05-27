# 窗口管理

窗口管理基于 `window.TWindow` 基类，提供窗口生命周期、状态控制、事件处理等功能。

## TWindow 基类

`window.TWindow` 继承自 LCL 的 `IForm`，是 ENERGY 框架的窗口基类。

```go
type TMainForm struct {
    window.TWindow
    // 自定义字段
}

var MainForm TMainForm
```

## 窗口生命周期

```
FormCreate → OnShow → [运行中] → OnCloseQuery → OnClose
```

### FormCreate - 窗口创建

窗口初始化时调用，用于创建组件和设置初始状态：

```go
func (m *TMainForm) FormCreate(sender lcl.IObject) {
    m.InternalBeforeFormCreate()  // 必须调用：初始化内部状态

    m.SetCaption("My Window")
    m.SetPosition(types.PoScreenCenter)
    m.SetWidth(1200)
    m.SetHeight(800)

    // 创建 WebView
    m.Webview1 = wv.NewWebview(m)
    m.Webview1.SetParent(m)
    m.Webview1.SetAlign(types.AlClient)
    m.Webview1.SetWindow(m)

    m.TWindow.FormCreate(sender)  // 必须调用：触发父类初始化
}
```

**重要**：`FormCreate` 中必须调用 `m.InternalBeforeFormCreate()` 和 `m.TWindow.FormCreate(sender)`。

### OnShow - 窗口显示

窗口显示时调用，适合创建浏览器和居中窗口：

```go
func (m *TMainForm) OnShow(sender lcl.IObject) {
    m.WorkAreaCenter()        // 窗口居中
    m.Webview1.CreateBrowser() // 创建浏览器
}
```

### OnCloseQuery - 关闭询问

用户请求关闭时调用，可阻止关闭：

```go
func (m *TMainForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
    // 设置 *canClose = false 可阻止关闭
    *canClose = true
}
```

### OnClose - 窗口关闭

窗口关闭时调用：

```go
func (m *TMainForm) OnClose(sender lcl.IObject, action *types.TCloseAction) {
    *action = types.CaFree // 释放窗口
}
```

## 窗口状态控制

```go
m.Minimize()        // 最小化
m.Maximize()        // 最大化
m.FullScreen()      // 全屏
m.ExitFullScreen()  // 退出全屏
m.Restore()         // 恢复
m.Close()           // 关闭
```

## 窗口状态查询

```go
m.IsMinimize()    // 是否最小化
m.IsMaximize()    // 是否最大化
m.IsFullScreen()  // 是否全屏
```

## 窗口属性设置

```go
m.SetCaption("标题")              // 设置标题
m.SetWidth(800)                   // 设置宽度
m.SetHeight(600)                  // 设置高度
m.SetLeft(100)                    // 设置左边距
m.SetTop(100)                     // 设置上边距
m.SetPosition(types.PoScreenCenter) // 居中显示
m.SetBorderStyleToFormBorderStyle(types.BsSingle)  // 边框样式
m.SetVisible(true)                // 可见性
m.SetEnabled(true)                // 启用状态
```

## 窗口事件

```go
m.SetOnClick(func(sender lcl.IObject) { /* 点击 */ })
m.SetOnResize(func(sender lcl.IObject) { /* 大小改变 */ })
m.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) { /* 鼠标按下 */ })
m.SetOnMouseMove(func(sender lcl.IObject, shift types.TShiftState, x, y int32) { /* 鼠标移动 */ })
m.SetOnKeyUp(func(sender lcl.IObject, key *uint16, shift types.TShiftState) { /* 键盘抬起 */ })
m.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) { /* 键盘按下 */ })
```

## 线程安全

在非主线程中操作窗口，必须使用 `lcl.RunOnMainThreadAsync`：

```go
ipc.On("close", func(context ipc.IContext) {
    tempWindow := getWindow(context.BrowserId())
    if tempWindow != nil {
        lcl.RunOnMainThreadAsync(func(id uint32) {
            tempWindow.Close()
        })
    }
})
```

## 多窗口管理

```go
var app.Forms []interface{} // 窗体列表

// 通过 BrowserId 查找窗口
getWindow := func(browserId uint32) *TMainForm {
    for _, form := range app.Forms {
        if tempWindow, ok := form.(*TMainForm); ok && tempWindow.BrowserId() == browserId {
            return tempWindow
        }
    }
    return nil
}
```

## 完整示例

```go
package main

import (
    "github.com/energye/energy/v3/application"
    "github.com/energye/energy/v3/window"
    "github.com/energye/energy/v3/wv"
    "github.com/energye/lcl/lcl"
    "github.com/energye/lcl/types"
)

type TMainForm struct {
    window.TWindow
    Webview1 wv.IWebview
}

var MainForm TMainForm

func main() {
    wvApp := wv.Init()
    wvApp.SetOptions(application.Options{
        Caption: "Window Demo",
        Width:   1024,
        Height:  768,
    })
    wv.Run(&MainForm)
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
    m.InternalBeforeFormCreate()
    m.SetCaption("Window Demo")
    m.SetPosition(types.PoScreenCenter)

    m.Webview1 = wv.NewWebview(m)
    m.Webview1.SetParent(m)
    m.Webview1.SetAlign(types.AlClient)
    m.Webview1.SetWindow(m)

    m.TWindow.FormCreate(sender)
}

func (m *TMainForm) OnShow(sender lcl.IObject) {
    m.WorkAreaCenter()
    m.Webview1.CreateBrowser()
}

func (m *TMainForm) OnCloseQuery(sender lcl.IObject, canClose *bool) {
    *canClose = true
}
```
