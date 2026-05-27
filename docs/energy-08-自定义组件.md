# 自定义组件

ENERGY 提供自定义绘制的 UI 组件，位于 `lcl/wg` 包下，包括 TButton（自定义按钮）和 TTab/TPage（标签页控件）。

## TButton 自定义按钮

TButton 是完全自定义绘制的按钮组件，支持渐变背景、图标、圆角、多种状态样式。

### 创建按钮

```go
import "github.com/energye/lcl/widget/wg"

button := wg.NewButton(parent)
button.SetCaption("点击我")
button.SetLeft(50)
button.SetTop(50)
button.SetWidth(120)
button.SetHeight(36)
```

### 按钮状态

| 状态 | 说明 | 触发条件 |
|------|------|----------|
| Default | 默认状态 | 正常显示 |
| Enter | 悬停状态 | 鼠标进入 |
| Down | 按下状态 | 鼠标按下 |
| Disabled | 禁用状态 | SetDisable(true) |

### 颜色配置

颜色方法接受两个参数：渐变起始色和渐变结束色。

```go
// 默认状态渐变
button.SetDefaultColor(startColor, endColor)

// 悬停状态渐变
button.SetEnterColor(startColor, endColor)

// 按下状态渐变
button.SetDownColor(startColor, endColor)

// 禁用状态渐变
button.SetDisabledColor(startColor, endColor)

// 设置单色（无渐变）
button.SetColor(color)

// 设置渐变色
button.SetColorGradient(startColor, endColor)
```

### 边框

边框按方向设置：

```go
// 设置指定方向边框颜色
button.SetBorderColor(direction, color)

// 设置指定方向边框宽度
button.SetBorderWidth(direction, width)

// 设置边框方向
button.SetBorderDirections(directions)
```

### 圆角半径

```go
button.SetRadius(8) // 设置圆角半径
```

### 透明度

```go
button.SetAlpha(200) // 0-255，255 为不透明
```

### 图标

```go
button.SetIcon("path/to/icon.png")              // 从文件设置图标
button.SetIconFormBytes(pngData)                  // 从字节数据设置图标
button.SetIconFavorite("path/to/favorite.png")    // 收藏图标
button.SetIconFavoriteFormBytes(pngData)          // 从字节数据设置收藏图标
button.SetIconClose("path/to/close.png")          // 关闭按钮图标
button.SetIconCloseFormBytes(pngData)             // 从字节数据设置关闭图标
button.SetIconCloseHighlight("path/to/highlight.png") // 关闭按钮高亮图标
button.SetIconCloseHighlightFormBytes(pngData)    // 从字节数据设置高亮图标
```

### 文本

```go
button.SetCaption("按钮文本")   // 设置按钮文本
button.Caption() string         // 获取按钮文本
button.SetText("文本")          // 设置文本（同 SetCaption）
button.Text() string            // 获取文本
button.AutoSizeWidth()          // 根据文本自动调整宽度
button.SetAutoSize(true)        // 设置自动大小
```

### 自动大小

```go
button.SetAutoSize(true) // 根据文本自动调整大小
```

### 提示文本

```go
button.ShowHint("这是提示文本")  // 显示提示
button.HideHint()                // 隐藏提示
button.SetCloseHintText("关闭")  // 设置关闭按钮提示文本
```

### 禁用状态

```go
button.SetDisable(true)  // 设置禁用
button.Disable() bool    // 查询是否禁用
```

### 事件

```go
button.SetOnMouseEnter(func(sender lcl.IObject) {
    // 鼠标进入
})

button.SetOnMouseLeave(func(sender lcl.IObject) {
    // 鼠标离开
})

button.SetOnMouseDown(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
    // 鼠标按下
})

button.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
    // 鼠标抬起
})

button.SetOnCloseClick(func(sender lcl.IObject) {
    // 关闭按钮点击
})

button.SetOnPaint(func(sender lcl.IObject) {
    // 自定义绘制
})
```

## TTab 标签页控件

TTab 是标签页容器，支持多页面切换、滚动、关闭等功能。

### 创建标签页

```go
tab := wg.NewTab(parent)
tab.SetLeft(10)
tab.SetTop(10)
tab.SetWidth(400)
tab.SetHeight(300)
```

### 添加页面

```go
page := tab.NewPage() // 返回 *TPage
page.SetCaption("标签标题")
page.SetActiveColor(activeColor)
page.SetDefaultColor(defaultColor)
```

### 页面管理

```go
// 获取所有页面
pages := tab.Pages() // 返回 []*TPage

// 隐藏所有活动页面
tab.HideAllActivated()

// 移除页面（传入 *TPage 指针）
tab.RemovePage(page)

// 设置变更回调
tab.SetOnChange(func(sender lcl.IObject) {
    // 标签变更时触发
})

// 滚动按钮
scrollLeft := tab.ScrollLeft()   // 获取左滚动按钮
scrollRight := tab.ScrollRight() // 获取右滚动按钮
tab.EnableScrollButton(true)     // 启用/禁用滚动按钮
```

### TPage 属性

```go
page.SetCaption("标题")         // 设置标签标题
page.SetActive(true)            // 设置为活动页
page.Active() bool              // 查询是否活动页
page.Hide()                     // 隐藏页面
page.Show()                     // 显示页面
page.Close()                    // 关闭页面
page.Remove()                   // 从标签页移除
page.SetActiveColor(color)      // 设置活动状态颜色
page.SetDefaultColor(color)     // 设置默认状态颜色
page.Button() *TButton          // 获取标签头按钮
page.IsEnterClose() bool        // 查询是否悬停关闭
```

### 页面事件

```go
page.SetOnShow(func(sender lcl.IObject) {
    // 页面显示时触发
})

page.SetOnHide(func(sender lcl.IObject) {
    // 页面隐藏时触发
})

page.SetOnClose(func(page *TPage, canClose *bool) {
    // 页面关闭时触发，设置 *canClose = false 可阻止关闭
})

page.SetOnClick(func(sender lcl.IObject) {
    // 标签头点击时触发
})
```

### 滚动支持

当标签数量超过容器宽度时，自动显示滚动按钮：

```go
// 滚动按钮自动管理，无需手动配置
```

## 使用场景

### TButton 适用场景

- 需要自定义外观的按钮
- 关闭按钮、收藏按钮等特殊功能按钮
- 需要渐变、圆角、透明效果的按钮
- 工具栏、标题栏按钮

### TTab 适用场景

- 多文档界面（MDI）
- 设置面板
- 代码编辑器标签
- 需要动态添加/删除页面的场景

## 完整示例

```go
package main

import (
    "fmt"
    "github.com/energye/lcl/lcl"
    "github.com/energye/lcl/types"
    "github.com/energye/lcl/types/colors"
    "github.com/energye/lcl/widget/wg"
)

type TMainForm struct {
    lcl.TEngForm
}

var mainForm TMainForm

func main() {
    lcl.Init()
    lcl.Run(&mainForm)
}

func (m *TMainForm) FormCreate(sender lcl.IObject) {
    m.SetCaption("自定义组件示例")
    m.SetPosition(types.PoScreenCenter)
    m.SetWidth(600)
    m.SetHeight(500)

    // 创建自定义按钮
    btn := wg.NewButton(m)
    btn.SetCaption("自定义按钮")
    btn.SetLeft(20)
    btn.SetTop(20)
    btn.SetWidth(150)
    btn.SetHeight(40)
    btn.SetDefaultColor(colors.TColor(0x3498DB), colors.TColor(0x2980B9))
    btn.SetEnterColor(colors.TColor(0x2980B9), colors.TColor(0x2471A3))
    btn.SetRadius(6)
    btn.SetOnMouseUp(func(sender lcl.IObject, button types.TMouseButton, shift types.TShiftState, x, y int32) {
        fmt.Println("按钮点击")
    })

    // 创建标签页
    tab := wg.NewTab(m)
    tab.SetLeft(20)
    tab.SetTop(80)
    tab.SetWidth(500)
    tab.SetHeight(350)

    // 添加页面
    page1 := tab.NewPage()
    page1.SetCaption("首页")

    page2 := tab.NewPage()
    page2.SetCaption("设置")

    page3 := tab.NewPage()
    page3.SetCaption("关于")
}
```
