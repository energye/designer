// Monaco editor initialization and content change handling

let editor = null;
let monacoRef = null;
let currentFilePath = '';
let saveTimeout = null;
let applyingEdit = false;
let organizeImportsTimeout = null;
// 保存 / 请求按 reqID 存 Map：旧实现用单个变量，连续请求会互相覆盖，
// 先发出的 Promise 永远不 resolve，直到超时才被丢弃。
let pendingOrganizeSave = new Map();
let files = new Map();
let currentDiagnostics = new Map();
let definitionRequestID = 0;
let pendingDefinition = new Map();
let formattingRequestID = 0;
let pendingFormatting = new Map();
let saveInProgress = false;
// gopls 首次索引期间由 Go 侧推送，决定 pending 是否延长等待。
let plsIndexing = false;

function hasActiveManualSave() {
    let active = false;
    pendingOrganizeSave.forEach(function (v) {
        if (v.phase === 'organize') active = true;
    });
    return active;
}

require.config({paths: {vs: './vs'}});
require(['vs/editor/editor.main'], function (monaco) {
    monacoRef = monaco;

    editor = monaco.editor.create(document.getElementById('container'), {
        value: '', language: 'go',
        // theme: 'vs-dark',
        fontSize: 14,
        automaticLayout: true, minimap: {enabled: true},
        scrollBeyondLastLine: false, renderWhitespace: 'selection',
        tabSize: 4, insertSpaces: true,
        multiCursorModifier: 'altKey',
        quickSuggestions: {other: true, comments: false, strings: false},
        suggestOnTriggerCharacters: true, acceptSuggestionOnEnter: 'on',
        wordBasedSuggestions: 'currentDocument',
        suggest: {showMethods: true, showFunctions: true, showFields: true, showVariables: true, showWords: true},
        parameterHints: {enabled: true}
    });

    registerPLSProviders();

    // Content change handler
    editor.onDidChangeModelContent(function () {
        if (applyingEdit) return;

        let changedFilePath = getFilePathByModel(editor.getModel());
        if (changedFilePath && files.has(changedFilePath)) {
            let info = files.get(changedFilePath);
            if (!info.isDirty) {
                info.isDirty = true;
                ipc.emit('set-file-dirty', [{file: changedFilePath, isDirty: true}]);
            }
            info.version++;
            ipc.emit('gopls-didChange', [{
                file: changedFilePath,
                content: info.model.getValue(),
                version: info.version
            }]);
        }
        clearTimeout(saveTimeout);
        saveTimeout = setTimeout(autoSave, 500);
    });

    // Ctrl+S save with format + organize imports
    editor.addCommand(monacoRef.KeyMod.CtrlCmd | monacoRef.KeyCode.KeyS, function () {
        manualSave();
    });

    // Go to Definition (Ctrl+Click) - async to avoid UI freeze
    editor.onMouseDown(function (e) {
        if (e.event.ctrlKey && e.event.leftButton) {
            let position = e.target.position;
            if (!position) return;
            let model = editor.getModel();
            if (!model) return;
            let filePath = getFilePathByModel(model);
            if (!filePath) return;
            requestDefinition(filePath, position.lineNumber - 1, position.column - 1);
        }
    });
});

// 请求 gopls 定义位置；hover 里的定义链接和 Ctrl+Click 共用。
function requestDefinition(filePath, line, column) {
    definitionRequestID++;
    let reqID = definitionRequestID;
    pendingDefinition.set(reqID, {sourceFile: filePath});
    schedulePLSTimeout(function () { return pendingDefinition.has(reqID); }, function () {
        pendingDefinition.delete(reqID);
    });

    ipc.emit('gopls-definition', [{
        requestID: reqID, file: filePath,
        line: line, column: column
    }]);
}

// Helper: create Monaco model URI from file path
function modelUri(path) {
    return monacoRef.Uri.parse('file:///' + path.replace(/\\/g, '/'));
}

// Request gopls formatting for a file
function requestFormatting(filePath, callback) {
    let model = getModelByFilePath(filePath);
    if (!model) {
        callback([]);
        return;
    }

    formattingRequestID++;
    let reqID = formattingRequestID;
    pendingFormatting.set(reqID, {callback: callback});
    schedulePLSTimeout(function () { return pendingFormatting.has(reqID); }, function () {
        pendingFormatting.delete(reqID);
        callback([]);
    });

    ipc.emit('gopls-formatting', [{
        requestID: reqID, file: filePath
    }]);
}

// Manual save: organize imports -> format -> save file
// Order matters: gopls already has current content for organizeImports,
// then after imports are cleaned we format, then save.
function manualSave() {
    if (!currentFilePath || !editor) return;
    let info = files.get(currentFilePath);
    if (!info || info.readOnly) return;

    // If a previous save is stuck (e.g. pendingOrganizeSave was overwritten),
    // allow a new manual save to proceed
    if (saveInProgress) {
        if (hasActiveManualSave()) {
            return; // active save in progress, wait
        }
        // Stuck state - reset and allow new save
        saveInProgress = false;
        pendingOrganizeSave.clear();
    }

    saveInProgress = true;
    clearTimeout(saveTimeout);
    clearTimeout(organizeImportsTimeout);
    pendingOrganizeSave.clear();

    let filePath = currentFilePath;
    let model = getModelByFilePath(filePath);
    if (!model) {
        saveInProgress = false;
        return;
    }

    // Step 1: organize imports first (gopls already has current file state)
    let fullRange = model.getFullModelRange();
    codeActionRequestID++;
    let organizeReqID = codeActionRequestID;
    pendingOrganizeSave.set(organizeReqID, {filePath: filePath, phase: 'organize'});
    schedulePLSTimeout(function () { return pendingOrganizeSave.has(organizeReqID); }, function () {
        pendingOrganizeSave.delete(organizeReqID);
        saveInProgress = false;
    });

    ipc.emit('gopls-codeAction', [{
        requestID: organizeReqID, file: filePath,
        startLine: fullRange.startLineNumber - 1, startChar: fullRange.startColumn - 1,
        endLine: fullRange.endLineNumber - 1, endChar: fullRange.endColumn - 1,
        kinds: 'source.organizeImports', diagnostics: currentDiagnostics.get(filePath) || []
    }]);
}
