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

!define MUI_ICON "${INFO_Icon}" ;"..\icon.ico"
!define MUI_UNICON "${INFO_UnIcon}" ;"..\icon.ico"

!define MUI_FINISHPAGE_NOAUTOCLOSE # Wait on the INSTFILES page so the user can take a look into the details of the installation steps
!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.


!insertmacro MUI_PAGE_WELCOME # Welcome to the installer page.

!ifdef NSIS_PAGE_LICENSE
    !insertmacro MUI_PAGE_LICENSE "${NSIS_PAGE_LICENSE}" # Add a LICENSE page to the installer
!endif

!insertmacro MUI_PAGE_DIRECTORY # In which folder install page.
!insertmacro MUI_PAGE_INSTFILES # Installing page.
!insertmacro MUI_PAGE_FINISH # Finished installation page.

!insertmacro MUI_UNPAGE_INSTFILES # Uinstalling page

!insertmacro MUI_LANGUAGE "${INFO_LANGUAGE}" # Set the Language of the installer


Name "${INFO_ProductName}"
OutFile ".\${INFO_InstallFileName}"
InstallDir "$PROGRAMFILES64\${INFO_CompanyName}\${INFO_ProductName}"
ShowInstDetails show # This will always show the installation details.

Function .onInit
FunctionEnd

Section
    !insertmacro energy.setShellContext

    SetOutPath $INSTDIR
    
    !insertmacro energy.files

    CreateShortcut "$SMPROGRAMS\${INFO_ShortCutName}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_ShortCutName}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro energy.writeUninstaller
SectionEnd

Section "uninstall" 
    !insertmacro energy.setShellContext

    RMDir /r "$AppData\${PRODUCT_EXECUTABLE}"

    RMDir /r $INSTDIR

    Delete "$SMPROGRAMS\${INFO_ShortCutName}.lnk"
    Delete "$DESKTOP\${INFO_ShortCutName}.lnk"

    !insertmacro energy.deleteUninstaller
SectionEnd
