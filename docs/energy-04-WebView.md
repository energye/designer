# WebView

WebView 模块提供跨平台 WebView 抽象，封装 Windows WebView2 / Linux WebKit2GTK / macOS WKWebView。

## 跨平台 WebView（推荐）

使用 `wv.IWebview` 接口，一套代码跨平台运行。

### 创建 WebView

```go
type TMainForm struct {
    window.TWindow
    Webview1 wv.IWebview
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
    m.InternalBeforeFormCreate()

    m.Webview1 = wv.NewWebview(m)     // 创建 WebView
    m.Webview1.SetParent(m)           // 设置父容器
    m.Webview1.SetAlign(types.AlClient) // 填满父容器
    m.Webview1.SetWindow(m)           // 绑定窗口

    m.TWindow.FormCreate(sender)
}

func (m *TMainForm) OnShow(sender lcl.IObject) {
    m.WorkAreaCenter()
    m.Webview1.CreateBrowser() // 创建浏览器实例
}
```

### IWebview 方法

#### 导航

```go
m.Webview1.LoadURL("https://example.com") // 加载 URL
```

#### JavaScript 执行

```go
// 同步执行
m.Webview1.ExecuteScript("document.title")

// 异步执行带回调
m.Webview1.ExecuteScriptCallback("document.title", func(result, err string) {
    if err == "" {
        fmt.Println("Title:", result)
    }
})
```

#### 浏览器信息

```go
m.Webview1.BrowserId() // 获取浏览器 ID
```

## 事件回调

### 加载状态变更

```go
m.Webview1.SetOnLoadChange(func(url, title string, load wv.TLoadChange) {
    switch load {
    case wv.LcStart:
        fmt.Println("开始加载:", url)
    case wv.LcLoading:
        fmt.Println("加载中:", url)
    case wv.LcFinish:
        fmt.Println("加载完成:", title)
    }
})
```

### 右键菜单

```go
m.Webview1.SetOnContextMenu(func(contextMenu *wv.TContextMenuItem) {
    contextMenu.Add("刷新页面", wv.CmkCommand)
    contextMenu.Add("", wv.CmkSeparator)
    contextMenu.Add("开发者工具", wv.CmkCommand)
})

m.Webview1.SetOnContextMenuCommand(func(commandId int32) {
    fmt.Println("菜单命令:", commandId)
})
```

### 弹出窗口

```go
m.Webview1.SetOnPopupWindow(func(targetURL string) bool {
    fmt.Println("弹出窗口:", targetURL)
    return true // 返回 true 拦截弹出
})
```

### 拖拽

```go
m.Webview1.SetOnDragEnter(func(dragType wv.TDragType, x, y int32) {
    fmt.Println("拖拽进入:", dragType)
})

m.Webview1.SetOnDragOver(func(data *wv.TDragData, x, y int32) {
    fmt.Println("拖拽悬停:", x, y)
})

m.Webview1.SetOnDragLeave(func() {
    fmt.Println("拖拽离开")
})
```

## 拖拽类型

| 类型 | 说明 |
|------|------|
| DragTypeNo | 无拖拽 |
| DragTypeFile | 文件拖拽 |
| DragTypeData | 数据拖拽 |

## TDragData 结构

```go
type TDragData struct {
    Type      TDragType // 拖拽类型
    Data      []byte    // 拖拽数据
    Filenames []string  // 文件名列表
}
```

## 平台底层 API

如需更精细的控制，可直接使用平台特定 API。

### Windows WebView2

```go
import wvWindows "github.com/energye/wv/windows"

// 创建 Loader
loader := wvWindows.NewLoader(nil)
loader.SetLoaderDllPath("WebView2Loader.dll")
loader.StartWebView2()

// 创建浏览器
browser := wvWindows.NewBrowser(m)
windowParent := wvWindows.NewWindowParent(m)
windowParent.SetBrowser(browser)

browser.SetOnAfterCreated(func(sender lcl.IObject) {
    windowParent.UpdateSize()
})

browser.SetOnNavigationStarting(func(sender, webView wvWindows.ICoreWebView2, args wvWindows.ICoreWebView2NavigationStartingEventArgs) {
    // 导航开始
})

browser.SetOnNavigationCompleted(func(sender, webView wvWindows.ICoreWebView2, args wvWindows.ICoreWebView2NavigationCompletedEventArgs) {
    // 导航完成
})

browser.SetOnWebMessageReceived(func(sender, webView wvWindows.ICoreWebView2, args wvWindows.ICoreWebView2WebMessageReceivedEventArgs) {
    message := args.WebMessageAsString()
    // 处理 JS 消息
})

browser.CreateBrowserWithHandleBool(windowParent.Handle(), true)
```

### Linux WebKit2GTK

```go
import wvLinux "github.com/energye/wv/linux"

webview := wvLinux.NewWebview(m)
webviewParent := wvLinux.NewWebviewParent(m)

// 配置设置
setting := wvLinux.NewSettings()
setting.SetEnableDeveloperExtras(true)
setting.SetUserAgentWithApplicationDetails("energy.io", "3.0")
setting.SetHardwareAccelerationPolicy(wvTypes.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
webview.SetSettings(setting)

// 注册 IPC 消息处理
webview.RegisterScriptMessageHandler("processMessage")

// 注册自定义 Scheme
wkContext := wvLinux.WebContext.Default()
wkContext.RegisterURIScheme("energy", webview.AsSchemeRequestDelegate())

// 事件
webview.SetOnLoadChange(func(sender lcl.IObject, loadEvent wvTypes.WebKitLoadEvent) {
    // 加载状态变更
})

webview.SetOnProcessMessage(func(sender lcl.IObject, jsValue wvLinux.IWkJSValue, processId wvTypes.TWkProcessId) {
    value := jsValue.StringValue()
    // 处理 JS 消息
})

webview.SetOnURISchemeRequest(func(sender lcl.IObject, wkURISchemeRequest wvTypes.WebKitURISchemeRequest) {
    // 自定义 Scheme 资源加载
})

webview.CreateBrowser()
webviewParent.SetWebview(webview)
webview.LoadURL("https://example.com")
```

### macOS WKWebView

```go
import wvDarwin "github.com/energye/wv/darwin"

webview := wvDarwin.NewWebview(m)
webviewParent := wvDarwin.NewWebviewParent(m)

// 配置
userContentController := wvDarwin.UserContentController.New()
scriptMessageHandler := wvDarwin.NewScriptMessageHandler(webview.AsReceiveScriptMessageDelegate())
userContentController.AddScriptMessageHandlerName(scriptMessageHandler, "processMessage")

configuration := wvDarwin.WebViewConfiguration.New()
configuration.SetUserContentController(userContentController.Data())

// 注册自定义 Scheme
URLSchemeHandler := wvDarwin.NewURLSchemeHandler(webview.AsWKURLSchemeHandlerDelegate())
configuration.SetURLSchemeHandlerForURLScheme(URLSchemeHandler.Data(), "energy")

// 初始化
frame := types.TRect{0, 0, 800, 600}
webview.InitWithFrameConfiguration(frame, configuration.Data())

// 事件
webview.SetOnProcessMessage(func(sender lcl.IObject, userContentController wvTypes.WKUserContentController, name string, data string) {
    // 处理 JS 消息
})

webview.SetOnFinishNavigation(func(sender lcl.IObject, navigation wvTypes.WKNavigation) {
    // 导航完成
})

webviewParent.SetWebview(webview.Data())
webview.LoadURL("https://example.com")
```

## WebView 设置

### Windows WebView2 设置

```go
settings := browser.Settings()
settings.SetAreDevToolsEnabled(true)
settings.SetAreDefaultContextMenusEnabled(true)
settings.SetIsScriptEnabled(true)
settings.SetIsWebMessageEnabled(true)
settings.SetUserAgent("Custom Agent")
```

### Linux WebKit2GTK 设置

```go
setting := wvLinux.NewSettings()
setting.SetEnableDeveloperExtras(true)
setting.SetEnableJavascript(true)
setting.SetAutoLoadImages(true)
setting.SetUserAgent("Custom Agent")
setting.SetHardwareAccelerationPolicy(wvTypes.WEBKIT_HARDWARE_ACCELERATION_POLICY_ALWAYS)
webview.SetSettings(setting)
```

## 完整示例

```go
package main

import (
    "embed"
    "fmt"
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
        Caption:    "WebView Demo",
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
    m.SetCaption("WebView Demo")
    m.SetPosition(types.PoScreenCenter)

    m.Webview1 = wv.NewWebview(m)
    m.Webview1.SetParent(m)
    m.Webview1.SetAlign(types.AlClient)
    m.Webview1.SetWindow(m)

    m.Webview1.SetOnLoadChange(func(url, title string, load wv.TLoadChange) {
        switch load {
        case wv.LcStart:
            fmt.Println("开始加载:", url)
        case wv.LcFinish:
            fmt.Println("加载完成:", title)
        }
    })

    m.Webview1.SetOnContextMenu(func(contextMenu *wv.TContextMenuItem) {
        contextMenu.Add("刷新", wv.CmkCommand)
    })

    m.TWindow.FormCreate(sender)
}

func (m *TMainForm) OnShow(sender lcl.IObject) {
    m.WorkAreaCenter()
    m.Webview1.CreateBrowser()
}
```
