# 系统托盘

系统托盘模块提供跨平台的系统托盘图标和菜单功能。

## 创建托盘图标

```go
import "github.com/energye/energy/v3/application"

tray := application.NewTrayIcon()
```

## 基本操作

### 显示/隐藏

```go
tray.Show()     // 显示托盘图标
tray.Hide()     // 隐藏托盘图标
tray.Visible()  // 查询是否可见
```

### 设置图标

```go
tray.SetIcon(icon)             // 设置图标对象
tray.SetIconBytes(data []byte) // 从字节数据设置图标
```

### 设置提示文本

```go
tray.SetHint("我的应用 - 运行中")
```

## 鼠标事件

```go
// 左键单击
tray.SetOnClick(func(sender lcl.IObject) {
    // 显示主窗口
})

// 左键双击
tray.SetOnDblClick(func(sender lcl.IObject) {
    // 显示/隐藏窗口
})

// 右键单击
tray.SetOnRightClick(func(sender lcl.IObject) {
    // 显示菜单
})
```

## 托盘菜单

### 获取菜单对象

```go
menu := tray.Menu()
```

### 添加菜单项

```go
item := menu.AddMenuItem("显示主窗口")
item.SetOnClick(func(sender lcl.IObject) {
    form.Show()
})
```

### 添加分隔线

```go
menu.AddSeparator()
```

### 子菜单

```go
item := menu.AddMenuItem("选项")
subItem := item.AddSubMenuItem("子选项")
subItem.SetOnClick(func(sender lcl.IObject) {
    // 子菜单点击
})
```

### 菜单项属性

```go
item.SetChecked(true)      // 设置选中状态
item.SetRadio(true)        // 设置为单选样式
item.SetEnabled(false)     // 设置为禁用
item.SetImage(image)       // 设置图标
item.SetBitmap(bitmap)     // 设置位图
```

### 图片列表

```go
menu.SetImageListEmbed(embedFS, []string{"resources/icon.png"})
```

## 完整示例

```go
package main

import (
    "github.com/energye/energy/v3/application"
    "github.com/energye/lcl/lcl"
)

func main() {
    form := lcl.Application.CreateForm()
    form.SetCaption("托盘示例")

    // 创建托盘图标
    tray := application.NewTrayIcon()
    tray.SetHint("我的应用")
    tray.SetIconBytes(iconData)

    // 创建菜单
    menu := tray.Menu()

    showItem := menu.AddMenuItem("显示窗口")
    showItem.SetOnClick(func(sender lcl.IObject) {
        form.Show()
    })

    menu.AddSeparator()

    exitItem := menu.AddMenuItem("退出")
    exitItem.SetOnClick(func(sender lcl.IObject) {
        lcl.Application.Terminate()
    })

    // 双击显示窗口
    tray.SetOnDblClick(func(sender lcl.IObject) {
        form.Show()
    })

    // 关闭时隐藏到托盘
    form.SetOnCloseQuery(func(sender lcl.IObject, canClose *bool) {
        *canClose = false
        form.Hide()
    })

    tray.Show()

    lcl.Run(form)
}
```
