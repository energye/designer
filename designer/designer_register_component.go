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

package designer

// 组件设计注册
// 所有要实现设计的组件都在此处注册

// 创建设计组件回调函数
type TNewComponent func(designerForm *FormTab, x, y int32) *TDesigningComponent

// 注册设计组件
// key: 组件类名, value: 组件创建函数
var registerComponents = make(map[string]TNewComponent)

func initRegisterComponent() {
	registerComponents["TButton"] = NewButtonDesigner
	registerComponents["TEdit"] = NewEditDesigner
	registerComponents["TCheckBox"] = NewCheckBoxDesigner
	registerComponents["TPanel"] = NewPanelDesigner
	registerComponents["TMainMenu"] = NewMainMenuDesigner
	registerComponents["TPopupMenu"] = NewPopupMenuDesigner
	registerComponents["TLabel"] = NewLabelDesigner
	registerComponents["TMemo"] = NewMemoDesigner
	registerComponents["TToggleBox"] = NewToggleBoxDesigner
	registerComponents["TLazVirtualStringTree"] = NewLazVirtualStringTreeDesigner
}

// 获取注册的设计组件
func GetRegisterComponent(name string) TNewComponent {
	if cb, ok := registerComponents[name]; ok {
		return cb
	}
	return nil
}
