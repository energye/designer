package designer

import (
	"testing"

	"github.com/energye/lcl/types"
)

func TestDesignerAlignRect(t *testing.T) {
	ref := types.TRect{Left: 100, Top: 100, Right: 300, Bottom: 200}
	target := types.TRect{Left: 20, Top: 30, Right: 80, Bottom: 70}
	cases := []struct {
		kind designerAlignKind
		want types.TRect
	}{
		{designerAlignLeft, types.TRect{Left: 100, Top: 30, Right: 160, Bottom: 70}},
		{designerAlignHCenter, types.TRect{Left: 170, Top: 30, Right: 230, Bottom: 70}},
		{designerAlignRight, types.TRect{Left: 240, Top: 30, Right: 300, Bottom: 70}},
		{designerAlignTop, types.TRect{Left: 20, Top: 100, Right: 80, Bottom: 140}},
		{designerAlignVCenter, types.TRect{Left: 20, Top: 130, Right: 80, Bottom: 170}},
		{designerAlignBottom, types.TRect{Left: 20, Top: 160, Right: 80, Bottom: 200}},
	}
	for _, c := range cases {
		if got := designerAlignRect(c.kind, ref, target); got != c.want {
			t.Fatalf("designerAlignRect(%v) = %+v, want %+v", c.kind, got, c.want)
		}
	}
}

func TestDesignerSameSizeRect(t *testing.T) {
	ref := types.TRect{Left: 0, Top: 0, Right: 200, Bottom: 100}
	target := types.TRect{Left: 20, Top: 30, Right: 80, Bottom: 70}
	cases := []struct {
		kind designerSameSizeKind
		want types.TRect
	}{
		{designerSameWidth, types.TRect{Left: 20, Top: 30, Right: 220, Bottom: 70}},
		{designerSameHeight, types.TRect{Left: 20, Top: 30, Right: 80, Bottom: 130}},
		{designerSameSize, types.TRect{Left: 20, Top: 30, Right: 220, Bottom: 130}},
	}
	for _, c := range cases {
		if got := designerSameSizeRect(c.kind, ref, target); got != c.want {
			t.Fatalf("designerSameSizeRect(%v) = %+v, want %+v", c.kind, got, c.want)
		}
	}
}

func TestDesignerDistributePositions(t *testing.T) {
	rects := []types.TRect{
		{Left: 0, Top: 0, Right: 100, Bottom: 20},
		{Left: 130, Top: 0, Right: 150, Bottom: 20},
		{Left: 300, Top: 0, Right: 360, Bottom: 20},
	}
	got := designerDistributePositions(rects, true)
	want := []int32{0, 190, 300}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("水平分布位置 = %v, want %v", got, want)
		}
	}

	vertical := []types.TRect{
		{Left: 0, Top: 10, Right: 20, Bottom: 30},
		{Left: 0, Top: 40, Right: 20, Bottom: 60},
		{Left: 0, Top: 100, Right: 20, Bottom: 140},
	}
	got = designerDistributePositions(vertical, false)
	want = []int32{10, 55, 100}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("垂直分布位置 = %v, want %v", got, want)
		}
	}
}

func TestDesignerDistributePositionsKeepEnds(t *testing.T) {
	// 少于 3 个组件不改变位置
	rects := []types.TRect{
		{Left: 10, Top: 0, Right: 50, Bottom: 20},
		{Left: 80, Top: 0, Right: 100, Bottom: 20},
	}
	got := designerDistributePositions(rects, true)
	if got[0] != 10 || got[1] != 80 {
		t.Fatalf("两个组件不应分布: %v", got)
	}
	// 重叠时允许负间距，但首尾仍保持不动
	overlap := []types.TRect{
		{Left: 0, Top: 0, Right: 100, Bottom: 20},
		{Left: 10, Top: 0, Right: 110, Bottom: 20},
		{Left: 20, Top: 0, Right: 120, Bottom: 20},
	}
	got = designerDistributePositions(overlap, true)
	if got[0] != 0 || got[2] != 20 {
		t.Fatalf("重叠分布首尾应保持不动: %v", got)
	}
}
