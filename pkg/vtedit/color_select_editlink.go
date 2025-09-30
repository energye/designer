package vtedit

import (
	"fmt"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
	"github.com/energye/lcl/types/keys"
)

// 颜色选择器

type TColorSelectEditLink struct {
	*TBaseEditLink
	colorText lcl.IEdit
	colorBtn  lcl.IColorButton
	bounds    types.TRect
	text      string
	stopping  bool
}

func NewColorSelectEditLink(bindData *TEditNodeData) *TColorSelectEditLink {
	link := new(TColorSelectEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.CreateEdit()
	return link
}

func (m *TColorSelectEditLink) CreateEdit() {
	logs.Debug("TColorSelectEditLink CreateEdit")
	m.colorText = lcl.NewEdit(nil)
	m.colorText.SetVisible(false)
	m.colorText.SetBorderStyle(types.BsSingle)
	m.colorText.SetAutoSize(false)
	m.colorText.SetDoubleBuffered(true)
	m.colorText.SetReadOnly(true)
	m.colorText.SetText(fmt.Sprintf("0x%X", m.BindData.EditNodeData.IntValue))
	clrFont := m.colorText.Font()
	clrFont.SetStyle(clrFont.Style().Include(types.FsBold))
	clrFont.SetColor(colors.TColor(m.BindData.EditNodeData.IntValue))
	m.colorText.SetOnKeyDown(func(sender lcl.IObject, key *uint16, shift types.TShiftState) {
		logs.Debug("TColorSelectEditLink OnKeyDown key:", *key)
		if *key == keys.VkReturn {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.VTree.EndEditNode()
			})
		}
	})

	m.colorBtn = lcl.NewColorButton(nil)
	m.colorBtn.SetVisible(false)
	m.colorBtn.SetAutoSize(false)
	m.colorBtn.SetButtonColor(colors.TColor(m.BindData.EditNodeData.IntValue))
	m.colorBtn.SetOnColorChanged(func(sender lcl.IObject) {
		color := m.colorBtn.ButtonColor()
		logs.Debug("TColorSelectEditLink OnColorChanged color:", color)
		m.BindData.EditNodeData.IntValue = int(color)
		m.colorText.SetText(fmt.Sprintf("0x%X", color))
		clrFont.SetColor(color)
		lcl.RunOnMainThreadAsync(func(id uint32) {
			m.VTree.EndEditNode()
		})
	})
}

func (m *TColorSelectEditLink) BeginEdit() bool {
	logs.Debug("TColorSelectEditLink BeginEdit")
	if !m.stopping {
		m.colorText.Show()
		m.colorText.SetFocus()
		m.colorBtn.Show()
	}
	return true
}

func (m *TColorSelectEditLink) CancelEdit() bool {
	logs.Debug("TColorSelectEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.colorText.Hide()
		m.colorBtn.Hide()
		m.VTree.CancelEditNode()
	}
	return true
}

func (m *TColorSelectEditLink) EndEdit() bool {
	color := m.colorBtn.ButtonColor()
	logs.Debug("TColorSelectEditLink EndEdit", "color:", color, "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.BindData.EditNodeData.IntValue = int(color)
		m.VTree.EndEditNode()
		m.colorText.Hide()
		m.colorBtn.Hide()
	}
	return true
}

func (m *TColorSelectEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TColorSelectEditLink PrepareEdit")
	if m.colorText == nil || !m.colorText.IsValid() {
		m.CreateEdit()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column

	m.colorText.SetParent(m.VTree)
	m.colorBtn.SetParent(m.VTree)
	return true
}

func (m *TColorSelectEditLink) GetBounds() types.TRect {
	return m.colorText.BoundsRect()
}

func (m *TColorSelectEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TColorSelectEditLink ProcessMessage")
	lcl.ControlHelper.WindowProc(m.colorText, msg)
	lcl.ControlHelper.WindowProc(m.colorBtn, msg)
}

func (m *TColorSelectEditLink) SetBounds(R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetWidth(columnRect.Width() - 26)
	m.colorText.SetBoundsRect(R)
	logs.Debug("TColorSelectEditLink SetBounds", R)

	R.Left = R.Left + R.Width()
	R.SetWidth(24)
	m.colorBtn.SetBoundsRect(R)
	logs.Debug("TColorSelectEditLink SetBounds", R)
}

func (m *TColorSelectEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TColorSelectEditLink Destroy")
	m.colorBtn.SetOnColorChanged(nil) // 解除注册的事件
	m.colorText.Free()
	m.colorBtn.Free()
}
