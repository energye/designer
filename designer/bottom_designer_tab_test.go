package designer

import (
	"testing"

	"github.com/energye/designer/consts"
	"github.com/energye/energy/v3/lcl/wg"
	"github.com/energye/lcl/types"
)

func TestFormTabNeedsRestore(t *testing.T) {
	cases := []struct {
		state FormTabState
		want  bool
	}{
		{FtsNone, false},
		{FtsClose, false},
		{FtsHide, true},
		{FtsShow, false},
	}
	for _, c := range cases {
		if got := formTabNeedsRestore(c.state); got != c.want {
			t.Fatalf("formTabNeedsRestore(%v) = %v, want %v", c.state, got, c.want)
		}
	}
}

func TestHiddenFormTabForPage(t *testing.T) {
	hiddenPage := &wg.TPage{}
	visiblePage := &wg.TPage{}
	hiddenForm := &FormTab{mainPage: hiddenPage, State: FtsHide}
	forms := map[int]*FormTab{
		1: hiddenForm,
		2: {mainPage: visiblePage, State: FtsShow},
		3: nil,
	}

	if got := hiddenFormTabForPage(forms, hiddenPage); got != hiddenForm {
		t.Fatalf("应命中隐藏窗体页, got %v", got)
	}
	if got := hiddenFormTabForPage(forms, visiblePage); got != nil {
		t.Fatalf("已显示的窗体页不应恢复, got %v", got)
	}
	if got := hiddenFormTabForPage(forms, &wg.TPage{}); got != nil {
		t.Fatalf("未知页不应命中, got %v", got)
	}
	if got := hiddenFormTabForPage(forms, nil); got != nil {
		t.Fatalf("nil 页不应命中, got %v", got)
	}
}

func TestSnapDesignerGrid(t *testing.T) {
	cases := []struct {
		in, want int32
	}{
		{0, 0}, {3, 0}, {4, 8}, {8, 8}, {13, 16},
		{-3, 0}, {-4, -8}, {-13, -16},
	}
	for _, c := range cases {
		if got := snapDesignerGrid(c.in); got != c.want {
			t.Fatalf("snapDesignerGrid(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestSnapGridEnabled(t *testing.T) {
	formTab := &FormTab{gridVisible: true}
	d := &drag{relation: &TDesigningComponent{FormTab: formTab}}
	if !d.snapGridEnabled(0) {
		t.Fatal("网格显示时应启用吸附")
	}
	if d.snapGridEnabled(types.NewSet(types.SsAlt)) {
		t.Fatal("Alt 应临时关闭吸附")
	}
	formTab.gridVisible = false
	if d.snapGridEnabled(0) {
		t.Fatal("网格隐藏时不应吸附")
	}
}

func TestGetMouseMsgShiftAlt(t *testing.T) {
	d := &TEngFormDesigner{}
	shift, _ := d.GetMouseMsgShift(&types.TLMMouse{Keys: consts.MK_ALT})
	if !shift.In(types.SsAlt) {
		t.Fatal("Alt 键应写入 ShiftState")
	}
	shift, _ = d.GetMouseMsgShift(&types.TLMMouse{Keys: consts.MK_ALT | consts.MK_CONTROL})
	if !shift.In(types.SsAlt) || !shift.In(types.SsCtrl) {
		t.Fatal("Alt+Ctrl 组合应同时生效")
	}
}

func TestRulerClosedGuards(t *testing.T) {
	page := &TFormDesignPage{rulerClosed: true}
	// 关闭后这些入口必须直接返回，不能访问已释放的 LCL 页面/窗体。
	page.syncRulers()
	page.invalidateRulers()
	page.startRulerSync()
	page.closeRuler()
	page.closeRuler()
}
