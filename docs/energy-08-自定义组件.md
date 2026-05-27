# 自定义组件

ENERGY 提供自定义绘制的 UI 组件，位于 `lcl/wg` 包下，包括 TButton（自定义按钮）和 TTab/TPage（标签页控件）。

## TButton 自定义按钮

TButton 是完全自定义绘制的按钮组件，支持渐变背景、图标、圆角、多种状态样式。

### 创建按钮

```go
import "github.com/energye/lcl/widget/wg"

button := wg.NewTButton(parent)
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
| Disabled | 禁用状态 | SetEnabled(false) |

### 颜色配置

```go
// 默认状态渐变
button.SetColor(clWhite)
button.SetColorEnd(clSilver)

// 悬停状态渐变
button.SetEnterColor(clBlue)
button.SetEnterColorEnd(clDarkBlue)

// 按下状态渐变
button.SetDownColor(clNavy)
button.SetDownColorEnd(clNavy)

// 禁用状态
button.SetDisabledColor(clGray)
button.SetDisabledColorEnd(clLightGray)
```

### 边框

```go
button.SetBorder(true)
button.SetBorderColor(clBlack)
button.SetBorderWidth(1)
```

### 圆角

```go
button.SetRounded(true)
button.SetCornerRadius(8)
```

### 透明度

```go
button.SetAlpha(200) // 0-255，255 为不透明
```

### 图标

```go
button.SetIcon(image)           // 设置图标
button.SetIconFavorite(true)    // 收藏图标样式
button.SetIconClose(true)       // 关闭按钮样式
button.SetIconCloseHighlight(true) // 关闭按钮高亮样式
button.SetIconCenter(true)      // 图标居中
```

### 文本对齐

```go
button.SetTextAlign(wg.AlignCenter)  // 居中
button.SetTextAlign(wg.AlignLeft)    // 左对齐
button.SetTextAlign(wg.AlignRight)   // 右对齐
```

### 自动大小

```go
button.SetAutoSize(true) // 根据文本自动调整大小
```

### 提示文本

```go
button.SetHint("这是提示文本")
```

### 点击事件

```go
button.SetOnClick(func(sender lcl.IObject) {
    fmt.Println("按钮被点击")
})
```

## TTab 标签页控件

TTab 是标签页容器，支持多页面切换、滚动、关闭等功能。

### 创建标签页

```go
tab := wg.NewTTab(parent)
tab.SetLeft(10)
tab.SetTop(10)
tab.SetWidth(400)
tab.SetHeight(300)
```

### 添加页面

```go
page := tab.NewPage("标签标题")
page.SetColor(clWhite)
```

### 页面管理

```go
// 隐藏所有页面
tab.HideAllActivated()

// 移除页面
tab.RemovePage(pageIndex)
```

### TPage 属性

```go
page.SetActive(true)    // 设置为活动页
page.Hide()             // 隐藏页面
page.Show()             // 显示页面
page.Close()            // 关闭页面
page.SetColor(clWhite)  // 设置背景色
```

### 页面按钮

每个页面有一个关联的按钮（标签头）：

```go
button := page.Button()
button.SetCaption("标签名")
```

### 滚动支持

当标签数量超过容器宽度时，自动显示滚动按钮：

```go
// 滚动按钮自动管理，无需手动配置
```

### 页面事件

```go
page.SetOnActivate(func(sender lcl.IObject) {
    // 页面激活时触发
})

page.SetOnClose(func(sender lcl.IObject) {
    // 页面关闭时触发
})
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
    "github.com/energye/lcl/widget/wg"
)

func main() {
    lcl.SetOnBeforeRun(func() {
        form := lcl.Application.CreateForm()

        // 创建自定义按钮
        btn := wg.NewTButton(form)
        btn.SetCaption("自定义按钮")
        btn.SetLeft(20)
        btn.SetTop(20)
        btn.SetWidth(150)
        btn.SetHeight(40)
        btn.SetColor(0x3498DB)
        btn.SetColorEnd(0x2980B9)
        btn.SetEnterColor(0x2980B9)
        btn.SetEnterColorEnd(0x2471A3)
        btn.SetRounded(true)
        btn.SetCornerRadius(6)
        btn.SetOnClick(func(sender lcl.IObject) {
            fmt.Println("按钮点击")
        })

        // 创建标签页
        tab := wg.NewTTab(form)
        tab.SetLeft(20)
        tab.SetTop(80)
        tab.SetWidth(500)
        tab.SetHeight(350)

        // 添加页面
        page1 := tab.NewPage("首页")
        page1.SetColor(clWhite)

        page2 := tab.NewPage("设置")
        page2.SetColor(clWhite)

        page3 := tab.NewPage("关于")
        page3.SetColor(clWhite)
    })

    lcl.Run(lcl.Application.Forms()...)
}
```
