# CLI 命令行工具

ENERGY Designer 提供独立命令行工具 `energy`，用于创建项目、编译、运行、打包，以及维护本机 ENERGY 配置和 CEF 运行环境。GUI 设计器中的 CEF 版本安装、切换能力也复用 CLI 的同一套核心逻辑。

## 安装

在源码目录编译：

```bash
cd designer/cmd/energy
go build -o energy
```

将生成的 `energy` 可执行文件放入系统 `PATH` 后，即可在任意目录使用。

## 基本用法

```bash
energy <command> [options]
```

不带参数运行 `energy` 会显示命令帮助。

参数解析规则：

| 写法 | 说明 |
|---|---|
| `-path /app/demo` | 短横线参数和值用空格分隔 |
| `-path=/app/demo` | 参数和值用 `=` 分隔 |
| `--all` | 布尔开关参数 |
| `chromium.version` | 非 `-` 开头的参数会作为位置参数使用，主要用于 `energy env` 查询 |

通用项目路径参数：

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-path <project-dir>` | 项目目录，目录内应包含 `.egp` 项目文件 | 当前工作目录 |

## 交互界面

CLI 在交互终端中默认使用 Bubble Tea/Bubbles 风格的终端界面，主要用于项目创建、选择 CEF 版本、确认覆盖、下载进度等场景。

以下情况会自动退回简单文本交互：

| 条件 | 说明 |
|---|---|
| `ENERGY_TUI=0` | 禁用 TUI，也支持 `false`、`off`、`no` |
| `CI` 环境变量存在 | CI 环境默认不使用 TUI |
| 非交互输入/输出 | 例如管道、重定向、`TERM=dumb` |
| Windows 7 等旧版 Windows | 避免终端虚拟控制能力不兼容 |

示例：

```bash
# 使用简单文本交互
ENERGY_TUI=0 energy init
```

## 命令列表

| 命令 | 说明 |
|---|---|
| `energy help` | 显示命令帮助 |
| `energy init` | 创建 ENERGY 项目 |
| `energy run` | 编译并运行当前项目 |
| `energy build` | 编译项目 |
| `energy package` | 编译并生成安装包 |
| `energy env` | 读取、写入本机 `~/.energy/config.json` |

## energy help

显示所有已注册命令。

```bash
energy help
```

## energy init

创建新的 ENERGY 项目。

```bash
energy init [-name <project-name>] [-path <project-dir>] [-ui <LCL|WV|CEF>]
```

CEF 项目可额外指定 CEF 安装参数：

```bash
energy init -name demo -ui CEF -cef-version 127 -os linux -arch amd64
```

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-name` | 项目名称；未传入时会交互输入 | 无 |
| `-path` | 项目创建目录 | 当前工作目录 |
| `-ui` | UI 框架，可选 `LCL`、`WV`、`CEF` | 未传入时交互选择 |
| `-framework` | `-ui` 的别名 | 无 |
| `-cef-version` | CEF 版本，支持内置完整版本或 major，例如 `127` | 未传入时交互选择 |
| `-os` | CEF 目标系统：`windows`、`linux`、`darwin` | 当前系统 |
| `-arch` | CEF 目标架构：`amd64`、`386`、`arm64`、`arm` | 当前架构 |
| `-cef-dir` | CEF 根目录 | `~/.energy/chromium` |

执行流程：

1. 收集项目名称、目录和 UI 框架。
2. 如果选择 `CEF`，先确认目标 CEF 版本、系统和架构。
3. 调用 CEF 安装/切换逻辑，确保 CEF 框架和对应 libenergy 运行时库可用。
4. 创建项目文件和模板代码。
5. 如目标目录已有冲突文件，会询问是否覆盖。
6. 创建完成后执行 Go 依赖更新。

示例：

```bash
# 交互式创建项目
energy init

# 创建 LCL 项目
energy init -name demo -path /home/user/demo -ui LCL

# 创建 CEF 项目，并使用 CEF 127
energy init -name cefdemo -path /home/user/cefdemo -ui CEF -cef-version 127
```

## energy run

编译并运行项目。

```bash
energy run [-path <project-dir>]
```

执行流程：

1. 加载项目目录下的 `.egp` 项目文件。
2. 执行当前平台构建。
3. 构建成功后运行生成的可执行文件。

示例：

```bash
energy run
energy run -path /home/user/myproject
```

## energy build

编译项目，生成可执行文件。

```bash
energy build [-path <project-dir>] [--all]
```

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-path` | 项目目录 | 当前工作目录 |
| `--all` | 按项目配置执行多平台构建 | 只构建当前平台 |

当前平台构建会根据运行系统进入对应构建流程：

| 当前系统 | 输出形式 |
|---|---|
| Windows | `build/<build_file_name>.exe` |
| Linux | `build/<build_file_name>` |
| macOS | `build/<package_name>.app/Contents/MacOS/<executable>` |

`--all` 依赖项目构建配置：

| 条件 | 说明 |
|---|---|
| `build_cgo_enabled=false` | 多平台构建会禁用 CGO |
| `build_other_platform=true` | 允许构建其他平台 |

满足条件时会依次构建：

| 系统 | 架构 |
|---|---|
| Windows | `amd64`、`386` |
| macOS | `amd64`、`arm64` |
| Linux | `amd64`、`386`、`arm`、`arm64` |

多平台构建输出文件会追加系统和架构后缀，例如 `demo_windows_amd64.exe`。

构建参数来源于项目 `.egp` 的 `build_option`：

| 字段 | 说明 |
|---|---|
| `output` | 构建输出目录 |
| `build_file_name` | 可执行文件名 |
| `build_mode_debug` | Debug 构建模式 |
| `build_mode_release` | Release 构建模式 |
| `go_args` | 传递给 `go build` 的额外参数，支持提取 `-tags`、`-ldflags` 和其他构建参数 |
| `mac_common_lib` | macOS 是否合并 Universal Binary |

示例：

```bash
# 构建当前平台
energy build

# 指定项目目录
energy build -path /home/user/myproject

# 按项目配置构建多平台产物
energy build --all
```

## energy package

编译并打包项目。

```bash
energy package [-path <project-dir>]
```

执行流程：

1. 加载 `.egp` 项目配置。
2. 强制启用 Release 构建模式。
3. 执行当前平台构建。
4. 根据当前系统和项目打包配置生成安装包。

各平台支持：

| 平台 | 支持产物 | 主要依赖 |
|---|---|---|
| Windows | NSIS `.exe`、MSIX/AppX | `makensis`、`MakeAppx.exe`、可选签名工具 |
| macOS | `.app`、`.pkg`、`.dmg` | `pkgbuild`、可选 `create-dmg`、可选 `codesign` |
| Linux | `.deb`、`.rpm`、`.AppImage` | `dpkg-deb`、`rpmbuild`、AppImage 相关工具 |

是否生成某种安装包由项目 `.egp` 的 `build_option` 控制，例如：

| 字段 | 说明 |
|---|---|
| `win_exe` | Windows NSIS 安装包 |
| `win_msi` | Windows MSIX/AppX 安装包 |
| `win_appx` | Windows MSIX/AppX 资源配置 |
| `mac_pkg` | macOS PKG |
| `mac_dmg` | macOS DMG |
| `linux_deb` | Linux DEB |
| `linux_rpm` | Linux RPM |
| `linux_app_image` | Linux AppImage |

示例：

```bash
energy package
energy package -path /home/user/myproject
```

## energy env

读取或写入本机配置文件：

```bash
~/.energy/config.json
```

基本用法：

```bash
energy env
energy env <key|json.path>
energy env -w <key|json.path>=<value>
energy env -l <module>
```

| 用法 | 说明 |
|---|---|
| `energy env` | 输出完整配置；读取时会合并当前 CLI 内置默认配置中缺失的字段 |
| `energy env chromium.version` | 按 JSON path 查询配置值 |
| `energy env version` | 按 key 名称模糊查询所有匹配项 |
| `energy env -w chromium.source=sourceforge` | 写入配置 |
| `energy env -l cef` | 列出已安装的 CEF 版本 |

### 查询配置

```bash
# 输出完整配置
energy env

# 查询 chromium 配置
energy env chromium

# 查询当前 CEF 版本
energy env chromium.version

# 查询所有名为 version 的字段
energy env version
```

JSON path 支持对象字段和数组索引：

```bash
energy env history_project[0]
```

### 写入配置

```bash
energy env -w registry=https://example.com/config
energy env -w proxy=http://127.0.0.1:7890
energy env -w chromium.dir=/home/user/.energy/chromium
```

写入规则：

| 规则 | 说明 |
|---|---|
| 类型保持 | 原字段是字符串、数字、布尔值时，会按原类型转换写入 |
| 模糊 key | 如果 key 匹配到多个字段，会交互选择目标字段 |
| JSON path | 可以直接写入明确路径，例如 `chromium.dir=...` |
| 新增字段 | 默认只允许创建 `chromium.source` |
| 废弃字段 | 写入默认配置时会移除旧的顶层 `cef_runtime` 字段 |

### 列出 CEF 版本

```bash
energy env -l cef
```

输出为 `chromium.dir` 下已安装的 CEF 目录。当前使用版本前会显示 `*`。

示例：

```text
  linux_amd64_109.1.18
* linux_amd64_127.3.5
  linux_amd64_147.0.14
```

### 安装或切换 CEF

写入 `chromium.version` 时，CLI 会执行安装/切换逻辑：

```bash
# 使用 major，自动解析为内置 CEF 完整版本
energy env -w chromium.version=127

# 使用完整 CEF 版本
energy env -w chromium.version=127.3.5

# 使用完整 os_arch_version
energy env -w chromium.version=linux_amd64_127.3.5
```

处理规则：

1. 如果目标 CEF 已安装且对应 libenergy 运行时库有效，直接切换。
2. 如果 CEF 未安装，先下载并解压 CEF 框架。
3. 如果目标 CEF 目录内的 libenergy 缺失、发行版本不匹配，或 API 读取到的 CEF major 不匹配，则重新下载 libenergy。
4. CEF 和 libenergy 准备完成后，更新 `chromium.version`。

`chromium.version` 保存格式为：

```text
<os>_<arch>_<cef-version>
```

例如：

```text
linux_amd64_127.3.5
windows_amd64_127.3.5
darwin_arm64_147.0.14
```

### 切换 libenergy 下载源

libenergy 下载源选择写入 `chromium.source`：

```bash
energy env -w chromium.source=sourceforge
```

当前内置配置中 `sourceforge` 已可用，`github` 预留为空。实际下载时优先使用：

1. 环境变量 `ENERGY_CEF_RUNTIME_SOURCE`
2. 用户配置 `chromium.source`
3. 内置配置中的默认 `source`

## CEF 和 libenergy 管理

CEF 根目录默认为：

```text
~/.energy/chromium
```

目录结构示例：

```text
~/.energy/chromium/
  .versions
  cef_binary_127.3.5+..._linux64_minimal.tar.bz2
  v3.0.1_libenergy-linux-amd64-gtk3-127.zip
  linux_amd64_127.3.5/
    libcef.so
    libenergy-amd64-gtk3.so
```

说明：

| 文件/目录 | 说明 |
|---|---|
| `<os>_<arch>_<cef-version>/` | CEF 版本目录 |
| `.versions` | CEF 文件清单和 libenergy 发行版本清单 |
| `cef_binary_*.tar.bz2` | CEF 下载包缓存 |
| `<runtime-release>_libenergy-*.zip` | libenergy 下载包缓存 |
| `libenergy-*.dll/.so/.dylib` | 解压到对应 CEF 版本目录中的运行时库 |

libenergy 下载包来自 CLI 内置配置 `resources/config.json` 的 `cef_runtime`，用户配置文件中不再持久化顶层 `cef_runtime`。

内置下载模板支持占位符：

| 占位符 | 含义 |
|---|---|
| `{version}` | libenergy 发行版本，例如 `v3.0.1` |
| `{major}` | CEF major，例如 `109`、`127`、`147` |
| `{os}` | 目标系统，例如 `windows`、`linux`、`darwin` |
| `{arch}` | 目标架构，例如 `amd64`、`386`、`arm64` |
| `{ws}` | Linux 窗口系统，通常为 `gtk2` 或 `gtk3` |

非 Linux 平台没有 `{ws}` 时，模板中的 `{ws}` 相关分隔符会自动移除。

SourceForge 下载地址会自动兼容重定向。如果配置为 SourceForge 项目文件地址且缺少 `/download`，CLI 会在实际请求时补齐。

## 环境变量

| 环境变量 | 说明 |
|---|---|
| `ENERGY_TUI` | 设置为 `0`、`false`、`off`、`no` 时禁用 TUI |
| `ENERGY_CEF_RUNTIME_SOURCE` | 临时覆盖 libenergy 下载源名称 |
| `ENERGY_CEF_RUNTIME_URL` | 临时追加 libenergy 下载 URL 模板，多个值用 `,` 或 `;` 分隔，并优先尝试 |
| `ENERGY_WS` | Linux 下选择 `gtk2` 或 `gtk3` libenergy |
| `CI` | 存在时自动使用简单文本交互 |
| `GOOS` | 多平台构建时由 CLI 设置目标系统 |
| `GOARCH` | 多平台构建时由 CLI 设置目标架构 |
| `CGO_ENABLED` | 多平台构建时禁用 CGO |
| `MACOSX_DEPLOYMENT_TARGET` | macOS 构建时可能设置最低系统版本 |
| `CGO_CFLAGS` / `CGO_LDFLAGS` | macOS 或 CGO 构建时使用 |

示例：

```bash
# 临时使用自定义 libenergy 下载模板
ENERGY_CEF_RUNTIME_URL='https://example.com/{version}/libenergy-{os}-{arch}-{ws}-{major}.zip' \
energy env -w chromium.version=127

# Linux 临时使用 GTK2 运行时库
ENERGY_WS=gtk2 energy env -w chromium.version=109
```

## 常用示例

```bash
# 查看帮助
energy help

# 创建 CEF 项目
energy init -name cefdemo -ui CEF -cef-version 127

# 查看本机配置
energy env

# 查看当前 CEF 版本
energy env chromium.version

# 列出已安装 CEF
energy env -l cef

# 切换到 CEF 147
energy env -w chromium.version=147

# 切换 libenergy 下载源
energy env -w chromium.source=sourceforge

# 构建并运行
energy run -path /home/user/cefdemo

# 构建安装包
energy package -path /home/user/cefdemo
```

## 排查

### `energy env -w chromium.version=127` 没有变化

可能原因：

| 原因 | 处理 |
|---|---|
| 当前内置配置没有唯一匹配的 CEF 版本 | 使用完整版本，例如 `127.3.5` |
| 目标版本不在内置支持列表中 | 检查 `resources/config.json` 的 `chromium` 配置 |
| 下载失败 | 检查网络、代理或下载源 |

### 下载进度不显示 TUI 样式

可能处于简单文本模式。检查：

```bash
echo $ENERGY_TUI
echo $CI
echo $TERM
```

如果在非交互终端、CI、Windows 7 或 `ENERGY_TUI=0` 下运行，这是预期行为。

### SourceForge 下载失败

确认内置或自定义 URL 是否指向 SourceForge 文件地址。CLI 会自动补齐 `/download` 并跟随重定向，但仍需要网络环境可以访问 SourceForge。

### 打包工具缺失

不同平台需要安装对应工具：

| 平台 | 常见工具 |
|---|---|
| Windows | NSIS、Windows SDK |
| macOS | Xcode Command Line Tools、`create-dmg` |
| Linux | `dpkg-deb`、`rpmbuild`、`file`、`fakeroot` |
