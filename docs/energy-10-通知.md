# 通知

通知模块提供跨平台的系统通知功能，支持自定义分类和操作按钮。

## 创建通知实例

```go
import "github.com/energye/energy/v3/platform/notification"

notify := notification.New() // 返回 INotification 接口
```

### 接口定义（platform/notification/types 包）

```go
type INotification interface {
    Initialize() error
    RequestNotificationAuthorization() (bool, error)
    CheckNotificationAuthorization() (bool, error)
    SendNotification(options Options) error
    SendNotificationWithActions(options Options) error
    RegisterNotificationCategory(category Category) error
    RemoveNotificationCategory(categoryID string) error
    RemoveAllPendingNotifications() error
    RemovePendingNotification(identifier string) error
    RemoveAllDeliveredNotifications() error
    RemoveDeliveredNotification(identifier string) error
    RemoveNotification(identifier string) error
    SetOnNotificationResponse(callback TNotificationResponseEvent)
}
```

## 请求权限

```go
// 请求通知权限（macOS 需要）
authorized, err := notify.RequestNotificationAuthorization()

// 检查权限状态
authorized, err := notify.CheckNotificationAuthorization()
```

## 发送通知

### 基本通知

```go
err := notify.SendNotification(types.Options{
    ID:    "unique-id",
    Title: "通知标题",
    Body:  "通知内容",
    Data: map[string]interface{}{
        "key": "value",
    },
})
```

### 带操作的通知

```go
// 注册分类
category := types.Category{
    ID: "message",
    Actions: []types.Action{
        {ID: "reply", Title: "回复"},
        {ID: "dismiss", Title: "忽略", Destructive: true},
    },
    HasReplyField:    true,
    ReplyPlaceholder: "输入回复...",
    ReplyButtonTitle: "发送",
}
err := notify.RegisterNotificationCategory(category)

// 发送带操作的通知
err = notify.SendNotificationWithActions(types.Options{
    ID:         "msg-001",
    Title:      "新消息",
    Body:       "您有一条新消息",
    CategoryID: "message",
    Data: map[string]interface{}{
        "sender": "张三",
    },
})
```

## 处理通知响应

```go
notify.SetOnNotificationResponse(func(result types.Result) {
    if result.Error != nil {
        fmt.Println("通知错误:", result.Error)
        return
    }

    response := result.Response
    fmt.Println("通知 ID:", response.ID)
    fmt.Println("操作:", response.ActionIdentifier)
    fmt.Println("分类:", response.CategoryID)
    fmt.Println("用户输入:", response.UserText)
    fmt.Println("自定义数据:", response.UserInfo)
})
```

## 清理通知

```go
err := notify.RemoveAllDeliveredNotifications()  // 移除所有已送达通知
err := notify.RemoveDeliveredNotification("id")  // 移除指定已送达通知
err := notify.RemoveAllPendingNotifications()     // 移除所有待发送通知
err := notify.RemovePendingNotification("id")     // 移除指定待发送通知
err := notify.RemoveNotification("id")            // 移除通知
```

## Options 配置

```go
type Options struct {
    ID         string                 // 通知唯一标识
    Title      string                 // 通知标题
    Subtitle   string                 // 通知副标题（macOS）
    Body       string                 // 通知内容
    CategoryID string                 // 通知分类 ID
    Data       map[string]interface{} // 自定义数据
}
```

## Category 分类

```go
type Category struct {
    ID               string   // 分类唯一标识
    Actions          []Action // 操作按钮列表
    HasReplyField    bool     // 是否显示回复输入框
    ReplyPlaceholder string   // 回复输入框占位文本
    ReplyButtonTitle string   // 回复按钮文本
}
```

## Action 操作

```go
type Action struct {
    ID          string // 操作唯一标识
    Title       string // 按钮文本
    Destructive bool   // 是否为破坏性操作（显示为红色）
}
```

## Response 响应

```go
type Response struct {
    ID               string                 // 通知 ID
    ActionIdentifier string                 // 用户点击的操作 ID
    CategoryID       string                 // 通知分类 ID
    Title            string                 // 通知标题
    Subtitle         string                 // 通知副标题
    Body             string                 // 通知内容
    UserText         string                 // 用户输入的文本
    UserInfo         map[string]interface{} // 自定义数据
}
```

## Result 结果

```go
type Result struct {
    Response Response // 通知响应
    Error    error    // 错误信息
}
```

## 常量

```go
const DefaultActionIdentifier = "DEFAULT_ACTION" // 默认操作标识
```

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/energye/energy/v3/platform/notification/types"
)

func setupNotification(notify types.INotification) {
    // 初始化
    notify.Initialize()

    // 请求权限
    authorized, _ := notify.RequestNotificationAuthorization()
    if !authorized {
        fmt.Println("通知权限未授权")
        return
    }

    // 注册分类
    category := types.Category{
        ID: "chat",
        Actions: []types.Action{
            {ID: "reply", Title: "回复"},
            {ID: "mark_read", Title: "标为已读"},
        },
        HasReplyField:    true,
        ReplyPlaceholder: "输入回复...",
    }
    notify.RegisterNotificationCategory(category)

    // 处理响应
    notify.SetOnNotificationResponse(func(result types.Result) {
        if result.Error != nil {
            fmt.Println("错误:", result.Error)
            return
        }
        switch result.Response.ActionIdentifier {
        case "reply":
            fmt.Println("用户回复:", result.Response.UserText)
        case "mark_read":
            fmt.Println("标为已读")
        }
    })

    // 发送通知
    notify.SendNotificationWithActions(types.Options{
        ID:         "chat-001",
        Title:      "新消息",
        Body:       "张三: 你好！",
        CategoryID: "chat",
        Data: map[string]interface{}{
            "sender": "张三",
        },
    })
}
```
