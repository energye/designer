# 系统托盘

系统托盘模块提供跨平台的系统托盘图标和菜单功能（Windows 和 macOS，Linux 使用系统通知）。

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
tray.Close()    // 关闭托盘图标
```

### 设置图标

```go
tray.SetIcon("path/to/icon.png")  // 从文件设置图标
tray.SetIconBytes(data []byte)    // 从字节数据设置图标
```

### 设置提示文本

```go
tray.SetHint("我的应用 - 运行中")
```

## 鼠标事件

```go
// 左键单击
tray.SetOnClick(func() {
    // 显示主窗口
})

// 左键双击
tray.SetOnDblClick(func() {
    // 显示/隐藏窗口
})

// 鼠标按下
tray.SetOnMouseDown(func(button types.TMouseButton, shift types.TShiftState, x, y int32) {
    // 鼠标按下事件
})

// 鼠标抬起
tray.SetOnMouseUp(func(button types.TMouseButton, shift types.TShiftState, x, y int32) {
    // 鼠标抬起事件
})

// 鼠标移动
tray.SetOnMouseMove(func(shift types.TShiftState, x, y int32) {
    // 鼠标移动事件
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
item.SetOnClick(func() {
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
subItem.SetOnClick(func() {
    // 子菜单点击
})
```

### 菜单项属性

```go
item.SetChecked(true)      // 设置选中状态
item.SetRadio(true)        // 设置为单选样式
item.SetEnabled(false)     // 设置为禁用
item.SetImage(imageName)   // 设置图标（通过图片名）
item.SetBitmap(imageData)  // 设置位图（字节数据）
item.SetImageIndex(index)  // 设置图片索引
item.Enabled() bool        // 查询是否启用
item.Checked() bool        // 查询是否选中
item.Clear()               // 清空子菜单
```

### 图片列表

```go
// 从文件路径设置图片列表
menu.SetImageList([]string{"path/to/icon1.png", "path/to/icon2.png"})

// 从嵌入 FS 设置图片列表
menu.SetImageListEmbed(embedFS, []string{"resources/icon1.png", "resources/icon2.png"})

// 从字节数据设置图片列表
menu.SetImageListDataBytes([][]byte{icon1Data, icon2Data})
```

### 菜单项绘制事件

```go
item.SetOnMeasureItem(fn lcl.TMenuMeasureItemEvent) // 自定义测量
item.SetOnDrawItem(fn lcl.TMenuDrawItemEvent)       // 自定义绘制
```

## 完整示例

```go
package main

import (
    "github.com/energye/energy/v3/application"
    "github.com/energye/lcl/lcl"
    "github.com/energye/lcl/types"
)

type TMainForm struct {
    lcl.TEngForm
    tray *application.TTrayIcon
}

var mainForm TMainForm

func main() {
    lcl.Init()
    lcl.Run(&mainForm)
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
    m.SetCaption("托盘示例")
    m.SetPosition(types.PoScreenCenter)

    // 创建托盘图标
    m.tray = application.NewTrayIcon()
    m.tray.SetHint("我的应用")
    m.tray.SetIconBytes(iconData)

    // 创建菜单
    menu := m.tray.Menu()

    showItem := menu.AddMenuItem("显示窗口")
    showItem.SetOnClick(func() {
        m.Show()
    })

    menu.AddSeparator()

    exitItem := menu.AddMenuItem("退出")
    exitItem.SetOnClick(func() {
        lcl.Application.Terminate()
    })

    // 双击显示窗口
    m.tray.SetOnDblClick(func() {
        m.Show()
    })

    // 关闭时隐藏到托盘
    m.SetOnCloseQuery(func(sender lcl.IObject, canClose *bool) {
        *canClose = false
        m.Hide()
    })

    m.tray.Show()
}
```
