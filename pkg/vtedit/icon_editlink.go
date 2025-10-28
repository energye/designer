package vtedit

import (
	"github.com/energye/designer/pkg/editorform"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 文本编辑框

type TIconEditLink struct {
	*TBaseEditLink
	btn       lcl.IBitBtn
	bounds    types.TRect
	alignment types.TAlignment
	stopping  bool
}

func NewIconEditLink(bindData *TEditNodeData) *TIconEditLink {
	link := new(TIconEditLink)
	link.TBaseEditLink = NewEditLink(link)
	link.BindData = bindData
	link.Create()
	return link
}

func (m *TIconEditLink) Create() {
	logs.Debug("TIconEditLink Create")
	m.btn = lcl.NewBitBtn(nil)
	m.btn.SetVisible(false)
	m.btn.SetAutoSize(false)
	m.btn.SetDoubleBuffered(true)
	m.btn.SetImages(tool.LoadImageList(nil, []string{"button/icon.png"}, 37, 26))
	m.btn.SetImageIndex(0)
	m.btn.SetCaption(m.BindData.EditStringValue())
	m.btn.SetOnClick(func(sender lcl.IObject) {
		editorform.NewGraphicPropertyEditor(func(filePath string, ok bool) {
			logs.Debug("TIconEditLink callback 图片目录:", filePath, ok)
		}).ShowModal()
	})
	m.btn.SetLayout(types.BlGlyphRight)
	textFont := m.btn.Font()
	textFont.SetStyle(textFont.Style().Include(types.FsBold))
	textFont.SetColor(0x2D5BC4)
}

func (m *TIconEditLink) BeginEdit() bool {
	logs.Debug("TIconEditLink BeginEdit")
	if !m.stopping {
		m.btn.Show()
		m.btn.SetFocus()
	}
	return true
}

func (m *TIconEditLink) CancelEdit() bool {
	logs.Debug("TIconEditLink CancelEdit")
	if !m.stopping {
		m.stopping = true
		m.btn.Hide()
		if m.VTree != nil {
			m.VTree.CancelEditNode()
		}
	}
	return true
}

func (m *TIconEditLink) EndEdit() bool {
	logs.Debug("TIconEditLink EndEdit", "m.stopping:", m.stopping)
	if !m.stopping {
		m.stopping = true
		m.btn.Hide()
		if m.VTree != nil {
			m.VTree.EndEditNode()
		}
	}
	return true
}

func (m *TIconEditLink) PrepareEdit(tree lcl.ILazVirtualStringTree, node types.PVirtualNode, column int32) bool {
	logs.Debug("TIconEditLink PrepareEdit")
	if m.btn == nil || !m.btn.IsValid() {
		m.Create()
	}
	m.VTree = tree
	m.Node = node
	m.Column = column
	m.btn.SetParent(m.VTree)
	m.btn.HandleNeeded()
	if column <= -1 {
		m.btn.SetBiDiMode(m.VTree.BiDiMode())
		m.alignment = m.VTree.Alignment()
	} else {
		columns := m.VTree.Header().Columns()
		m.btn.SetBiDiMode(columns.ItemsWithColumnIndexToVirtualTreeColumn(column).BiDiMode())
		m.alignment = columns.ItemsWithColumnIndexToVirtualTreeColumn(column).Alignment()
	}

	if m.btn.BiDiMode() != types.BdLeftToRight {
		switch m.alignment {
		case types.TaLeftJustify:
			m.alignment = types.TaRightJustify
		case types.TaRightJustify:
			m.alignment = types.TaLeftJustify
		}
	}
	return true
}

func (m *TIconEditLink) GetBounds() types.TRect {
	return m.btn.BoundsRect()
}

func (m *TIconEditLink) ProcessMessage(msg *types.TLMessage) {
	logs.Debug("TIconEditLink ProcessMessage")
}

func (m *TIconEditLink) SetBounds(R types.TRect) {
	columnRect := m.VTree.GetDisplayRect(m.Node, m.Column, false, false, true)
	R.Left = columnRect.Left
	R.Top = columnRect.Top
	R.SetHeight(columnRect.Height())
	R.SetWidth(columnRect.Width())
	m.btn.SetBoundsRect(R)
	logs.Debug("TIconEditLink SetBounds", R)
}

func (m *TIconEditLink) Destroy(sender lcl.IObject) {
	logs.Debug("TIconEditLink Destroy")
	m.btn.Free()
}
