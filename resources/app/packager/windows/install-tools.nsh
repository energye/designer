; 工具宏定义脚本, 用于 energy designer 制作应用程序安装包
; 脚本功能, 定义宏, 打包应用二进制文件, 自动安装运行时环境

!include "x64.nsh"
!include "WinVer.nsh"
!include "FileFunc.nsh"

; 脚本宏定义

!define INFO_EXECUTE_BINARY "{{.BuildName}}"
!define INFO_EXECUTE_BINARY_PATH "{{.BuildFileNamePath}}"
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

{{if .DefaultInstall}}
!define INFO_DefaultInstall "{{.DefaultInstall}}"
{{else}}
!define INFO_DefaultInstall "$PROGRAMFILES64\${INFO_CompanyName}\${INFO_ProductName}"
{{end}}

!define INFO_RuntimeWebView2Setup "{{.RuntimeWebView2Setup}}"

; 唯一 key
!define UNINST_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\${INFO_UNINST_KEY_NAME}"

; 授权文件信息
{{if .NSISLicense}}
    !define NSIS_PAGE_LICENSE "{{.NSISLicense}}"
{{end}}

; 执行等级
{{if .NSISRequestExecutionLevel}}
    !define REQUEST_EXECUTION_LEVEL "{{.NSISRequestExecutionLevel}}"
    RequestExecutionLevel "${REQUEST_EXECUTION_LEVEL}"
{{end}}

; 打包进的文件
!macro energy.files
    File "/oname=${INFO_EXECUTE_BINARY}" "${INFO_EXECUTE_BINARY_PATH}"

    ; File "file xxx"
    ; File -r "file xxx"

    File "{{.RuntimeLibEnergy}}" ; libenergy-xxx.dll
    {{if .RuntimeWebView2Loader}}
    File "{{.RuntimeWebView2Loader}}" ; WebView2Loader-xxx.dll
    {{end}}

    {{range $i,$path := .NSISInclude }}
        File /r "{{$path}}"{{end}}
!macroend

; 写入卸载信息
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

; 上下文切换
!macro energy.setShellContext
    ${If} ${REQUEST_EXECUTION_LEVEL} == "admin"
        SetShellVarContext all
    ${else}
        SetShellVarContext current
    ${EndIf}
!macroend

; Webview2
!macro energy.webview2Install
    SetRegView 64
    !define WEBVIEW2_CLSID "{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"

    ; 检查 WebView2 是否已经安装

    ; 检查 32位 系统级（WOW6432NODE）
    ReadRegStr $0 HKLM "SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\${WEBVIEW2_CLSID}" "pv"
    ${If} $0 != ""
        Goto webview2_done ; 已安装
    ${EndIf}

    ; 检查 当前用户级（无管理员权限时）
    ${If} ${REQUEST_EXECUTION_LEVEL} == "user"
        ReadRegStr $0 HKCU "Software\Microsoft\EdgeUpdate\Clients\${WEBVIEW2_CLSID}" "pv"
        ${If} $0 != ""
            Goto webview2_done ; 已安装
        ${EndIf}
    ${EndIf}

    ; 开始安装
    SetDetailsPrint both
    DetailPrint "Installing Microsoft Edge WebView2 Runtime..."
    SetDetailsPrint listonly

    ; 释放安装程序
    InitPluginsDir
    StrCpy $1 "$pluginsdir\webview2boot"
    CreateDirectory "$1"
    SetOutPath "$1"
    File "MicrosoftEdgeWebview2Setup.exe"

    ; 静默安装
    ExecWait '"$1\MicrosoftEdgeWebview2Setup.exe" /silent /install'

    ; 清理临时文件
    Delete "$1\MicrosoftEdgeWebview2Setup.exe"
    RMDir "$1"

    ; 安装后必须再次校验
    SetRegView 64
    ReadRegStr $2 HKLM "SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\${WEBVIEW2_CLSID}" "pv"

    ${If} $2 == ""
        ${If} ${REQUEST_EXECUTION_LEVEL} == "user"
            ReadRegStr $2 HKCU "Software\Microsoft\EdgeUpdate\Clients\${WEBVIEW2_CLSID}" "pv"
        ${EndIf}
    ${EndIf}

    ${If} $2 == ""
        SetDetailsPrint both
        MessageBox MB_ICONSTOP "WebView2 Runtime installation failed!"
        Abort ; 安装失败，直接终止安装程序
    ${EndIf}

    ; 已安装或完成安装
    webview2_done:
    SetDetailsPrint both
    SetRegView lastused
!macroend

; 关联文件
!macro APP_ASSOCIATE EXT FILECLASS DESCRIPTION ICON COMMANDTEXT COMMAND
  ReadRegStr $R0 SHELL_CONTEXT "Software\Classes\.${EXT}" ""
  WriteRegStr SHELL_CONTEXT "Software\Classes\.${EXT}" "${FILECLASS}_backup" "$R0"
  WriteRegStr SHELL_CONTEXT "Software\Classes\.${EXT}" "" "${FILECLASS}"
  WriteRegStr SHELL_CONTEXT "Software\Classes\${FILECLASS}" "" `${DESCRIPTION}`
  WriteRegStr SHELL_CONTEXT "Software\Classes\${FILECLASS}\DefaultIcon" "" `${ICON}`
  WriteRegStr SHELL_CONTEXT "Software\Classes\${FILECLASS}\shell" "" "open"
  WriteRegStr SHELL_CONTEXT "Software\Classes\${FILECLASS}\shell\open" "" `${COMMANDTEXT}`
  WriteRegStr SHELL_CONTEXT "Software\Classes\${FILECLASS}\shell\open\command" "" `${COMMAND}`
!macroend

!macro APP_UNASSOCIATE EXT FILECLASS
  ReadRegStr $R0 SHELL_CONTEXT "Software\Classes\.${EXT}" `${FILECLASS}_backup`
  WriteRegStr SHELL_CONTEXT "Software\Classes\.${EXT}" "" "$R0"
  DeleteRegKey SHELL_CONTEXT `Software\Classes\${FILECLASS}`
!macroend

!macro energy.associateFiles
    {{range .AssociateFiles}}
    !insertmacro APP_ASSOCIATE "{{.Ext}}" "{{.FileClass}}" "{{.Description}}" "$INSTDIR\{{.Icon}}" "{{.CommandText}}" "$INSTDIR\${INFO_EXECUTE_BINARY} $\"%1$\""
    File "{{.SrcIcon}}"
    {{end}}
!macroend

!macro energy.unAssociateFiles
    {{range .AssociateFiles}}
    !insertmacro APP_UNASSOCIATE "{{.Ext}}" "{{.FileClass}}"
    Delete "$INSTDIR\{{.Icon}}"
    {{end}}
!macroend


; 自定义打开软件协议
!macro CUSTOM_PROTOCOL_ASSOCIATE PROTOCOL DESCRIPTION ICON COMMAND
  DeleteRegKey SHELL_CONTEXT "Software\Classes\${PROTOCOL}"
  WriteRegStr SHELL_CONTEXT "Software\Classes\${PROTOCOL}" "" "${DESCRIPTION}"
  WriteRegStr SHELL_CONTEXT "Software\Classes\${PROTOCOL}" "URL Protocol" ""
  WriteRegStr SHELL_CONTEXT "Software\Classes\${PROTOCOL}\DefaultIcon" "" "${ICON}"
  WriteRegStr SHELL_CONTEXT "Software\Classes\${PROTOCOL}\shell" "" ""
  WriteRegStr SHELL_CONTEXT "Software\Classes\${PROTOCOL}\shell\open" "" ""
  WriteRegStr SHELL_CONTEXT "Software\Classes\${PROTOCOL}\shell\open\command" "" "${COMMAND}"
!macroend

!macro CUSTOM_PROTOCOL_UNASSOCIATE PROTOCOL
  DeleteRegKey SHELL_CONTEXT "Software\Classes\${PROTOCOL}"
!macroend

!macro energy.customAssociateProtocols
    {{range .AssociateProtocols}}
      !insertmacro CUSTOM_PROTOCOL_ASSOCIATE "{{.Scheme}}" "{{.Description}}" "$INSTDIR\${INFO_EXECUTE_BINARY},0" "$INSTDIR\${INFO_EXECUTE_BINARY} $\"%1$\""
    {{end}}
!macroend

!macro energy.unCustomAssociateProtocols
    {{range .AssociateProtocols}}
      !insertmacro CUSTOM_PROTOCOL_UNASSOCIATE "{{.Scheme}}"
    {{end}}
!macroend