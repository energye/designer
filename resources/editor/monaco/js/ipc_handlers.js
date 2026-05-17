// Inbound IPC event listeners from Go

ipc.on('open-file', function (data, readOnly) {
    openFile(data, readOnly);
});

ipc.on('close-file', function (data) {
    closeFile(data);
});

ipc.on('save-current-file', function () {
    manualSave();
});

ipc.on('goto-position', function (filePath, line, character) {
    if (!editor) return;
    // If the file is already open and active, position immediately
    if (currentFilePath === filePath && files.has(filePath)) {
        editor.revealLineInCenter(line + 1);
        editor.setPosition({lineNumber: line + 1, column: character + 1});
        editor.focus();
    } else {
        // Store pending position - will be applied when openFile completes
        pendingGotoPosition = {filePath: filePath, line: line, character: character};
    }
});

ipc.on('file-changed-externally', function (filePath) {
    let info = files.get(filePath);
    if (!info || info.isDirty) return;
    reloadFileFromDisk(filePath, '文件已重新加载: ' + getFileName(filePath));
});

ipc.on('file-conflict-detected', function (filePath) {
    reloadFileFromDisk(filePath);
});

