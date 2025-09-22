package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
	"log"
)

type TOnNewData func(node types.PVirtualNode, column int32, value string)

type IStringEditLink interface {
	lcl.ICustomVTEditLink
	AsIVTEditLink() lcl.IVTEditLink
	SetOnNewData(fn TOnNewData)
}

type TStringEditLink struct {
	lcl.ICustomVTEditLink
	newData    TOnNewData
	edit       lcl.IEdit
	tree       lcl.ILazVirtualStringTree
	node       types.PVirtualNode
	column     int32
	textBounds types.TRect
	text       string
	alignment  types.TAlignment
	stopping   bool
}

func NewStringEditLink() IStringEditLink {
	m := new(TStringEditLink)
	m.ICustomVTEditLink = lcl.NewCustomVTEditLink()
	m.ICustomVTEditLink.SetOnBeginEdit(m.BeginEdit)
	m.ICustomVTEditLink.SetOnCancelEdit(m.CancelEdit)
	m.ICustomVTEditLink.SetOnEndEdit(m.EndEdit)
	m.ICustomVTEditLink.SetOnPrepareEdit(m.PrepareEdit)
	m.ICustomVTEditLink.SetOnGetBounds(m.GetBounds)
	m.ICustomVTEditLink.SetOnProcessMessage(m.ProcessMessage)
	m.ICustomVTEditLink.SetOnSetBounds(m.SetBounds)
	m.ICustomVTEditLink.SetOnDestroy(m.Destroy)
	m.CreateEdit()
	return m
}

func (m *TStringEditLink) SetOnNewData(fn TOnNewData) {
	m.newData = fn
}

func (m *TStringEditLink) AsIVTEditLink() lcl.IVTEditLink {
	return lcl.AsVTEditLink(m.ICustomVTEditLink.AsIntfVTEditLink())
}

func (m *TStringEditLink) CreateEdit() {
	m.edit = lcl.NewEdit(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		if *key == keys.VkReturn {
			//m.tree.SetText(m.node, m.column, m.edit.Text())
			//m.tree.EndEditNode()
			//m.edit.Hide()
		}
	})
}

// 通知编辑链接现在可以开始编辑。后代可以通过返回False来取消节点编辑。
func (m *TStringEditLink) BeginEdit() bool {
	log.Println("BeginEdit")
	if !m.stopping {
		m.edit.Show()
		m.edit.SelectAll()
		m.edit.SetFocus()
	}
	return true
}

func (m *TStringEditLink) CancelEdit() bool {
	log.Println("CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.Hide()
		m.tree.CancelEditNode()
	}
	return true
}

func (m *TStringEditLink) EndEdit() bool {
	text := m.edit.Text()
	log.Println("EndEdit Modified:", m.edit.Modified(), text, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		if m.edit.Modified() {
			//m.tree.SetText(m.node, m.column, text)
			if m.newData != nil {
				m.newData(m.node, m.column, text)
			}
		}
		m.tree.EndEditNode()
		m.edit.Hide()
	}
	return true
}

func (m *TStringEditLink) PrepareEdit(tree lcl.IBaseVirtualTree, node types.PVirtualNode, column int32) bool {
	log.Println("PrepareEdit")
	if !m.edit.IsValid() {
		m.CreateEdit()
	}
	m.tree = lcl.AsLazVirtualStringTree(tree)
	m.node = node
	m.column = column
	// 节点的初始大小、字体和文本。
	m.tree.GetTextInfo(node, column, m.edit.Font(), &m.textBounds, &m.text)
	log.Println("PrepareEdit GetTextInfo:", m.textBounds, m.text)
	m.edit.Font().SetColor(colors.ClWindowText)
	m.edit.SetParent(m.tree)
	m.edit.HandleNeeded()
	m.edit.SetText(m.text)
	if column <= -1 {
		m.edit.SetBiDiMode(m.tree.BiDiMode())
		m.alignment = m.tree.Alignment()
	} else {
		columns := m.tree.Header().Columns()
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
	log.Println("GetBounds", m.edit.BoundsRect().Width())
	return m.edit.BoundsRect()
}

func (m *TStringEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TStringEditLink) SetBounds(R types.TRect) {
	log.Println("SetBounds", R)
	columnRect := m.tree.GetDisplayRect(m.node, m.column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)
}

func (m *TStringEditLink) Destroy(sender lcl.IObject) {
	log.Println("Destroy")
	m.edit.Free()
}
