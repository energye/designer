# ENERGY GUI Designer

![go-version](https://img.shields.io/github/go-mod/go-version/energye/designer?logo=git&logoColor=green)
[![github](https://img.shields.io/github/last-commit/energye/energy/main.svg?logo=github&logoColor=green&label=commit)](https://github.com/energye/designer)
![repo](https://img.shields.io/github/repo-size/energye/designer.svg?logo=github&logoColor=green&label=repo-size)

---

## Project Overview

ENERGY Designer is a visual designer specifically built for the ENERGY cross-platform GUI framework, developed using the Go LCL component library. It provides a what-you-see-is-what-you-get (WYSIWYG) design experience with simplified GUI design capabilities. Developers can intuitively create and edit GUI interfaces by dragging controls and setting properties, while automatically generating corresponding Go source code.

## Core Features

### Designer Functionality

- Visual Design: WYSIWYG GUI interface design experience
- Drag-and-Drop Controls: Support for standard control drag-and-drop layout
- Property Editing: Real-time property and event configuration
- Component Management: Complete component tree viewing and management
- Preview and Run: Preview design effects
- Code Generation: Automatically convert UI designs to Go source code
- Multi-Form Support: Support for multi-tab form design
- Undo/Redo: Support for undoing and redoing design operations

### Supported Component Types

- LCL Native Controls: Standard desktop application controls
- CEF Controls: Chromium Embedded Framework browser controls
- WebView Controls: System-native Webview controls (Windows WebView2, Linux WebKit2, macOS WKWebView)

### Cross-Platform Support

- Windows: Supports 386 and amd64 architectures
- macOS: Supports Intel (amd64) and Apple Silicon (arm64) architectures
- Linux: Supports amd64, 386, arm64, arm architectures

### Use Cases

1. Rapid Prototyping: Quickly build GUI application prototypes
2. Interface Design: Visually design complex user interfaces
3. Code Generation: Automatically generate Go code compliant with ENERGY standards
4. Learning Tool: Help beginners understand GUI programming concepts

*ENERGY Designer - Making GUI Development Easier*

### Installing Designer

#### Requirements

- Install Git
- Install Golang
- Execute the following commands

```cmd
Clone the repository
git clone https://github.com/energye/designer.git
Enter the designer directory and update module dependencies
go mod tidy
```

- Run from the designer directory

`go run main.go`

If you have any requirements, please submit an issue or join the QQ group for discussion (541258627).

### Screenshots

![ENERGY-designer-create-project.png](docs/image/ENERGY-designer-create-project.png)
![ENERGY-designer-home.png](docs/image/ENERGY-designer-home.png)
![ENERGY-designer-widget.png](docs/image/ENERGY-designer-widget.png)
![ENERGY-designer-run.png](docs/image/ENERGY-designer-run.png)
![ENERGY-designer-config-app.png](docs/image/ENERGY-designer-config-app.png)
![ENERGY-designer-config-build-1.png](docs/image/ENERGY-designer-config-build-1.png)
![ENERGY-designer-config-build-2.png](docs/image/ENERGY-designer-config-build-2.png)
- menu

![ENERGY-designer-menu-run.png](docs/image/ENERGY-designer-menu-run.png)
![ENERGY-designer-menu-options.png](docs/image/ENERGY-designer-menu-options.png)
![ENERGY-designer-menu-file.png](docs/image/ENERGY-designer-menu-file.png)
- main.go

![ENERGY-project-main.go.png](docs/image/ENERGY-project-main.go.png)
- form1.go

![ENERGY-project-form1.go.png](docs/image/ENERGY-project-form1.go.png)