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

!define MUI_ICON "${INFO_Icon}"
!define MUI_UNICON "${INFO_UnIcon}"

!define MUI_FINISHPAGE_NOAUTOCLOSE
!define MUI_ABORTWARNING


!insertmacro MUI_PAGE_WELCOME

!ifdef NSIS_PAGE_LICENSE
    !insertmacro MUI_PAGE_LICENSE "${NSIS_PAGE_LICENSE}"
!endif

!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

!insertmacro MUI_UNPAGE_INSTFILES

!insertmacro MUI_LANGUAGE "${INFO_LANGUAGE}"


Name "${INFO_ProductName}"
OutFile ".\${INFO_InstallFileName}"
InstallDir "$PROGRAMFILES64\${INFO_CompanyName}\${INFO_ProductName}"
ShowInstDetails show # This will always show the installation details.

Function .onInit
FunctionEnd

Function un.onInit
    MessageBox MB_YESNO|MB_ICONQUESTION "您确定要完全卸载吗？所有程序文件都会被删除" IDYES noabort
    Abort
noabort:
FunctionEnd

Section
    !insertmacro energy.setShellContext

    SetOutPath $INSTDIR
    
    !insertmacro energy.files

    CreateShortcut "$SMPROGRAMS\${INFO_ShortCutName}.lnk" "$INSTDIR\${INFO_EXECUTE_BINARY}"
    CreateShortCut "$DESKTOP\${INFO_ShortCutName}.lnk" "$INSTDIR\${INFO_EXECUTE_BINARY}"

    !insertmacro energy.writeUninstaller
SectionEnd

Section "uninstall" 
    !insertmacro energy.setShellContext

    RMDir /r "$AppData\${INFO_EXECUTE_BINARY}"

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_ShortCutName}.lnk"
    Delete "$DESKTOP\${INFO_ShortCutName}.lnk"
    Delete "$INSTDIR\uninstall.exe"
    SetRegView 64
    DeleteRegKey HKLM "${UNINST_KEY}"
SectionEnd
