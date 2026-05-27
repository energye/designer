# 菜单编辑

菜单编辑模块提供标准的编辑操作菜单（撤销、重做、剪切、复制、粘贴等），适用于 LCL 控件、SynEdit 编辑器和 WebView。

## TMenuEditing 组件

```go
import engLCL "github.com/energye/energy/v3/lcl"

editing := engLCL.NewMenuEditing(parent)
```

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
menuBar := lcl.NewMainMenu(form)

// 创建编辑菜单
editMenu := lcl.NewMenuItem(form)
editMenu.SetCaption("编辑(&E)")
menuBar.Items().Add(editMenu)

// 添加撤销菜单项，绑定到 UndoAction
undoItem := lcl.NewMenuItem(form)
undoItem.SetCaption("撤销")
undoItem.SetAction(editing.UndoAction)
editMenu.Add(undoItem)

// 添加剪切菜单项，绑定到 CutAction
cutItem := lcl.NewMenuItem(form)
cutItem.SetCaption("剪切")
cutItem.SetAction(editing.CutAction)
editMenu.Add(cutItem)

// 添加复制菜单项，绑定到 CopyAction
copyItem := lcl.NewMenuItem(form)
copyItem.SetCaption("复制")
copyItem.SetAction(editing.CopyAction)
editMenu.Add(copyItem)
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
control := engLCL.PlatformControl()
// Windows/Linux: "Ctrl"
// macOS: "Meta"
```

## 完整示例

```go
package main

import (
    engLCL "github.com/energye/energy/v3/lcl"
    "github.com/energye/lcl/lcl"
    "github.com/energye/lcl/types"
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
    m.SetCaption("编辑菜单示例")
    m.SetPosition(types.PoScreenCenter)

    // 创建编辑操作管理器
    editing := engLCL.NewMenuEditing(m)

    // 创建菜单栏
    menuBar := lcl.NewMainMenu(m)

    // 编辑菜单
    editMenu := lcl.NewMenuItem(m)
    editMenu.SetCaption("编辑(&E)")
    menuBar.Items().Add(editMenu)

    // 撤销
    undoItem := lcl.NewMenuItem(m)
    undoItem.SetCaption("撤销")
    undoItem.SetAction(editing.UndoAction)
    editMenu.Add(undoItem)

    // 重做
    redoItem := lcl.NewMenuItem(m)
    redoItem.SetCaption("重做")
    redoItem.SetAction(editing.RedoAction)
    editMenu.Add(redoItem)

    // 分隔线
    sep1 := lcl.NewMenuItem(m)
    sep1.SetCaption("-")
    editMenu.Add(sep1)

    // 剪切
    cutItem := lcl.NewMenuItem(m)
    cutItem.SetCaption("剪切")
    cutItem.SetAction(editing.CutAction)
    editMenu.Add(cutItem)

    // 复制
    copyItem := lcl.NewMenuItem(m)
    copyItem.SetCaption("复制")
    copyItem.SetAction(editing.CopyAction)
    editMenu.Add(copyItem)

    // 粘贴
    pasteItem := lcl.NewMenuItem(m)
    pasteItem.SetCaption("粘贴")
    pasteItem.SetAction(editing.PasteAction)
    editMenu.Add(pasteItem)

    // 删除
    deleteItem := lcl.NewMenuItem(m)
    deleteItem.SetCaption("删除")
    deleteItem.SetAction(editing.DeleteAction)
    editMenu.Add(deleteItem)

    // 分隔线
    sep2 := lcl.NewMenuItem(m)
    sep2.SetCaption("-")
    editMenu.Add(sep2)

    // 全选
    selectAllItem := lcl.NewMenuItem(m)
    selectAllItem.SetCaption("全选")
    selectAllItem.SetAction(editing.SelectAllAction)
    editMenu.Add(selectAllItem)
}
```
