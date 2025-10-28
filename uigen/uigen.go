package uigen

// UI 布局文件生成, JSON 格式, 自动时时生成
// 依赖 designer 生成 JSON UI文件(form[n].ui)
// 只存放被修改的组件属性
// xxx.ui 文件内容是 tree JSON 结构, 数据格式为组件[变更的属性列表]
// 生成触发条件: 即时触发 防抖
