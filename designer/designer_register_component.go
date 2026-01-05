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
	"github.com/energye/energy/v3/wv"
	"github.com/energye/lcl/lcl"
)

// 组件设计注册
// 所有要实现设计的组件都在此处注册
//
//// 创建设计组件回调函数
//type TNewComponent func(designerForm *FormTab, x, y int32) *TDesigningComponent
//
//// 注册设计组件
//// key: 组件类名, value: 组件创建函数
//var registerComponents = make(map[string]TNewComponent)
//
//func initRegisterComponent() {
//	logs.Println("初始化注册组件")
//	registerComponents["TButton"] = NewButtonDesigner
//	registerComponents["TEdit"] = NewEditDesigner
//	registerComponents["TCheckBox"] = NewCheckBoxDesigner
//	registerComponents["TPanel"] = NewPanelDesigner
//	registerComponents["TMainMenu"] = NewMainMenuDesigner
//	registerComponents["TPopupMenu"] = NewPopupMenuDesigner
//	registerComponents["TLabel"] = NewLabelDesigner
//	registerComponents["TMemo"] = NewMemoDesigner
//	registerComponents["TToggleBox"] = NewToggleBoxDesigner
//	registerComponents["TLazVirtualStringTree"] = NewLazVirtualStringTreeDesigner
//}

// TRegisterComponent 注册组件信息
type TRegisterComponent struct {
	ClassName string                    // 类名
	Func      any                       // 创建函数
	AsFunc    any                       // 转换函数
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
func NewRegisterComponent(func_, asFunc any, type_ consts.ComponentType, mod consts.Mod) *TRegisterComponent {
	return &TRegisterComponent{
		Func:   func_,
		AsFunc: asFunc,
		Type:   type_,
		Mod:    mod,
	}
}

// NewLCLVisualRegisterComponent 创建一个新的LCL可视化组件注册器
// 该函数用于创建专门针对LCL可视组件的注册器实例
//
//   - createFunc: 组件创建函数，类型为any，用于实际创建组件实例
//   - TRegisterComponent: 返回指向TRegisterComponent的指针，表示新创建的组件注册器
func NewLCLVisualRegisterComponent(func_, asFunc any) *TRegisterComponent {
	return NewRegisterComponent(func_, asFunc, consts.CtVisual, consts.ModLCL)
}

// NewEnergyCustomVisualRegisterComponent 创建一个新的 energy 自定义视觉注册组件
// 该函数用于注册具有自定义视觉效果的能量相关组件
//
//   - func_ - 要注册的函数对象，通常为组件的构造函数或初始化函数
//   - asFunc - 作为函数使用的对象，用于类型转换或接口实现
//   - *TRegisterComponent - 返回注册组件的指针，用于后续的组件管理操作
func NewEnergyCustomVisualRegisterComponent(func_, asFunc any) *TRegisterComponent {
	return NewRegisterComponent(func_, asFunc, consts.CtCustomVisual, consts.ModEnergy)
}

// NewLCLNonVisualRegisterComponent 创建一个新的非可视化组件注册器
// 该函数用于创建一个专门用于注册非可视化组件的注册器实例
//
//   - createFunc: 组件创建函数，可以是任何类型的函数
//   - TRegisterComponent: 返回一个新的组件注册器指针，用于后续的组件注册操作
func NewLCLNonVisualRegisterComponent(func_, asFunc any) *TRegisterComponent {
	return NewRegisterComponent(func_, asFunc, consts.CtNonVisual, consts.ModLCL)
}

// initRegisterComponent2 初始化并注册所有可用的组件。
// 此函数在程序启动时调用，用于将标准、附加、通用、对话框、杂项、系统以及 Lazarus 特有组件 和 Web 组件
// 注册到全局组件映射中。
func initRegisterComponent() {
	logs.Println("初始化注册组件")
	// 标准组件
	AddRegisterComponent("TMainMenu", NewLCLNonVisualRegisterComponent(lcl.TMainMenuClass, lcl.AsMainMenu))
	AddRegisterComponent("TPopupMenu", NewLCLNonVisualRegisterComponent(lcl.TPopupMenuClass, lcl.AsPopupMenu))
	AddRegisterComponent("TButton", NewLCLVisualRegisterComponent(lcl.TButtonClass, lcl.AsButton))
	AddRegisterComponent("TLabel", NewLCLVisualRegisterComponent(lcl.TLabelClass, lcl.AsLabel))
	AddRegisterComponent("TEdit", NewLCLVisualRegisterComponent(lcl.TEditClass, lcl.AsEdit))
	AddRegisterComponent("TMemo", NewLCLVisualRegisterComponent(lcl.TMemoClass, lcl.AsMemo))
	AddRegisterComponent("TToggleBox", NewLCLVisualRegisterComponent(lcl.TToggleBoxClass, lcl.AsToggleBox))
	AddRegisterComponent("TCheckBox", NewLCLVisualRegisterComponent(lcl.TCheckBoxClass, lcl.AsCheckBox))
	AddRegisterComponent("TRadioButton", NewLCLVisualRegisterComponent(lcl.TRadioButtonClass, lcl.AsRadioButton))
	AddRegisterComponent("TListBox", NewLCLVisualRegisterComponent(lcl.TListBoxClass, lcl.AsListBox))
	AddRegisterComponent("TComboBox", NewLCLVisualRegisterComponent(lcl.TComboBoxClass, lcl.AsComboBox))
	AddRegisterComponent("TScrollBar", NewLCLVisualRegisterComponent(lcl.TScrollBarClass, lcl.AsScrollBar))
	AddRegisterComponent("TGroupBox", NewLCLVisualRegisterComponent(lcl.TGroupBoxClass, lcl.AsGroupBox))
	AddRegisterComponent("TRadioGroup", NewLCLVisualRegisterComponent(lcl.TRadioGroupClass, lcl.AsRadioGroup))
	AddRegisterComponent("TCheckGroup", NewLCLVisualRegisterComponent(lcl.TCheckGroupClass, lcl.AsCheckGroup))
	AddRegisterComponent("TPanel", NewLCLVisualRegisterComponent(lcl.TPanelClass, lcl.AsPanel))
	AddRegisterComponent("TFrame", NewLCLVisualRegisterComponent(lcl.TFrameClass, lcl.AsFrame))
	AddRegisterComponent("TActionList", NewLCLNonVisualRegisterComponent(lcl.TActionListClass, lcl.AsActionList))

	// 附加组件
	AddRegisterComponent("TBitBtn", NewLCLVisualRegisterComponent(lcl.TBitBtnClass, lcl.AsBitBtn))
	AddRegisterComponent("TSpeedButton", NewLCLVisualRegisterComponent(lcl.TSpeedButtonClass, lcl.AsSpeedButton))
	AddRegisterComponent("TStaticText", NewLCLVisualRegisterComponent(lcl.TStaticTextClass, lcl.AsStaticText))
	AddRegisterComponent("TImage", NewLCLVisualRegisterComponent(lcl.TImageClass, lcl.AsImage))
	AddRegisterComponent("TShape", NewLCLVisualRegisterComponent(lcl.TShapeClass, lcl.AsShape))
	AddRegisterComponent("TBevel", NewLCLVisualRegisterComponent(lcl.TBevelClass, lcl.AsBevel))
	AddRegisterComponent("TPaintBox", NewLCLVisualRegisterComponent(lcl.TPaintBoxClass, lcl.AsPaintBox))
	AddRegisterComponent("TLabeledEdit", NewLCLVisualRegisterComponent(lcl.TLabeledEditClass, lcl.AsLabeledEdit))
	AddRegisterComponent("TSplitter", NewLCLVisualRegisterComponent(lcl.TSplitterClass, lcl.AsSplitter))
	AddRegisterComponent("TTrayIcon", NewLCLNonVisualRegisterComponent(lcl.TTrayIconClass, lcl.AsTrayIcon))
	AddRegisterComponent("TControlBar", NewLCLVisualRegisterComponent(lcl.TControlBarClass, lcl.AsControlBar))
	AddRegisterComponent("TFlowPanel", NewLCLVisualRegisterComponent(lcl.TFlowPanelClass, lcl.AsFlowPanel))
	AddRegisterComponent("TMaskEdit", NewLCLVisualRegisterComponent(lcl.TMaskEditClass, lcl.AsMaskEdit))
	AddRegisterComponent("TCheckListBox", NewLCLVisualRegisterComponent(lcl.TCheckListBoxClass, lcl.AsCheckListBox))
	AddRegisterComponent("TScrollBox", NewLCLVisualRegisterComponent(lcl.TScrollBoxClass, lcl.AsScrollBox))
	AddRegisterComponent("TApplicationProperties", NewLCLNonVisualRegisterComponent(lcl.TApplicationPropertiesClass, lcl.AsApplicationProperties))
	AddRegisterComponent("TStringGrid", NewLCLVisualRegisterComponent(lcl.TStringGridClass, lcl.AsStringGrid))
	AddRegisterComponent("TDrawGrid", NewLCLVisualRegisterComponent(lcl.TDrawGridClass, lcl.AsDrawGrid))
	AddRegisterComponent("TPairSplitter", NewLCLVisualRegisterComponent(lcl.TPairSplitterClass, lcl.AsPairSplitter))
	AddRegisterComponent("TColorBox", NewLCLVisualRegisterComponent(lcl.TColorBoxClass, lcl.AsColorBox))
	AddRegisterComponent("TColorListBox", NewLCLVisualRegisterComponent(lcl.TColorListBoxClass, lcl.AsColorListBox))
	AddRegisterComponent("TValueListEditor", NewLCLVisualRegisterComponent(lcl.TValueListEditorClass, lcl.AsValueListEditor))

	// 通用组件
	AddRegisterComponent("TTrackBar", NewLCLVisualRegisterComponent(lcl.TTrackBarClass, lcl.AsTrackBar))
	AddRegisterComponent("TProgressBar", NewLCLVisualRegisterComponent(lcl.TProgressBarClass, lcl.AsProgressBar))
	AddRegisterComponent("TTreeView", NewLCLVisualRegisterComponent(lcl.TTreeViewClass, lcl.AsTreeView))
	AddRegisterComponent("TListView", NewLCLVisualRegisterComponent(lcl.TListViewClass, lcl.AsListView))
	AddRegisterComponent("TStatusBar", NewLCLVisualRegisterComponent(lcl.TStatusBarClass, lcl.AsStatusBar))
	AddRegisterComponent("TToolBar", NewLCLVisualRegisterComponent(lcl.TToolBarClass, lcl.AsToolBar))
	AddRegisterComponent("TCoolBar", NewLCLVisualRegisterComponent(lcl.TCoolBarClass, lcl.AsCoolBar))
	AddRegisterComponent("TUpDown", NewLCLVisualRegisterComponent(lcl.TUpDownClass, lcl.AsUpDown))
	AddRegisterComponent("TPageControl", NewLCLVisualRegisterComponent(lcl.TPageControlClass, lcl.AsPageControl))
	AddRegisterComponent("THeaderControl", NewLCLVisualRegisterComponent(lcl.THeaderControlClass, lcl.AsHeaderControl))
	AddRegisterComponent("TImageList", NewLCLNonVisualRegisterComponent(lcl.TImageListClass, lcl.AsImageList))
	AddRegisterComponent("TPopupNotifier", NewLCLNonVisualRegisterComponent(lcl.TPopupNotifierClass, lcl.AsPopupNotifier))
	AddRegisterComponent("TDateTimePicker", NewLCLVisualRegisterComponent(lcl.TDateTimePickerClass, lcl.AsDateTimePicker))
	AddRegisterComponent("TRichMemo", NewLCLVisualRegisterComponent(lcl.TRichMemoClass, lcl.AsRichMemo))

	// 对话框组件
	AddRegisterComponent("TOpenDialog", NewLCLNonVisualRegisterComponent(lcl.TOpenDialogClass, lcl.AsOpenDialog))
	AddRegisterComponent("TSaveDialog", NewLCLNonVisualRegisterComponent(lcl.TSaveDialogClass, lcl.AsSaveDialog))
	AddRegisterComponent("TSelectDirectoryDialog", NewLCLNonVisualRegisterComponent(lcl.TSelectDirectoryDialogClass, lcl.AsSelectDirectoryDialog))
	AddRegisterComponent("TColorDialog", NewLCLNonVisualRegisterComponent(lcl.TColorDialogClass, lcl.AsColorDialog))
	AddRegisterComponent("TFontDialog", NewLCLNonVisualRegisterComponent(lcl.TFontDialogClass, lcl.AsFontDialog))
	AddRegisterComponent("TFindDialog", NewLCLNonVisualRegisterComponent(lcl.TFindDialogClass, lcl.AsFindDialog))
	AddRegisterComponent("TReplaceDialog", NewLCLNonVisualRegisterComponent(lcl.TReplaceDialogClass, lcl.AsReplaceDialog))
	AddRegisterComponent("TTaskDialog", NewLCLNonVisualRegisterComponent(lcl.TTaskDialogClass, lcl.AsTaskDialog))
	AddRegisterComponent("TOpenPictureDialog", NewLCLNonVisualRegisterComponent(lcl.TOpenPictureDialogClass, lcl.AsOpenPictureDialog))
	AddRegisterComponent("TSavePictureDialog", NewLCLNonVisualRegisterComponent(lcl.TSavePictureDialogClass, lcl.AsSavePictureDialog))
	AddRegisterComponent("TCalendarDialog", NewLCLNonVisualRegisterComponent(lcl.TCalendarDialogClass, lcl.AsCalendarDialog))
	AddRegisterComponent("TCalculatorDialog", NewLCLNonVisualRegisterComponent(lcl.TCalculatorDialogClass, lcl.AsCalculatorDialog))
	AddRegisterComponent("TPrinterSetupDialog", NewLCLNonVisualRegisterComponent(lcl.TPrinterSetupDialogClass, lcl.AsPrinterSetupDialog))
	AddRegisterComponent("TPrintDialog", NewLCLNonVisualRegisterComponent(lcl.TPrintDialogClass, lcl.AsPrintDialog))
	AddRegisterComponent("TPageSetupDialog", NewLCLNonVisualRegisterComponent(lcl.TPageSetupDialogClass, lcl.AsPageSetupDialog))

	// 杂项组件
	AddRegisterComponent("TColorButton", NewLCLVisualRegisterComponent(lcl.TColorButtonClass, lcl.AsColorButton))
	AddRegisterComponent("TSpinEdit", NewLCLVisualRegisterComponent(lcl.TSpinEditClass, lcl.AsSpinEdit))
	AddRegisterComponent("TFloatSpinEdit", NewLCLVisualRegisterComponent(lcl.TFloatSpinEditClass, lcl.AsFloatSpinEdit))
	AddRegisterComponent("TCalendar", NewLCLVisualRegisterComponent(lcl.TCalendarClass, lcl.AsCalendar))
	AddRegisterComponent("TEditButton", NewLCLVisualRegisterComponent(lcl.TEditButtonClass, lcl.AsEditButton))
	AddRegisterComponent("TFileNameEdit", NewLCLVisualRegisterComponent(lcl.TFileNameEditClass, lcl.AsFileNameEdit))
	AddRegisterComponent("TDirectoryEdit", NewLCLVisualRegisterComponent(lcl.TDirectoryEditClass, lcl.AsDirectoryEdit))
	AddRegisterComponent("TDateEdit", NewLCLVisualRegisterComponent(lcl.TDateEditClass, lcl.AsDateEdit))
	AddRegisterComponent("TTimeEdit", NewLCLVisualRegisterComponent(lcl.TTimeEditClass, lcl.AsTimeEdit))
	AddRegisterComponent("TComboBoxEx", NewLCLVisualRegisterComponent(lcl.TComboBoxExClass, lcl.AsComboBoxEx))
	AddRegisterComponent("TButtonPanel", NewLCLVisualRegisterComponent(lcl.TButtonPanelClass, lcl.AsButtonPanel))
	AddRegisterComponent("TCheckComboBox", NewLCLVisualRegisterComponent(lcl.TCheckComboBoxClass, lcl.AsCheckComboBox))
	AddRegisterComponent("TLinkLabel", NewLCLVisualRegisterComponent(lcl.TLinkLabelClass, lcl.AsLinkLabel))
	AddRegisterComponent("TXButton", NewLCLVisualRegisterComponent(lcl.TXButtonClass, lcl.AsXButton))
	AddRegisterComponent("TImageButton", NewLCLVisualRegisterComponent(lcl.TImageButtonClass, lcl.AsImageButton))
	AddRegisterComponent("TATGauge", NewLCLVisualRegisterComponent(lcl.TATGaugeClass, lcl.AsATGauge))
	AddRegisterComponent("TOpenGLControl", NewLCLVisualRegisterComponent(lcl.TOpenGLControlClass, lcl.AsOpenGLControl))

	// 系统组件
	AddRegisterComponent("TTimer", NewLCLNonVisualRegisterComponent(lcl.TTimerClass, lcl.AsTimer))

	// Laz组件
	AddRegisterComponent("TLazVirtualDrawTree", NewLCLVisualRegisterComponent(lcl.TLazVirtualDrawTreeClass, lcl.AsLazVirtualDrawTree))
	AddRegisterComponent("TLazVirtualStringTree", NewLCLVisualRegisterComponent(lcl.TLazVirtualStringTreeClass, lcl.AsLazVirtualStringTree))
	AddRegisterComponent("TDividerBevel", NewLCLVisualRegisterComponent(lcl.TDividerBevelClass, lcl.AsDividerBevel))
	AddRegisterComponent("TCheckBoxThemed", NewLCLVisualRegisterComponent(lcl.TCheckBoxThemedClass, lcl.AsCheckBoxThemed))
	AddRegisterComponent("TExtendedNotebook", NewLCLVisualRegisterComponent(lcl.TExtendedNotebookClass, lcl.AsExtendedNotebook))
	AddRegisterComponent("TListFilterEdit", NewLCLVisualRegisterComponent(lcl.TListFilterEditClass, lcl.AsListFilterEdit))
	AddRegisterComponent("TTreeFilterEdit", NewLCLVisualRegisterComponent(lcl.TTreeFilterEditClass, lcl.AsTreeFilterEdit))
	AddRegisterComponent("TVTHeaderPopupMenu", NewLCLNonVisualRegisterComponent(lcl.TVTHeaderPopupMenuClass, lcl.AsVTHeaderPopupMenu))

	// Web组件
	AddRegisterComponent("TWebview", NewEnergyCustomVisualRegisterComponent(wv.NewWebview, nil))
	//AddRegisterComponent("TWebview", NewRegisterComponent(wv.NewBrowserWindow, nil, consts.CtVisual, consts.ModEnergy))
}

//
//// 获取注册的设计组件
//func GetRegisterComponent(name string) TNewComponent {
//	if cb, ok := registerComponents[name]; ok {
//		return cb
//	}
//	return nil
//}
