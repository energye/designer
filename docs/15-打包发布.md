# 打包发布

ENERGY Designer 支持多平台安装包制作，可将编译后的应用程序打包为各平台标准的安装格式。通过 **设置 > 构建选项**（快捷键 `Ctrl+F9` / macOS `Cmd+F9`）切换到"打包配置"标签页进行设置。

## 通用配置

| 配置项 | 说明 | 输入类型 | 默认值 |
|--------|------|----------|--------|
| 安装包名称 | 生成的安装包文件名 | 文本输入 | 项目名称 |

## macOS 打包

### 打包格式

| 格式 | 说明 | 默认值 |
|------|------|--------|
| PKG 安装包 | macOS 标准安装包格式 | 开启 |
| DMG 镜像 | macOS 磁盘映像文件 | 关闭 |
| 通用二进制 (Universal Binary) | 同时支持 Intel (amd64) 和 Apple Silicon (arm64) 架构 | 关闭 |

**通用二进制说明：** 启用后，构建会自动分别编译 amd64 和 arm64 架构的临时二进制，然后使用 `lipo` 工具合并为 Universal Binary。需安装 Xcode Command Line Tools。

### 签名配置

签名功能通过启用/禁用按钮控制开关。配置内容为 `codesign` 命令列表，每行一条命令。

**默认签名命令：**
```
codesign -f -s "-" "$APP_NAME/Contents/Frameworks/$ENERGY.DYLIB"
codesign -f -s "-" --options runtime "$APP_NAME"
```

**说明：**
- `-s "-"` 表示 ad-hoc 签名（预览运行时使用）
- 发布时可替换为开发者证书签名
- `--options runtime` 启用 hardened runtime（公证要求）

### 关联文件

配置应用程序关联的文件类型，每行一条记录，使用 `|` 分隔字段：

```
扩展名 | 类名 | 类型描述 | 角色 | 优先级 | 图标 | MIME类型
```

**字段说明：**

| 字段 | 说明 | 可选值 |
|------|------|--------|
| 扩展名 | 文件扩展名（不含点号） | 如 `txt`、`pdf` |
| 类名 | 唯一的文件类型标识符 | 如 `AppTxtFile` |
| 类型描述 | 文件类型的描述文字 | 如 `Text Document` |
| 角色 | 应用对该文件类型的角色 | `Editor`（编辑器）/ `Viewer`（查看器） |
| 优先级 | 关联优先级 | `Owner`（所有者）/ `Default`（默认） |
| 图标 | 关联图标文件名 | 如 `MyIcon.icns` |
| MIME 类型 | MIME 类型标识 | 如 `application/x-gproj` |

**示例：**
```
txt | AppTxtFile | Text Document | Editor | Owner | MyIcon.icns | text/plain
```

### 关联协议

配置应用程序关联的 URL 协议，每行一条记录，使用 `|` 分隔：

```
协议名称 | 协议描述
```

**示例：**
```
myapp | Open My App
```

## Windows 打包

### 打包格式

| 格式 | 说明 | 默认值 |
|------|------|--------|
| EXE 安装包 (MakeNsis) | 使用 NSIS 脚本生成 Windows 安装程序 | 开启 |
| MSIX 安装包 (MakeAppx) | Windows 应用商店格式 | 关闭 |

### 默认安装路径

| 配置项 | 说明 | 输入类型 | 示例 |
|--------|------|----------|------|
| 默认安装路径 | 应用程序的默认安装目录 | 文本输入 | `C:\Program Files` |

### 签名配置

签名功能通过启用/禁用按钮控制开关。配置内容为 `signtool` 命令，格式为 `key=value`，每行一条。

**签名模式：**

| 模式 | 格式 | 说明 |
|------|------|------|
| auto | `auto=signtool sign /a /fd SHA256 /tr http://timestamp.digicert.com /td SHA256` | 自动选择证书 |
| file | `file=signtool sign /f cert.pfx /p 密码 /fd SHA256 /tr http://timestamp.digicert.com /td SHA256` | 指定证书文件路径和密码 |

**说明：**
- `/a` - 自动选择最佳证书
- `/fd SHA256` - 文件摘要算法
- `/tr` - 时间戳服务器 URL
- `/td SHA256` - 时间戳摘要算法
- `/f` - 证书文件路径
- `/p` - 证书密码

### 关联文件

配置应用程序关联的文件类型，每行一条记录，使用 `|` 分隔：

```
扩展名 | 唯一类名 | 类型描述 | 图标 | 右键菜单文本
```

**字段说明：**

| 字段 | 说明 | 示例 |
|------|------|------|
| 扩展名 | 文件扩展名（不含点号） | `txt` |
| 唯一类名 | 注册表中的文件类型标识 | `AppTxtFile` |
| 类型描述 | 文件类型的描述文字 | `My Project File` |
| 图标 | 关联图标文件名 | `MyFile.ico` |
| 右键菜单文本 | 右键菜单中显示的操作文本 | `Open with Your App` |

**示例：**
```
txt | AppTxtFile | My Project File | MyFile.ico | Open with Your App
```

### 关联协议

配置应用程序关联的 URL 协议，每行一条记录，使用 `|` 分隔：

```
协议名称 | 协议描述
```

**示例：**
```
myapp | Open My App
```

### NSIS 安装包资源

配置 NSIS 安装程序的视觉资源，格式为 `key=value`，每行一条：

| 键名 | 说明 | 支持格式 |
|------|------|----------|
| welcome | 欢迎页面横幅图片 | `.png` / `.bmp` |
| header | 安装程序头部横幅图片 | `.png` / `.bmp` |
| icon | 安装程序图标 | `.png` / `.ico` |
| unicon | 卸载程序图标 | `.png` / `.ico` |

**示例：**
```
welcome=assets/welcome.png
header=assets/header.bmp
icon=assets/app.ico
unicon=assets/uninstall.ico
```

### NSIS 许可证

在文本框中输入许可证协议内容，打包时自动生成 UTF-8 BOM 编码的 `resources/<项目名>-license.txt` 文件并嵌入安装程序。

### AppX 资源

配置 MSIX/AppX 打包所需的资源文件，格式为 `key=value`，每行一条：

| 键名 | 说明 |
|------|------|
| propertiesLogo | 属性页面 Logo |
| square44x44Logo | 44×44 方形图标（任务栏等） |
| square150x150Logo | 150×150 方形图标（开始菜单） |
| wide310x150Logo | 310×150 宽幅图标（磁贴） |
| splashScreen | 启动画面图片 |
| associateFileIcon | 文件关联图标 |
| associateProtocolLogo | 协议关联 Logo |

**示例：**
```
propertiesLogo=assets/StoreLogo.png
square44x44Logo=assets/Square44x44Logo.png
square150x150Logo=assets/Square150x150Logo.png
wide310x150Logo=assets/Wide310x150Logo.png
splashScreen=assets/SplashScreen.png
```

## Linux 打包

### 打包格式

| 格式 | 说明 | 默认值 |
|------|------|--------|
| deb | Debian/Ubuntu 安装包 | 开启 |
| rpm | Red Hat/CentOS/Fedora 安装包 | 关闭 |
| AppImage | 通用 Linux 应用格式（免安装） | 关闭 |

### 配置参数

| 配置项 | 说明 | 输入类型 | 默认值 | 示例 |
|--------|------|----------|--------|------|
| 依赖包列表 | 应用运行所需的系统依赖包 | 文本输入 | 空 | `libc6 (>= 2.17), libgtk-3-0` |
| 应用分类 | 应用在桌面环境中的分类 | 文本输入 | `Utility;` | `Development;Utility;` |
| 项目主页 | 应用的官方网站地址 | 文本输入 | 空 | `https://example.com` |
| 维护者 | 安装包维护者信息 | 文本输入 | 空 | `Name <email@example.com>` |
| 开源协议 | 应用的开源许可证 | 文本输入 | `MIT` | `GPL-3.0` |

### 依赖包格式规则

| 格式 | 说明 | 示例 |
|------|------|------|
| 基本格式 | 包名 + 版本约束 | `libc6 (>= 2.17)` |
| 多个依赖 | 逗号分隔 | `libc6 (>= 2.17), libgtk-3-0` |
| 备选依赖 | 管道符分隔（RPM 仅取第一个） | `libssl1.1 \| libssl3` |

**运行时库自动选择规则：**

| GUI 渲染框架 | UI 配置 | 运行时库选择 |
|--------------|---------|-------------|
| LCL | GTK2 | GTK2 |
| LCL | GTK3 | GTK3 |
| LCL + WV | - | GTK3 |
| LCL + CEF | - | GTK3 |

## 打包产物汇总

| 平台 | 格式 | 生成工具 | 前置要求 |
|------|------|----------|----------|
| macOS | `.app` | 内置 | 无 |
| macOS | `.pkg` | `pkgbuild` | 系统内置 |
| macOS | `.dmg` | `create-dmg` | `brew install create-dmg` |
| Windows | `.exe` 安装包 | NSIS | 需安装 NSIS |
| Windows | AppX | `makeappx.exe` | 需安装 Windows SDK |
| Linux | `.deb` | `dpkg-deb` | 需安装 `dpkg-dev` |
| Linux | `.rpm` | `rpmbuild` | 需安装 `rpm-build` |
| Linux | `.AppImage` | 内置 | 需安装 `file`、`fakeroot` |

## 资源嵌入

打包过程中自动处理以下资源：

| 资源类型 | 说明 |
|----------|------|
| 应用图标 | 自动转换为各平台格式（`.ico` / `.icns` / `.png`） |
| 版本信息 | 自动嵌入到可执行文件中 |
| Manifest 文件 | Windows 下自动生成并嵌入 |
| 框架运行时 | 自动打包必要的框架库文件（`libenergy.dylib` / `libenergy.so`） |
