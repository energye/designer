# 菜单编辑

菜单编辑模块提供标准的编辑操作菜单（撤销、重做、剪切、复制、粘贴等），适用于 LCL 控件、SynEdit 编辑器和 WebView。

## TMenuEditing 组件

```go
import "github.com/energye/lcl/lcl"

editing := lcl.NewTMenuEditing(parent)
```

## 标准编辑操作

TMenuEditing 提供以下标准操作：

| 操作 | 方法 | 快捷键 | 说明 |
|------|------|--------|------|
| 撤销 | Undo | Ctrl+Z | 撤销上一步操作 |
| 重做 | Redo | Ctrl+Y | 重做上一步操作 |
| 剪切 | Cut | Ctrl+X | 剪切选中内容 |
| 复制 | Copy | Ctrl+C | 复制选中内容 |
| 粘贴 | Paste | Ctrl+V | 粘贴剪贴板内容 |
| 删除 | Delete | Delete | 删除选中内容 |
| 全选 | SelectAll | Ctrl+A | 选中全部内容 |

## 可用性检查

每个操作都有对应的 `Can*` 方法检查是否可用：

```go
if editing.CanUndo() {
    editing.Undo()
}

if editing.CanCopy() {
    editing.Copy()
}
```

| 检查方法 | 说明 |
|----------|------|
| CanUndo() | 是否可撤销 |
| CanRedo() | 是否可重做 |
| CanCut() | 是否可剪切 |
| CanCopy() | 是否可复制 |
| CanPaste() | 是否可粘贴 |
| CanDelete() | 是否可删除 |
| CanSelectAll() | 是否可全选 |

## 平台适配

快捷键自动适配平台：

| 平台 | 修饰键 | 说明 |
|------|--------|------|
| Windows | Ctrl | Ctrl+Z, Ctrl+C 等 |
| Linux | Ctrl | Ctrl+Z, Ctrl+C 等 |
| macOS | Meta (Cmd) | Cmd+Z, Cmd+C 等 |

### 获取平台修饰键

```go
func platformControl() string {
    if desTool.IsDarwin {
        return "Meta"
    }
return "Ctrl"
}
```

## 在菜单中使用

```go
// 创建菜单栏
menuBar := lcl.NewTMainMenu(form)

// 创建编辑菜单
editMenu := lcl.NewTMenuItem(menuBar)
editMenu.SetCaption("编辑(&E)")

// 添加编辑操作
editing := lcl.NewTMenuEditing(form)

undoItem := lcl.NewTMenuItem(editMenu)
undoItem.SetCaption("撤销")
undoItem.SetShortcut("Ctrl+Z")
undoItem.SetOnClick(func(sender lcl.IObject) {
    editing.Undo()
})

copyItem := lcl.NewTMenuItem(editMenu)
copyItem.SetCaption("复制")
copyItem.SetShortcut("Ctrl+C")
copyItem.SetOnClick(func(sender lcl.IObject) {
    editing.Copy()
})
```

## ActionList

TMenuEditing 内部使用 ActionList 管理标准操作：

```go
actionList := editing.ActionList()
```

ActionList 包含预定义的标准编辑动作，可以绑定到菜单项或工具栏按钮。

## 完整示例

```go
package main

import (
    "github.com/energye/lcl/lcl"
)

func main() {
    lcl.SetOnBeforeRun(func() {
        form := lcl.Application.CreateForm()
        form.SetCaption("编辑菜单示例")

        // 创建编辑操作管理器
        editing := lcl.NewTMenuEditing(form)

        // 创建菜单栏
        menuBar := lcl.NewTMainMenu(form)

        // 文件菜单
        fileMenu := lcl.NewTMenuItem(menuBar)
        fileMenu.SetCaption("文件(&F)")

        // 编辑菜单
        editMenu := lcl.NewTMenuItem(menuBar)
        editMenu.SetCaption("编辑(&E)")

        // 添加编辑操作到菜单
        items := []struct {
            caption  string
            shortcut string
            action   func()
        }{
            {"撤销", "Ctrl+Z", func() { editing.Undo() }},
            {"重做", "Ctrl+Y", func() { editing.Redo() }},
            {"", "", nil}, // 分隔线
            {"剪切", "Ctrl+X", func() { editing.Cut() }},
            {"复制", "Ctrl+C", func() { editing.Copy() }},
            {"粘贴", "Ctrl+V", func() { editing.Paste() }},
            {"删除", "Delete", func() { editing.Delete() }},
            {"", "", nil}, // 分隔线
            {"全选", "Ctrl+A", func() { editing.SelectAll() }},
        }

        for _, item := range items {
            if item.caption == "" {
                sep := lcl.NewTMenuItem(editMenu)
                sep.SetCaption("-")
                editMenu.Add(sep)
                continue
            }

            menuItem := lcl.NewTMenuItem(editMenu)
            menuItem.SetCaption(item.caption)
            menuItem.SetShortcut(item.shortcut)
            action := item.action
            menuItem.SetOnClick(func(sender lcl.IObject) {
                action()
            })
            editMenu.Add(menuItem)
        }
    })

    lcl.Run(lcl.Application.Forms()...)
}
```
