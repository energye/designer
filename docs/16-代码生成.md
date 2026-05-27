# 代码生成

ENERGY Designer 的代码生成引擎负责将可视化设计转换为可维护的 Go 源码。代码生成在设计过程中自动触发，无需手动操作。

## 生成的文件类型

### UI 布局文件 (.ui)

JSON 格式的布局描述文件，存储在 `layouts/` 目录下，记录窗体中所有组件的属性和层级关系。

- 实时自动更新：每次属性变更、组件增删、位置调整都会自动同步
- 用于布局恢复：重新打开项目时，从 `.ui` 文件恢复设计状态
- 文件名格式：`{窗体名小写}.ui`
- 存储位置：`layouts/` 目录

### UI 代码文件 (.ui.go)

根据 `.ui` 文件自动生成的 Go 源码，存储在 `app/` 目录下，包含窗体和组件的创建、属性设置代码。

- **自动生成，禁止手动修改** - 每次设计变更会被重新生成
- 包含组件创建、属性设置、父子关系建立等代码
- 文件名格式：`{窗体名小写}.ui.go`
- 存储位置：`app/` 目录

### 用户代码文件 (.go)

用户的业务逻辑代码文件，存储在 `app/` 目录下，包含事件处理函数等。

- **用户可自由编辑** - 设计器使用 AST 技术保护用户代码
- 事件绑定自动生成函数骨架
- 窗体名称修改时自动更新 self 引用
- 文件名格式：`{窗体名小写}.go`
- 存储位置：`app/` 目录

### 窗体列表文件 (app.go)

自动生成的窗体维护列表，存储在 `app/` 目录下。

- 包含所有窗体的注册列表
- 包含平台特定的初始化代码（如 Linux GTK 版本选择、macOS Universal Binary 设置）
- 文件名：`app.go`
- 存储位置：`app/` 目录

## 文件关系链

三个文件形成完整的代码生成链：

```
layouts/form1.ui (JSON 布局)
       │
       ▼
app/form1.ui.go (自动生成的 Go 代码)
       │
       ▼
app/form1.go (用户业务逻辑代码)
```

### 结构体组合关系

用户代码文件通过结构体嵌入组合自动生成的 UI 结构体：

```go
// form1.ui.go - 自动生成
type TForm1UI struct {
    window.TWindow
    Button1  lcl.IButton
    Panel1   lcl.IPanel
}

// form1.go - 用户代码
type TForm1 struct {
    TForm1UI  // 嵌入自动生成的 UI 结构体
}
```

这种设计使得：
- 自动生成的代码可以随时重新生成，不影响用户代码
- 用户代码通过嵌入自动获得所有组件的访问权限
- 用户可以安全地添加自定义字段和方法

## 窗体生命周期

自动生成的 `.ui.go` 文件包含完整的窗体生命周期事件：

| 事件 | 说明 | 触发时机 |
|------|------|----------|
| `FormCreate` | 窗体创建 | 窗体初始化时自动调用 |
| `OnShow` | 窗体显示 | 窗体显示时触发 |
| `OnCloseQuery` | 关闭前询问 | 用户请求关闭窗体时触发，可阻止关闭 |
| `OnClose` | 窗体关闭 | `OnCloseQuery` 允许关闭后触发 |

用户可以在 `.go` 文件中重写这些事件：

```go
func (m *TForm1) OnFormCreate(sender lcl.IObject) {
    // 窗体初始化代码
}

func (m *TForm1) OnShow(sender lcl.IObject) {
    // 窗体显示时的处理
}
```

## .ui 文件格式

`.ui` 文件使用 JSON 格式描述窗体布局：

```json
{
  "name": "Form1",
  "class_name": "TEngForm",
  "mod": "",
  "type": 0,
  "properties": [
    {"name": "Caption", "value": "Form1", "type": 0},
    {"name": "Height", "value": 400, "type": 1},
    {"name": "Width", "value": 600, "type": 1}
  ],
  "child": [
    {
      "name": "Button1",
      "class_name": "TButton",
      "mod": "lcl",
      "type": 2,
      "properties": [
        {"name": "Caption", "value": "Button1", "type": 0},
        {"name": "Left", "value": 48, "type": 1},
        {"name": "Top", "value": 50, "type": 1}
      ]
    }
  ]
}
```

### 字段说明

| 字段 | 说明 |
|------|------|
| `name` | 组件名称 |
| `class_name` | 组件类名（如 `TButton`、`TPanel`） |
| `mod` | 组件来源模块（`"lcl"` = LCL 原生，`"wvEnergy"` = WebView，空 = 窗体本身） |
| `type` | 组件类型（0 = 窗体，2 = 可视组件） |
| `properties` | 属性数组 |
| `child` | 子组件数组（递归结构） |

### 属性类型

| type 值 | 类型 | 说明 |
|---------|------|------|
| 0 | 文本 (String) | 字符串值 |
| 1 | 整数 (Integer) | 数值 |
| 2 | 布尔 (Boolean) | true/false |
| 3 | 枚举 (Enum) | 预定义选项值 |
| 4 | 集合 (Set) | 多选值 |
| 5 | 颜色 (Color) | 颜色值 |
| 6 | 类实例 (Class) | 复合对象（如 Icon） |

## 代码生成流程

```
组件属性变更 → UI 布局文件 (.ui) 更新 → Go 代码生成 (.ui.go) → 运行预览
```

### 触发时机

以下操作会触发代码生成：
- 放置新组件到设计画布
- 修改组件属性
- 删除组件
- 调整组件位置或大小
- 修改窗体属性
- 绑定或修改事件

## 用户代码保护

ENERGY Designer 使用 AST（抽象语法树）技术分析用户的 Go 代码文件，确保：

- 用户手写的代码不会被覆盖
- 自动生成的事件处理函数骨架正确插入
- 窗体名称修改时自动更新 self 相关引用
- 新增事件绑定只追加不覆盖

## 依赖管理

代码生成引擎自动管理项目依赖：
- 根据使用的组件自动添加所需的 import 包
- 根据 GUI 渲染框架（LCL / WebView / CEF）使用对应的包路径
- 跨平台时自动区分各平台的包路径（如 WebView 在 Windows/Linux/macOS 使用不同包）

## 平台特定初始化

自动生成的 `app.go` 包含平台特定的初始化代码：

```go
func init() {
    if runtime.GOOS == "linux" {
        libname.UseWS = "gtk3"  // 根据 UI 框架选择 GTK2 或 GTK3
    }
    if runtime.GOOS == "darwin" {
        libname.EnableUniversalBinary = false  // macOS Universal Binary 设置
    }
}
```

## main.go 入口差异

不同 GUI 框架的 `main.go` 入口代码不同：

### LCL 框架

```go
func main() {
    lcl.Run(app.Forms...)
}
```

### WebView 框架

```go
func main() {
    wvApp := wv.Init()
    wvApp.SetOptions(application.Options{DefaultURL: "https://example.com"})
    wv.Run(app.Forms...)
}
```

### CEF 框架

```go
暂时未添加
```
