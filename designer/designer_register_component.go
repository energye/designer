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

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
)

// 组件设计注册
// 所有要实现设计的组件都在此处注册

// 创建设计组件回调函数
type TNewComponent func(designerForm *FormTab, x, y int32) *TDesigningComponent

// 注册设计组件
// key: 组件类名, value: 组件创建函数
var registerComponents = make(map[string]TNewComponent)

func initRegisterComponent() {
	logs.Println("初始化注册组件")
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

// TRegisterComponent 注册组件信息
type TRegisterComponent struct {
	ClassName  string                    // 类名
	CreateFunc any                       // 创建函数
	Type       consts.ComponentType      // 类型
	Mod        consts.Mod                // 模块
	Default    tool.HashMap[string, any] // 默认属性, 反射调用API. key: 属性名, value: 属性值
}

// 组件全局注册表
var registerComponentsTest = tool.NewHashMap[string, *TRegisterComponent]()

// AddRegisterComponent 注册组件到全局注册表中
//
//   - key: 组件的唯一标识符，将作为组件的类名
//   - component: 要注册的组件实例指针
func AddRegisterComponent(key string, component *TRegisterComponent) {
	component.ClassName = key
	registerComponentsTest.Add(key, component)
}

// NewRegisterComponent 创建一个新的注册组件实例
//   - createFunc: 组件创建函数，用于实例化具体组件
//   - type_: 组件类型，标识组件的分类
//   - mod: 组件所属模块，定义组件的作用域
//   - 返回值: 指向TRegisterComponent结构体的指针
func NewRegisterComponent(createFunc any, type_ consts.ComponentType, mod consts.Mod) *TRegisterComponent {
	return &TRegisterComponent{
		CreateFunc: createFunc,
		Type:       type_,
		Mod:        mod,
	}
}

// NewLCLVisualRegisterComponent 创建一个新的LCL可视化组件注册器
// 该函数用于创建专门针对LCL可视组件的注册器实例
//
//   - createFunc: 组件创建函数，类型为any，用于实际创建组件实例
//   - TRegisterComponent: 返回指向TRegisterComponent的指针，表示新创建的组件注册器
func NewLCLVisualRegisterComponent(createFunc any) *TRegisterComponent {
	return NewRegisterComponent(createFunc, consts.CtVisual, consts.ModLCL)
}

// NewLCLNonVisualRegisterComponent 创建一个新的非可视化组件注册器
// 该函数用于创建一个专门用于注册非可视化组件的注册器实例
//
//   - createFunc: 组件创建函数，可以是任何类型的函数
//   - TRegisterComponent: 返回一个新的组件注册器指针，用于后续的组件注册操作
func NewLCLNonVisualRegisterComponent(createFunc any) *TRegisterComponent {
	return NewRegisterComponent(createFunc, consts.CtNonVisual, consts.ModLCL)
}

// initRegisterComponent2 初始化并注册所有可用的组件。
// 此函数在程序启动时调用，用于将标准、附加、通用、对话框、杂项、系统以及 Lazarus 特有组件 和 Web 组件
// 注册到全局组件映射中。
func initRegisterComponent2() {
	logs.Println("初始化注册组件")
	// 标准组件
	AddRegisterComponent("TMainMenu", NewLCLNonVisualRegisterComponent(lcl.NewMainMenu))
	AddRegisterComponent("TPopupMenu", NewLCLNonVisualRegisterComponent(lcl.NewPopupMenu))
	AddRegisterComponent("TButton", NewLCLVisualRegisterComponent(lcl.NewButton))
	AddRegisterComponent("TLabel", NewLCLVisualRegisterComponent(lcl.NewLabel))
	AddRegisterComponent("TEdit", NewLCLVisualRegisterComponent(lcl.NewEdit))
	AddRegisterComponent("TMemo", NewLCLVisualRegisterComponent(lcl.NewMemo))
	AddRegisterComponent("TToggleBox", NewLCLVisualRegisterComponent(lcl.NewToggleBox))
	AddRegisterComponent("TCheckBox", NewLCLVisualRegisterComponent(lcl.NewCheckBox))
	AddRegisterComponent("TRadioButton", NewLCLVisualRegisterComponent(lcl.NewRadioButton))
	AddRegisterComponent("TListBox", NewLCLVisualRegisterComponent(lcl.NewListBox))
	AddRegisterComponent("TComboBox", NewLCLVisualRegisterComponent(lcl.NewComboBox))
	AddRegisterComponent("TScrollBar", NewLCLVisualRegisterComponent(lcl.NewScrollBar))
	AddRegisterComponent("TGroupBox", NewLCLVisualRegisterComponent(lcl.NewGroupBox))
	AddRegisterComponent("TRadioGroup", NewLCLVisualRegisterComponent(lcl.NewRadioGroup))
	AddRegisterComponent("TCheckGroup", NewLCLVisualRegisterComponent(lcl.NewCheckGroup))
	AddRegisterComponent("TPanel", NewLCLVisualRegisterComponent(lcl.NewPanel))
	AddRegisterComponent("TFrame", NewLCLVisualRegisterComponent(lcl.NewFrame))
	AddRegisterComponent("TActionList", NewLCLNonVisualRegisterComponent(lcl.NewActionList))

	// 附加组件
	AddRegisterComponent("TBitBtn", NewLCLVisualRegisterComponent(lcl.NewBitBtn))
	AddRegisterComponent("TSpeedButton", NewLCLVisualRegisterComponent(lcl.NewSpeedButton))
	AddRegisterComponent("TStaticText", NewLCLVisualRegisterComponent(lcl.NewStaticText))
	AddRegisterComponent("TImage", NewLCLVisualRegisterComponent(lcl.NewImage))
	AddRegisterComponent("TShape", NewLCLVisualRegisterComponent(lcl.NewShape))
	AddRegisterComponent("TBevel", NewLCLVisualRegisterComponent(lcl.NewBevel))
	AddRegisterComponent("TPaintBox", NewLCLVisualRegisterComponent(lcl.NewPaintBox))
	AddRegisterComponent("TLabeledEdit", NewLCLVisualRegisterComponent(lcl.NewLabeledEdit))
	AddRegisterComponent("TSplitter", NewLCLVisualRegisterComponent(lcl.NewSplitter))
	AddRegisterComponent("TTrayIcon", NewLCLNonVisualRegisterComponent(lcl.NewTrayIcon))
	AddRegisterComponent("TControlBar", NewLCLVisualRegisterComponent(lcl.NewControlBar))
	AddRegisterComponent("TFlowPanel", NewLCLVisualRegisterComponent(lcl.NewFlowPanel))
	AddRegisterComponent("TMaskEdit", NewLCLVisualRegisterComponent(lcl.NewMaskEdit))
	AddRegisterComponent("TCheckListBox", NewLCLVisualRegisterComponent(lcl.NewCheckListBox))
	AddRegisterComponent("TScrollBox", NewLCLVisualRegisterComponent(lcl.NewScrollBox))
	AddRegisterComponent("TApplicationProperties", NewLCLNonVisualRegisterComponent(lcl.NewApplicationProperties))
	AddRegisterComponent("TStringGrid", NewLCLVisualRegisterComponent(lcl.NewStringGrid))
	AddRegisterComponent("TDrawGrid", NewLCLVisualRegisterComponent(lcl.NewDrawGrid))
	AddRegisterComponent("TPairSplitter", NewLCLVisualRegisterComponent(lcl.NewPairSplitter))
	AddRegisterComponent("TColorBox", NewLCLVisualRegisterComponent(lcl.NewColorBox))
	AddRegisterComponent("TColorListBox", NewLCLVisualRegisterComponent(lcl.NewColorListBox))
	AddRegisterComponent("TValueListEditor", NewLCLVisualRegisterComponent(lcl.NewValueListEditor))

	// 通用组件
	AddRegisterComponent("TTrackBar", NewLCLVisualRegisterComponent(lcl.NewTrackBar))
	AddRegisterComponent("TProgressBar", NewLCLVisualRegisterComponent(lcl.NewProgressBar))
	AddRegisterComponent("TTreeView", NewLCLVisualRegisterComponent(lcl.NewTreeView))
	AddRegisterComponent("TListView", NewLCLVisualRegisterComponent(lcl.NewListView))
	AddRegisterComponent("TStatusBar", NewLCLVisualRegisterComponent(lcl.NewStatusBar))
	AddRegisterComponent("TToolBar", NewLCLVisualRegisterComponent(lcl.NewToolBar))
	AddRegisterComponent("TCoolBar", NewLCLVisualRegisterComponent(lcl.NewCoolBar))
	AddRegisterComponent("TUpDown", NewLCLVisualRegisterComponent(lcl.NewUpDown))
	AddRegisterComponent("TPageControl", NewLCLVisualRegisterComponent(lcl.NewPageControl))
	AddRegisterComponent("THeaderControl", NewLCLVisualRegisterComponent(lcl.NewHeaderControl))
	AddRegisterComponent("TImageList", NewLCLNonVisualRegisterComponent(lcl.NewImageList))
	AddRegisterComponent("TPopupNotifier", NewLCLNonVisualRegisterComponent(lcl.NewPopupNotifier))
	AddRegisterComponent("TDateTimePicker", NewLCLVisualRegisterComponent(lcl.NewDateTimePicker))
	AddRegisterComponent("TRichMemo", NewLCLVisualRegisterComponent(lcl.NewRichMemo))

	// 对话框组件
	AddRegisterComponent("TOpenDialog", NewLCLNonVisualRegisterComponent(lcl.NewOpenDialog))
	AddRegisterComponent("TSaveDialog", NewLCLNonVisualRegisterComponent(lcl.NewSaveDialog))
	AddRegisterComponent("TSelectDirectoryDialog", NewLCLNonVisualRegisterComponent(lcl.NewSelectDirectoryDialog))
	AddRegisterComponent("TColorDialog", NewLCLNonVisualRegisterComponent(lcl.NewColorDialog))
	AddRegisterComponent("TFontDialog", NewLCLNonVisualRegisterComponent(lcl.NewFontDialog))
	AddRegisterComponent("TFindDialog", NewLCLNonVisualRegisterComponent(lcl.NewFindDialog))
	AddRegisterComponent("TReplaceDialog", NewLCLNonVisualRegisterComponent(lcl.NewReplaceDialog))
	AddRegisterComponent("TTaskDialog", NewLCLNonVisualRegisterComponent(lcl.NewTaskDialog))
	AddRegisterComponent("TOpenPictureDialog", NewLCLNonVisualRegisterComponent(lcl.NewOpenPictureDialog))
	AddRegisterComponent("TSavePictureDialog", NewLCLNonVisualRegisterComponent(lcl.NewSavePictureDialog))
	AddRegisterComponent("TCalendarDialog", NewLCLNonVisualRegisterComponent(lcl.NewCalendarDialog))
	AddRegisterComponent("TCalculatorDialog", NewLCLNonVisualRegisterComponent(lcl.NewCalculatorDialog))
	AddRegisterComponent("TPrinterSetupDialog", NewLCLNonVisualRegisterComponent(lcl.NewPrinterSetupDialog))
	AddRegisterComponent("TPrintDialog", NewLCLNonVisualRegisterComponent(lcl.NewPrintDialog))
	AddRegisterComponent("TPageSetupDialog", NewLCLNonVisualRegisterComponent(lcl.NewPageSetupDialog))

	// 杂项组件
	AddRegisterComponent("TColorButton", NewLCLVisualRegisterComponent(lcl.NewColorButton))
	AddRegisterComponent("TSpinEdit", NewLCLVisualRegisterComponent(lcl.NewSpinEdit))
	AddRegisterComponent("TFloatSpinEdit", NewLCLVisualRegisterComponent(lcl.NewFloatSpinEdit))
	AddRegisterComponent("TCalendar", NewLCLVisualRegisterComponent(lcl.NewCalendar))
	AddRegisterComponent("TEditButton", NewLCLVisualRegisterComponent(lcl.NewEditButton))
	AddRegisterComponent("TFileNameEdit", NewLCLVisualRegisterComponent(lcl.NewFileNameEdit))
	AddRegisterComponent("TDirectoryEdit", NewLCLVisualRegisterComponent(lcl.NewDirectoryEdit))
	AddRegisterComponent("TDateEdit", NewLCLVisualRegisterComponent(lcl.NewDateEdit))
	AddRegisterComponent("TTimeEdit", NewLCLVisualRegisterComponent(lcl.NewTimeEdit))
	AddRegisterComponent("TComboBoxEx", NewLCLVisualRegisterComponent(lcl.NewComboBoxEx))
	AddRegisterComponent("TButtonPanel", NewLCLVisualRegisterComponent(lcl.NewButtonPanel))
	AddRegisterComponent("TCheckComboBox", NewLCLVisualRegisterComponent(lcl.NewCheckComboBox))
	AddRegisterComponent("TLinkLabel", NewLCLVisualRegisterComponent(lcl.NewLinkLabel))
	AddRegisterComponent("TXButton", NewLCLVisualRegisterComponent(lcl.NewXButton))
	AddRegisterComponent("TImageButton", NewLCLVisualRegisterComponent(lcl.NewImageButton))
	AddRegisterComponent("TATGauge", NewLCLVisualRegisterComponent(lcl.NewATGauge))
	AddRegisterComponent("TOpenGLControl", NewLCLVisualRegisterComponent(lcl.NewOpenGLControl))

	// 系统组件
	AddRegisterComponent("TTimer", NewLCLNonVisualRegisterComponent(lcl.NewTimer))

	// Laz组件
	AddRegisterComponent("TLazVirtualDrawTree", NewLCLVisualRegisterComponent(lcl.NewLazVirtualDrawTree))
	AddRegisterComponent("TLazVirtualStringTree", NewLCLVisualRegisterComponent(lcl.NewLazVirtualStringTree))
	AddRegisterComponent("TDividerBevel", NewLCLVisualRegisterComponent(lcl.NewDividerBevel))
	AddRegisterComponent("TCheckBoxThemed", NewLCLVisualRegisterComponent(lcl.NewCheckBoxThemed))
	AddRegisterComponent("TExtendedNotebook", NewLCLVisualRegisterComponent(lcl.NewExtendedNotebook))
	AddRegisterComponent("TListFilterEdit", NewLCLVisualRegisterComponent(lcl.NewListFilterEdit))
	AddRegisterComponent("TTreeFilterEdit", NewLCLVisualRegisterComponent(lcl.NewTreeFilterEdit))
	AddRegisterComponent("TVTHeaderPopupMenu", NewLCLNonVisualRegisterComponent(lcl.NewVTHeaderPopupMenu))

}

// 获取注册的设计组件
func GetRegisterComponent(name string) TNewComponent {
	if cb, ok := registerComponents[name]; ok {
		return cb
	}
	return nil
}
