# 菜单编辑

菜单编辑模块提供标准的编辑操作菜单（撤销、重做、剪切、复制、粘贴等），适用于 LCL 控件、SynEdit 编辑器和 WebView。

## TMenuEditing 组件

```go
import "github.com/energye/lcl/lcl"

editing := lcl.NewMenuEditing(parent)
```

**注意**：构造函数是 `NewMenuEditing`，不是 `NewTMenuEditing`。

## 内置 Action

TMenuEditing 包含以下预定义的 Action：

| Action 字段 | 类型 | 说明 |
|-------------|------|------|
| ActionList | lcl.IActionList | Action 列表 |
| UndoAction | lcl.IEditUndo | 撤销 |
| RedoAction | lcl.IAction | 重做 |
| CutAction | lcl.IEditCut | 剪切 |
| CopyAction | lcl.IEditCopy | 复制 |
| PasteAction | lcl.IEditPaste | 粘贴 |
| DeleteAction | lcl.IEditDelete | 删除 |
| SelectAllAction | lcl.IEditSelectAll | 全选 |

## 使用方式

TMenuEditing 通过 Action 模式工作，每个 Action 自带 Execute 和 Update 回调：

```go
editing := lcl.NewMenuEditing(form)

// Action 的 Execute 回调在操作执行时触发
// Action 的 Update 回调在菜单显示时触发，用于更新可用状态
```

### 在菜单中使用

```go
// 创建菜单栏
menuBar := lcl.NewTMainMenu(form)

// 创建编辑菜单
editMenu := lcl.NewTMenuItem(menuBar)
editMenu.SetCaption("编辑(&E)")

// 添加撤销菜单项，绑定到 UndoAction
undoItem := lcl.NewTMenuItem(editMenu)
undoItem.SetCaption("撤销")
undoItem.SetAction(editing.UndoAction)

// 添加剪切菜单项，绑定到 CutAction
cutItem := lcl.NewTMenuItem(editMenu)
cutItem.SetCaption("剪切")
cutItem.SetAction(editing.CutAction)

// 添加复制菜单项，绑定到 CopyAction
copyItem := lcl.NewTMenuItem(editMenu)
copyItem.SetCaption("复制")
copyItem.SetAction(editing.CopyAction)
```

### 使用 ActionList

```go
// 获取 ActionList，可用于绑定到工具栏
actionList := editing.ActionList
```

## 平台适配

快捷键自动适配平台：

| 平台 | 修饰键 | 说明 |
|------|--------|------|
| Windows | Ctrl | Ctrl+Z, Ctrl+C 等 |
| Linux | Ctrl | Ctrl+Z, Ctrl+C 等 |
| macOS | Meta (Cmd) | Cmd+Z, Cmd+C 等 |

### 获取平台修饰键

```go
control := lcl.PlatformControl()
// Windows/Linux: "Ctrl"
// macOS: "Meta"
```

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
        editing := lcl.NewMenuEditing(form)

        // 创建菜单栏
        menuBar := lcl.NewTMainMenu(form)

        // 编辑菜单
        editMenu := lcl.NewTMenuItem(menuBar)
        editMenu.SetCaption("编辑(&E)")

        // 撤销
        undoItem := lcl.NewTMenuItem(editMenu)
        undoItem.SetCaption("撤销")
        undoItem.SetAction(editing.UndoAction)

        // 重做
        redoItem := lcl.NewTMenuItem(editMenu)
        redoItem.SetCaption("重做")
        redoItem.SetAction(editing.RedoAction)

        // 分隔线
        sep1 := lcl.NewTMenuItem(editMenu)
        sep1.SetCaption("-")

        // 剪切
        cutItem := lcl.NewTMenuItem(editMenu)
        cutItem.SetCaption("剪切")
        cutItem.SetAction(editing.CutAction)

        // 复制
        copyItem := lcl.NewTMenuItem(editMenu)
        copyItem.SetCaption("复制")
        copyItem.SetAction(editing.CopyAction)

        // 粘贴
        pasteItem := lcl.NewTMenuItem(editMenu)
        pasteItem.SetCaption("粘贴")
        pasteItem.SetAction(editing.PasteAction)

        // 删除
        deleteItem := lcl.NewTMenuItem(editMenu)
        deleteItem.SetCaption("删除")
        deleteItem.SetAction(editing.DeleteAction)

        // 分隔线
        sep2 := lcl.NewTMenuItem(editMenu)
        sep2.SetCaption("-")

        // 全选
        selectAllItem := lcl.NewTMenuItem(editMenu)
        selectAllItem.SetCaption("全选")
        selectAllItem.SetAction(editing.SelectAllAction)
    })

    lcl.Run(lcl.Application.Forms()...)
}
```
