# 应用配置

应用配置用于设置应用程序的元数据信息，如标题、图标、版本号、权限等。通过 **设置 > 应用配置**（快捷键 `Ctrl+F11` / macOS `Cmd+F11`）打开。

## 通用配置

| 配置项 | 说明 | 输入类型 | 默认值/提示 |
|--------|------|----------|-------------|
| 标题 | 应用程序的显示名称 | 文本输入 | 提示：`my energy app` |
| 标识 | 应用程序的唯一标识符（反向域名格式） | 文本输入 | 提示：`company.product.app` |
| 描述 | 应用程序的描述信息 | 文本输入 | 提示：`your application description.` |
| 版本 | 应用程序版本号 | 文本输入 | 提示：`1.2.3.4` |
| 版权 | 应用程序的版权声明 | 文本输入 | 提示：`Copyright (C)` |
| 应用图标 | 应用程序的图标文件 | 图标选择器（点击上传） | 支持 PNG 格式，最大 1024×1024 |

### 图标设置

- 支持上传 PNG 格式的图标文件
- 设计器自动将 PNG 转换为各平台所需格式：
  - **Windows** - 自动生成 `.ico` 格式
  - **macOS** - 自动生成 `.icns` 格式
  - **Linux** - 使用 `.png` 格式
- 图标自动缩放到各平台要求的标准尺寸

## Windows 配置

### Manifest 设置

Windows Manifest 用于声明应用程序的兼容性、权限和系统行为。配置分为下拉选择项和复选框两类。

#### 下拉选择项

| 配置项 | 说明 | 可选值 | 默认值 |
|--------|------|--------|--------|
| 兼容性操作系统 | 目标 Windows 最低版本 | `Windows Vista`、`Windows 7`、`Windows 8`、`Windows 8.1`、`Windows 10`、`Windows 11` | `Windows Vista` |
| DPI 感知 | 高 DPI 显示支持模式 | `System (开启)` - 系统级 DPI 感知<br>`UnAware (关闭)` - 不感知 DPI<br>`PerMonitor (true/PM)` - 每显示器感知<br>`PerMonitorV2 (true/PM-V2)` - 每显示器感知 V2 | `System (开启)` |
| 运行级别 | 应用程序权限运行级别 | `AsInvoker (当前用户)` - 以当前用户权限运行<br>`HighestAvailable (最高可用权限)` - 以最高可用权限运行<br>`RequireAdministrator (要求管理员)` - 要求管理员权限 | `AsInvoker (当前用户)` |

#### 复选框选项

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| uiAccess (用户界面访问) | 允许应用程序访问受保护的用户界面元素 | 关闭 |
| autoElevate (自动提权) | 自动请求提升权限（需配合管理员运行级别） | 关闭 |
| DisableTheming (禁用主题) | 禁用 Windows 视觉主题渲染 | 关闭 |
| disableWindowFiltering (禁用窗口过滤) | 禁用窗口过滤机制 | 关闭 |
| highResolutionScrollingAware (高分辨率滚动) | 启用高分辨率显示器的平滑滚动支持 | 开启 |
| ultraHighResolutionScrollingAware (超高分辨率滚动) | 启用超高分辨率显示器的滚动支持 | 开启 |
| longPathAware (启用长路径支持) | 支持超过 260 字符的文件路径 | 开启 |
| printerDriverIsolation (打印驱动隔离) | 将打印机驱动程序隔离在独立进程中运行 | 关闭 |
| gdiScaling (GDI 自动缩放) | 启用 GDI 自动缩放以改善高 DPI 显示效果 | 开启 |
| segmentHeap (分段堆) | 使用段堆内存管理方式，可降低内存占用 | 关闭 |
| useCommonControlsV6 (视觉样式) | 使用 Common Controls v6，启用现代视觉样式 | 开启 |

## macOS 配置

### Info.plist 设置

Info.plist 是 macOS 应用包的核心配置文件，定义了应用的元数据和系统行为。

#### 下拉选择项

| 配置项 | 说明 | 可选值 | 默认值 |
|--------|------|--------|--------|
| LSUIElement (运行模式) | 控制应用是否在 Dock 中显示图标 | `false (常规前台应用)` - 正常显示在 Dock<br>`true (后台应用, 无 Dock 图标)` - 后台运行，不显示 Dock 图标 | `false (常规前台应用)` |
| LS Minimum System Version (最低系统版本) | 应用支持的最低 macOS 版本 | `10.15 (Intel)` - 支持 Intel Mac<br>`11.0 (Apple Silicon)` - 支持 Apple Silicon 芯片 | `10.15 (Intel)` |

#### 文本输入项

| 配置项 | 说明 | 提示信息 |
|--------|------|----------|
| CFBundleName (短标题) | 应用的短显示名称 | `应用的短显示名称, 默认: 构建二进制文件名` |
| CFBundleLocalizations (本地化语言列表) | 支持的语言本地化列表，逗号分隔 | `本地化语言列表, 豆号分隔 zh_CN,en` |

#### 自动生成的 Info.plist 字段

以下字段在打包时由设计器根据通用配置自动生成，无需手动设置：

| 字段 | 说明 | 数据来源 |
|------|------|----------|
| CFBundleExecutable | 可执行文件名 | 项目构建文件名 |
| CFBundleDisplayName | 应用显示名称 | 通用配置中的"标题" |
| CFBundleIdentifier | 应用唯一标识符 | 通用配置中的"标识" |
| CFBundleVersion | 内部版本号 | 通用配置中的"版本" |
| CFBundleShortVersionString | 短版本号 | 通用配置中的"版本" |
| CFBundleGetInfoString | 简介信息 | 通用配置中的"描述" |
| CFBundleIconFile | 图标文件名 | 自动生成的 `.icns` 文件名 |
| NSHumanReadableCopyright | 版权信息 | 通用配置中的"版权" |

## Linux 配置

Linux 平台目前无特定的应用配置选项（相关配置在打包配置中设置）。

## 保存配置

修改配置后点击"保存"按钮，配置信息写入项目的 `.egp` 配置文件。同时自动更新应用图标文件到项目的资源目录。
