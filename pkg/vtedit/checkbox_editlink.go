package vtedit

import (
	"bytes"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"strconv"
)

// CheckBox

type TCheckBoxEditLink struct {
	*TBaseEditLink
	checkbox  lcl.ICheckBox
	alignment types.TAlignment
	stopping  bool
}

func NewCheckBoxEditLink(bindData *TEditLinkNodeData) *TCheckBoxEditLink {
	link := new(TCheckBoxEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TCheckBoxEditLink) Create() {
	logs.Debug("TCheckBoxEditLink CreateEdit")
	m.checkbox = lcl.NewCheckBox(nil)
	m.checkbox.SetVisible(false)
	m.checkbox.SetCaption("(false)")
	m.checkbox.SetDoubleBuffered(true)
	m.checkbox.SetOnChange(func(sender lcl.IObject) {
		m.checkbox.SetCaption("(" + strconv.FormatBool(m.checkbox.Checked()) + ")")
		m.BindData.Checked = m.checkbox.Checked()
		logs.Debug("TCheckBoxEditLink OnChange checked:", m.BindData.Checked)
		node := m.Node.ToGo()
		parentNode := node.Parent
		if pData := GetPropertyNodeData(parentNode); pData != nil {
			dataList := pData.CheckBoxValue
			buf := bytes.Buffer{}
			buf.WriteString("[")
			i := 0
			for _, item := range dataList {
				if item.Checked {
					if i > 0 {
						buf.WriteString(",")
					}
					buf.WriteString(item.Name)
					i++
				}
			}
			buf.WriteString("]")
			pData.StringValue = buf.String()
			logs.Debug("TCheckBoxEditLink OnChange ParentNode-text:", pData.StringValue)
			m.VTree.InvalidateNode(parentNode)
		}
	})
}

func (m *TCheckBoxEditLink) BeginEdit() bool {
	if !m.stopping {
		m.checkbox.SetVisible(true)
	}
	return true
}

func (m *TCheckBoxEditLink) CancelEdit() bool {
	logs.Debug("TCheckBoxEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.checkbox.SetVisible(false)
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TCheckBoxEditLink) EndEdit() bool {
	value := m.checkbox.Checked()
	logs.Debug("TCheckBoxEditLink EndEdit", "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.BoolValue = value
		m.VTree.EndEditNode()
		m.checkbox.SetVisible(false)
	}
	return true
}

func (m *TCheckBoxEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TCheckBoxEditLink PrepareEdit")
	if m.checkbox == nil || m.checkbox.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	value := m.VTree.Text(m.Node, m.Column)
	if v, err := strconv.ParseBool(value); err == nil {
		m.BindData.Checked = v
	}
	m.checkbox.SetParent(m.VTree)
	m.checkbox.HandleNeeded()
	m.checkbox.SetChecked(m.BindData.Checked)
	return true
}

func (m *TCheckBoxEditLink) GetBounds() (R types.TRect) {
	logs.Debug("TCheckBoxEditLink GetBounds")
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R = columnRect
	return
}

func (m *TCheckBoxEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TCheckBoxEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.checkbox, msg)
}

func (m *TCheckBoxEditLink) SetBounds(R types.TRect) {
	logs.Debug("TCheckBoxEditLink SetBounds", R)
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left + 5
	R.Top = columnRect.Top + 3
	//R.SetHeight(columnRect.Height())
	//R.SetWidth(columnRect.Width())
	m.checkbox.SetBoundsRect(R)

}

func (m *TCheckBoxEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TCheckBoxEditLink Destroy")
	m.checkbox.Free()
}
