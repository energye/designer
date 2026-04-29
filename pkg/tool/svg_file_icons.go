package tool

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/energye/designer/pkg/svg/oksvg"
	"github.com/energye/designer/pkg/svg/rasterx"
	"github.com/energye/designer/resources"
	"image"
	"image/png"
	"path/filepath"
	"strings"
)

// SvgIconConfig 图标配置文件结构
type SvgIconConfig struct {
	IconDefinitions map[string]struct {
		IconPath string `json:"iconPath"`
	} `json:"iconDefinitions"`
	FileExtensions map[string]string `json:"fileExtensions"`
	LanguageIds    map[string]string `json:"languageIds"`
	File           string            `json:"file"`
	Folder         string            `json:"folder"`
}

var svgIconConfig *SvgIconConfig

func init() {
	configData := resources.Images("icons/svg_symbol-icon.json")
	if configData == nil {
		return
	}
	var config SvgIconConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return
	}
	svgIconConfig = &config

}

// GetSVGIconPath 根据文件扩展名获取图标路径
func GetSVGIconPath(fileName string) string {
	if svgIconConfig == nil {
		return ""
	}
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
	if ext == "" {
		ext = fileName
	}
	var (
		path string
		ok   bool
	)
	definit, ok := svgIconConfig.IconDefinitions[ext]
	if ok {
		return definit.IconPath
	}
	path, ok = svgIconConfig.FileExtensions[ext]
	if ok {
		definit, ok = svgIconConfig.IconDefinitions[path]
		if ok {
			return definit.IconPath
		}
	}
	path, ok = svgIconConfig.LanguageIds[ext]
	if ok {
		definit, ok = svgIconConfig.IconDefinitions[path]
		if ok {
			return definit.IconPath
		}
	}
	definit, ok = svgIconConfig.IconDefinitions[svgIconConfig.File]
	if ok {
		return definit.IconPath
	}
	return ""
}

// GetSVGIconData 根据文件名获取 svg 图标数据
func GetSVGIconData(fileName string) []byte {
	path := GetSVGIconPath(fileName)
	if path == "" {
		return nil
	}
	return resources.Images(path)
}

// SVGToPNG svg 转 png
func SVGToPNG(svgData []byte, targetWidth, targetHeight int) ([]byte, error) {
	if svgData == nil {
		return nil, errors.New("svgData is nil")
	}
	icon, err := oksvg.ReadIconStream(bytes.NewBuffer(svgData))
	if err != nil {
		return nil, err
	}
	w, h := float64(targetWidth), float64(targetHeight)
	icon.SetTarget(0, 0, w, h)
	rgba := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	scannerGV := rasterx.NewScannerGV(targetWidth, targetHeight, rgba, rgba.Bounds())
	rasterizer := rasterx.NewDasher(targetWidth, targetHeight, scannerGV)
	icon.Draw(rasterizer, 1.0)
	pngBuf := bytes.NewBuffer(nil)
	err = png.Encode(pngBuf, rgba)
	return pngBuf.Bytes(), err
}
