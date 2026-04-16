package icns

import (
	"image"
	"os"
)

func PngToIcns(srcPngFile, destIcnsFile string) error {
	iconPngFile, err := os.Open(srcPngFile)
	if err != nil {
		return err
	}
	defer iconPngFile.Close()
	iconPngSrcImg, _, err := image.Decode(iconPngFile)
	if err != nil {
		return err
	}
	iconPngIcnsDest, err := os.Create(destIcnsFile)
	if err != nil {
		return err
	}
	defer iconPngIcnsDest.Close()
	if err := Encode(iconPngIcnsDest, iconPngSrcImg); err != nil {
		return err
	}
	return nil
}
