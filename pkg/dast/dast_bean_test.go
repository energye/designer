package dast

import (
	"github.com/energye/lcl/lcl"
	lclTypes "github.com/energye/lcl/types"
)

type TTestEvent func(pam1 int32, param2 string) string
type TContextPopupEvent func(sender lcl.IObject, mousePos lclTypes.TPoint, handled *bool)
type TStartDockEvent func(sender lcl.IObject, dragObject *lcl.IDragDockObject)

// Method2
func (m *TTestStruct) Method3(pam1 lcl.IObject, param2 string, param3 TTestStruct, param4 lclTypes.TRect, param5 ITestInterface) lcl.IPanel {
	return &lcl.TPanel{}
}
