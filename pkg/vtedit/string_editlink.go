package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
	"log"
)

// 文本编辑框

type TStringEditLink struct {
	*TBaseEditLink
	edit      lcl.IEdit
	bounds    types.TRect
	text      string
	alignment types.TAlignment
	stopping  bool
}

func NewStringEditLink(bindData *TEditLinkNodeData) *TStringEditLink {
	link := new(TStringEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TStringEditLink) Create() {
	log.Println("TStringEditLink Create")
	m.edit = lcl.NewEdit(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetDoubleBuffered(true)
	m.edit.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})
}

// 通知编辑链接现在可以开始编辑。后代可以通过返回False来取消节点编辑。
func (m *TStringEditLink) BeginEdit() bool {
	log.Println("TStringEditLink BeginEdit")
	if !m.stopping {
		m.edit.Show()
		m.edit.SelectAll()
		m.edit.SetFocus()
	}
	return true
}

func (m *TStringEditLink) CancelEdit() bool {
	log.Println("TStringEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.Hide()
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TStringEditLink) EndEdit() bool {
	value := m.edit.Text()
	log.Println("TStringEditLink EndEdit Modified:", m.edit.Modified(), "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		if m.edit.Modified() {
			m.BindData.StringValue = value
		}
		m.VTree.EndEditNode()
		m.edit.Hide()
	}
	return true
}

func (m *TStringEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TStringEditLink PrepareEdit")
	if m.edit == nil || !m.edit.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	// 节点的初始大小、字体和文本。
	m.VTree.GetTextInfo(node, column, m.edit.Font(), &m.bounds, &m.text)
	log.Println("  PrepareEdit GetTextInfo:", m.bounds, m.text)
	m.edit.Font().SetColor(colors.ClWindowText)
	m.edit.SetParent(m.VTree)
	m.edit.HandleNeeded()
	m.edit.SetText(m.text)
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

func (m *TStringEditLink) GetBounds() types.TRect {
	log.Println("TStringEditLink GetBounds")
	return m.edit.BoundsRect()
}

func (m *TStringEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TStringEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TStringEditLink) SetBounds(R types.TRect) {
	log.Println("TStringEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)
}

func (m *TStringEditLink) Destroy(sender lcl.IObject) {
	log.Println("TStringEditLink Destroy")
	m.edit.Free()
}
