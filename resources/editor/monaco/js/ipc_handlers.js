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
    let key = resolveFilePath(filePath);
    // If the file is already open and active, position immediately
    if (key !== null && currentFilePath === key) {
        editor.revealLineInCenter(line + 1);
        editor.setPosition({lineNumber: line + 1, column: character + 1});
        editor.focus();
    } else {
        // Store pending position - will be applied when openFile completes
        pendingGotoPosition = {filePath: key === null ? filePath : key, line: line, character: character};
    }
});

ipc.on('file-changed-externally', function (filePath) {
    let key = resolveFilePath(filePath);
    if (key === null) return;
    let info = files.get(key);
    // Go 侧与前端的脏状态偶尔不同步（例如生成前刚落盘），此时不静默跳过。
    if (info.isDirty) {
        showFileConflict(key);
        return;
    }
    reloadFileFromDisk(key, '文件已重新加载: ' + getFileName(key));
});

// 外部变更且本地有未保存内容：绝不直接重载，交给用户选择，避免丢失编辑内容。
ipc.on('file-conflict-detected', function (filePath) {
    let key = resolveFilePath(filePath);
    if (key !== null) showFileConflict(key);
});

// 设计器要写盘但编辑器里有未保存内容：按 Monaco 当前 model 保存所有脏文件。
// 保存内容由前端给出，Go 侧不持有文本副本；每个 save-file 完成后 Go 侧
// 会确认脏状态清空并重放被暂缓的代码生成。
ipc.on('save-all-request', function (requestID) {
    let targets = [];
    files.forEach(function (info, path) {
        if (info.isDirty && !info.readOnly) targets.push(path);
    });
    targets.forEach(function (path) {
        let info = files.get(path);
        ipc.emit('save-file', [{
            file: path, content: info.model.getValue(), readOnly: false
        }], function (result) {
            if (result === 'ok' && files.has(path)) {
                files.get(path).isDirty = false;
                clearFileConflict(path);
            } else {
                showNotification('自动保存失败: ' + getFileName(path));
            }
        });
    });
});

// gopls 不可用 / 正在启动 / 首次索引时在代码编辑页给出明确提示。
let plsStatusOK = false;
let plsStatusReason = '';

function renderGoplsStatus(ok, reason) {
    plsStatusOK = (ok === true);
    plsStatusReason = reason || '';
    let bar = document.getElementById('gopls-status-bar');
    if (!bar) return;
    if (plsIndexing && plsStatusOK) {
        bar.textContent = 'Go 语言服务正在建立索引，首次补全 / 诊断 / 跳转可能较慢，请稍候…';
        bar.classList.add('show');
        return;
    }
    if (plsStatusOK) {
        bar.classList.remove('show');
        bar.textContent = '';
        return;
    }
    if (plsStatusReason.indexOf('正在启动') !== -1) {
        bar.textContent = 'Go 语言服务正在启动，请稍候…';
        bar.classList.add('show');
        return;
    }
    bar.textContent = 'Go 语言服务不可用：' + (plsStatusReason || '未知原因') +
        '。补全、诊断、跳转定义暂不可用；请安装 gopls（go install golang.org/x/tools/gopls@latest）' +
        '或设置 DESIGNER_GOPLS_PATH。';
    bar.classList.add('show');
}

ipc.on('gopls-status', function (ok, reason) {
    renderGoplsStatus(ok, reason);
});

// 首次索引期间 Go 侧会把请求超时从 10s 放宽到 120s，前端同步延长 pending。
ipc.on('gopls-indexing', function (indexing) {
    plsIndexing = !!indexing;
    renderGoplsStatus(plsStatusOK, plsStatusReason);
});

// 页面脚本就绪后主动同步一次，防止状态上报早于 WebView IPC 注册而丢失。
ipc.emit('gopls-status-request', [], function (result) {
    if (!result) return;
    let data = result;
    if (typeof result === 'string') {
        try { data = JSON.parse(result); } catch (e) { return; }
    }
    if (data && typeof data.indexing === 'boolean') plsIndexing = data.indexing;
    renderGoplsStatus(data.ok, data.reason);
});

let conflictReloadButton = document.getElementById('conflict-reload');
if (conflictReloadButton) {
    conflictReloadButton.addEventListener('click', function () {
        resolveFileConflict('reload');
    });
}
let conflictKeepButton = document.getElementById('conflict-keep');
if (conflictKeepButton) {
    conflictKeepButton.addEventListener('click', function () {
        resolveFileConflict('keep');
    });
}
