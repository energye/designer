package tool

import (
	"github.com/energye/designer/resources"
	"github.com/energye/lcl/lcl"
	"os"
	"strings"
)

// 加载图像到列表
func ImageListAddPng(imageList lcl.IImageList, filePath string) {
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
	} else {
		data = resources.Images("components/default_150.png")
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

func Equal(s1, s2 string) bool {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)
	return s1 == s2
}

func FirstToUpper(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		} else if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
