package config

import "github.com/energye/lcl/types"

// StorageWindowLayout 设计器窗口布局
type StorageWindowLayout struct {
	WindowBoundsRect types.TRect           `json:"bounds_rect"`
	WindowState      types.TWindowState    `json:"window_state"`
	MenuView         *StorageMenuView      `json:"menu_view"`
	ContentLayout    *StorageContentLayout `json:"content_layout"`
}

func (m *StorageWindowLayout) InitDefaultMenuView() {
	if m.MenuView == nil {
		m.MenuView = &StorageMenuView{
			WidgetsChecked:   true,
			ProjectChecked:   true,
			InspectorChecked: true,
			ConsoleChecked:   true,
			StatusbarChecked: true,
		}
	}
}

func (m *StorageWindowLayout) InitDefaultContentLayout() {
	if m.ContentLayout == nil {
		m.ContentLayout = &StorageContentLayout{
			WidgetPanelWidth:      170,
			ProjectPanelWidth:     150,
			InspectorPanelWidth:   225,
			ConsoleLogPanelHeight: 150,
		}
	}
}

// StorageMenuView 菜单视图
type StorageMenuView struct {
	WidgetsChecked   bool `json:"widgets_checked"`
	ProjectChecked   bool `json:"project_checked"`
	InspectorChecked bool `json:"inspector_checked"`
	ConsoleChecked   bool `json:"console_checked"`
	StatusbarChecked bool `json:"statusbar_checked"`
}

// StorageContentLayout 内容布局
type StorageContentLayout struct {
	WidgetPanelWidth      int32 `json:"widget_panel_width"`
	ProjectPanelWidth     int32 `json:"project_panel_width"`
	InspectorPanelWidth   int32 `json:"inspector_panel_width"`
	ConsoleLogPanelHeight int32 `json:"console_log_panel_height"`
}
