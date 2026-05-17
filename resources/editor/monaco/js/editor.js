// Monaco editor initialization and content change handling

let editor = null;
let monacoRef = null;
let currentFilePath = '';
let saveTimeout = null;
let applyingEdit = false;
let organizeImportsTimeout = null;
let pendingOrganizeSave = null;
let files = new Map();
let currentDiagnostics = new Map();
let definitionRequestID = 0;
let pendingDefinition = null;
let formattingRequestID = 0;
let pendingFormatting = null;
let saveInProgress = false;

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

            definitionRequestID++;
            let reqID = definitionRequestID;
            pendingDefinition = {reqID: reqID, sourceFile: filePath};

            ipc.emit('gopls-definition', [{
                requestID: reqID, file: filePath,
                line: position.lineNumber - 1, column: position.column - 1
            }]);
        }
    });
});

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
    pendingFormatting = {reqID: reqID, callback: callback};

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
    if (!info) return;

    // If a previous save is stuck (e.g. pendingOrganizeSave was overwritten),
    // allow a new manual save to proceed
    if (saveInProgress) {
        if (pendingOrganizeSave && pendingOrganizeSave.phase === 'organize') {
            return; // active save in progress, wait
        }
        // Stuck state - reset and allow new save
        saveInProgress = false;
        pendingOrganizeSave = null;
    }

    saveInProgress = true;
    clearTimeout(saveTimeout);
    clearTimeout(organizeImportsTimeout);
    pendingOrganizeSave = null;

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
    pendingOrganizeSave = {filePath: filePath, reqID: organizeReqID, phase: 'organize'};

    ipc.emit('gopls-codeAction', [{
        requestID: organizeReqID, file: filePath,
        startLine: fullRange.startLineNumber - 1, startChar: fullRange.startColumn - 1,
        endLine: fullRange.endLineNumber - 1, endChar: fullRange.endColumn - 1,
        kinds: 'source.organizeImports', diagnostics: currentDiagnostics.get(filePath) || []
    }]);
}
