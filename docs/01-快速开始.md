# 快速开始

ENERGY Designer 是专为 ENERGY 跨平台 GUI 框架打造的可视化设计工具，以原生桌面 UI 设计为核心，生成可维护的 Go 源码。同时支持通过 WebView / CEF 组件嵌入现代 Web 内容，轻松构建原生 + Web 的混合桌面应用。

## 环境要求

- [Git](https://git-scm.com/)
- [Golang](https://go.dev/) 1.20 或更高版本

## 安装运行

### 从源码运行

```bash
git clone https://github.com/energye/designer.git
cd designer
go mod tidy
go run main.go
```

### 下载发行版

前往 [GitHub Releases](https://github.com/energye/designer/releases) 下载已编译版本，直接运行即可。

## 1 分钟创建第一个应用

1. **启动设计器** - 运行 ENERGY Designer
2. **新建项目** - 按 `Ctrl+P`（macOS 为 `Cmd+P`），填写项目名称和目录，选择 GUI 渲染框架
3. **拖拽组件到画布** - 从左侧"组件库"面板选择组件，点击后在设计画布上点击放置
4. **配置属性和事件** - 在右侧"对象检查器"面板修改组件属性、绑定事件
5. **运行预览** - 按 `F9` 编译并预览应用

## 支持的 UI 技术

| 引擎       | 描述          | 适用场景           |
|----------|-------------|----------------|
| LCL 原生控件 | 系统原生桌面控件    | 企业管理软件 / 工具类应用 |
| WebView  | 系统自带网页引擎    | 轻量混合桌面应用       |
| CEF      | Chromium 内核 | 复杂现代化桌面应用      |

## 平台支持

| 平台      | 架构                        | 说明               |
|---------|---------------------------|------------------|
| Windows | 386 / amd64               | Win7+            |
| macOS   | amd64 / arm64             | Universal Binary |
| Linux   | 386 / amd64 / arm / arm64 | 主流发行版            |

## 项目文件结构

创建项目后，ENERGY Designer 会自动生成以下目录结构：

```
myproject/
├── app/                        # 应用代码目录
│   ├── app.go                  # 窗体列表和平台初始化（自动生成）
│   ├── myform.go               # 用户业务逻辑代码（可自由编辑）
│   ├── myform.ui.go            # UI 代码文件（自动生成，勿手动修改）
├── layouts/                    # UI 布局文件目录
│   └── myform.ui               # UI 布局文件（JSON 格式，勿手动修改）
├── resources/                  # 资源目录
│   ├── embed/                  # 嵌入资源（icon.png、icon.ico、icon.icns）
│   ├── metadata/               # 平台元数据（Windows .syso 文件、macOS 本地化）
│   ├── resources.go            # 资源引用代码（自动生成）
│   └── resources_windows.go    # Windows 平台资源引用（自动生成）
├── go.mod                      # Go 模块文件
├── go.sum                      # Go 依赖校验文件
└── main.go                     # 程序入口
```

### 文件说明

| 文件 | 说明 | 是否可编辑           |
|------|------|-----------------|
| `main.go` | 程序入口，根据 GUI 框架自动选择初始化方式 | 是               |
| `app/app.go` | 窗体维护列表，包含所有窗体的注册和平台特定初始化 | 否         |
| `app/myform.go` | 用户业务逻辑代码，事件处理函数写在此文件中 | 是 |
| `app/myform.ui.go` | 根据布局文件自动生成的 Go 代码，包含组件创建和属性设置 | 否         |
| `layouts/myform.ui` | JSON 格式的 UI 布局描述文件，记录组件属性和层级关系 | 否        |
| `resources/` | 应用资源目录，包含图标、平台元数据等 | 部分可编辑           |

## 社区交流

- Issues: https://github.com/energye/designer/issues
- QQ 群：541258627
