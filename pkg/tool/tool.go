package tool

import (
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
)

func ImageListAddPng(imageList lcl.IImageList, filePath string) {
	if filePath == "" {
		return
	}
	data := resources.Images(filePath)
	if data != nil {
		pic := lcl.NewPicture()
		defer pic.Free()

		mem := lcl.NewMemoryStream()
		defer mem.Free()
		lcl.StreamHelper.WriteBuffer(mem, data)
		mem.SetPosition(0)
		pic.LoadFromStream(mem)
		imageList.Add(pic.Bitmap(), nil)
	}
}
