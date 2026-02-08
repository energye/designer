// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"encoding/json"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
)

type TDesignerWebview struct {
	lcl.ICustomPanel
}

func (m *TDesignerWebview) Published() (props []lcl.ComponentProperties) {
	propStrList := tool.NewArray[string]()
	propStrList.Add(`{"name":"Align","value":"alCustom","kind":"tkEnumeration","type":"TAlign","options":"alNone,alTop,alBottom,alLeft,alRight,alClient,alCustom"}`)
	propStrList.Add(`{"name":"Anchors","value":"akTop,akLeft","kind":"tkSet","type":"TAnchors","options":"akTop,akLeft,akRight,akBottom"}`)
	propStrList.Add(`{"name":"Caption","value":"` + m.Caption() + `","kind":"tkAString","type":"TTranslateString","options":""}`)
	propStrList.Add(`{"name":"Width","value":"170","kind":"tkInteger","type":"LongInt","options":""}`)
	propStrList.Add(`{"name":"Height","value":"50","kind":"tkInteger","type":"LongInt","options":""}`)
	propStrList.Add(`{"name":"Top","value":"0","kind":"tkInteger","type":"LongInt","options":""}`)
	propStrList.Add(`{"name":"Left","value":"0","kind":"tkInteger","type":"LongInt","options":""}`)
	propStrList.Add(`{"name":"Name","value":"` + m.Name() + `","kind":"tkAString","type":"AnsiString","options":""}`)
	//propStrList.Add(`{"name":"Visible","value":"1","kind":"tkBool","type":"Boolean","options":""}`)

	props = make([]lcl.ComponentProperties, propStrList.Len())
	for i, prop := range propStrList.Values() {
		var propItem lcl.ComponentProperties
		if err := json.Unmarshal([]byte(prop), &propItem); err == nil {
			props[i] = propItem
		}
	}
	return
}

func NewDesignerWebview(owner lcl.IComponent) *TDesignerWebview {
	m := &TDesignerWebview{}
	m.ICustomPanel = lcl.NewCustomPanel(owner)
	m.SetParentColor(true)
	m.SetParentDoubleBuffered(true)
	m.SetBevelInner(types.BvNone)
	m.SetBevelOuter(types.BvNone)
	m.SetBorderStyleToBorderStyle(types.BsSingle)
	return m
}
