// File open/close/save/reload logic

var pendingGotoPosition = null; // {line, character} - applied after file open completes

function switchToFile(filePath) {
    var info = files.get(filePath);
    if (currentFilePath && files.has(currentFilePath)) {
        files.get(currentFilePath).viewState = editor.saveViewState();
    }
    editor.setModel(info.model);
    if (info.viewState) editor.restoreViewState(info.viewState);
    editor.updateOptions({readOnly: !!info.readOnly});
    currentFilePath = filePath;
    updateFilePathBar(filePath);

    // Apply pending goto position if this is the target file
    if (pendingGotoPosition && pendingGotoPosition.filePath === filePath) {
        var pos = pendingGotoPosition;
        pendingGotoPosition = null;
        editor.revealLineInCenter(pos.line + 1);
        editor.setPosition({lineNumber: pos.line + 1, column: pos.character + 1});
        editor.focus();
    }
}

function updateFilePathBar(filePath) {
    var bar = document.getElementById('file-path-bar');
    if (!bar) return;
    var info = files.get(filePath);
    var displayPath = filePath.replace(/\\/g, '/');
    bar.style.color = '';
    if (info && info.readOnly) {
        displayPath = displayPath + '  [Read Only]';
        bar.style.color = '#b8860b';
    }
    if (goplsStatus === 'unavailable' && filePath.replace(/\\/g, '/').endsWith('.go')) {
        displayPath = '⚠ gopls missing  |  ' + displayPath;
        bar.style.color = '#b8860b';
    }
    bar.textContent = displayPath;
}

function openFile(filePath, readOnly) {
    if (!editor) return;

    ipc.emit('open-file-request', [filePath], function (response) {
        if (!response) return;

        var data = JSON.parse(response);
        var content = data.content;
        var language = data.language;
        var modTime = data.modTime;
        var serverReadOnly = !!data.readOnly;
        var effectiveReadOnly = readOnly || serverReadOnly;

        if (files.has(filePath)) {
            var info = files.get(filePath);
            info.readOnly = effectiveReadOnly;
            if (info.modTime !== modTime && !info.isDirty) {
                info.model.setValue(content);
                info.modTime = modTime;
                info.version++;
                ipc.emit('gopls-didChange', [{
                    file: filePath, content: content, version: info.version
                }]);
            }
        } else {
            var uri = modelUri(filePath);
            files.set(filePath, {
                model: monacoRef.editor.createModel(content, language, uri),
                modTime: modTime, isDirty: false, viewState: null, readOnly: effectiveReadOnly, version: 1
            });
            ipc.emit('register-opened-file', [{file: filePath, modTime: modTime}]);
            ipc.emit('gopls-didOpen', [{
                file: filePath, languageId: language, content: content, version: 1
            }]);
        }

        switchToFile(filePath);
    });
}

function closeFile(filePath) {
    if (!files.has(filePath)) return;
    var info = files.get(filePath);

    ipc.emit('gopls-didClose', [filePath]);
    info.model.dispose();
    files.delete(filePath);
    ipc.emit('unregister-opened-file', [filePath]);

    if (currentFilePath === filePath) {
        currentFilePath = '';
        if (files.size > 0) {
            var firstKey = files.keys().next().value;
            switchToFile(firstKey);
        } else {
            editor.setModel(monacoRef.editor.createModel('', 'plaintext'));
        }
    }
}

function autoSave() {
    if (!currentFilePath || !editor) return;
    var model = editor.getModel();
    if (!model) return;

    requestOrganizeImports(currentFilePath, model, true);
}

// Reload file from disk and re-sync with gopls
function reloadFileFromDisk(filePath, notifyMessage) {
    if (!files.has(filePath)) return;

    ipc.emit('reload-file-request', [filePath], function (response) {
        if (!response) return;
        var data = JSON.parse(response);
        var info = files.get(filePath);
        applyingEdit = true;
        info.model.setValue(data.content);
        applyingEdit = false;
        info.modTime = data.modTime;
        info.version = 1;
        info.isDirty = false;
        if (currentFilePath === filePath) editor.setModel(info.model);
        ipc.emit('gopls-didClose', [filePath]);
        ipc.emit('gopls-didOpen', [{
            file: filePath, languageId: data.language, content: data.content, version: 1
        }]);
        ipc.emit('set-file-dirty', [{file: filePath, isDirty: false}]);
        ipc.emit('register-opened-file', [{file: filePath, modTime: data.modTime}]);
        if (notifyMessage) showNotification(notifyMessage);
    });
}
