// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package event

// 数据载体类型
type Type int32

// 项目管理模块, 数据载体类型
const (
	ProjectCreate     Type = iota // 项目创建
	ProjectLoad                   // 项目加载
	ProjectUpdateForm             // 项目配置更新 Form 信息
	ProjectConfig                 // 项目(应用)配置
	EnvConfig                     // 项目(环境)配置
	BuildConfig                   // 项目(构建)配置
)

// 代码生成模块, 数据载体类型
const (
	CodeGenUI    Type = iota // 生成 xxx.ui.go UI 布局代码
	CodeGenEvent             // 生成 xxx.go 用户代码, 事件
	CodeGenSelf              // 生成 xxx.go 用户代码和 self 自引用(窗体 name 被修改时)
)

// 控制台, 数据载体类型
const (
	ConsoleInfo  Type = iota // 普通消息类型
	ConsoleClear             // 清空消息类型
)

// TPayload 通用的事件数据载体
type TPayload struct {
	Type Type
	Data any
}
