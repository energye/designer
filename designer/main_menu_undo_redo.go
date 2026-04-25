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
	m.undoAction.SetOnExecute(m.undoActionOnExecute)
	m.undoAction.SetOnUpdate(m.undoActionOnUpdate)

	m.redoAction.SetOnExecute(m.redoActionOnExecute)
	m.redoAction.SetOnUpdate(m.redoActionOnUpdate)
}

func (m *TUndoRedo) undoActionOnExecute(sender lcl.IObject) {
	activeControl := lcl.Screen.ActiveControl()
	if activeControl != nil {
		clsName := activeControl.ClassName()
		switch clsName {
		case "TSynEdit":
			lcl.AsSynEdit(activeControl).Undo()
		default:
			m.undoAction.ExecuteTarget(activeControl)
		}
	}
}

func (m *TUndoRedo) undoActionOnUpdate(sender lcl.IObject) {
	activeControl := lcl.Screen.ActiveControl()
	if activeControl != nil {
		if activeControl.IsObjectInstanceOf(lcl.TCustomEditClass()) {
			m.undoAction.SetEnabled(lcl.AsCustomEdit(activeControl).CanUndo())
		} else if activeControl.IsObjectInstanceOf(lcl.TSynEditClass()) {
			m.undoAction.SetEnabled(lcl.AsSynEdit(activeControl).CanUndo())
		} else {
			m.undoAction.SetEnabled(false)
		}
	}
}

func (m *TUndoRedo) redoActionOnExecute(sender lcl.IObject) {
	activeControl := lcl.Screen.ActiveControl()
	if activeControl != nil {
		clsName := activeControl.ClassName()
		switch clsName {
		case "TSynEdit":
			lcl.AsSynEdit(activeControl).Redo()
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
