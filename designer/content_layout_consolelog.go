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
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strings"
	"time"
)

type ContentLayoutConsoleLog struct {
	tab *wg.TTab

	popupMenu lcl.IPopupMenu

	designer     lcl.ISynEdit // 设计器日志
	designerPage *wg.TPage    // 设计器日志

	build     lcl.ISynEdit // 构建输出日志
	buildPage *wg.TPage    // 构建输出日志
}

func initContentLayoutConsoleLog(owner *ContentLayout) *ContentLayoutConsoleLog {
	m := &ContentLayoutConsoleLog{}

	m.designer = lcl.NewSynEdit(owner.consoleLogPanel)
	m.designer.SetBorderStyleToBorderStyle(types.BsNone)
	m.designer.SetAlign(types.AlClient)
	m.designer.SetReadOnly(true)
	m.designer.SetWantTabs(false)
	m.designer.Gutter().SetVisible(false)
	m.designer.Gutter().SetWidth(0)
	m.designer.SetRightEdge(-1)
	//m.designer.SetHighlighter(m.synAnySyn)
	m.designer.SetOnSpecialLineMarkup(func(sender lcl.IObject, line int32, special *bool,
		markup lcl.ISynSelectedColor) {
		lines := m.designer.Lines()
		lineText := lines.Strings(line - 1)
		markup.SetBackground(0xFFFFFF)
		*special = true
		//markup.SetForeground(0x1E1F22)
		if strings.Contains(lineText, "[ERROR]") {
			markup.SetForeground(0x0000CC)
		} else if strings.Contains(lineText, "[WARN]") {
			markup.SetForeground(0x0066CC)
		} else if strings.Contains(lineText, "[INFO]") {
			markup.SetForeground(0x008800)
		} else if strings.Contains(lineText, "[DEBUG]") {
			markup.SetForeground(0x666666)
		} else {
			//markup.SetForeground(0x008800)
		}
	})

	font := m.designer.Font()
	font.SetHeight(-13)
	if tool.IsDarwin {
		font.SetQuality(types.FqAntialiased)
	} else if tool.IsWindows {
		font.SetQuality(types.FqCleartypeNatural)
		//font.SetName("Microsoft YaHei UI")
	} else {
		font.SetQuality(types.FqAntialiased)
	}
	m.designer.SetParent(owner.consoleLogPanel)

	m.CreatePopupMenu()
	return m
}

func (m *ContentLayoutConsoleLog) WriteDesignerLog(s ...string) {
	if m == nil || m.designer == nil {
		return
	}
	lines := m.designer.Lines()
	for _, text := range s {
		text = "[" + time.Now().Format("15:04:05.000") + "] " + text
		lines.Add(text)
		m.TrimLines(lines, 500)
	}
	m.designer.SetLeftChar(1)
	m.designer.SetTopLine(m.designer.Lines().Count())
}

func (m *ContentLayoutConsoleLog) TrimLines(lines lcl.IStrings, maxLines int32) {
	if m == nil || m.designer == nil {
		return
	}
	if lines.Count() > maxLines {
		lines.Delete(0)
	}
}

func (m *ContentLayoutConsoleLog) ClearConsole() {
	if m == nil || m.designer == nil {
		return
	}
	m.designer.Lines().Clear()
}

func (m *ContentLayoutConsoleLog) CreatePopupMenu() {
	popupMenu := lcl.NewPopupMenu(m.designer)
	m.popupMenu = popupMenu

	copyItem := lcl.NewMenuItem(m.designer)
	copyItem.SetCaption("复制")
	copyItem.SetOnClick(m.onCopyItemClick)
	popupMenu.Items().Add(copyItem)

	selectAllItem := lcl.NewMenuItem(m.designer)
	selectAllItem.SetCaption("全选")
	selectAllItem.SetOnClick(m.onSelectAllItemClick)
	popupMenu.Items().Add(selectAllItem)

	separator1 := lcl.NewMenuItem(m.designer)
	separator1.SetCaption("-")
	popupMenu.Items().Add(separator1)

	separator2 := lcl.NewMenuItem(m.designer)
	separator2.SetCaption("-")
	popupMenu.Items().Add(separator2)

	clearItem := lcl.NewMenuItem(m.designer)
	clearItem.SetCaption("清空")
	clearItem.SetOnClick(m.onClearItemClick)
	popupMenu.Items().Add(clearItem)

	m.designer.SetPopupMenu(popupMenu)
}

func (m *ContentLayoutConsoleLog) onCopyItemClick(sender lcl.IObject) {
	if m.designer != nil && m.designer.SelText() != "" {
		lcl.Clipboard.SetAsText(m.designer.SelText())
	}
}

func (m *ContentLayoutConsoleLog) onSelectAllItemClick(sender lcl.IObject) {
	if m.designer != nil {
		m.designer.SelectAll()
	}
}

func (m *ContentLayoutConsoleLog) onClearItemClick(sender lcl.IObject) {
	m.ClearConsole()
}
