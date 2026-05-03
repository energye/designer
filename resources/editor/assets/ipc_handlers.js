// Inbound IPC event listeners from Go

ipc.on('open-file', function (data, readOnly) { openFile(data, readOnly); });
ipc.on('close-file', function (data) { closeFile(data); });
ipc.on('save-current-file', function () { autoSave(); });

ipc.on('file-changed-externally', function (filePath) {
    var info = files.get(filePath);
    if (!info || info.isDirty) return;
    reloadFileFromDisk(filePath, '文件已重新加载: ' + getFileName(filePath));
});

ipc.on('file-conflict-detected', function (filePath) {
    reloadFileFromDisk(filePath);
});
