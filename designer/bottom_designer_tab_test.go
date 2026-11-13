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

func TestSnapEdgeToTargets(t *testing.T) {
	targets := []designerAlignTarget{
		{pos: 100, from: 10, to: 40},
		{pos: 200, from: 0, to: 20},
	}
	edge, target, hit := snapEdgeToTargets(103, targets, designerAlignThreshold)
	if !hit || edge != 100 || target.pos != 100 || target.from != 10 || target.to != 40 {
		t.Fatalf("edge snap failed: edge=%d target=%+v hit=%v", edge, target, hit)
	}
	edge, _, hit = snapEdgeToTargets(110, targets, designerAlignThreshold)
	if hit || edge != 110 {
		t.Fatalf("超出阈值不应吸附: edge=%d hit=%v", edge, hit)
	}
}

func TestSnapMoveAxis(t *testing.T) {
	target := func(pos int32) []designerAlignTarget {
		return []designerAlignTarget{{pos: pos, from: 0, to: 10}}
	}
	// 左边对齐
	start, matched, hit := snapMoveAxis(98, 50, target(100), designerAlignThreshold)
	if !hit || start != 100 || matched.pos != 100 {
		t.Fatalf("左边吸附失败: start=%d matched=%+v", start, matched)
	}
	// 右边对齐
	start, matched, hit = snapMoveAxis(100, 50, target(150), designerAlignThreshold)
	if !hit || start != 100 || matched.pos != 150 {
		t.Fatalf("右边吸附失败: start=%d matched=%+v", start, matched)
	}
	// 中心对齐
	start, matched, hit = snapMoveAxis(100, 50, target(127), designerAlignThreshold)
	if !hit || start != 102 || matched.pos != 127 {
		t.Fatalf("中心吸附失败: start=%d matched=%+v", start, matched)
	}
}

func TestGuideFor(t *testing.T) {
	// 拖动控件 [20,80] 与目标 [10,40] 取联合区间
	guide := guideFor(20, 60, designerAlignTarget{pos: 100, from: 10, to: 40})
	if guide.pos != 100 || guide.from != 10 || guide.to != 80 {
		t.Fatalf("联合区间错误: %+v", guide)
	}
	// 拖动控件 [0,20] 与目标 [10,40] 取联合区间
	guide = guideFor(0, 20, designerAlignTarget{pos: 5, from: 10, to: 40})
	if guide.from != 0 || guide.to != 40 {
		t.Fatalf("联合区间错误: %+v", guide)
	}
}
