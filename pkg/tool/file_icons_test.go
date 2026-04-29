package tool

import (
	"github.com/energye/designer/resources"
	"testing"
)

func TestSVGToPNG(t *testing.T) {
	svgData := resources.Images("icons/files_svg/angular.svg")
	pngData, err := SVGToPNG(svgData, 24, 24)
	if err != nil {
		t.Errorf("SVGToPNG failed: %v", err)
	}

	if len(pngData) == 0 {
		t.Error("SVGToPNG returned empty PNG data")
	}

	if len(pngData) > 0 && pngData[0] != 0x89 {
		t.Error("Returned data does not appear to be valid PNG format")
	}
}
