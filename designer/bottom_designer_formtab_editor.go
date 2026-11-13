// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/energye/designer/designer/editor"
	"github.com/energye/designer/designer/editor/webview"
	projBean "github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/metadata"
	"github.com/energye/energy/v3/core"
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"github.com/energye/lcl/types/colors"
)

const (
	designerRulerSize  = 24
	designerRulerMinor = 8
	designerRulerMajor = 50
	designerRulerLabel = 100
	designerRulerBase  = 2048
)

var (
	designerRulerBg   = colors.RGBToColor(245, 245, 245)
	designerRulerLine = colors.RGBToColor(160, 160, 160)
	designerRulerText = colors.RGBToColor(80, 80, 80)
)

var gFromEditor editor.IEditor

type IFormDesignCode interface {
	UIFile() string
	GOFile() string
	GOUserFile() string
}

// 窗体设计页, 只针对设计窗体和对应的代码 tab page
// 如果非设计窗口代码, 常规代码文件直接在 mainPage 初始化编辑器
type TFormDesignPage struct {
	//formTab            *wg.TTab  // 设计窗体Tab
	//formDesignPage     *wg.TPage // 设计窗体
	//formUserEditorPage *wg.TPage // 用户代码 tab page
	////formUIEditorPage   *wg.TPage      // UI 代码 tab page

	formDesignScroll lcl.IScrollBox // 设计窗体滚动条

	formTab *FormTab

	rulerTop        lcl.IPaintBox
	rulerLeft       lcl.IPaintBox
	rulerTopBitmap  lcl.IBitmap
	rulerLeftBitmap lcl.IBitmap
	rulerTopLen     int32
	rulerLeftLen    int32
	rulerStop       chan struct{}
	lastRulerX      int32
	lastRulerY      int32
	lastRulerW      int32
	lastRulerH      int32

	//
	formPageControl    lcl.IPageControl
	formDesignPage     lcl.ITabSheet // "窗体"
	formUserEditorPage lcl.ITabSheet // "代码"
	formUIEditorPage   lcl.ITabSheet // "UI代码" - 只读

	wvWindowParent core.WindowParent

	formDesignCode IFormDesignCode
}

func NewFormDesignPage(formTab *FormTab) *TFormDesignPage {
	m := &TFormDesignPage{formDesignCode: formTab, formTab: formTab}
	//m.formTab = wg.NewTab(formTab.mainPage)
	//m.formTab.SetBounds(0, 0, formTab.mainPage.Width(), formTab.mainPage.Height())
	//m.formTab.SetAlign(types.AlClient)
	//m.formTab.EnableScrollButton(false)
	//m.formTab.SetParent(formTab.mainPage)
	//
	//m.formDesignPage = m.formTab.NewPage()
	//m.formDesignPage.SetDoubleBuffered(true)
	//m.formDesignPage.Button().SetCaption("窗体")
	//m.formDesignPage.SetOnHide(func(sender lcl.IObject) {
	//	m.formDesignScroll.SetVisible(false)
	//})
	//m.formDesignPage.SetOnShow(func(sender lcl.IObject) {
	//	m.formDesignScroll.SetVisible(true)
	//})
	//setFormDesignPageStyle(m.formDesignPage, resources.Images("components/tform.png"))
	//
	//m.formUserEditorPage = m.formTab.NewPage()
	//m.formUserEditorPage.SetDoubleBuffered(true)
	//m.formUserEditorPage.Button().SetCaption("代码")
	//m.formUserEditorPage.SetOnShow(m.UserEditorPageOnShow)
	//name := tool.GetSVGIconPath("go")
	//codeIconPngData, _ := tool.SVGToPNG(resources.Images(name), 24, 24)
	//setFormDesignPageStyle(m.formUserEditorPage, codeIconPngData)

	m.formPageControl = lcl.NewPageControl(formTab.mainPage)
	m.formPageControl.SetAlign(types.AlClient)
	m.formPageControl.SetParent(formTab.mainPage)

	m.formDesignPage = lcl.NewTabSheet(formTab.mainPage)
	m.formDesignPage.SetName("FormDesignPage")
	m.formDesignPage.SetPageControl(m.formPageControl)
	m.formDesignPage.SetCaption(metadata.GI18n.Dict("FormDesignPage.Caption"))
	m.formDesignPage.SetOnShow(func(sender lcl.IObject) {
		m.invalidateRulers()
	})

	m.formUserEditorPage = lcl.NewTabSheet(formTab.mainPage)
	m.formUserEditorPage.SetName("FormUserEditorPage")
	m.formUserEditorPage.SetPageControl(m.formPageControl)
	m.formUserEditorPage.SetCaption(metadata.GI18n.Dict("FormUserEditorPage.Caption"))
	m.formUserEditorPage.SetOnShow(m.UserEditorPageOnShow)

	//m.formUIEditorPage = lcl.NewTabSheet(formTab.mainPage)
	//m.formUIEditorPage.SetPageControl(m.formPageControl)
	//m.formUIEditorPage.SetCaption("UI代码")
	//m.formUIEditorPage.SetOnShow(m.UIEditorPageOnShow)

	// 标尺在滚动框外层：横向 AlTop、纵向 AlLeft，滚动框最后 AlClient 占剩余区域。
	m.rulerTop = lcl.NewPaintBox(formTab.mainPage)
	m.rulerTop.SetParent(m.formDesignPage)
	m.rulerTop.SetAlign(types.AlTop)
	m.rulerTop.SetHeight(designerRulerSize)
	m.rulerTop.SetOnPaint(func(sender lcl.IObject) { m.paintRuler(true) })

	m.rulerLeft = lcl.NewPaintBox(formTab.mainPage)
	m.rulerLeft.SetParent(m.formDesignPage)
	m.rulerLeft.SetAlign(types.AlLeft)
	m.rulerLeft.SetWidth(designerRulerSize)
	m.rulerLeft.SetOnPaint(func(sender lcl.IObject) { m.paintRuler(false) })

	m.formDesignScroll = lcl.NewScrollBox(formTab.mainPage)
	m.formDesignScroll.SetAlign(types.AlClient)
	m.formDesignScroll.SetAutoScroll(true)
	if tool.IsDarwin {
		// fix: laz MacOS bug 默认隐藏滚动条, 手动控制显示
		// 该bug体现为当同时出现横竖滚动条时, UI 锁死崩溃, 在laz 4.6 复现
		hBar := m.formDesignScroll.HorzScrollBar()
		vBar := m.formDesignScroll.VertScrollBar()
		hBar.SetVisible(false)
		vBar.SetVisible(false)
		m.formDesignScroll.SetOnResize(func(sender lcl.IObject) {
			// 这里使用 < 小于做判断
			// 当滚动box宽或高小于设计窗体大小时显示或隐藏滚动条
			vBar.SetVisible(m.formDesignScroll.Height() < formTab.formDesigner.Form.Height())
			hBar.SetVisible(m.formDesignScroll.Width() < formTab.formDesigner.Form.Width())
		})
	}
	m.formDesignScroll.SetBorderStyleToBorderStyle(types.BsNone)
	m.formDesignScroll.SetDoubleBuffered(true)
	m.formDesignScroll.SetParent(m.formDesignPage)
	m.startRulerSync()

	//m.formTab.HideAllActivated()
	//m.formDesignPage.SetActive(true)
	//m.formTab.RecalculatePosition()

	return m
}

//func (m *TFormDesignPage) ActiveDesignPage() {
//	m.formTab.HideAllActivated()
//	m.formDesignPage.SetActive(true)
//}
//func (m *TFormDesignPage) ActiveEditorPage() {
//	m.formTab.HideAllActivated()
//	m.formUserEditorPage.SetActive(true)
//}

func (m *TFormDesignPage) UserEditorPageOnShow(sender lcl.IObject) {
	m.initEditor()
	m.SwitchTabPageEditor(false)
}

func (m *TFormDesignPage) UIEditorPageOnShow(sender lcl.IObject) {
	m.initEditor()
	m.SwitchTabPageEditor(true)
}

func (m *TFormDesignPage) initEditor() {
	if gFromEditor == nil {
		gFromEditor = editor.NewEditor(designer.tab)
		if gFromEditor != nil {
			if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
				m.wvWindowParent = wvEditor.Webview().WindowParent()
			}
		}
	} else {
		if m.wvWindowParent == nil {
			m.wvWindowParent = webview.NewWebviewWindowParent(designer.tab)
		}
	}
	if gFromEditor != nil {
		webview.SetOnGoToDefinition(goToDefinition)
	}
}

func (m *TFormDesignPage) SwitchTabPageEditor(uiCode bool) {
	if gFromEditor != nil && m.wvWindowParent != nil {
		if wvEditor, ok := gFromEditor.(webview.IWebviewEditor); ok {
			canLoad := make(chan error, 1)
			wvEditor.SetCanLoadChan(canLoad)
			go func() {
				_ = <-canLoad
				wvEditor.SetCanLoadChan(nil)
				close(canLoad)
				var (
					filePath string
					readOnly bool
				)
				if uiCode {
					filePath = filepath.Join(projBean.CodePath(), m.formDesignCode.GOFile())
					readOnly = true
				} else {
					filePath = filepath.Join(projBean.CodePath(), m.formDesignCode.GOUserFile())
					readOnly = false
				}
				if wvEditor.Initialized() {
					lcl.RunOnMainThreadAsync(func(id uint32) {
						editor.OpenFileInEditor(filePath, readOnly)
					})
				}
			}()
			var targetOwner lcl.IWinControl
			if uiCode && m.formUIEditorPage != nil {
				targetOwner = m.formUIEditorPage
			} else {
				targetOwner = m.formUserEditorPage
			}
			wvEditor.SwitchTabPage(targetOwner, m.wvWindowParent)
			wvEditor.CreateBrowser()
		}
	}
}

func (m *TFormDesignPage) ActiveCodeEditorTab() {
	activeIdx := m.formPageControl.ActivePageIndex()
	m.initEditor()
	m.SwitchTabPageEditor(activeIdx == 2)
	if activeIdx >= 1 {
		go func() {
			time.AfterFunc(time.Millisecond*200, func() {
				lcl.RunOnMainThreadAsync(func(id uint32) {
					m.formPageControl.SetActivePageIndex(activeIdx)
				})
			})
		}()
	}
}

// ==== 设计画布标尺 ====

// paintRuler 标尺每帧只绘制预渲染的刻度位图，并按滚动偏移平移。
func (m *TFormDesignPage) paintRuler(horizontal bool) {
	var pb lcl.IPaintBox
	if horizontal {
		pb = m.rulerTop
	} else {
		pb = m.rulerLeft
	}
	if pb == nil || !pb.IsValid() {
		return
	}
	canvas := pb.Canvas()
	if canvas == nil || !canvas.IsValid() {
		return
	}
	canvas.BrushToBrush().SetColor(designerRulerBg)
	canvas.FillRectWithRect(types.Rect(0, 0, pb.ClientWidth(), pb.ClientHeight()))

	m.ensureRulerBitmaps()
	if horizontal {
		if m.rulerTopBitmap != nil {
			canvas.DrawWithIntX2Graphic(m.rulerOriginX(), 0, m.rulerTopBitmap)
		}
	} else {
		if m.rulerLeftBitmap != nil {
			canvas.DrawWithIntX2Graphic(0, m.rulerOriginY(), m.rulerLeftBitmap)
		}
	}
}

func (m *TFormDesignPage) rulerOriginX() int32 {
	if m.formTab == nil || m.formTab.formDesigner == nil || m.formTab.formDesigner.Form == nil {
		return 0
	}
	// 顶部标尺从页面 x=0 开始，而滚动框在左侧标尺右边，需要补上滚动框 Left。
	return m.formDesignScroll.Left() + m.formTab.formDesigner.Form.Left() -
		m.formDesignScroll.HorzScrollBar().Position()
}

func (m *TFormDesignPage) rulerOriginY() int32 {
	if m.formTab == nil || m.formTab.formDesigner == nil || m.formTab.formDesigner.Form == nil {
		return 0
	}
	return m.formTab.formDesigner.Form.Top() - m.formDesignScroll.VertScrollBar().Position()
}

func rulerBucketLength(required int32) int32 {
	length := int32(designerRulerBase)
	for length < required {
		length *= 2
	}
	return length
}

func (m *TFormDesignPage) ensureRulerBitmaps() {
	if m.formTab == nil || m.formTab.formDesigner == nil || m.formTab.formDesigner.Form == nil {
		return
	}
	w, h := m.formTab.formDesigner.Form.Width(), m.formTab.formDesigner.Form.Height()
	if w <= 0 || h <= 0 {
		return
	}
	topLen := rulerBucketLength(w)
	if m.rulerTopBitmap == nil || m.rulerTopLen != topLen {
		m.buildRulerBitmap(true, topLen)
	}
	leftLen := rulerBucketLength(h)
	if m.rulerLeftBitmap == nil || m.rulerLeftLen != leftLen {
		m.buildRulerBitmap(false, leftLen)
	}
}

func (m *TFormDesignPage) buildRulerBitmap(horizontal bool, length int32) {
	bmp := lcl.NewBitmap()
	if horizontal {
		bmp.SetSize(length, designerRulerSize)
	} else {
		bmp.SetSize(designerRulerSize, length)
	}
	canvas := bmp.Canvas()
	canvas.BrushToBrush().SetColor(designerRulerBg)
	canvas.FillRectWithRect(types.Rect(0, 0, bmp.Width(), bmp.Height()))
	pen := canvas.PenToPen()
	pen.SetColor(designerRulerLine)
	pen.SetWidth(1)
	pen.SetStyle(types.PsSolid)
	font := canvas.FontToFont()
	font.SetColor(designerRulerText)
	font.SetSize(7)

	if horizontal {
		bottom := int32(designerRulerSize - 1)
		canvas.MoveToWithIntX2(0, bottom)
		canvas.LineToWithIntX2(length-1, bottom)
		for pos := int32(0); pos < length; pos += designerRulerMinor {
			if pos%designerRulerMajor == 0 {
				continue
			}
			canvas.MoveToWithIntX2(pos, bottom)
			canvas.LineToWithIntX2(pos, bottom-3)
		}
		for pos := int32(0); pos < length; pos += designerRulerMajor {
			if pos%designerRulerLabel == 0 {
				continue
			}
			canvas.MoveToWithIntX2(pos, bottom)
			canvas.LineToWithIntX2(pos, bottom-7)
		}
		for pos := int32(0); pos < length; pos += designerRulerLabel {
			canvas.MoveToWithIntX2(pos, bottom)
			canvas.LineToWithIntX2(pos, bottom-11)
			canvas.TextOutWithIntX2Str(pos+2, 1, strconv.Itoa(int(pos)))
		}
		if m.rulerTopBitmap != nil {
			m.rulerTopBitmap.Free()
		}
		m.rulerTopBitmap, m.rulerTopLen = bmp, length
		return
	}

	right := int32(designerRulerSize - 1)
	canvas.MoveToWithIntX2(right, 0)
	canvas.LineToWithIntX2(right, length-1)
	for pos := int32(0); pos < length; pos += designerRulerMinor {
		if pos%designerRulerMajor == 0 {
			continue
		}
		canvas.MoveToWithIntX2(right, pos)
		canvas.LineToWithIntX2(right-3, pos)
	}
	for pos := int32(0); pos < length; pos += designerRulerMajor {
		if pos%designerRulerLabel == 0 {
			continue
		}
		canvas.MoveToWithIntX2(right, pos)
		canvas.LineToWithIntX2(right-7, pos)
	}
	for pos := int32(0); pos < length; pos += designerRulerLabel {
		canvas.MoveToWithIntX2(right, pos)
		canvas.LineToWithIntX2(right-11, pos)
		canvas.TextOutWithIntX2Str(1, pos+1, strconv.Itoa(int(pos)))
	}
	if m.rulerLeftBitmap != nil {
		m.rulerLeftBitmap.Free()
	}
	m.rulerLeftBitmap, m.rulerLeftLen = bmp, length
}

func (m *TFormDesignPage) freeRulerBitmaps() {
	if m.rulerTopBitmap != nil {
		m.rulerTopBitmap.Free()
		m.rulerTopBitmap = nil
	}
	if m.rulerLeftBitmap != nil {
		m.rulerLeftBitmap.Free()
		m.rulerLeftBitmap = nil
	}
	m.rulerTopLen, m.rulerLeftLen = 0, 0
}

// startRulerSync 启动轻量滚动同步：仅滚动位置/窗体尺寸变化时重绘标尺。
func (m *TFormDesignPage) startRulerSync() {
	if m.rulerStop != nil {
		return
	}
	stop := make(chan struct{})
	m.rulerStop = stop
	go func() {
		ticker := time.NewTicker(60 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				lcl.RunOnMainThreadAsync(func(id uint32) {
					m.syncRulers()
				})
			}
		}
	}()
	m.invalidateRulers()
}

func (m *TFormDesignPage) stopRulerSync() {
	if m.rulerStop != nil {
		close(m.rulerStop)
		m.rulerStop = nil
	}
}

// closeRuler 释放标尺资源；停止后已排队的 UI 回调会因 rulerStop=nil 直接返回。
func (m *TFormDesignPage) closeRuler() {
	m.stopRulerSync()
	m.freeRulerBitmaps()
	m.formTab = nil
}

func (m *TFormDesignPage) invalidateRulers() {
	if m.rulerTop != nil && m.rulerTop.IsValid() {
		m.rulerTop.Invalidate()
	}
	if m.rulerLeft != nil && m.rulerLeft.IsValid() {
		m.rulerLeft.Invalidate()
	}
}

func (m *TFormDesignPage) syncRulers() {
	if m.rulerStop == nil || m.formTab == nil || m.formTab.mainPage == nil || !m.formTab.mainPage.Visible() {
		return
	}
	if m.formDesignPage == nil || !m.formDesignPage.Visible() {
		return
	}
	if m.formTab.formDesigner == nil || m.formTab.formDesigner.Form == nil {
		return
	}
	form := m.formTab.formDesigner.Form
	x := m.formDesignScroll.HorzScrollBar().Position()
	y := m.formDesignScroll.VertScrollBar().Position()
	w, h := form.Width(), form.Height()
	if x == m.lastRulerX && y == m.lastRulerY && w == m.lastRulerW && h == m.lastRulerH {
		return
	}
	m.lastRulerX, m.lastRulerY, m.lastRulerW, m.lastRulerH = x, y, w, h
	m.invalidateRulers()
}
