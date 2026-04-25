//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License 2.0
//
//----------------------------------------

package designer

import "github.com/energye/lcl/lcl"

type TUndoRedo struct {
	undoAction lcl.IEditUndo
	redoAction lcl.IAction
}

func (m *TUndoRedo) init() {
	m.redoAction.SetOnExecute(m.redoActionOnExecute)
	m.redoAction.SetOnUpdate(m.redoActionOnUpdate)
}

func (m *TUndoRedo) redoActionOnExecute(sender lcl.IObject) {
	activeControl := lcl.Screen.ActiveControl()
	if activeControl != nil {
		clsName := activeControl.ClassName()
		switch clsName {
		case "TSynEdit":
			m.redoAction.ExecuteTarget(activeControl)
		default:
			m.undoAction.Execute()
		}
	}
}

func (m *TUndoRedo) redoActionOnUpdate(sender lcl.IObject) {
	activeControl := lcl.Screen.ActiveControl()
	if activeControl != nil {
		clsName := activeControl.ClassName()
		// activeControl.IsObjectInstanceOf(lcl.TSynEditClass())
		//fmt.Println("IsObjectInstanceOf:", activeControl.IsObjectInstanceOf(lcl.TSynEditClass()))
		switch clsName {
		case "TSynEdit":
			m.redoAction.SetEnabled(lcl.AsSynEdit(activeControl).CanRedo())
		default:
			m.redoAction.SetEnabled(m.undoAction.Enabled())
		}
	}
}
