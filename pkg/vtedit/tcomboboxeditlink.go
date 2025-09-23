package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"log"
)

// 下拉框

type TComboBoxEditLink struct {
	*TBaseEditLink
	edit      lcl.IComboBox
	bounds    types.TRect
	text      string
	alignment types.TAlignment
	stopping  bool
}

func NewComboBoxEditLink() *TComboBoxEditLink {
	m := new(TComboBoxEditLink)
	m.TBaseEditLink = NewEditLink(m)
	m.CreateEdit()
	return m
}

func (m *TComboBoxEditLink) CreateEdit() {
	log.Println("TComboBoxEditLink CreateEdit")
	m.edit = lcl.NewComboBox(nil)
	m.edit.SetVisible(false)
	m.edit.SetBorderStyle(types.BsSingle)
	m.edit.SetAutoSize(false)
	m.edit.SetDoubleBuffered(true)
	m.edit.SetOnChange(func(sender lcl.IObject) {
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.VTree.EndEditNode()
		})
	})
}

// 通知编辑链接现在可以开始编辑。后代可以通过返回False来取消节点编辑。
func (m *TComboBoxEditLink) BeginEdit() bool {
	log.Println("TComboBoxEditLink BeginEdit")
	if !m.stopping {
		m.edit.Show()
		m.edit.SelectAll()
		m.edit.SetFocus()
	}
	return true
}

func (m *TComboBoxEditLink) CancelEdit() bool {
	log.Println("TComboBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.edit.Hide()
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TComboBoxEditLink) EndEdit() bool {
	value := m.edit.Text()
	log.Println("TComboBoxEditLink EndEdit", "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		if m.OnNewData != nil {
			m.OnNewData(m.Node, m.Column, value)
		}
		m.VTree.EndEditNode()
		m.edit.Hide()
	}
	return true
}

func (m *TComboBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TComboBoxEditLink PrepareEdit")
	if m.edit == nil || !m.edit.IsValid() {
		m.CreateEdit()
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
	return true
}

func (m *TComboBoxEditLink) GetBounds() types.TRect {
	log.Println("TComboBoxEditLink GetBounds")
	return m.edit.BoundsRect()
}

func (m *TComboBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TComboBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.edit, msg)
}

func (m *TComboBoxEditLink) SetBounds(R types.TRect) {
	log.Println("TComboBoxEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetWidth(columnRect.Width())
	m.edit.SetBoundsRect(R)
}

func (m *TComboBoxEditLink) Destroy(sender lcl.IObject) {
	log.Println("TComboBoxEditLink Destroy")
	m.edit.Free()
}
