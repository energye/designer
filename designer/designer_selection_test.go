package designer

import (
	"testing"

	"github.com/energye/lcl/types"
)

func TestNormalizeRect(t *testing.T) {
	r := normalizeRect(types.Point(120, 90), types.Point(20, 10))
	if r.Left != 20 || r.Top != 10 || r.Right != 120 || r.Bottom != 90 {
		t.Fatalf("normalizeRect 错误: %+v", r)
	}
	if r.Width() != 100 || r.Height() != 80 {
		t.Fatalf("normalizeRect 尺寸错误: %v x %v", r.Width(), r.Height())
	}
	// 反向拖动
	r = normalizeRect(types.Point(0, 0), types.Point(50, 60))
	if r.Left != 0 || r.Top != 0 || r.Right != 50 || r.Bottom != 60 {
		t.Fatalf("normalizeRect 正向错误: %+v", r)
	}
}

func TestRectsIntersect(t *testing.T) {
	base := types.TRect{Left: 10, Top: 10, Right: 110, Bottom: 60}
	cases := []struct {
		rect types.TRect
		want bool
	}{
		{types.TRect{Left: 20, Top: 20, Right: 40, Bottom: 40}, true},    // 完全包含
		{types.TRect{Left: 0, Top: 0, Right: 20, Bottom: 20}, true},      // 部分重叠
		{types.TRect{Left: 110, Top: 10, Right: 200, Bottom: 60}, false}, // 贴右边不算相交
		{types.TRect{Left: 0, Top: 60, Right: 100, Bottom: 90}, false},   // 贴下边不算相交
		{types.TRect{Left: 200, Top: 200, Right: 300, Bottom: 300}, false},
	}
	for _, c := range cases {
		if got := rectsIntersect(base, c.rect); got != c.want {
			t.Fatalf("rectsIntersect(%+v) = %v, want %v", c.rect, got, c.want)
		}
	}
}
