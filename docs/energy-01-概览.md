# ENERGY 框架概览

ENERGY 是基于 Go 语言的跨平台桌面 GUI 框架，整合原生控件（LCL）和 WebView 引擎（Windows WebView2 / Linux WebKit2GTK / macOS WKWebView），一套代码运行在 Windows、macOS、Linux 上。

## 包结构

```
github.com/energye/energy/v3      # ENERGY 框架高层封装
├── application/                   # 应用配置、本地资源加载、系统托盘、通知
├── window/                        # 窗口基类 TWindow
├── wv/                            # 跨平台 WebView 抽象层
├── ipc/                           # Go ↔ JavaScript 进程间通信
├── logger/                        # 异步日志
└── platform/                      # 平台特定功能（通知等）

github.com/energye/lcl            # LCL 原生控件库（底层）
github.com/energye/wv             # 平台 WebView 绑定（底层）
```

## 核心包导入

```go
import (
    "github.com/energye/energy/v3/application" // 应用配置
    "github.com/energye/energy/v3/ipc"         // IPC 通信
    "github.com/energye/energy/v3/wv"          // WebView
    "github.com/energye/energy/v3/window"      // 窗口基类
    "github.com/energye/lcl/lcl"               // LCL 核心
    "github.com/energye/lcl/types"             // UI 类型常量
    "github.com/energye/lcl/types/colors"      // 颜色常量
    "github.com/energye/lcl/api"               // API 工具
    "github.com/energye/lcl/api/libname"       // 平台库配置
)
```

## 两种初始化模式

### 模式一：高层 API（推荐）

使用 `wv.Init()` + `wv.Run()` 一行代码启动：

```go
package main

import (
    "embed"
    "github.com/energye/energy/v3/application"
    "github.com/energye/energy/v3/ipc"
    "github.com/energye/energy/v3/wv"
    "github.com/energye/energy/v3/window"
    "github.com/energye/lcl/lcl"
    "github.com/energye/lcl/types"
)

//go:embed resources/*
var resources embed.FS

type TMainForm struct {
    window.TWindow
    Webview1 wv.IWebview
}

var MainForm TMainForm

func main() {
    wvApp := wv.Init()
    wvApp.SetOptions(application.Options{
        Caption:    "My App",
        DefaultURL: "app://custom/index.html",
        Width:      1200,
        Height:     800,
    })
    wvApp.SetLocalLoad(application.LocalLoad{
        Scheme:     "app",
        Domain:     "custom",
        ResRootDir: "resources",
        FS:         resources,
    })

    ipc.On("hello", func(context ipc.IContext) {
        context.Result("world")
    })

    wv.Run(&MainForm)
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
    m.InternalBeforeFormCreate()
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
```

### 模式二：底层 API（完全控制）

手动初始化 LCL 和 WebView，适合需要精细控制的场景：

```go
package main

import (
    "github.com/energye/lcl/lcl"
    "github.com/energye/lcl/api/libname"
    "github.com/energye/wv"
)

func main() {
    libname.UseWS = "gtk3" // Linux 必须
    lcl.Init()
    wv.Init()

    lcl.Application.Initialize()
    lcl.Application.SetMainFormOnTaskBar(true)
    lcl.Application.NewForm(&mainForm)
    lcl.Application.Run()
}
```

## 平台要求

| 平台 | 要求 |
|------|------|
| Windows | WebView2 运行时（Edge 自带或单独安装） |
| Linux | WebKit2GTK 4.0 或 4.1，GTK3 |
| macOS | 系统内置 WebKit2 |

## 文档索引

| 文档 | 内容 |
|------|------|
| [应用配置](energy-02-应用配置.md) | Options 全部字段、平台配置 |
| [窗口管理](energy-03-窗口管理.md) | TWindow、窗口状态、事件 |
| [WebView](energy-04-WebView.md) | WebView 创建、配置、事件、导航 |
| [IPC 通信](energy-05-IPC通信.md) | Go ↔ JS 双向通信 |
| [本地资源加载](energy-06-本地资源加载.md) | LocalLoad、自定义 Scheme、XHRProxy |
| [日志](energy-07-日志.md) | 异步日志系统 |
| [自定义组件](energy-08-自定义组件.md) | TButton、TTab/TPage |
| [系统托盘](energy-09-系统托盘.md) | 托盘图标、菜单 |
| [通知](energy-10-通知.md) | 跨平台通知 |
| [菜单编辑](energy-11-菜单编辑.md) | 标准编辑操作 |
