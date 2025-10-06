package message

import (
	"github.com/energye/lcl/lcl"
	"github.com/energye/lcl/types"
	"time"
)

// 消息弹框

type TMessage struct {
	form        lcl.IForm
	alphaStep   byte
	showTimer   lcl.ITimer
	afterTime   *time.Timer
	displayTime int32
	content     lcl.ILabel
}

var message *TMessage

func mustMessage() {
	if message == nil {
		message = new(TMessage)
		message.alphaStep = 5
		form := lcl.NewForm(nil)
		message.form = form
		form.SetBorderStyleToFormBorderStyle(types.BsNone)
		form.Canvas().SetAntialiasingMode(types.AmOn)
		form.SetControlStyle(form.ControlStyle().Include(types.CsParentBackground))
		form.SetFormStyle(types.FsStayOnTop)
		form.SetAlphaBlend(true)
		form.SetAlphaBlendValue(0)

		form.SetOnPaint(message.OnPaint)
		form.SetOnClick(message.OnClick)

		showTimer := lcl.NewTimer(form)
		showTimer.SetEnabled(false)
		showTimer.SetInterval(15)
		showTimer.SetOnTimer(message.OnShowTimer)
		message.showTimer = showTimer

		content := lcl.NewLabel(form)
		content.SetParent(form)
		message.content = content
	}
}

func rect(width, height int32) types.TRect {
	rect := lcl.Application.MainForm().BoundsRect()
	rect.Left = rect.Left + rect.Width()/2 - width/2
	rect.Top = rect.Top + rect.Height()/2 - height/2
	rect.SetWidth(width)
	rect.SetHeight(height)
	return rect
}

func Info(title, content string, width, height int32) {
	mustMessage()
	message.Hide()
	message.displayTime = 3 // 秒
	message.form.SetBoundsRect(rect(width, height))
	message.content.SetCaption(title + "\n  " + content)
	message.showTimer.SetEnabled(true)
	message.form.Show()
}

func (m *TMessage) OnShowTimer(sender lcl.IObject) {
	abv := m.form.AlphaBlendValue()
	if abv >= 255 {
		m.showTimer.SetEnabled(false)
		// displayTime 秒后关闭
		m.afterTime = time.AfterFunc(time.Second*time.Duration(m.displayTime), func() {
			lcl.RunOnMainThreadAsync(func(id uint32) {
				m.Hide()
			})
		})
		return
	}
	m.form.SetAlphaBlendValue(abv + m.alphaStep)
}

func (m *TMessage) OnPaint(sender lcl.IObject) {

}

func (m *TMessage) Hide() {
	if m.afterTime != nil {
		m.afterTime.Stop()
		m.afterTime = nil
	}
	m.showTimer.SetEnabled(false)
	m.form.SetAlphaBlendValue(0)
	m.form.Hide()
}

func (m *TMessage) OnClick(sender lcl.IObject) {
	m.Hide()
}
