# IPC 通信

IPC（Inter-Process Communication）模块提供 Go 与 JavaScript 之间的双向通信，支持事件注册、消息发送、自动方法绑定。

## Go 端 API

### On - 注册事件监听

```go
ipc.On("eventName", func(context ipc.IContext) {
    data := context.Data()          // 获取参数（any）
    browserId := context.BrowserId() // 获取浏览器 ID（uint32）
    context.Result(result)           // 返回结果给 JS（支持多个参数）
})
```

### Emit - 向 JS 发送事件

```go
ipc.Emit("eventName", arg1, arg2, ...)
```

### EmitOptions - 带选项发送

```go
ipc.EmitOptions(&ipc.OptionsEvent{
    BrowserId: browserId, // uint32
    Name:      "eventName",
    Data:      "data",    // any
    Callback:  nil,       // EventCallback
})
```

### RemoveOn - 移除事件监听

```go
ipc.RemoveOn("eventName")
```

### BindEvent - 自动绑定结构体方法

将结构体的所有导出方法自动绑定到 IPC 事件，事件名为 `{类型名}.{方法名}`：

```go
type DemoBind struct {
    Field1 string
    Field2 string
    Field3 int
}

func (m *DemoBind) Test1(data string, data2 int, browserId int) {
    fmt.Println("Test1 called:", data, data2, browserId)
}

func (m DemoBind) Test2(datas string, datai int, dataf32 float32, datab bool) string {
    return "Test2 result"
}

func (m *DemoBind) TestResult() *DemoBind {
    return &DemoBind{Field1: "a", Field2: "b", Field3: 1}
}

// 绑定，事件名为 DemoBind.Test1、DemoBind.Test2、DemoBind.TestResult
ipc.BindEvent(&DemoBind{})
```

JS 端调用：
```javascript
ipc.emit('DemoBind.Test1', ['datastr', 1009, energy.getEnv("browser").id], function(result) {
    console.log(result);
});
```

### BindEventPrefix - 带前缀绑定

```go
ipc.BindEventPrefix("demo", &DemoBind{})
```

绑定后事件名为 `demo.Test1`、`demo.Test2`、`demo.TestResult`。

## IContext 接口

```go
type IContext interface {
    Data() any               // 获取事件参数
    Result(data ...any)      // 返回结果给 JS（支持多个返回值）
    BrowserId() uint32       // 触发事件的浏览器 ID
}
```

### 参数类型

JS 传递的参数会被自动转换：

```go
ipc.On("calculate", func(context ipc.IContext) {
    data := context.Data().([]interface{})
    params := data[0].(map[string]interface{})
    a := params["a"].(float64)
    b := params["b"].(float64)
    operator := params["operator"].(string)

    var result float64
    switch operator {
    case "+":
        result = a + b
    case "-":
        result = a - b
    }

    context.Result(map[string]interface{}{
        "result": result,
    })
})
```

### 多个返回值

```go
ipc.On("test", func(context ipc.IContext) {
    context.Result("ResultData", 123, 888.99, true, time.Now().String())
})
```

## JavaScript 端 API

### 发送事件到 Go

```javascript
// 基本发送，参数为数组，最后一个参数为回调函数
ipc.emit('eventName', [arg1, arg2], function(result) {
    console.log(result);
});

// 无参数
ipc.emit('eventName', [], function(result) {
    console.log(result);
});

// 带参数
ipc.emit('calculate', [{a: 10, b: 20, operator: '+'}], function(result) {
    console.log(result.result); // 30
});
```

### 监听 Go 事件

```javascript
ipc.on('eventName', function(arg1, arg2, arg3) {
    console.log(arg1, arg2, arg3);
});
```

### 获取环境信息

```javascript
energy.getEnv("browser").id  // 获取当前浏览器 ID
```

## IPC 消息流程

```
JS: ipc.emit('name', [data], callback)
    ↓
Go: ipc.On('name', func(ctx) { ctx.Result(...) })
    ↓
JS: callback(result)
```

```
Go: ipc.Emit('name', data)
    ↓
JS: ipc.on('name', function(data) { ... })
```

## 平台差异

IPC 底层传输机制因平台而异：

| 平台 | JS 端 API | 说明 |
|------|-----------|------|
| Windows | `window.chrome.webview.postMessage()` | WebView2 原生 API |
| Linux | `window.webkit.messageHandlers.external.postMessage()` | WebKit2GTK 原生 API |
| macOS | `window.webkit.messageHandlers.external.postMessage()` | WKWebView 原生 API |

框架自动处理这些差异，用户无需关心。

## 完整示例

### Go 端

```go
package main

import (
    "fmt"
    "runtime"
    "github.com/energye/energy/v3/application"
    "github.com/energye/energy/v3/ipc"
    "github.com/energye/energy/v3/wv"
    "github.com/energye/energy/v3/window"
    "github.com/energye/lcl/lcl"
)

type TMainForm struct {
    window.TWindow
}

var MainForm TMainForm

func main() {
    wvApp := wv.Init()
    wvApp.SetOptions(application.Options{
        Caption: "IPC Demo",
        Width:   800,
        Height:  600,
    })

    // 注册事件
    ipc.On("get-system-info", func(context ipc.IContext) {
        context.Result(map[string]interface{}{
            "platform": runtime.GOOS,
            "arch":     runtime.GOARCH,
        })
    })

    ipc.On("calculate", func(context ipc.IContext) {
        data := context.Data().([]interface{})
        params := data[0].(map[string]interface{})
        a := params["a"].(float64)
        b := params["b"].(float64)
        context.Result(map[string]interface{}{
            "result": a + b,
        })
    })

    // 自动绑定
    ipc.BindEvent(&MyHandler{})

    wv.Run(&MainForm)
}

type MyHandler struct{}

func (h *MyHandler) Save(data string) {
    fmt.Println("Save:", data)
}
```

### JavaScript 端

```javascript
// 调用 Go 方法
ipc.emit('get-system-info', [], function(result) {
    document.getElementById('platform').innerText = result.platform;
});

ipc.emit('calculate', [{a: 10, b: 20}], function(result) {
    console.log('计算结果:', result.result); // 30
});

// 监听 Go 推送
ipc.on('update', function(data) {
    console.log('收到更新:', data);
});

// 调用绑定方法
ipc.emit('MyHandler.Save', ['hello'], function(result) {
    console.log('保存成功');
});
```
