package vtedit

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

// 下拉框

type IComboBoxEditLink interface {
	lcl.ICustomVTEditLink
	AsIVTEditLink() lcl.IVTEditLink
	SetOnNewData(fn TOnNewData)
}

type TComboBoxEditLink struct {
	lcl.ICustomVTEditLink
	newData    TOnNewData
	edit       lcl.IComboBox
	tree       lcl.ILazVirtualStringTree
	node       types.PVirtualNode
	column     int32
	textBounds types.TRect
	text       string
	alignment  types.TAlignment
	stopping   bool
}
