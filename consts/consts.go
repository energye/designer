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

package consts

const AppPackageName = "app"

// 项目配置文件扩展名
const EGPExt = ".egp"

// UI布局文件扩展名
const UIExt = ".ui"

// UI而已Go文件扩展名
const UIGoExt = ".ui.go"

// 应用代码目录的窗体列表文件名
const FormListFileName = "app.go"

// PreviewState 预览状态
type PreviewState int

const (
	PsStop     PreviewState = iota // 停止
	PsStarting                     // 启动中
	PsStarted                      // 启动完成
)

// DragShowStatus 拖拽显示状态
type DragShowStatus int32

const (
	DsAll         DragShowStatus = iota // 显示所有
	DsRightBottom                       // 显示 右 右下 下
)

const (
	// Mouse message key states
	MK_LBUTTON  = 1
	MK_RBUTTON  = 2
	MK_SHIFT    = 4
	MK_CONTROL  = 8
	MK_MBUTTON  = 0x10
	MK_XBUTTON1 = 0x20
	MK_XBUTTON2 = 0x40
	// following are "virtual" key states
	MK_DOUBLECLICK = 0x80
	MK_TRIPLECLICK = 0x100
	MK_QUADCLICK   = 0x200
	MK_ALT         = 0x20000000
)

const (
	DLeft = iota
	DTop
	DRight
	DBottom
	DLeftTop
	DRightTop
	DLeftBottom
	DRightBottom
)

// PropertyDataType 属性数据组件类型
type PropertyDataType int32

const (
	PdtText PropertyDataType = iota
	PdtInt
	PdtUint16
	PdtInt64
	PdtFloat
	PdtRadiobutton
	PdtCheckBox
	PdtCheckBoxDraw
	PdtCheckBoxList
	PdtComboBox
	PdtClassDialog
	PdtColorSelect
	PdtClass
)

type PropertyKind = string

const (
	TkClass       PropertyKind = "tkClass"
	TkEnumeration PropertyKind = "tkEnumeration"
	TkSet         PropertyKind = "tkSet"
	TkBool        PropertyKind = "tkBool"
	TkAString     PropertyKind = "tkAString"
	TkChar        PropertyKind = "tkChar"
	TkInteger     PropertyKind = "tkInteger"
	TkInt64       PropertyKind = "tkInt64"
	TkMethod      PropertyKind = "tkMethod"
)

// 组件类型
type ComponentType int32

const (
	CtForm      ComponentType = iota // 窗体
	CtNonVisual                      // 非可视组件
	CtVisual                         // 可视组件
)

// 改变 Z 序
type ChangeLevel int32

const (
	CLevelFront      ChangeLevel = iota //
	CLevelBack                          //
	CLevelForwardOne                    //
	CLevelBackOne                       //
)

// 组件所属模块 LCL, CEF, Webview2, Webkit2, Webkit2
type ComponentModule int32

const (
	CmLCL       ComponentModule = iota // LCL
	CmCEF                              // CEF
	CmWVWindows                        // Webview2 Windows
	CmWVMacOS                          // Webkit2 MacOS
	CmWVLinux                          // Webkit2 Linux
)
