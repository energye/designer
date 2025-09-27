# designer

#### 介绍
ENERGY GUI 设计器

#### 软件架构
- Golang
- lcl + lib


#### 安装教程
- clone
- golang

#### 使用说明
- go run


### 参考

lazarus 设计器 

### 对象查看器

#### 组件树

#### 属性设置

##### 类型
- edit: 文本
1. string
2. int
3. float
4. bool
- combobox: 下拉框
1. 文本
- sub-tree: 子菜单
1. set 集合
2. class 类
3. checkbox
- checkbox: 复选框
1. checkbox

#### 事件设置

### 文件目录

- resources/component_property.json
```json
设计时组件属性配置, 在加载到设计列表使用
{
  "common": { 通用属性配置
    "exclude": ["Action"], 排除的属性
    "include": [ 包含的属性
      {"name": "","value": "", "kind": "", "type": ""}
    ]
  },
  "custom": { 自定义属性配置
    "TButton": {"name": "","value": "", "kind": "", "type": ""} 组件名 : 组件属性信息
    ... 更多其它组件的自定义属性配置
  }
}
```
- resources/config.json
