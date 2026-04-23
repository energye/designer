package designer

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/widget/wg"
	"time"
)

type ContentLayoutConsoleLog struct {
	tab *wg.TTab

	//designer     lcl.IMemo // 设计器日志
	designer     lcl.ISynEdit // 设计器日志
	designerPage *wg.TPage    // 设计器日志

	build     lcl.ISynEdit // 构建输出日志
	buildPage *wg.TPage    // 构建输出日志
}

func initContentLayoutConsoleLog(owner *ContentLayout) *ContentLayoutConsoleLog {
	m := &ContentLayoutConsoleLog{}
	//m.designer = lcl.NewMemo(owner.consoleLogPanel)
	//m.designer.SetBorderStyle(types.BsNone)
	//m.designer.SetAlign(types.AlClient)
	//m.designer.SetWantTabs(false)
	//m.designer.SetWordWrap(false)
	//m.designer.SetScrollBars(types.SsAutoBoth)
	m.designer = lcl.NewSynEdit(owner.consoleLogPanel)
	m.designer.SetBorderStyleToBorderStyle(types.BsNone)
	m.designer.SetAlign(types.AlClient)
	m.designer.SetWantTabs(false)
	m.designer.Gutter().SetVisible(false)
	m.designer.Gutter().SetWidth(0)
	m.designer.SetRightEdge(0)
	//m.designer.SetScrollBars(types.SsAutoBoth)
	font := m.designer.Font()
	font.SetName("Courier New")
	//font.SetSize(8)
	font.SetHeight(-13)
	font.SetQuality(types.FqAntialiased)
	m.designer.SetParent(owner.consoleLogPanel)
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
		m.designer.SetLeftChar(1)
		m.TrimMemoLines(500)
		m.designer.SetSelStart(m.designer.SelStart() + int32(len(text)))
	}
}

func (m *ContentLayoutConsoleLog) TrimMemoLines(maxLines int32) {
	if m == nil || m.designer == nil {
		return
	}
	if m.designer.Lines().Count() > maxLines {
		m.designer.Lines().Delete(0)
	}
}

func (m *ContentLayoutConsoleLog) ClearConsole() {
	if m == nil || m.designer == nil {
		return
	}
	m.designer.Lines().Clear()
}

func (m *ContentLayoutConsoleLog) CreatePopupMenu() {
	//procedure TLogViewer.CreatePopupMenu;
	//var
	//		Item: TMenuItem;
	//begin
	//FPopupMenu := TPopupMenu.Create(Self);
	//
	//// 复制
	//Item := TMenuItem.Create(FPopupMenu);
	//Item.Caption := '复制 (&C)';
	//Item.ShortCut := TextToShortCut('Ctrl+C');
	//Item.OnClick := @CopyItemClick;
	//FPopupMenu.Items.Add(Item);
	//
	//// 全选
	//Item := TMenuItem.Create(FPopupMenu);
	//Item.Caption := '全选 (&A)';
	//Item.ShortCut := TextToShortCut('Ctrl+A');
	//Item.OnClick := @SelectAllItemClick;
	//FPopupMenu.Items.Add(Item);
	//
	//// 分隔线
	//Item := TMenuItem.Create(FPopupMenu);
	//Item.Caption := '-';
	//FPopupMenu.Items.Add(Item);
	//
	//// 保存到文件
	//Item := TMenuItem.Create(FPopupMenu);
	//Item.Caption := '保存到文件 (&S)...';
	//Item.OnClick := @SaveToFileClick;
	//FPopupMenu.Items.Add(Item);
	//
	//// 分隔线
	//Item := TMenuItem.Create(FPopupMenu);
	//Item.Caption := '-';
	//FPopupMenu.Items.Add(Item);
	//
	//// 清空
	//Item := TMenuItem.Create(FPopupMenu);
	//Item.Caption := '清空 (&L)';
	//Item.OnClick := @ClearItemClick;
	//FPopupMenu.Items.Add(Item);
	//
	//// 关联菜单
	//SynEdit1.PopupMenu := FPopupMenu;
	//end;
}
