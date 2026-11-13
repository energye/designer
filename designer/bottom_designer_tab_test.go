package designer

import (
	"testing"

	"github.com/energye/energy/v3/lcl/wg"
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
