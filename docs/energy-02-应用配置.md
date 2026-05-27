# 应用配置

应用配置通过 `application.Options` 结构体设置，控制窗口外观、行为和平台特定选项。

## 基本用法

```go
wvApp := wv.Init()
wvApp.SetOptions(application.Options{
    Caption:    "My App",
    DefaultURL: "app://custom/index.html",
    Width:      1200,
    Height:     800,
})
```

## Options 完整字段

### 窗口基本属性

| 字段 | 类型 | 说明 |
|------|------|------|
| Caption | string | 窗口标题 |
| DefaultURL | string | 默认加载的 URL |
| X | int32 | 窗口 X 坐标 |
| Y | int32 | 窗口 Y 坐标 |
| Width | int32 | 窗口宽度 |
| Height | int32 | 窗口高度 |
| MinWidth | int32 | 最小宽度（0 = 无限制） |
| MinHeight | int32 | 最小高度（0 = 无限制） |
| MaxWidth | int32 | 最大宽度（0 = 无限制） |
| MaxHeight | int32 | 最大高度（0 = 无限制） |
| DefaultWindowStatus | types.TWindowState | 默认窗口状态 |

### 功能开关

| 字段 | 类型 | 说明 |
|------|------|------|
| DisableDevTools | bool | 禁用开发者工具 |
| DisableContextMenu | bool | 禁用右键菜单 |
| DisableWebkitAppRegionDClk | bool | 禁用 Webkit AppRegion 双击 |
| DisableResize | bool | 禁用窗口调整大小 |
| DisableMinimize | bool | 禁用最小化按钮 |
| DisableMaximize | bool | 禁用最大化按钮 |
| DisableSystemMenu | bool | 禁用系统菜单 |

### 外观

| 字段 | 类型 | 说明 |
|------|------|------|
| Frameless | bool | 无边框窗口 |
| WindowTransparent | bool | 窗口透明 |
| WebviewTransparent | bool | WebView 背景透明 |
| BackgroundColor | *colors.TARGB | 背景颜色（使用 `colors.NewARGB()` 创建） |

### 平台配置

| 字段 | 类型 | 说明 |
|------|------|------|
| Windows | Windows | Windows 平台配置 |
| MacOS | MacOS | macOS 平台配置 |
| Linux | Linux | Linux 平台配置 |

## Windows 平台配置

```go
Windows: application.Windows{
    Theme:           application.Dark,
    BackdropType:    application.BtAcrylic,
    WindowProtected: true,
}
```

### Theme 主题

| 值 | 说明 |
|-----|------|
| SystemDefault | 跟随系统（默认） |
| Dark | 深色主题 |
| Light | 浅色主题 |

### BackdropType 背景材质

| 值 | 说明 |
|-----|------|
| BtAuto | 自动选择 |
| BtNone | 无材质 |
| BtMica | Mica 材质（Win11） |
| BtAcrylic | 亚克力材质（Win10+） |
| BtTabbed | 标签页材质 |

### ThemeSetting 主题颜色

```go
Windows: application.Windows{
    ThemeSetting: &application.ThemeSetting{
        DarkTitleBar:           0x000000,
        DarkTitleBarInactive:   0x2B2B2B,
        DarkTitleText:          0xFFFFFF,
        DarkTitleTextInactive:  0x999999,
        DarkBorder:             0x333333,
        DarkBorderInactive:     0x222222,
        LightTitleBar:          0xFFFFFF,
        LightTitleBarInactive:  0xF0F0F0,
        LightTitleText:         0x000000,
        LightTitleTextInactive: 0x999999,
        LightBorder:            0xCCCCCC,
        LightBorderInactive:    0xE0E0E0,
    },
}
```

ThemeSetting 字段（均为 int32 颜色值）：

| 字段 | 说明 |
|------|------|
| DarkTitleBar | 深色主题标题栏颜色 |
| DarkTitleBarInactive | 深色主题非活动标题栏颜色 |
| DarkTitleText | 深色主题标题文字颜色 |
| DarkTitleTextInactive | 深色主题非活动标题文字颜色 |
| DarkBorder | 深色主题边框颜色 |
| DarkBorderInactive | 深色主题非活动边框颜色 |
| LightTitleBar | 浅色主题标题栏颜色 |
| LightTitleBarInactive | 浅色主题非活动标题栏颜色 |
| LightTitleText | 浅色主题标题文字颜色 |
| LightTitleTextInactive | 浅色主题非活动标题文字颜色 |
| LightBorder | 浅色主题边框颜色 |
| LightBorderInactive | 浅色主题非活动边框颜色 |

### WindowProtected 窗口保护

启用后窗口内容受系统保护（如截图保护）。

## macOS 平台配置

```go
MacOS: application.MacOS{
    AppearanceName:    application.NSAppearanceNameDarkAqua,
    FullSizeContent:   true,
    TitleTransparent:  true,
    TitleHideText:     true,
    WindowRadius:      8,
    UseWindowDelegate: true,
    ToolBar:           &application.ToolBar{},
}
```

### AppearanceName 外观

| 值 | 说明 |
|-----|------|
| NSAppearanceNameAqua | 标准浅色系统外观 |
| NSAppearanceNameDarkAqua | 标准深色系统外观 |
| NSAppearanceNameVibrantLight | 浅色生动外观 |
| NSAppearanceNameAccessibilityHighContrastAqua | 高对比度浅色外观 |
| NSAppearanceNameAccessibilityHighContrastDarkAqua | 高对比度深色外观 |
| NSAppearanceNameAccessibilityHighContrastVibrantLight | 高对比度浅色生动外观 |
| NSAppearanceNameAccessibilityHighContrastVibrantDark | 高对比度深色生动外观 |

### 其他 macOS 选项

| 字段 | 类型 | 说明 |
|------|------|------|
| WindowRadius | float32 | 窗口圆角半径（Frameless=true 时有效） |
| FullSizeContent | bool | 全尺寸内容（延伸到标题栏） |
| TitleTransparent | bool | 标题栏透明 |
| TitleHideText | bool | 隐藏标题文字 |
| UseWindowDelegate | bool | 使用窗口委托（用于窗口事件回调） |
| ToolBar | *ToolBar | 工具栏配置 |

## Linux 平台配置

```go
Linux: application.Linux{
    HardwareGPU: application.HGPUEnable,
}
```

### HardwareGPU 硬件加速

| 值 | 说明 |
|-----|------|
| HGPUDefault | 默认设置 |
| HGPUEnable | 启用硬件 GPU |
| HGPUDisable | 禁用硬件 GPU（虚拟机推荐） |

## 完整示例

```go
wvApp := wv.Init()
wvApp.SetOptions(application.Options{
    Caption:            "ENERGY App",
    DefaultURL:         "fs://energy/index.html",
    Width:              1200,
    Height:             800,
    Frameless:          true,
    WindowTransparent:  true,
    WebviewTransparent: true,
    BackgroundColor:    colors.NewARGB(0, 0, 0, 0),
    DisableDevTools:    false,
    DisableContextMenu: false,
    DisableResize:      false,
    Windows: application.Windows{
        Theme:           application.Dark,
        BackdropType:    application.BtAcrylic,
        WindowProtected: true,
    },
    MacOS: application.MacOS{
        AppearanceName:    application.NSAppearanceNameDarkAqua,
        FullSizeContent:   true,
        TitleTransparent:  true,
        TitleHideText:     true,
        WindowRadius:      8,
        UseWindowDelegate: true,
        ToolBar:           &application.ToolBar{},
    },
    Linux: application.Linux{
        HardwareGPU: application.HGPUEnable,
    },
})
```

## 本地资源加载

通过 `SetLocalLoad` 配置本地资源加载，详见 [本地资源加载](energy-06-本地资源加载.md)。

```go
wvApp.SetLocalLoad(application.LocalLoad{
    Scheme:     "app",
    Domain:     "custom",
    ResRootDir: "resources",
    FS:         resources, // embed.FS
})
```
