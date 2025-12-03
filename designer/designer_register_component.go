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
	ClassName string                    // 类名
	Func      any                       // 创建函数
	Type      consts.ComponentType      // 类型
	Mod       consts.Mod                // 模块
	Default   tool.HashMap[string, any] // 默认属性, 反射调用API. key: 属性名, value: 属性值
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
func NewRegisterComponent(func_ any, type_ consts.ComponentType, mod consts.Mod) *TRegisterComponent {
	return &TRegisterComponent{
		Func: func_,
		Type: type_,
		Mod:  mod,
	}
}

// NewLCLVisualRegisterComponent 创建一个新的LCL可视化组件注册器
// 该函数用于创建专门针对LCL可视组件的注册器实例
//
//   - createFunc: 组件创建函数，类型为any，用于实际创建组件实例
//   - TRegisterComponent: 返回指向TRegisterComponent的指针，表示新创建的组件注册器
func NewLCLVisualRegisterComponent(func_ any) *TRegisterComponent {
	return NewRegisterComponent(func_, consts.CtVisual, consts.ModLCL)
}

// NewLCLNonVisualRegisterComponent 创建一个新的非可视化组件注册器
// 该函数用于创建一个专门用于注册非可视化组件的注册器实例
//
//   - createFunc: 组件创建函数，可以是任何类型的函数
//   - TRegisterComponent: 返回一个新的组件注册器指针，用于后续的组件注册操作
func NewLCLNonVisualRegisterComponent(func_ any) *TRegisterComponent {
	return NewRegisterComponent(func_, consts.CtNonVisual, consts.ModLCL)
}

// initRegisterComponent2 初始化并注册所有可用的组件。
// 此函数在程序启动时调用，用于将标准、附加、通用、对话框、杂项、系统以及 Lazarus 特有组件 和 Web 组件
// 注册到全局组件映射中。
func initRegisterComponent2() {
	logs.Println("初始化注册组件")
	// 标准组件
	AddRegisterComponent("TMainMenu", NewLCLNonVisualRegisterComponent(lcl.TMainMenuClass))
	AddRegisterComponent("TPopupMenu", NewLCLNonVisualRegisterComponent(lcl.TPopupMenuClass))
	AddRegisterComponent("TButton", NewLCLVisualRegisterComponent(lcl.TButtonClass))
	AddRegisterComponent("TLabel", NewLCLVisualRegisterComponent(lcl.TLabelClass))
	AddRegisterComponent("TEdit", NewLCLVisualRegisterComponent(lcl.TEditClass))
	AddRegisterComponent("TMemo", NewLCLVisualRegisterComponent(lcl.TMemoClass))
	AddRegisterComponent("TToggleBox", NewLCLVisualRegisterComponent(lcl.TToggleBoxClass))
	AddRegisterComponent("TCheckBox", NewLCLVisualRegisterComponent(lcl.TCheckBoxClass))
	AddRegisterComponent("TRadioButton", NewLCLVisualRegisterComponent(lcl.TRadioButtonClass))
	AddRegisterComponent("TListBox", NewLCLVisualRegisterComponent(lcl.TListBoxClass))
	AddRegisterComponent("TComboBox", NewLCLVisualRegisterComponent(lcl.TComboBoxClass))
	AddRegisterComponent("TScrollBar", NewLCLVisualRegisterComponent(lcl.TScrollBarClass))
	AddRegisterComponent("TGroupBox", NewLCLVisualRegisterComponent(lcl.TGroupBoxClass))
	AddRegisterComponent("TRadioGroup", NewLCLVisualRegisterComponent(lcl.TRadioGroupClass))
	AddRegisterComponent("TCheckGroup", NewLCLVisualRegisterComponent(lcl.TCheckGroupClass))
	AddRegisterComponent("TPanel", NewLCLVisualRegisterComponent(lcl.TPanelClass))
	AddRegisterComponent("TFrame", NewLCLVisualRegisterComponent(lcl.TFrameClass))
	AddRegisterComponent("TActionList", NewLCLNonVisualRegisterComponent(lcl.TActionListClass))

	// 附加组件
	AddRegisterComponent("TBitBtn", NewLCLVisualRegisterComponent(lcl.TBitBtnClass))
	AddRegisterComponent("TSpeedButton", NewLCLVisualRegisterComponent(lcl.TSpeedButtonClass))
	AddRegisterComponent("TStaticText", NewLCLVisualRegisterComponent(lcl.TStaticTextClass))
	AddRegisterComponent("TImage", NewLCLVisualRegisterComponent(lcl.TImageClass))
	AddRegisterComponent("TShape", NewLCLVisualRegisterComponent(lcl.TShapeClass))
	AddRegisterComponent("TBevel", NewLCLVisualRegisterComponent(lcl.TBevelClass))
	AddRegisterComponent("TPaintBox", NewLCLVisualRegisterComponent(lcl.TPaintBoxClass))
	AddRegisterComponent("TLabeledEdit", NewLCLVisualRegisterComponent(lcl.TLabeledEditClass))
	AddRegisterComponent("TSplitter", NewLCLVisualRegisterComponent(lcl.TSplitterClass))
	AddRegisterComponent("TTrayIcon", NewLCLNonVisualRegisterComponent(lcl.TTrayIconClass))
	AddRegisterComponent("TControlBar", NewLCLVisualRegisterComponent(lcl.TControlBarClass))
	AddRegisterComponent("TFlowPanel", NewLCLVisualRegisterComponent(lcl.TFlowPanelClass))
	AddRegisterComponent("TMaskEdit", NewLCLVisualRegisterComponent(lcl.TMaskEditClass))
	AddRegisterComponent("TCheckListBox", NewLCLVisualRegisterComponent(lcl.TCheckListBoxClass))
	AddRegisterComponent("TScrollBox", NewLCLVisualRegisterComponent(lcl.TScrollBoxClass))
	AddRegisterComponent("TApplicationProperties", NewLCLNonVisualRegisterComponent(lcl.TApplicationPropertiesClass))
	AddRegisterComponent("TStringGrid", NewLCLVisualRegisterComponent(lcl.TStringGridClass))
	AddRegisterComponent("TDrawGrid", NewLCLVisualRegisterComponent(lcl.TDrawGridClass))
	AddRegisterComponent("TPairSplitter", NewLCLVisualRegisterComponent(lcl.TPairSplitterClass))
	AddRegisterComponent("TColorBox", NewLCLVisualRegisterComponent(lcl.TColorBoxClass))
	AddRegisterComponent("TColorListBox", NewLCLVisualRegisterComponent(lcl.TColorListBoxClass))
	AddRegisterComponent("TValueListEditor", NewLCLVisualRegisterComponent(lcl.TValueListEditorClass))

	// 通用组件
	AddRegisterComponent("TTrackBar", NewLCLVisualRegisterComponent(lcl.TTrackBarClass))
	AddRegisterComponent("TProgressBar", NewLCLVisualRegisterComponent(lcl.TProgressBarClass))
	AddRegisterComponent("TTreeView", NewLCLVisualRegisterComponent(lcl.TTreeViewClass))
	AddRegisterComponent("TListView", NewLCLVisualRegisterComponent(lcl.TListViewClass))
	AddRegisterComponent("TStatusBar", NewLCLVisualRegisterComponent(lcl.TStatusBarClass))
	AddRegisterComponent("TToolBar", NewLCLVisualRegisterComponent(lcl.TToolBarClass))
	AddRegisterComponent("TCoolBar", NewLCLVisualRegisterComponent(lcl.TCoolBarClass))
	AddRegisterComponent("TUpDown", NewLCLVisualRegisterComponent(lcl.TUpDownClass))
	AddRegisterComponent("TPageControl", NewLCLVisualRegisterComponent(lcl.TPageControlClass))
	AddRegisterComponent("THeaderControl", NewLCLVisualRegisterComponent(lcl.THeaderControlClass))
	AddRegisterComponent("TImageList", NewLCLNonVisualRegisterComponent(lcl.TImageListClass))
	AddRegisterComponent("TPopupNotifier", NewLCLNonVisualRegisterComponent(lcl.TPopupNotifierClass))
	AddRegisterComponent("TDateTimePicker", NewLCLVisualRegisterComponent(lcl.TDateTimePickerClass))
	AddRegisterComponent("TRichMemo", NewLCLVisualRegisterComponent(lcl.TRichMemoClass))

	// 对话框组件
	AddRegisterComponent("TOpenDialog", NewLCLNonVisualRegisterComponent(lcl.TOpenDialogClass))
	AddRegisterComponent("TSaveDialog", NewLCLNonVisualRegisterComponent(lcl.TSaveDialogClass))
	AddRegisterComponent("TSelectDirectoryDialog", NewLCLNonVisualRegisterComponent(lcl.TSelectDirectoryDialogClass))
	AddRegisterComponent("TColorDialog", NewLCLNonVisualRegisterComponent(lcl.TColorDialogClass))
	AddRegisterComponent("TFontDialog", NewLCLNonVisualRegisterComponent(lcl.TFontDialogClass))
	AddRegisterComponent("TFindDialog", NewLCLNonVisualRegisterComponent(lcl.TFindDialogClass))
	AddRegisterComponent("TReplaceDialog", NewLCLNonVisualRegisterComponent(lcl.TReplaceDialogClass))
	AddRegisterComponent("TTaskDialog", NewLCLNonVisualRegisterComponent(lcl.TTaskDialogClass))
	AddRegisterComponent("TOpenPictureDialog", NewLCLNonVisualRegisterComponent(lcl.TOpenPictureDialogClass))
	AddRegisterComponent("TSavePictureDialog", NewLCLNonVisualRegisterComponent(lcl.TSavePictureDialogClass))
	AddRegisterComponent("TCalendarDialog", NewLCLNonVisualRegisterComponent(lcl.TCalendarDialogClass))
	AddRegisterComponent("TCalculatorDialog", NewLCLNonVisualRegisterComponent(lcl.TCalculatorDialogClass))
	AddRegisterComponent("TPrinterSetupDialog", NewLCLNonVisualRegisterComponent(lcl.TPrinterSetupDialogClass))
	AddRegisterComponent("TPrintDialog", NewLCLNonVisualRegisterComponent(lcl.TPrintDialogClass))
	AddRegisterComponent("TPageSetupDialog", NewLCLNonVisualRegisterComponent(lcl.TPageSetupDialogClass))

	// 杂项组件
	AddRegisterComponent("TColorButton", NewLCLVisualRegisterComponent(lcl.TColorButtonClass))
	AddRegisterComponent("TSpinEdit", NewLCLVisualRegisterComponent(lcl.TSpinEditClass))
	AddRegisterComponent("TFloatSpinEdit", NewLCLVisualRegisterComponent(lcl.TFloatSpinEditClass))
	AddRegisterComponent("TCalendar", NewLCLVisualRegisterComponent(lcl.TCalendarClass))
	AddRegisterComponent("TEditButton", NewLCLVisualRegisterComponent(lcl.TEditButtonClass))
	AddRegisterComponent("TFileNameEdit", NewLCLVisualRegisterComponent(lcl.TFileNameEditClass))
	AddRegisterComponent("TDirectoryEdit", NewLCLVisualRegisterComponent(lcl.TDirectoryEditClass))
	AddRegisterComponent("TDateEdit", NewLCLVisualRegisterComponent(lcl.TDateEditClass))
	AddRegisterComponent("TTimeEdit", NewLCLVisualRegisterComponent(lcl.TTimeEditClass))
	AddRegisterComponent("TComboBoxEx", NewLCLVisualRegisterComponent(lcl.TComboBoxExClass))
	AddRegisterComponent("TButtonPanel", NewLCLVisualRegisterComponent(lcl.TButtonPanelClass))
	AddRegisterComponent("TCheckComboBox", NewLCLVisualRegisterComponent(lcl.TCheckComboBoxClass))
	AddRegisterComponent("TLinkLabel", NewLCLVisualRegisterComponent(lcl.TLinkLabelClass))
	AddRegisterComponent("TXButton", NewLCLVisualRegisterComponent(lcl.TXButtonClass))
	AddRegisterComponent("TImageButton", NewLCLVisualRegisterComponent(lcl.TImageButtonClass))
	AddRegisterComponent("TATGauge", NewLCLVisualRegisterComponent(lcl.TATGaugeClass))
	AddRegisterComponent("TOpenGLControl", NewLCLVisualRegisterComponent(lcl.TOpenGLControlClass))

	// 系统组件
	AddRegisterComponent("TTimer", NewLCLNonVisualRegisterComponent(lcl.TTimerClass))

	// Laz组件
	AddRegisterComponent("TLazVirtualDrawTree", NewLCLVisualRegisterComponent(lcl.TLazVirtualDrawTreeClass))
	AddRegisterComponent("TLazVirtualStringTree", NewLCLVisualRegisterComponent(lcl.TLazVirtualStringTreeClass))
	AddRegisterComponent("TDividerBevel", NewLCLVisualRegisterComponent(lcl.TDividerBevelClass))
	AddRegisterComponent("TCheckBoxThemed", NewLCLVisualRegisterComponent(lcl.TCheckBoxThemedClass))
	AddRegisterComponent("TExtendedNotebook", NewLCLVisualRegisterComponent(lcl.TExtendedNotebookClass))
	AddRegisterComponent("TListFilterEdit", NewLCLVisualRegisterComponent(lcl.TListFilterEditClass))
	AddRegisterComponent("TTreeFilterEdit", NewLCLVisualRegisterComponent(lcl.TTreeFilterEditClass))
	AddRegisterComponent("TVTHeaderPopupMenu", NewLCLNonVisualRegisterComponent(lcl.TVTHeaderPopupMenuClass))

}

// 获取注册的设计组件
func GetRegisterComponent(name string) TNewComponent {
	if cb, ok := registerComponents[name]; ok {
		return cb
	}
	return nil
}
