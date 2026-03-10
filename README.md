## ENERGY GUI Designer

![go-version](https://img.shields.io/github/go-mod/go-version/energye/designer?logo=git&logoColor=green)
[![github](https://img.shields.io/github/last-commit/energye/energy/main.svg?logo=github&logoColor=green&label=commit)](https://github.com/energye/designer)
![repo](https://img.shields.io/github/repo-size/energye/designer.svg?logo=github&logoColor=green&label=repo-size)
---

### 🌟 项目简介

ENERGY Designer 是专为 ENERGY 跨平台 GUI 框架打造、且基于该框架开发的可视化设计器，它采用 Go LCL 组件库实现，提供所见即所得的设计体验与简化的 GUI 设计功能，开发者可通过拖拽控件、设置属性直观操作，快速创建和编辑 GUI 界面的同时，还能自动生成对应的 Go 代码。

### 核心特性

####  设计器功能
- **可视化设计**：所见即所得的 GUI 界面设计体验
- **控件拖拽**：支持标准控件的拖拽布局
- **属性编辑**：实时属性和事件设置
- **组件管理**：完整的组件树查看和管理
- **预览运行**：实时预览设计效果

###  使用场景

1. **快速原型开发**：快速构建 GUI 应用
2. **界面设计**：可视化设计复杂的用户界面
3. **代码生成**：自动生 ENERGY GUI 的代码
4. **丰富的组件**：原生控件, CEF控件, Webview控件

*ENERGY Designer - 让 GUI 开发更简单*

### 安装 Designer

- 创建工作目录. 如: `C:\app\workspace`
- 克隆项目依赖模块. 进入 `workspace` 目录，执行如下克隆命令

| 项目模块     | 克隆命令                                              | 描述                                |
|----------|---------------------------------------------------|-----------------------------------|
| LCL      | git clone https://github.com/energye/lcl.git      | ENERGY GUI 框架基础库, 需切换分枝`3.0-beta` |
| Widget   | git clone https://github.com/energye/widget.git   | ENERGY GUI 框架自定义组件库               |
| Designer | git clone https://github.com/energye/designer.git | ENERGY GUI 设计器库                   |

- 初始化工作空间, 进入 `workspace` 目录，执行如下命令
```shell
  go work init
  go work use lcl
  go work use widget
  go work use designer
```

- 运行. 进入 designer 目录

`go run main.go`

如果有需求，请提 issue, 或加入 QQ 群进行交流（541258627）

### 截图

![ENERGY-designer-create-project.png](docs/image/ENERGY-designer-create-project.png)


![ENERGY-designer-macOSrun.png.png](docs%2Fimage%2FENERGY-designer-macOSrun.png.png)

![ENERGY-designer.png](docs%2Fimage%2FENERGY-designer.png)

![ENERGY-designer-linux.png](docs/image/ENERGY-designer-linux.png)

![ENERGY-designer-preview.png](docs%2Fimage%2FENERGY-designer-preview.png)

![ENERGY-designer-linux-preview.png](docs/image/ENERGY-designer-linux-preview.png)