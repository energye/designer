# 本地资源加载

本地资源加载模块允许 WebView 从本地文件系统或嵌入的 `embed.FS` 加载内容，无需启动 HTTP 服务器。

## 高层 API（推荐）

### LocalLoad 配置

```go
wvApp := wv.Init()
wvApp.SetLocalLoad(application.LocalLoad{
    Scheme:     "app",       // URL Scheme
    Domain:     "custom",    // 域名
    ResRootDir: "resources", // 资源目录
    FS:         resources,   // embed.FS（可选）
})
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| Scheme | string | URL Scheme 名称（如 "app"、"fs"），默认 "fs" |
| Domain | string | 域名（如 "custom"、"energy"），默认 "energy" |
| ResRootDir | string | 资源根目录（资源读取目录） |
| FS | emfs.IEmbedFS | 嵌入的文件系统（可选，不设置则从本地目录读取） |
| Home | string | 默认首页，默认 "/index.html" |
| Proxy | IXHRProxy | XHR 代理配置（可选） |

### URL 格式

配置完成后，使用 `{Scheme}://{Domain}/{path}` 格式访问资源：

```go
wvApp.SetOptions(application.Options{
    DefaultURL: "app://custom/index.html", // 加载 resources/index.html
})
```

### 嵌入文件系统

```go
//go:embed resources/*
var resources embed.FS

wvApp.SetLocalLoad(application.LocalLoad{
    Scheme:     "app",
    Domain:     "custom",
    ResRootDir: "resources",
    FS:         resources,
})
```

### 本地文件系统

不设置 `FS` 字段，从本地目录读取：

```go
wvApp.SetLocalLoad(application.LocalLoad{
    Scheme:     "app",
    Domain:     "custom",
    ResRootDir: "./web", // 从 ./web 目录读取
})
```

## 自定义 Scheme

### Windows WebView2

通过 Loader 的 `SetOnGetCustomSchemes` 注册：

```go
loader := wvWindows.NewLoader(nil)
loader.SetOnGetCustomSchemes(func(sender lcl.IObject, customSchemes *wvWindows.IWVCustomSchemeInfoArrayWrap) {
    tmpCustomSchemes := wvWindows.NewCustomSchemeInfoArrayWrapWithInt(1)
    tmpCustomSchemes.SetValue(0, wvWindows.TWVCustomSchemeInfo{
        SchemeName:            "myscheme",
        TreatAsSecure:         1,
        HasAuthorityComponent: 1,
    })
    *customSchemes = tmpCustomSchemes
})

// 拦截资源请求
browser.AddWebResourceRequestedFilter("myscheme*", wvTypes.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)
browser.SetOnWebResourceRequested(func(sender, webView wvWindows.ICoreWebView2, args wvWindows.ICoreWebView2WebResourceRequestedEventArgs) {
    // 处理资源请求
    req := args.Request()
    reqUrl, _ := url.Parse(req.Uri())
    data, _ := resources.ReadFile("assets" + reqUrl.Path)
    stream := lcl.NewBytesStream(data)
    adapter := lcl.NewStreamAdapter(stream, types.SoOwned)
    environment.CreateWebResourceResponse(adapter, 200, "OK", "Content-Type: text/html", &response)
    args.SetResponse(response)
})
```

### Linux WebKit2GTK

通过 `WebContext.RegisterURIScheme` 注册：

```go
webview := wvLinux.NewWebview(m)

wkContext := wvLinux.WebContext.Default()
wkContext.RegisterURIScheme("energy", webview.AsSchemeRequestDelegate())

webview.SetOnURISchemeRequest(func(sender lcl.IObject, wkURISchemeRequest wvTypes.WebKitURISchemeRequest) {
    uriSchemeRequest := wvLinux.NewURISchemeRequest(wkURISchemeRequest)
    path := uriSchemeRequest.Path()
    data, _ := resources.ReadFile(filepath.Join("assets", path))
    ins := wvLinux.InputStream.New(uintptr(unsafe.Pointer(&data[0])), int64(len(data)))
    uriSchemeRequest.Finish(ins.Data(), int64(len(data)), "text/html")
})
```

### macOS WKWebView

通过 `URLSchemeHandler` 注册：

```go
webview := wvDarwin.NewWebview(m)

configuration := wvDarwin.WebViewConfiguration.New()
URLSchemeHandler := wvDarwin.NewURLSchemeHandler(webview.AsWKURLSchemeHandlerDelegate())
configuration.SetURLSchemeHandlerForURLScheme(URLSchemeHandler.Data(), "energy")

webview.SetOnStartURLSchemeTask(func(sender lcl.IObject, webView wvTypes.WKWebView, urlSchemeTask wvTypes.WKURLSchemeTask) {
    // 处理资源请求
})

webview.SetOnStopURLSchemeTask(func(sender lcl.IObject, webView wvTypes.WKWebView, urlSchemeTask wvTypes.WKURLSchemeTask) {
    // 停止资源请求
})
```

## XHR 代理

XHR 代理允许 WebView 中的 XMLHttpRequest 通过 Go 端发起 HTTP 请求。

### 配置

```go
wvApp.SetLocalLoad(application.LocalLoad{
    Scheme:     "app",
    Domain:     "custom",
    ResRootDir: "resources",
    FS:         resources,
    Proxy: &application.XHRProxy{
        Scheme: application.LpsHttp,
        IP:     "127.0.0.1",
        Port:   8080,
    },
})
```

### XHRProxy 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| Scheme | LocalProxyScheme | 代理协议（LpsHttp / LpsHttps） |
| IP | string | 代理 IP 地址，默认 "localhost" |
| Port | int | 代理端口，默认 80 |
| SSL | XHRProxySSL | HTTPS 证书配置 |
| HttpClient | *HttpClient | 自定义 HTTP 客户端 |

### HTTPS 代理

```go
Proxy: &application.XHRProxy{
    Scheme: application.LpsHttps,
    IP:     "127.0.0.1",
    Port:   8443,
    SSL: application.XHRProxySSL{
        RootDir: "certs",           // 证书文件目录
        Cert:    "server.crt",      // 证书文件路径
        Key:     "server.key",      // 私钥文件路径
        CARoots: []string{"ca.crt"}, // CA 根证书文件路径列表
    },
}
```

### XHRProxySSL 字段

| 字段 | 类型 | 说明 |
|------|------|------|
| FS | emfs.IEmbedFS | 证书文件嵌入 FS（可选） |
| RootDir | string | 证书文件根目录 |
| Cert | string | 证书文件路径 |
| Key | string | 私钥文件路径 |
| CARoots | []string | CA 根证书文件路径列表 |

### XHRProxyResponse 响应

| 字段 | 类型 | 说明 |
|------|------|------|
| Data | []byte | 响应数据 |
| DataSize | int | 数据大小 |
| StatusCode | int32 | HTTP 状态码 |
| Status | string | 状态文本 |
| Header | map[string][]string | 响应头 |

## JavaScript 端使用

### 加载本地资源

```html
<!-- 直接引用本地资源 -->
<script src="app://custom/scripts/app.js"></script>
<link rel="stylesheet" href="app://custom/styles/main.css">
<img src="app://custom/images/logo.png">
```

### 发起代理请求

```javascript
// 请求会通过 Go 端的 XHRProxy 代理
fetch("http://127.0.0.1:8080/api/data")
    .then(response => response.json())
    .then(data => {
        console.log(data);
    });
```

## 完整示例

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
        Caption:    "Local Load Demo",
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

### 资源目录结构

```
project/
├── resources/
│   ├── index.html
│   ├── app.js
│   ├── style.css
│   └── images/
│       └── logo.png
├── main.go
└── go.mod
```
