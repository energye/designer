; 工具宏定义脚本, 用于 energy designer 制作应用程序安装包
; 脚本功能, 定义宏, 打包应用二进制文件, 自动安装运行时环境

!include "x64.nsh"
!include "WinVer.nsh"
!include "FileFunc.nsh"

!define INFO_EXECUTE_BINARY "${.BuildName}"
!define INFO_InstallFileName "{{.InstallFileName}}"
!define INFO_CompanyName "{{.CompanyName}}"
!define INFO_ProductName "{{.ProductName}}"
!define INFO_ShortCutName "{{.ShortCutName}}"
!define INFO_FileVersion "{{.FileVersion}}"
!define INFO_ProductVersion "{{.ProductVersion}}"
!define INFO_FileDescription "{{.FileDescription}}"
!define INFO_Copyright "{{.Copyright}}"
!define INFO_UNINST_KEY_NAME "${INFO_CompanyName}${INFO_ProductName}"
!define INFO_Icon "{{.NSISIcon}}"
!define INFO_UnIcon "{{.NSISUnIcon}}"
!define INFO_LANGUAGE "{{.NSISLanguage}}"

!define UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${INFO_UNINST_KEY_NAME}"

{{if .NSISLicense}}
    !define NSIS_PAGE_LICENSE "{{.NSISLicense}}" ; license.txt path
{{end}}

{{if .NSISRequestExecutionLevel}}
    !define REQUEST_EXECUTION_LEVEL "{{.NSISRequestExecutionLevel}}"
    RequestExecutionLevel "${REQUEST_EXECUTION_LEVEL}" ; admin or ""
{{end}}


!macro energy.files
    File "/oname=${INFO_EXECUTE_BINARY}" "{{.ProjectPath}}\{{.ExeName}}.exe"

    ; File "file xxx"
    ; File -r "file xxx"

    {{range $i,$path := .NSISInclude }}
        File /r "{{$path}}"{{end}}
    !macroend

!macroend

!macro energy.writeUninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"

    SetRegView 64
    WriteRegStr HKLM "${UNINST_KEY}" "Publisher" "${INFO_CompanyName}"
    WriteRegStr HKLM "${UNINST_KEY}" "DisplayName" "${INFO_ProductName}"
    WriteRegStr HKLM "${UNINST_KEY}" "DisplayVersion" "${INFO_ProductVersion}"
    WriteRegStr HKLM "${UNINST_KEY}" "DisplayIcon" "$INSTDIR\${INFO_EXECUTE_BINARY}"
    WriteRegStr HKLM "${UNINST_KEY}" "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
    WriteRegStr HKLM "${UNINST_KEY}" "QuietUninstallString" "$\"$INSTDIR\uninstall.exe$\" /S"

    ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
    IntFmt $0 "0x%08X" $0
    WriteRegDWORD HKLM "${UNINST_KEY}" "EstimatedSize" "$0"
!macroend

!macro energy.deleteUninstaller
    Delete "$INSTDIR\uninstall.exe"

    SetRegView 64
    DeleteRegKey HKLM "${UNINST_KEY}"
!macroend

!macro energy.setShellContext
    ${If} ${REQUEST_EXECUTION_LEVEL} == "admin"
        SetShellVarContext all
    ${else}
        SetShellVarContext current
    ${EndIf}
!macroend
