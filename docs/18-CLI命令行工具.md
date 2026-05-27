# CLI 命令行工具

ENERGY Designer 提供独立的命令行工具 `energy`，支持在终端中直接执行项目的编译、运行和打包操作，无需打开设计器 GUI。

## 安装

CLI 工具随 ENERGY Designer 源码一起提供。编译方式：

```bash
cd designer/cmd/energy
go build -o energy
```

或将编译后的 `energy` 可执行文件放入系统 PATH 中，即可全局使用。

## 基本用法

```bash
energy <command> [options]
```

不带任何参数运行 `energy` 会显示帮助信息。

## 命令列表

### energy help

显示所有可用命令的帮助信息。

```bash
energy help
```

输出示例：
```
energy command
   run
     energy run, Run the application
   build
     energy build, Build the application binary
   package
     energy package, Build the application installer package
   help
     energy help
```

---

### energy run

编译并运行应用程序。

```bash
energy run [-path <project-path>]
```

| 参数     | 说明                       | 是否必需 | 默认值        |
|--------|--------------------------|------|------------|
| `-path` | 项目根目录路径（包含 `.egp` 文件的目录） | 否    | 当前工作目录     |

**执行流程：**

1. 加载项目配置文件（`.egp`）
2. 执行 `go build` 编译项目
3. macOS 平台自动创建 `.app` 应用包
4. 运行编译后的可执行文件
5. 标准输出实时显示在终端

**示例：**

```bash
# 在当前目录运行项目
energy run

# 指定项目路径运行
energy run -path /home/user/myproject

# 使用 = 号分隔参数
energy run -path=/home/user/myproject
```

---

### energy build

编译应用程序，生成可执行文件。

```bash
energy build [-path <project-path>] [--all]
```

| 参数     | 说明                       | 是否必需 | 默认值        |
|--------|--------------------------|------|------------|
| `-path` | 项目根目录路径（包含 `.egp` 文件的目录） | 否    | 当前工作目录     |
| `--all` | 构建所有平台和架构的可执行文件          | 否    | 仅构建当前平台    |

**编译模式：**

- **Debug 模式** - 包含调试信息，不优化
- **Release 模式** - 去除调试信息（`-s -w`），启用优化（`-trimpath`），Windows 下隐藏控制台窗口（`-H windowsgui`）

**平台特定行为：**

| 平台      | 输出文件格式              | 特殊处理                        |
|---------|---------------------|-----------------------------|
| Windows | `build/<name>.exe`  | Release 模式自动添加 `-H windowsgui` |
| macOS   | `build/<name>.app`  | 自动创建 `.app` 包结构              |
| Linux   | `build/<name>`      | Release 模式自动执行 `strip`        |

**跨平台构建（`--all`）：**

使用 `--all` 参数时，会禁用 CGO 并依次构建以下平台：

| 目标平台    | 目标架构                     |
|---------|--------------------------|
| Windows | amd64, 386               |
| macOS   | amd64, arm64             |
| Linux   | amd64, 386, arm, arm64   |

跨平台构建的输出文件名自动追加平台和架构后缀，如 `myapp_windows_amd64.exe`。

**macOS Universal Binary：**

当项目配置中启用了 `MacCommonLib` 选项时，构建会自动：
1. 分别编译 amd64 和 arm64 架构的临时二进制
2. 使用 `lipo` 工具合并为 Universal Binary
3. 清理临时文件

**示例：**

```bash
# 构建当前平台
energy build

# 构建当前平台（指定路径）
energy build -path /home/user/myproject

# 构建所有平台
energy build --all

# 构建所有平台（指定路径）
energy build -path /home/user/myproject --all
```

---

### energy package

编译并打包应用程序为安装包。

```bash
energy package [-path <project-path>]
```

| 参数     | 说明                       | 是否必需 | 默认值        |
|--------|--------------------------|------|------------|
| `-path` | 项目根目录路径（包含 `.egp` 文件的目录） | 否    | 当前工作目录     |

**执行流程：**

1. 加载项目配置文件（`.egp`）
2. 强制启用 Release 编译模式
3. 执行 `go build` 编译项目
4. 根据目标平台执行对应的打包流程
5. 生成安装包文件

**各平台打包产物：**

#### macOS

| 格式       | 生成工具        | 前置要求                      |
|----------|-------------|---------------------------|
| `.app`   | 内置          | 无                         |
| `.pkg`   | `pkgbuild`  | 系统内置                     |
| `.dmg`   | `create-dmg` | 需安装：`brew install create-dmg` |

macOS 打包流程：
1. 创建 `.app` 包结构（`Contents/MacOS`、`Contents/Frameworks`、`Contents/Resources`）
2. 复制可执行文件、运行时库（`libenergy.dylib`）、图标（`.icns`）
3. 生成 `Info.plist` 配置文件
4. 复制本地化资源（`.lproj`）
5. 如启用签名，执行 `codesign` 代码签名
6. 如启用 PKG，使用 `pkgbuild` 生成 `.pkg` 安装包
7. 如启用 DMG，使用 `create-dmg` 生成 `.dmg` 磁盘映像

#### Windows

| 格式        | 生成工具              | 前置要求                    |
|-----------|-------------------|-------------------------|
| `.exe` 安装包 | NSIS              | 需安装 NSIS                |
| AppX      | `makeappx.exe`    | 需安装 Windows SDK         |

Windows 打包流程：
1. 编译生成可执行文件
2. 如启用签名，使用 `signtool.exe` 签名
3. 如启用 NSEXE，生成 NSIS 脚本并编译为安装程序
4. 如启用 AppX，使用 `makeappx.exe` 打包

签名支持两种模式：
- `auto` - 自动选择证书
- `file` - 指定证书文件路径和密码

#### Linux

| 格式         | 生成工具       | 前置要求                |
|------------|------------|---------------------|
| `.deb`     | `dpkg-deb` | 需安装 `dpkg-dev`      |
| `.rpm`     | `rpmbuild` | 需安装 `rpm-build`     |
| `.AppImage` | 内置         | 需安装 `file`、`fakeroot` |

Linux 打包流程：
1. 编译生成可执行文件
2. 自动生成 `.desktop` 桌面快捷方式文件
3. 根据 GUI 渲染框架自动选择 GTK2/GTK3 运行时库（`libenergy.so`）
4. 自动生成依赖列表（GTK、WebKit2GTK 等）
5. 按目标格式生成安装包

**Linux 运行时库选择规则：**

| GUI 渲染框架 | UI 配置  | 运行时库选择 |
|-----------|--------|---------|
| LCL       | GTK3   | GTK3    |
| LCL       | GTK2   | GTK2    |
| LCL + WV  | -      | GTK3    |
| LCL + CEF | -      | GTK3    |

**示例：**

```bash
# 打包当前平台
energy package

# 指定项目路径打包
energy package -path /home/user/myproject
```

---

## 项目配置文件

CLI 工具通过读取项目根目录下的 `.egp` 文件获取项目配置。`.egp` 是 JSON 格式文件，包含：

- 项目基本信息（名称、版本等）
- 应用配置（图标、元数据）
- 构建配置（目标平台、架构、输出目录等）
- 打包配置（各平台安装包选项、签名配置等）

CLI 支持两种方式指定项目路径：
1. 使用 `-path` 参数显式指定
2. 在项目根目录下直接运行命令（自动检测当前目录）

如果 `-path` 指向的是目录而非 `.egp` 文件，CLI 会自动在该目录下搜索 `.egp` 文件。

## 环境变量

CLI 在编译过程中会自动设置以下环境变量：

| 环境变量                      | 说明               | 适用平台          |
|---------------------------|------------------|---------------|
| `GOOS`                    | 目标操作系统           | 跨平台构建         |
| `GOARCH`                  | 目标 CPU 架构        | 跨平台构建         |
| `CGO_ENABLED`             | 是否启用 CGO         | 跨平台构建/Linux   |
| `MACOSX_DEPLOYMENT_TARGET` | macOS 最低部署版本     | macOS          |
| `CGO_CFLAGS`              | CGO 编译标志         | macOS          |
| `CGO_LDFLAGS`             | CGO 链接标志         | macOS          |
| `ENERGY_WS`               | 窗口系统类型（gtk2/gtk3） | Linux          |

## 编译参数

CLI 会从项目的构建配置中读取以下参数并传递给 `go build`：

| 参数          | 说明                     |
|-------------|------------------------|
| `-v`        | 显示编译详细信息（始终启用）        |
| `-tags`     | 构建标签（如 `prod`、`dev`）   |
| `-ldflags`  | 链接器标志（如 `-s -w`）       |
| `-trimpath` | 去除编译路径信息（Release 模式）   |
| `-o`        | 输出文件路径                 |

自定义编译参数可通过项目构建配置中的 `GoArgs` 字段指定，CLI 会自动提取 `-tags` 和 `-ldflags` 参数并合并。
