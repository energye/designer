// Monaco editor initialization and content change handling

var editor = null;
var monacoRef = null;
var currentFilePath = '';
var saveTimeout = null;
var applyingEdit = false;
var organizeImportsTimeout = null;
var pendingOrganizeSave = null;
var saveAfterOrganizeTimeout = null;
var files = new Map();
var currentDiagnostics = new Map();

require.config({paths: {vs: './vs'}});
require(['vs/editor/editor.main'], function (monaco) {
    monacoRef = monaco;

    editor = monaco.editor.create(document.getElementById('container'), {
        value: '', language: 'go', theme: 'vs-dark', fontSize: 14,
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

    registerLSPProviders();

    // Content change handler
    editor.onDidChangeModelContent(function () {
        if (applyingEdit) return;

        var changedFilePath = getFilePathByModel(editor.getModel());
        if (changedFilePath && files.has(changedFilePath)) {
            var info = files.get(changedFilePath);
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

    ipc.emit("monaco-inited", []);

    // Go to Definition (Ctrl+Click)
    editor.onMouseDown(function (e) {
        if (e.event.ctrlKey && e.event.leftButton) {
            var position = e.target.position;
            if (!position) return;
            var model = editor.getModel();
            if (!model) return;
            var filePath = getFilePathByModel(model);
            if (!filePath) return;

            ipc.emit('gopls-definition', [{
                requestID: 0, file: filePath,
                line: position.lineNumber - 1, column: position.column - 1
            }], function (response) {
                if (!response || response === 'null') return;
                var data = JSON.parse(response);
                if (!data.file) return;

                var targetPath = data.file;
                var targetLine = data.range.start.line;
                var targetCol = data.range.start.character;

                if (targetPath === filePath) {
                    editor.revealLineInCenter(targetLine + 1);
                    editor.setPosition({
                        lineNumber: targetLine + 1,
                        column: targetCol + 1
                    });
                    editor.focus();
                } else {
                    ipc.emit('go-to-definition', [{file: targetPath, range: data.range}]);
                }
            });
        }
    });
});

// Helper: create Monaco model URI from file path
function modelUri(path) {
    return monacoRef.Uri.parse('file:///' + path.replace(/\\/g, '/'));
}
