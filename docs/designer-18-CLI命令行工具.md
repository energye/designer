# CLI 命令行工具

ENERGY Designer 提供命令行工具 `energy`，可在终端中创建项目、运行项目、构建程序、生成安装包，以及管理本机配置和 CEF 版本。

## 安装

在源码目录编译：

```bash
cd designer/cmd/energy
go install
```

## 基本用法

```bash
energy <command> [options]
```

不带参数运行 `energy` 会显示帮助信息。

参数支持以下写法：

```bash
energy run -path /home/user/demo
energy run -path=/home/user/demo
energy build --all
```

常用命令：

| 命令 | 说明 |
|---|---|
| `energy help` | 显示帮助 |
| `energy init` | 创建项目 |
| `energy run` | 编译并运行项目 |
| `energy build` | 编译项目 |
| `energy package` | 生成安装包 |
| `energy env` | 查看或修改本机配置 |

## 交互模式

部分命令在参数缺失时会进入交互模式，例如：

- 创建项目时输入项目名称
- 选择 UI 框架
- 选择 CEF 版本
- 确认是否覆盖文件
- 显示下载进度

如需使用简单文本模式，可设置：

```bash
ENERGY_TUI=0 energy init
```

在 CI、非交互终端或旧版 Windows 系统中，也会自动使用简单文本模式。

## energy help

显示命令帮助：

```bash
energy help
```

## energy init

创建 ENERGY 项目。

```bash
energy init [-name <project-name>] [-path <project-dir>] [-ui <LCL|WV|CEF>]
```

参数：

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-name` | 项目名称 | 未指定时交互输入 |
| `-path` | 项目目录 | 当前目录 |
| `-ui` | UI 框架，可选 `LCL`、`WV`、`CEF` | 未指定时交互选择 |
| `-framework` | `-ui` 的别名 | 无 |

创建 CEF 项目时还可使用：

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-cef-version` | CEF 版本，支持 `109`、`127`、`147` 这类 major 写法，也支持完整版本 | 未指定时交互选择 |
| `-os` | CEF 目标系统：`windows`、`linux`、`darwin` | 当前系统 |
| `-arch` | CEF 目标架构：`amd64`、`386`、`arm64`、`arm` | 当前架构 |
| `-cef-dir` | CEF 安装目录 | `~/.energy/chromium` |

示例：

```bash
# 交互式创建项目
energy init

# 创建 LCL 项目
energy init -name demo -path /home/user/demo -ui LCL

# 创建 WebView 项目
energy init -name webdemo -ui WV

# 创建 CEF 项目，并使用 CEF 127
energy init -name cefdemo -ui CEF -cef-version 127

# 指定 CEF 目录
energy init -name cefdemo -ui CEF -cef-version 127 -cef-dir /opt/energy/chromium
```

## energy run

编译并运行项目。

```bash
energy run [-path <project-dir>]
```

参数：

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-path` | 项目目录，目录内应包含 `.egp` 项目文件 | 当前目录 |

示例：

```bash
# 在当前项目目录运行
energy run

# 指定项目目录运行
energy run -path /home/user/demo
```

## energy build

编译项目并生成可执行文件。

```bash
energy build [-path <project-dir>] [--all]
```

参数：

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-path` | 项目目录 | 当前目录 |
| `--all` | 按项目构建配置生成多平台产物 | 只构建当前平台 |

输出位置由项目配置中的 `build_option.output` 决定，默认通常是项目目录下的 `build`。

示例：

```bash
# 构建当前平台
energy build

# 指定项目目录构建
energy build -path /home/user/demo

# 构建多平台产物
energy build --all
```

多平台构建需要项目配置中允许构建其他平台，并且适合非 CGO 构建。输出文件名会自动带上系统和架构后缀，例如：

```text
demo_windows_amd64.exe
demo_linux_amd64
demo_darwin_arm64
```

## energy package

编译并生成安装包。

```bash
energy package [-path <project-dir>]
```

参数：

| 参数 | 说明 | 默认值 |
|---|---|---|
| `-path` | 项目目录 | 当前目录 |

示例：

```bash
energy package
energy package -path /home/user/demo
```

支持的打包类型由项目配置决定：

| 平台 | 常见产物 |
|---|---|
| Windows | NSIS `.exe`、MSIX/AppX |
| macOS | `.app`、`.pkg`、`.dmg` |
| Linux | `.deb`、`.rpm`、`.AppImage` |

常见外部工具要求：

| 平台 | 常见工具 |
|---|---|
| Windows | NSIS、Windows SDK |
| macOS | Xcode Command Line Tools、`create-dmg` |
| Linux | `dpkg-deb`、`rpmbuild`、`file`、`fakeroot` |

## energy env

查看或修改本机配置文件：

```text
~/.energy/config.json
```

基本用法：

```bash
energy env
energy env <key|json.path>
energy env -w <key|json.path>=<value>
energy env -l <module>
```

### 查看配置

```bash
# 查看完整配置
energy env

# 查看 Chromium/CEF 配置
energy env chromium

# 查看当前 CEF 版本
energy env chromium.version

# 查看所有名称为 version 的配置项
energy env version
```

### 修改配置

```bash
# 修改 CEF 根目录
energy env -w chromium.dir=/home/user/.energy/chromium

# 修改代理
energy env -w proxy=http://127.0.0.1:7890

# 修改下载源
energy env -w chromium.source=sourceforge
```

写入布尔值和数字时直接使用文本值：

```bash
energy env -w some_bool=true
energy env -w some_number=123
```

如果一个名称匹配到多个配置项，CLI 会提示选择要修改的配置项。建议优先使用完整 JSON 路径，例如 `chromium.version`。

### 列出已安装 CEF

```bash
energy env -l cef
```

当前使用的版本前会显示 `*`：

```text
  linux_amd64_109.1.18
* linux_amd64_127.3.5
  linux_amd64_147.0.14
```

### 安装或切换 CEF

使用 `chromium.version` 切换 CEF 版本：

```bash
# 使用 major
energy env -w chromium.version=127

# 使用完整版本
energy env -w chromium.version=127.3.5

# 使用完整平台版本
energy env -w chromium.version=linux_amd64_127.3.5
```

说明：

- 如果目标版本未安装，会自动下载。
- 如果目标版本已安装，会直接切换。
- CEF 项目需要的 libenergy 运行时库会随 CEF 版本一起准备。
- `chromium.version` 最终保存格式为 `<os>_<arch>_<cef-version>`。

示例：

```text
windows_amd64_127.3.5
linux_amd64_127.3.5
darwin_arm64_147.0.14
```

### 切换下载源

切换 libenergy 下载源：

```bash
energy env -w chromium.source=sourceforge
```

临时指定下载源：

```bash
ENERGY_CEF_RUNTIME_SOURCE=sourceforge energy env -w chromium.version=127
```

临时指定自定义下载地址模板：

```bash
ENERGY_CEF_RUNTIME_URL='https://example.com/{version}/libenergy-{os}-{arch}-{ws}-{major}.zip' \
energy env -w chromium.version=127
```

下载地址模板支持：

| 占位符 | 说明 |
|---|---|
| `{version}` | libenergy 发行版本 |
| `{major}` | CEF major，例如 `127` |
| `{os}` | 系统，例如 `windows`、`linux`、`darwin` |
| `{arch}` | 架构，例如 `amd64`、`arm64` |
| `{ws}` | Linux 使用，通常为 `gtk2` 或 `gtk3` |

Linux 下可临时指定 GTK 运行时：

```bash
ENERGY_WS=gtk2 energy env -w chromium.version=109
ENERGY_WS=gtk3 energy env -w chromium.version=127
```

## 常用环境变量

| 环境变量 | 说明 |
|---|---|
| `ENERGY_TUI` | 设置为 `0`、`false`、`off`、`no` 时禁用交互式 TUI |
| `ENERGY_CEF_RUNTIME_SOURCE` | 临时指定 libenergy 下载源 |
| `ENERGY_CEF_RUNTIME_URL` | 临时指定 libenergy 下载地址模板，多个地址可用 `,` 或 `;` 分隔 |
| `ENERGY_WS` | Linux 下指定 `gtk2` 或 `gtk3` |
| `CI` | CI 环境中自动使用简单文本交互 |

## 常用示例

```bash
# 查看帮助
energy help

# 创建项目
energy init -name demo -ui LCL

# 创建 CEF 项目
energy init -name cefdemo -ui CEF -cef-version 127

# 运行项目
energy run -path /home/user/demo

# 构建项目
energy build -path /home/user/demo

# 打包项目
energy package -path /home/user/demo

# 查看配置
energy env

# 查看当前 CEF
energy env chromium.version

# 列出已安装 CEF
energy env -l cef

# 切换到 CEF 147
energy env -w chromium.version=147

# 设置 libenergy 下载源
energy env -w chromium.source=sourceforge
```

## 常见问题

### CEF 下载失败

检查网络、代理和下载源：

```bash
energy env proxy
energy env chromium.source
```

需要代理时：

```bash
energy env -w proxy=http://127.0.0.1:7890
```

### 切换 CEF 后版本没有变化

建议使用完整版本或完整平台版本：

```bash
energy env -w chromium.version=127.3.5
energy env -w chromium.version=linux_amd64_127.3.5
```

也可以先查看本机已安装版本：

```bash
energy env -l cef
```

### 交互界面显示异常

可切换为简单文本模式：

```bash
ENERGY_TUI=0 energy init
```

### 打包失败提示工具不存在

按当前平台安装对应工具：

| 平台 | 工具 |
|---|---|
| Windows | NSIS、Windows SDK |
| macOS | Xcode Command Line Tools、`create-dmg` |
| Linux | `dpkg-deb`、`rpmbuild`、`file`、`fakeroot` |
