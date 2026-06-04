// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package project

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/energye/designer/consts"
	projecttpl "github.com/energye/designer/internal/templates/project"
	"github.com/energye/designer/options/bean"
	"github.com/energye/designer/pkg/config"
	"github.com/energye/designer/pkg/icns"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/pkg/winicon"
	"github.com/energye/designer/pkg/winres"
	"github.com/energye/designer/pkg/winres/version"
	"github.com/energye/designer/resources"
)

type ConflictKind string

const (
	ConflictExistingProject ConflictKind = "existing_project"
	ConflictNotEmpty        ConflictKind = "not_empty"
)

type ConflictDecision int

const (
	ConflictCancel ConflictDecision = iota
	ConflictOverwrite
)

type Conflict struct {
	Kind    ConflictKind
	Dir     string
	File    string
	Message string
}

type CreateOptions struct {
	Name               string
	Dir                string
	GUIRenderFramework bean.GUIRenderFramework
	FrameworkVersion   string
	OnConflict         func(Conflict) ConflictDecision
}

type CreateResult struct {
	Project *bean.TProject
	Dir     string
	EGPPath string
}

func Create(options CreateOptions) (*CreateResult, error) {
	options.Name = strings.TrimSpace(options.Name)
	options.Dir = strings.TrimSpace(options.Dir)
	if options.Name == "" {
		return nil, errors.New("project name is empty")
	}
	if options.Dir == "" {
		return nil, errors.New("project directory is empty")
	}
	if options.GUIRenderFramework == "" {
		options.GUIRenderFramework = bean.GUIRenderFramework_LCL
	}
	if !isSupportedGUIRenderFramework(options.GUIRenderFramework) {
		return nil, fmt.Errorf("unsupported GUI render framework: %s", options.GUIRenderFramework)
	}
	absDir, err := resolveCreateDir(options.Name, options.Dir)
	if err != nil {
		return nil, err
	}
	if err = checkCreateDir(absDir, options.OnConflict); err != nil {
		return nil, err
	}

	newProject := &bean.TProject{
		Name:               options.Name,
		EGPName:            options.Name + consts.EGPExt,
		Main:               "main.go",
		GUIRenderFramework: options.GUIRenderFramework,
		FrameworkVersion:   options.FrameworkVersion,
		Package:            consts.AppPackageName,
	}
	newProject.InitAppOption()
	newProject.InitBuildOption()

	if err = WriteEGPConfig(absDir, newProject); err != nil {
		bean.GPath = ""
		bean.GProject = nil
		return nil, err
	}
	bean.GPath = absDir
	bean.GProject = newProject
	if err = createProjectDir(); err != nil {
		return nil, err
	}
	if err = saveOrUpdateWindowsManifest(); err != nil {
		return nil, err
	}
	if err = createAppLocalizations(); err != nil {
		return nil, err
	}
	if err = updateWindowICON(); err != nil {
		return nil, err
	}
	return &CreateResult{
		Project: newProject,
		Dir:     absDir,
		EGPPath: filepath.Join(absDir, newProject.EGPName),
	}, nil
}

func resolveCreateDir(projectName, selectedDir string) (string, error) {
	absDir, err := filepath.Abs(selectedDir)
	if err != nil {
		return "", err
	}
	targetDir := absDir
	if !tool.Equal(filepath.Base(filepath.Clean(absDir)), projectName) {
		targetDir = filepath.Join(absDir, projectName)
	}
	if tool.IsExist(targetDir) {
		if !tool.IsDir(targetDir) {
			return "", fmt.Errorf("project path is not a directory: %s", targetDir)
		}
		return targetDir, nil
	}
	if err = os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", err
	}
	return targetDir, nil
}

func WriteEGPConfig(path string, project *bean.TProject) error {
	if project == nil {
		return errors.New("project configuration is empty")
	}
	if project.CheckLinuxWSGTK3() {
		project.BuildOption.UIGtk3 = true
	}
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, project.EGPName), data, 0644)
}

func checkCreateDir(dir string, onConflict func(Conflict) ConflictDecision) error {
	if !tool.IsExist(dir) {
		return fmt.Errorf("directory does not exist: %s", dir)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var existingEGP string
	isNotEmpty := false
	for _, entry := range entries {
		if entry.IsDir() {
			isNotEmpty = true
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.HasSuffix(name, consts.EGPExt) {
			existingEGP = entry.Name()
			break
		}
		isNotEmpty = true
	}
	if existingEGP != "" {
		conflict := Conflict{
			Kind:    ConflictExistingProject,
			Dir:     dir,
			File:    existingEGP,
			Message: fmt.Sprintf("Project config already exists in current directory: %s. Overwrite?", existingEGP),
		}
		if onConflict == nil || onConflict(conflict) != ConflictOverwrite {
			return errors.New("project creation cancelled")
		}
		if err = os.Remove(filepath.Join(dir, existingEGP)); err != nil {
			return err
		}
		return nil
	}
	if isNotEmpty {
		conflict := Conflict{
			Kind:    ConflictNotEmpty,
			Dir:     dir,
			Message: "Directory is not empty. Create?",
		}
		if onConflict == nil || onConflict(conflict) != ConflictOverwrite {
			return errors.New("project creation cancelled")
		}
	}
	return nil
}

func isSupportedGUIRenderFramework(framework bean.GUIRenderFramework) bool {
	return framework == bean.GUIRenderFramework_LCL ||
		framework == bean.GUIRenderFramework_WV ||
		framework == bean.GUIRenderFramework_CEF
}

func createProjectDir() error {
	if bean.GProject == nil || bean.GPath == "" {
		return errors.New("project is not initialized")
	}
	appRoot := bean.GPath
	appCodePath := filepath.Join(appRoot, consts.AppPackageName)
	resourcesPath := bean.ResourcePath()
	resourcesEmbedPath := bean.ResourceEmbedPath()
	resourcesMetadataPath := bean.ResourceMetadataPath()
	for _, path := range []string{appCodePath, resourcesPath, resourcesEmbedPath, resourcesMetadataPath} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	data := *bean.GProject
	localModule := tool.Buffer{}
	localModule.WriteString("require (\n")
	for modName, modVersion := range config.DesignerConfig.Dependencies {
		localModule.WriteString(modName, " ", modVersion, "\n")
	}
	localModule.WriteString(")", "\n")
	data.Data = localModule.String()

	type createFile struct {
		path string
		name string
		data string
	}
	files := []createFile{
		{appCodePath, consts.FormListFileName, buildTemplateData(projecttpl.AppCodeTemplate, &data)},
		{resourcesPath, "resources.go", buildTemplateData(projecttpl.ResourcesGoTemplate, &data)},
		{resourcesPath, "resources_windows.go", buildTemplateData(projecttpl.ResourcesWindowsGoTemplate, &data)},
		{resourcesMetadataPath, "metadata_windows.go", projecttpl.MetadataWindowsGoTemplate},
		{appRoot, "go.mod", buildTemplateData(projecttpl.GoModTemplate, &data)},
	}
	switch data.GUIRenderFramework {
	case bean.GUIRenderFramework_LCL:
		files = append(files, createFile{appRoot, "main.go", buildTemplateData(projecttpl.RunLCLCodeTemplate, &data)})
	case bean.GUIRenderFramework_WV:
		webPath := filepath.Join(appRoot, "web")
		if err := os.MkdirAll(webPath, os.ModePerm); err != nil {
			return err
		}
		files = append(files,
			createFile{appRoot, "main.go", buildTemplateData(projecttpl.RunWVCodeTemplate, &data)},
			createFile{webPath, "index.html", buildTemplateData(projecttpl.WebIndexHTMLTemplate, &data)},
		)
	case bean.GUIRenderFramework_CEF:
		webPath := filepath.Join(appRoot, "web")
		if err := os.MkdirAll(webPath, os.ModePerm); err != nil {
			return err
		}
		files = append(files,
			createFile{appRoot, "main.go", buildTemplateData(projecttpl.RunLCLCodeTemplate, &data)},
			createFile{webPath, "index.html", buildTemplateData(projecttpl.WebIndexHTMLTemplate, &data)},
		)
	}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(file.path, file.name), []byte(file.data), 0644); err != nil {
			return err
		}
	}
	return nil
}

func buildTemplateData(templateData string, data any) string {
	tmpl, err := template.New("project").Parse(templateData)
	if err != nil {
		return ""
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func saveOrUpdateWindowsManifest() error {
	iconData := bean.GProject.AppOption.Icon
	if iconData.Data == nil {
		iconData.Data = resources.Images("icons/window-icon_256x256.png")
		iconData.W = 256
		iconData.H = 256
	}

	icoSetBuf := tool.Buffer{}
	if err := winicon.GenerateIcon(bytes.NewBuffer(iconData.Data), &icoSetBuf, []int{256, 128, 64, 48, 32, 16}); err != nil {
		return err
	}
	rs := &winres.ResourceSet{}
	ico, err := winres.LoadICO(bytes.NewReader(icoSetBuf.Bytes()))
	if err != nil {
		return err
	}
	if err = rs.SetIcon(winres.RT_ICON, ico); err != nil {
		return err
	}
	rs.SetManifest(newManifest())

	v := version.Info{}
	v.ProductVersion = appVersionNum(bean.GProject.AppOption.Version)
	v.FileVersion = appVersionNum(bean.GProject.AppOption.Version)
	v.Flags.SpecialBuild = true
	v.Timestamp = time.Now()
	v.Set(2052, version.CompanyName, bean.GProject.AppOption.Id)
	v.Set(2052, version.ProductName, bean.GProject.AppOption.Title)
	v.Set(2052, version.LegalCopyright, bean.GProject.AppOption.Copyright)
	v.Set(2052, version.FileDescription, bean.GProject.AppOption.Desc)
	v.Set(2052, version.ProductVersion, bean.GProject.AppOption.Version)
	v.Set(2052, version.FileVersion, bean.GProject.AppOption.Version)
	v.Set(2052, version.Comments, bean.GProject.AppOption.Desc)
	rs.SetVersionInfo(v)

	resourcesPath := bean.ResourceMetadataPath()
	for _, arch := range []winres.Arch{winres.ArchAMD64, winres.ArchI386} {
		sysoOutBuf := tool.Buffer{}
		if err = rs.WriteObject(&sysoOutBuf, arch); err != nil {
			return err
		}
		sysoOutFile := fmt.Sprintf("%s-windows_%v.syso", bean.GProject.Name, arch)
		if err = os.WriteFile(filepath.Join(resourcesPath, sysoOutFile), sysoOutBuf.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func newManifest() winres.AppManifest {
	return winres.AppManifest{
		Identity: winres.AssemblyIdentity{
			Name:    bean.GProject.AppOption.Id,
			Version: appVersionNum(bean.GProject.AppOption.Version),
		},
		Description:                       bean.GProject.AppOption.Desc,
		UIAccess:                          bean.GProject.AppOption.Windows.Manifest.UIAccess,
		AutoElevate:                       bean.GProject.AppOption.Windows.Manifest.AutoElevate,
		DisableTheming:                    bean.GProject.AppOption.Windows.Manifest.DisableTheming,
		DisableWindowFiltering:            bean.GProject.AppOption.Windows.Manifest.DisableWindowFiltering,
		HighResolutionScrollingAware:      bean.GProject.AppOption.Windows.Manifest.HighResolutionScrollingAware,
		UltraHighResolutionScrollingAware: bean.GProject.AppOption.Windows.Manifest.UltraHighResolutionScrollingAware,
		LongPathAware:                     bean.GProject.AppOption.Windows.Manifest.LongPathAware,
		PrinterDriverIsolation:            bean.GProject.AppOption.Windows.Manifest.PrinterDriverIsolation,
		GDIScaling:                        bean.GProject.AppOption.Windows.Manifest.GDIScaling,
		SegmentHeap:                       bean.GProject.AppOption.Windows.Manifest.SegmentHeap,
		UseCommonControlsV6:               bean.GProject.AppOption.Windows.Manifest.UseCommonControlsV6,
		ExecutionLevel:                    winres.ExecutionLevel(bean.GProject.AppOption.Windows.Manifest.RunLevel),
		Compatibility:                     winres.SupportedOS(bean.GProject.AppOption.Windows.Manifest.CompatibilityOS),
		DPIAwareness:                      winres.DPIAwareness(bean.GProject.AppOption.Windows.Manifest.DPI),
	}
}

func appVersionNum(versionText string) [4]uint16 {
	versionNum := [4]uint16{0, 0, 0, 0}
	for i, v := range tool.Split(versionText, ".") {
		if i < len(versionNum) {
			vn, _ := strconv.ParseUint(v, 10, 16)
			versionNum[i] = uint16(vn)
		}
	}
	return versionNum
}

func createAppLocalizations() error {
	resourcesMetadataPath := bean.ResourceMetadataPath()
	for _, local := range bean.GProject.AppOption.MacOS.PList.CFBundleLocalizations {
		resourcesLocal := filepath.Join(resourcesMetadataPath, local+".lproj")
		if tool.IsExist(resourcesLocal) {
			continue
		}
		if err := os.MkdirAll(resourcesLocal, 0755); err != nil {
			return err
		}
		localizations := `/* localizations */
CFBundleDisplayName = "{{CFBundleDisplayName}}";
CFBundleName = "{{CFBundleName}}";
`
		localizations = strings.Replace(localizations, "{{CFBundleDisplayName}}", bean.GProject.AppOption.MacOS.PList.CFBundleDisplayName, 1)
		localizations = strings.Replace(localizations, "{{CFBundleName}}", bean.GProject.AppOption.MacOS.PList.CFBundleName, 1)
		if err := os.WriteFile(filepath.Join(resourcesLocal, "InfoPlist.strings"), []byte(localizations), 0644); err != nil {
			return err
		}
	}
	return nil
}

func updateWindowICON() error {
	icon := bean.GProject.AppOption.Icon
	if icon.Data == nil {
		icon = bean.TAppIcon{
			Data: resources.Images("icons/window-icon_256x256.png"),
			W:    256,
			H:    256,
		}
	}
	if icon.Data == nil || icon.W <= 0 || icon.H <= 0 {
		return errors.New("icon data cannot be empty or the size is incorrect")
	}
	iconData := icon.Data
	if tool.IsDarwin {
		if icon.W > 1024 || icon.H > 1024 {
			iconData = tool.Scale(iconData, 1024, 1024)
		}
	} else if icon.W > 256 || icon.H > 256 {
		iconData = tool.Scale(iconData, 256, 256)
	}

	pngBuf := new(bytes.Buffer)
	pngBuf.Write(iconData)
	pngImage, err := png.Decode(pngBuf)
	if err != nil {
		return err
	}
	icoBuf := new(bytes.Buffer)
	if err = tool.Encode(icoBuf, pngImage); err != nil {
		return err
	}

	embedPath := bean.ResourceEmbedPath()
	iconIcoFilePath := filepath.Join(embedPath, "icon.ico")
	if err = os.WriteFile(iconIcoFilePath, icoBuf.Bytes(), 0644); err != nil {
		return err
	}
	iconPngFilePath := filepath.Join(embedPath, "icon.png")
	if err = os.WriteFile(iconPngFilePath, iconData, 0644); err != nil {
		return err
	}

	iconPngFile, err := os.Open(iconPngFilePath)
	if err != nil {
		return err
	}
	defer iconPngFile.Close()
	iconPngSrcImg, _, err := image.Decode(iconPngFile)
	if err != nil {
		return err
	}
	iconPngIcnsDest, err := os.Create(filepath.Join(embedPath, "icon.icns"))
	if err != nil {
		return err
	}
	defer iconPngIcnsDest.Close()
	return icns.Encode(iconPngIcnsDest, iconPngSrcImg)
}
