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

package codegen

import (
	"github.com/energye/designer/consts"
	codegentpl "github.com/energye/designer/internal/templates/codegen"
	"github.com/energye/designer/pkg/tool"
	"go/format"
	"strings"
	"testing"
	"text/template"
)

func newTestFormData(browser *TComponentData) *TFormData {
	form := &TComponentData{
		Name:      "Form1",
		ClassName: "TForm1",
		Mod:       consts.ModLCL,
		Type:      consts.CtForm,
		Children:  []*TComponentData{browser},
	}
	browser.Parent = form
	return &TFormData{
		BaseInfo:    &TBaseInfo{},
		PackageName: "app",
		Imports:     tool.NewHashSet[string](),
		Form:        form,
	}
}

func TestCEFChromiumComponentCode(t *testing.T) {
	chromium := &TComponentData{
		Name:      "Chromium1",
		ClassName: "TChromium",
		Mod:       consts.ModCEFEnergy,
		Type:      consts.CtCustomVisual,
	}
	data := newTestFormData(chromium)
	data.Imports.Add(lcl)
	data.Imports.Add(energyWindow)
	data.Imports.Add(lclTypes)
	data.Imports.Add(cefEnergy)

	if got, want := chromium.GoIntfName(), "cef.IEmbeddedBrowser"; got != want {
		t.Fatalf("GoIntfName() = %q, want %q", got, want)
	}
	if got, want := chromium.GoNewObject(), "m.Chromium1 = cef.NewEmbeddedBrowser(m)\n"; got != want {
		t.Fatalf("GoNewObject() = %q, want %q", got, want)
	}
	if got, want := chromium.GoSetObjectParent(), "m.Chromium1.SetParent(m)\n"; got != want {
		t.Fatalf("GoSetObjectParent() = %q, want %q", got, want)
	}
	if got, want := chromium.GoSetWindow(), "m.Chromium1.SetWindow(m)\n"; got != want {
		t.Fatalf("GoSetWindow() = %q, want %q", got, want)
	}

	tmpl, err := template.New("cef").Parse(codegentpl.CEFAutoCodeTemplate)
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	if err = tmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}
	code, err := format.Source([]byte(buf.String()))
	if err != nil {
		t.Fatalf("CEF generated code format error: %v\n%s", err, buf.String())
	}
	for _, want := range []string{
		"Chromium1 cef.IEmbeddedBrowser",
		"m.Chromium1 = cef.NewEmbeddedBrowser(m)",
		"m.Chromium1.SetParent(m)",
		"m.Chromium1.SetWindow(m)",
	} {
		if !strings.Contains(string(code), want) {
			t.Fatalf("CEF generated code missing %q:\n%s", want, code)
		}
	}
}

func TestWVWebviewComponentCode(t *testing.T) {
	webview := &TComponentData{
		Name:      "Webview1",
		ClassName: "TWebview",
		Mod:       consts.ModWVEnergy,
		Type:      consts.CtCustomVisual,
	}
	data := newTestFormData(webview)
	data.Imports.Add(lcl)
	data.Imports.Add(energyWindow)
	data.Imports.Add(lclTypes)
	data.Imports.Add(wvEnergy)

	if got, want := webview.GoIntfName(), "wv.IWebview"; got != want {
		t.Fatalf("GoIntfName() = %q, want %q", got, want)
	}
	if got, want := webview.GoSetObjectParent(), "m.Webview1.SetParent(m)\n"; got != want {
		t.Fatalf("GoSetObjectParent() = %q, want %q", got, want)
	}
	if got, want := webview.GoSetWindow(), "m.Webview1.SetWindow(m)\n"; got != want {
		t.Fatalf("GoSetWindow() = %q, want %q", got, want)
	}

	tmpl, err := template.New("wv").Parse(codegentpl.WVAutoCodeTemplate)
	if err != nil {
		t.Fatal(err)
	}
	var buf strings.Builder
	if err = tmpl.Execute(&buf, data); err != nil {
		t.Fatal(err)
	}
	code, err := format.Source([]byte(buf.String()))
	if err != nil {
		t.Fatalf("WV generated code format error: %v\n%s", err, buf.String())
	}
	for _, want := range []string{
		"Webview1 wv.IWebview",
		"m.Webview1 = wv.NewWebview(m)",
		"m.Webview1.SetParent(m)",
		"m.Webview1.SetWindow(m)",
	} {
		if !strings.Contains(string(code), want) {
			t.Fatalf("WV generated code missing %q:\n%s", want, code)
		}
	}
}
