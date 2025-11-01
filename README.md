# ENERGY GUI Designer


ENERGY Designer 是一个基于 ENERGY GUI 框架开发的跨平台 GUI 可视化设计器，专为 ENERGY GUI 框架打造。它提供了所见即所得的设计体验，让开发者能够快速创建和编辑 GUI 界面。

## 🌟 项目简介

本项目是 ENERGY 跨平台 GUI 框架的官方可视化设计器，采用 Go LCL 组件库实现，提供了简化的 GUI 设计功能。通过该设计器，开发者可以直观地拖拽控件、设置属性，并自动生成对应的 Go 代码。

## 🚀 核心特性

### 🔧 设计器功能
- **可视化设计**：所见即所得的 GUI 界面设计体验
- **控件拖拽**：支持标准控件的拖拽布局
- **属性编辑**：实时属性和事件设置
- **组件管理**：完整的组件树查看和管理
- **预览运行**：实时预览设计效果

### 📦 项目管理
- **项目文件**：`.egp`(ENERGY GUI Project) 格式的项目配置文件
- **UI 布局**：生成和恢复 UI 布局文件
- **代码生成**：根据设计自动生成 Go 代码

### 🎨 控件支持
- **标准控件**：Button、Label、Edit、CheckBox 等
- **附加控件**：Image、Panel、GroupBox 等
- **布局控件**：Splitter、ScrollBox 等


## 🛠 技术架构

### 核心组件
- **Go LCL**：基于 Lazarus Component Library 的 Go 绑定
- **Virtual TreeView**：高性能的树形控件用于属性编辑
- **模板引擎**：代码生成模板系统


## 🎯 使用场景

1. **快速原型开发**：快速构建 GUI 应用原型
2. **界面设计**：可视化设计复杂的用户界面
3. **代码生成**：自动生 ENERGY GUI 的代码


## 📖 文档说明

详细的功能设计和使用说明请参考：
- [功能设计文档](docs/功能设计.md)
- [待完成功能](docs/待完成功能.md)


*ENERGY Designer - 让 GUI 开发更简单*

![designer_window.png](docs%2Fimage%2Fdesigner_window.png)