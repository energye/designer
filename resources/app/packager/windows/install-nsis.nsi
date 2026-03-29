Unicode true

!include "install-tools.nsh"

VIProductVersion "${INFO_ProductVersion}.0"
VIFileVersion    "${INFO_FileVersion}.0"

VIAddVersionKey "ProductName"     "${INFO_ProductName}"
VIAddVersionKey "CompanyName"     "${INFO_CompanyName}"
VIAddVersionKey "FileDescription" "${INFO_FileDescription}"
VIAddVersionKey "ProductVersion"  "${INFO_ProductVersion}"
VIAddVersionKey "FileVersion"     "${INFO_FileVersion}"
VIAddVersionKey "LegalCopyright"  "${INFO_Copyright}"

!include "MUI2.nsh"

; 宏定义

!define MUI_ICON "${INFO_Icon}"
!define MUI_UNICON "${INFO_UnIcon}"

!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN
!define MUI_FINISHPAGE_RUN_TEXT "Launch Now ${INFO_ShortCutName}"
!define MUI_FINISHPAGE_RUN_FUNCTION LaunchApp

; 插入界面

!insertmacro MUI_PAGE_WELCOME

!ifdef NSIS_PAGE_LICENSE
    !define MUI_LICENSEPAGE_CHECKBOX
    !insertmacro MUI_PAGE_LICENSE "${NSIS_PAGE_LICENSE}"
!endif

!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "${INFO_LANGUAGE}"

Name "${INFO_ShortCutName} ${INFO_ProductVersion}"
OutFile ".\${INFO_InstallFileName}"
InstallDir "${INFO_DefaultInstall}"
ShowInstDetails show
ShowUnInstDetails show

Function .onInit
  !insertmacro energy.setShellContext
  SetRegView 64
  ReadRegStr $0 HKLM "${UNINST_KEY}" "UninstallString"
  StrCmp $0 "" done

  MessageBox MB_OKCANCEL|MB_ICONEXCLAMATION "Already installed. Reinstall now?" IDOK done
  Abort

done:
FunctionEnd

; 卸载提示
Function un.onInit
    MessageBox MB_YESNO|MB_ICONQUESTION "您确定要完全卸载吗？$\n$\n所有程序文件都会被删除" IDYES noabort
    Abort
noabort:
FunctionEnd

; 安装完提示立即运行
Function LaunchApp
    Exec "$INSTDIR\${INFO_EXECUTE_BINARY}"
FunctionEnd

; 创建桌面快捷
Section "Desktop" SEC_DESKTOP
    !insertmacro energy.setShellContext
    CreateShortCut "$DESKTOP\${INFO_ShortCutName}.lnk" "$INSTDIR\${INFO_EXECUTE_BINARY}"
SectionEnd

; 创建开始菜单快捷
Section "StartMenu" SEC_STARTMENU
    !insertmacro energy.setShellContext
    CreateDirectory "$SMPROGRAMS\${INFO_ShortCutName}"
    CreateShortcut "$SMPROGRAMS\${INFO_ShortCutName}\${INFO_ShortCutName}.lnk" "$INSTDIR\${INFO_EXECUTE_BINARY}" "" \
      "$INSTDIR\${INFO_EXECUTE_BINARY}" 0
    CreateShortcut "$SMPROGRAMS\${INFO_ShortCutName}\Uninstall.lnk" "$INSTDIR\uninstall.exe" "" \
      "$INSTDIR\uninstall.exe" 0
SectionEnd

Section
    !insertmacro energy.setShellContext

    ; 判断是否需要安装 webview2 runtime
    !ifdef INFO_RuntimeWebView2Setup
    !insertmacro energy.webview2Install
    !endif

    ; 安装目录
    SetOutPath $INSTDIR
    ; 打包进的文件
    !insertmacro energy.files
    ; 关联文件
    !insertmacro energy.associateFiles
    ; 关联协议
    !insertmacro energy.customAssociateProtocols

    ;CreateShortcut "$SMPROGRAMS\${INFO_ShortCutName}.lnk" "$INSTDIR\${INFO_EXECUTE_BINARY}"
    ;CreateShortCut "$DESKTOP\${INFO_ShortCutName}.lnk" "$INSTDIR\${INFO_EXECUTE_BINARY}"

    !insertmacro energy.writeUninstaller
SectionEnd

Section "uninstall" 
    !insertmacro energy.setShellContext
    !insertmacro energy.unAssociateFiles
    !insertmacro energy.unCustomAssociateProtocols

    Delete "$SMPROGRAMS\${INFO_ShortCutName}\${INFO_ShortCutName}.lnk"
    Delete "$SMPROGRAMS\${INFO_ShortCutName}\UnInstall.lnk"
    RMDir /r "$SMPROGRAMS\${INFO_ShortCutName}"

    Delete "$DESKTOP\${INFO_ShortCutName}.lnk"

    RMDir /r "$AppData\${INFO_ShortCutName}"

    Delete "$INSTDIR\uninstall.exe"
    RMDir /r $INSTDIR

    SetRegView 64
    DeleteRegKey HKLM "${UNINST_KEY}"
SectionEnd
