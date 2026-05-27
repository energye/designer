# 通知

通知模块提供跨平台的系统通知功能，支持自定义分类和操作按钮。

## 创建通知实例

```go
import "github.com/energye/energy/v3/platform/notification"

notify := notification.New()
```

## 请求权限

```go
// macOS 需要请求通知权限
notify.RequestNotificationAuthorization()

// 检查权限状态
notify.CheckNotificationAuthorization()
```

## 发送通知

### 基本通知

```go
notify.SendNotification(notification.Options{
    ID:    "unique-id",
    Title: "通知标题",
    Body:  "通知内容",
    Sound: true,
    Data: map[string]interface{}{
        "key": "value",
    },
})
```

### 带操作的通知

```go
// 注册分类
category := notification.Category{
    ID: "message",
    Actions: []notification.Action{
        {ID: "reply", Title: "回复"},
        {ID: "dismiss", Title: "忽略"},
    },
    HasReplyField: true,
}
notify.RegisterNotificationCategory(category)

// 发送带操作的通知
notify.SendNotificationWithActions(notification.Options{
    ID:         "msg-001",
    Title:      "新消息",
    Body:       "您有一条新消息",
    CategoryID: "message",
    Sound:      true,
})
```

## 处理通知响应

```go
notify.SetOnNotificationResponse(func(result notification.Result) {
    actionID := result.ActionID
    switch actionID {
    case "reply":
        // 用户点击了回复
    case "dismiss":
        // 用户点击了忽略
    }
})
```

## 清理通知

```go
notify.RemoveAllDeliveredNotifications()  // 移除所有已送达通知
notify.RemoveAllPendingNotifications()    // 移除所有待发送通知
```

## Options 配置

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 通知唯一标识 |
| Title | string | 通知标题 |
| Body | string | 通知内容 |
| Sound | bool | 是否播放声音 |
| CategoryID | string | 通知分类 ID |
| Data | map[string]interface{} | 自定义数据 |

## Category 分类

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 分类唯一标识 |
| Actions | []Action | 操作按钮列表 |
| HasReplyField | bool | 是否显示回复输入框 |

## Action 操作

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | string | 操作唯一标识 |
| Title | string | 按钮文本 |

## 完整示例

```go
package main

import (
    "github.com/energye/energy/v3/platform/notification"
)

func main() {
    notify := notification.New()

    // 请求权限（macOS）
    notify.RequestNotificationAuthorization()

    // 注册分类
    category := notification.Category{
        ID: "chat",
        Actions: []notification.Action{
            {ID: "reply", Title: "回复"},
            {ID: "mark_read", Title: "标为已读"},
        },
        HasReplyField: true,
    }
    notify.RegisterNotificationCategory(category)

    // 处理响应
    notify.SetOnNotificationResponse(func(result notification.Result) {
        switch result.ActionID {
        case "reply":
            openChatWindow()
        case "mark_read":
            markAsRead()
        }
    })

    // 发送通知
    notify.SendNotificationWithActions(notification.Options{
        ID:         "chat-001",
        Title:      "新消息",
        Body:       "张三: 你好！",
        CategoryID: "chat",
        Sound:      true,
        Data: map[string]interface{}{
            "sender": "张三",
        },
    })
}
```
