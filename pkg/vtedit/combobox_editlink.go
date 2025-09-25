package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
	"log"
)

// 下拉框

type TComboBoxEditLink struct {
	*TBaseEditLink
	combobox lcl.IComboBox
	bounds   types.TRect
	text     string
	stopping bool
}

func NewComboBoxEditLink(bindData *TEditLinkNodeData) *TComboBoxEditLink {
	link := new(TComboBoxEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.CreateEdit()
	return link
}

func (m *TComboBoxEditLink) CreateEdit() {
	log.Println("TComboBoxEditLink CreateEdit")
	m.combobox = lcl.NewComboBox(nil)
	m.combobox.SetVisible(false)
	m.combobox.SetBorderStyle(types.BsSingle)
	m.combobox.SetAutoSize(false)
	m.combobox.SetDoubleBuffered(true)
	m.combobox.SetOnChange(func(sender lcl.IObject) {
		m.BindData.Index = m.combobox.ItemIndex()
		m.BindData.StringValue = m.combobox.Text()
	})
	m.combobox.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})
	items := m.combobox.Items()
	for _, item := range m.BindData.ComboBoxValue {
		items.Add(item.StringValue)
	}
}

// 通知编辑链接现在可以开始编辑。后代可以通过返回False来取消节点编辑。
func (m *TComboBoxEditLink) BeginEdit() bool {
	log.Println("TComboBoxEditLink BeginEdit")
	if !m.stopping {
		m.combobox.Show()
		m.combobox.SelectAll()
		m.combobox.SetFocus()
	}
	return true
}

func (m *TComboBoxEditLink) CancelEdit() bool {
	log.Println("TComboBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.combobox.Hide()
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TComboBoxEditLink) EndEdit() bool {
	value := m.combobox.Text()
	log.Println("TComboBoxEditLink EndEdit", "value:", value, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.Index = m.combobox.ItemIndex()
		m.BindData.StringValue = m.combobox.Text()
		m.VTree.EndEditNode()
		m.combobox.Hide()
	}
	return true
}

func (m *TComboBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	log.Println("TComboBoxEditLink PrepareEdit")
	if m.combobox == nil || !m.combobox.IsValid() {
		m.CreateEdit()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	m.VTree.GetTextInfo(node, column, m.combobox.Font(), &m.bounds, &m.text)
	log.Println("  PrepareEdit GetTextInfo:", m.bounds, m.text)
	m.combobox.Font().SetColor(colors.ClWindowText)
	m.combobox.SetParent(m.VTree)
	m.combobox.HandleNeeded()
	m.combobox.SetText(m.text)
	return true
}

func (m *TComboBoxEditLink) GetBounds() types.TRect {
	log.Println("TComboBoxEditLink GetBounds")
	return m.combobox.BoundsRect()
}

func (m *TComboBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	log.Println("TComboBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.combobox, msg)
}

func (m *TComboBoxEditLink) SetBounds(R types.TRect) {
	log.Println("TComboBoxEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetWidth(columnRect.Width())
	m.combobox.SetBoundsRect(R)
}

func (m *TComboBoxEditLink) Destroy(sender lcl.IObject) {
	log.Println("TComboBoxEditLink Destroy")
	m.combobox.Free()
}
