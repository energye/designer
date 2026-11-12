// File open/close/save/reload logic

let pendingGotoPosition = null; // {line, character} - applied after file open completes
let conflictQueue = [];        // 待用户处理的文件冲突，一次显示一个

function switchToFile(filePath) {
    let info = files.get(filePath);
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
        let pos = pendingGotoPosition;
        pendingGotoPosition = null;
        editor.revealLineInCenter(pos.line + 1);
        editor.setPosition({lineNumber: pos.line + 1, column: pos.character + 1});
        editor.focus();
    }
}

function updateFilePathBar(filePath) {
    let bar = document.getElementById('file-path-bar');
    if (!bar) return;
    let info = files.get(filePath);
    let displayPath = filePath.replace(/\\/g, '/');
    bar.style.color = '';
    if (info && info.readOnly) {
        displayPath = displayPath + '  [Read Only]';
        bar.style.color = '#b8860b';
    }
    bar.textContent = displayPath;
}

function openFile(filePath, readOnly) {
    if (!editor) return;

    ipc.emit('open-file-request', [filePath], function (response) {
        if (!response) return;
        let data;
        if (typeof response === 'string') {
            data = JSON.parse(response);
        } else if (typeof response === 'object') {
            data = response;
        }
        let content = data.content;
        let language = data.language;
        let modTime = data.modTime;
        let serverReadOnly = !!data.readOnly;
        let effectiveReadOnly = readOnly || serverReadOnly;

        if (files.has(filePath)) {
            let info = files.get(filePath);
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
            // let uri = modelUri(filePath);
            let uri = filePath;
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
    let key = resolveFilePath(filePath);
    if (key === null) return;
    let info = files.get(key);

    ipc.emit('gopls-didClose', [key]);
    info.model.dispose();
    files.delete(key);
    ipc.emit('unregister-opened-file', [key]);

    if (currentFilePath === key) {
        currentFilePath = '';
        if (files.size > 0) {
            let firstKey = files.keys().next().value;
            switchToFile(firstKey);
        } else {
            editor.setModel(monacoRef.editor.createModel('', 'plaintext'));
        }
    }
    renderFileConflict();
}

function autoSave() {
    if (!currentFilePath || !editor) return;
    let info = files.get(currentFilePath);
    if (!info || info.readOnly) return;
    let model = editor.getModel();
    if (!model) return;

    requestOrganizeImports(currentFilePath, model, true);
}

// Reload file from disk and re-sync with gopls
function reloadFileFromDisk(filePath, notifyMessage) {
    if (!files.has(filePath)) return;

    ipc.emit('reload-file-request', [filePath], function (response) {
        if (!response) return;
        // WebView2 的 IPC 回调会把 Go 侧结构体直接解析成对象；只有字符串
        // 形式才需要 JSON.parse。旧实现无条件 parse 会抛异常，导致冲突提示
        // 里的「重新加载」静默失效。
        let data = response;
        if (typeof response === 'string') {
            try {
                data = JSON.parse(response);
            } catch (e) {
                return;
            }
        }
        if (!data || typeof data.content !== 'string') return;
        let info = files.get(filePath);
        applyingEdit = true;
        try {
            info.model.setValue(data.content);
        } finally {
            applyingEdit = false;
        }
        info.modTime = data.modTime;
        info.readOnly = !!data.readOnly;
        // version 不再重置为 1：Go 侧派发器按文档单调分配版本，重置会让
        // 后续 didChange 的 version 回退，与 gopls 已打开的文档状态冲突。
        info.version++;
        info.isDirty = false;
        if (currentFilePath === filePath) {
            editor.setModel(info.model);
            editor.updateOptions({readOnly: info.readOnly});
        }
        updateFilePathBar(filePath);
        ipc.emit('gopls-didClose', [filePath]);
        ipc.emit('gopls-didOpen', [{
            file: filePath, languageId: data.language, content: data.content, version: info.version
        }]);
        ipc.emit('set-file-dirty', [{file: filePath, isDirty: false}]);
        ipc.emit('register-opened-file', [{file: filePath, modTime: data.modTime}]);
        clearFileConflict(filePath);
        if (notifyMessage) showNotification(notifyMessage);
    });
}

// clearFileConflict 文件已保存或已重载时撤掉该文件的冲突提示。
function clearFileConflict(filePath) {
    let index = conflictQueue.indexOf(filePath);
    if (index === -1) return;
    conflictQueue.splice(index, 1);
    renderFileConflict();
}

// 外部变更冲突：Go 侧每秒检测，dirty 时通知一次，这里排队让用户决定。
// 「保留我的修改」不动 model 也不清 dirty；只有「重新加载」才读磁盘覆盖。
function showFileConflict(filePath) {
    if (!files.has(filePath)) return;
    if (conflictQueue.indexOf(filePath) === -1) conflictQueue.push(filePath);
    renderFileConflict();
}

function renderFileConflict() {
    let bar = document.getElementById('conflict-bar');
    let message = document.getElementById('conflict-message');
    if (!bar || !message) return;
    while (conflictQueue.length > 0 && !files.has(conflictQueue[0])) conflictQueue.shift();
    if (conflictQueue.length === 0) {
        bar.classList.remove('show');
        return;
    }
    let path = conflictQueue[0];
    message.textContent = '文件在外部被修改，但当前有未保存内容: '
        + getFileName(path) + ' — 请选择处理方式';
    bar.classList.add('show');
}

function resolveFileConflict(action) {
    if (conflictQueue.length === 0) return;
    let filePath = conflictQueue.shift();
    if (action === 'reload' && files.has(filePath)) {
        reloadFileFromDisk(filePath, '已重新加载: ' + getFileName(filePath));
    } else {
        showNotification('已保留本地修改: ' + getFileName(filePath));
    }
    renderFileConflict();
}
