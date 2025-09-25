package vtedit

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
	"strconv"
)

// 文本编辑框

type TIntEditLink struct {
	*TBaseEditLink
	edit      lcl.IEdit
	bounds    types.TRect
	alignment types.TAlignment
	stopping  bool
}

func NewIntEditLink(bindData *TEditLinkNodeData) *TIntEditLink {
	link := new(TIntEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TIntEditLink) Create() {
	logs.Debug("TIntEditLink Create")
	m.edit = lcl.NewEdit(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetDoubleBuffered(true)
	oldText := m.edit.Text()
	m.edit.SetOnKeyPress(func(sender lcl.IObject, key *uint16) {
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		} else if !((*key >= keys.Vk0 && *key <= keys.Vk9) || (*key == keys.VkBack)) {
			*key = 0
			return
		}
		oldText = m.edit.Text()
	})
	m.edit.SetOnChange(func(sender lcl.IObject) {
		text := m.edit.Text()
		if text == "" {
			return
		}
		if _, err := strconv.Atoi(text); err != nil {
			m.edit.SetText(oldText)
			m.edit.SetSelStart(int32(len(oldText)))
		}
	})
}

// 通知编辑链接现在可以开始编辑。后代可以通过返回False来取消节点编辑。
func (m *TIntEditLink) BeginEdit() bool {
	logs.Debug("TIntEditLink BeginEdit")
	if !m.stopping {
		m.edit.Show()
		m.edit.SelectAll()
		m.edit.SetFocus()
	}
	return true
}

func (m *TIntEditLink) CancelEdit() bool {
	logs.Debug("TIntEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.Hide()
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TIntEditLink) EndEdit() bool {
	value := m.edit.Text()
	logs.Debug("TIntEditLink EndEdit", "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		if v, err := strconv.Atoi(value); err == nil {
			m.BindData.IntValue = v
		}
		m.VTree.EndEditNode()
		m.edit.Hide()
	}
	return true
}

func (m *TIntEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TIntEditLink PrepareEdit")
	if m.edit == nil || !m.edit.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	m.edit.Font().SetColor(colors.ClWindowText)
	m.edit.SetParent(m.VTree)
	m.edit.HandleNeeded()
	m.edit.SetText(strconv.Itoa(m.BindData.IntValue))
	if column <= -1 {
		m.edit.SetBiDiMode(m.VTree.BiDiMode())
		m.alignment = m.VTree.Alignment()
	} else {
		columns := m.VTree.Header().Columns()
		m.edit.SetBiDiMode(columns.ItemsWithColumnIndexToVirtualTreeColumn(column).BiDiMode())
		m.alignment = columns.ItemsWithColumnIndexToVirtualTreeColumn(column).Alignment()
	}

	if m.edit.BiDiMode() != types.BdLeftToRight {
		switch m.alignment {
		case types.TaLeftJustify:
			m.alignment = types.TaRightJustify
		case types.TaRightJustify:
			m.alignment = types.TaLeftJustify
		}
	}
	return true
}

func (m *TIntEditLink) GetBounds() types.TRect {
	logs.Debug("TIntEditLink GetBounds")
	return m.edit.BoundsRect()
}

func (m *TIntEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TIntEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TIntEditLink) SetBounds(R types.TRect) {
	logs.Debug("TIntEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)
}

func (m *TIntEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TIntEditLink Destroy")
	m.edit.Free()
}
