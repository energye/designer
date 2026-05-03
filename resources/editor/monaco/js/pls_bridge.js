// PLS provider registration and response handlers

var completionRequestID = 0;
var pendingCompletion = null;
var signatureHelpRequestID = 0;
var pendingSignatureHelp = null;
var codeActionRequestID = 0;
var pendingCodeAction = null;

function registerPLSProviders() {
    // Completion provider
    monacoRef.languages.registerCompletionItemProvider('go', {
        triggerCharacters: ['.', '(', '"', "'", '/', '@'],
        provideCompletionItems: function (model, position, token, context) {
            var filePath = getFilePathByModel(model);
            if (!filePath) return {suggestions: []};

            var line = position.lineNumber - 1;
            var column = position.column - 1;
            completionRequestID++;
            var reqID = completionRequestID;

            var triggerKind = 1;
            var triggerChar = '';
            if (context && context.triggerKind === 2) {
                triggerKind = 2;
                triggerChar = context.triggerCharacter || '';
            } else if (context && context.triggerKind === 3) {
                triggerKind = 3;
                triggerChar = context.triggerCharacter || '';
            }

            return new Promise(function (resolve) {
                pendingCompletion = {reqID: reqID, resolve: resolve};
                // Override the timeout-based one since we set pendingCompletion after
                setTimeout(function () {
                    if (pendingCompletion && pendingCompletion.reqID === reqID) {
                        pendingCompletion = null;
                        resolve({suggestions: []});
                    }
                }, 3000);

                ipc.emit('gopls-completion', [{
                    requestID: reqID, file: filePath, line: line, column: column,
                    triggerKind: triggerKind, triggerChar: triggerChar
                }]);
            });
        }
    });

    // Signature Help provider
    monacoRef.languages.registerSignatureHelpProvider('go', {
        signatureHelpTriggerCharacters: ['(', ','],
        provideSignatureHelp: function (model, position, token, context) {
            var filePath = getFilePathByModel(model);
            if (!filePath) return null;

            var line = position.lineNumber - 1;
            var column = position.column - 1;
            signatureHelpRequestID++;
            var reqID = signatureHelpRequestID;

            return new Promise(function (resolve) {
                pendingSignatureHelp = {reqID: reqID, resolve: resolve};
                setTimeout(function () {
                    if (pendingSignatureHelp && pendingSignatureHelp.reqID === reqID) {
                        pendingSignatureHelp = null;
                        resolve(null);
                    }
                }, 3000);

                ipc.emit('gopls-signatureHelp', [{
                    requestID: reqID, file: filePath, line: line, column: column
                }]);
            });
        }
    });

    // Code Action provider
    monacoRef.languages.registerCodeActionProvider('go', {
        provideCodeActions: function (model, range, context) {
            var filePath = getFilePathByModel(model);
            if (!filePath) return {actions: [], dispose: function(){}};

            codeActionRequestID++;
            var reqID = codeActionRequestID;

            var kinds = 'quickfix,source.organizeImports';
            var diags = currentDiagnostics.get(filePath) || [];
            var rangeDiags = diags.filter(function (d) {
                var dStartLine = (d.range && d.range.start ? d.range.start.line + 1 : 1);
                var dEndLine = (d.range && d.range.end ? d.range.end.line + 1 : 1);
                return dEndLine >= range.startLineNumber && dStartLine <= range.endLineNumber;
            });

            return new Promise(function (resolve) {
                pendingCodeAction = {reqID: reqID, resolve: resolve};
                setTimeout(function () {
                    if (pendingCodeAction && pendingCodeAction.reqID === reqID) {
                        pendingCodeAction = null;
                        resolve({actions: [], dispose: function(){}});
                    }
                }, 3000);

                ipc.emit('gopls-codeAction', [{
                    requestID: reqID, file: filePath,
                    startLine: range.startLineNumber - 1, startChar: range.startColumn - 1,
                    endLine: range.endLineNumber - 1, endChar: range.endColumn - 1,
                    kinds: kinds, diagnostics: rangeDiags
                }]);
            });
        }
    });
}

// === Completion Response Handler ===
ipc.on('gopls-completion-response', function (reqID, resultJSON) {
    if (!pendingCompletion || pendingCompletion.reqID !== reqID) return;

    var resolve = pendingCompletion.resolve;
    pendingCompletion = null;

    var items = [];
    try {
        var parsed = JSON.parse(resultJSON);
        if (Array.isArray(parsed)) {
            items = parsed;
        }
    } catch (e) {}

    var K = monacoRef.languages.CompletionItemKind;
    var Snippet = monacoRef.languages.CompletionItemInsertTextRule;
    var suggestions = items.map(function (item) {
        var insertText = item.insertText || item.label;
        var kind = plsKindToMonaco(item.kind, K);
        var isSnippet = item.insertTextFormat === 2;

        // Enhance function/method completions with parentheses and parameter placeholders
        if ((kind === K.Function || kind === K.Method) && !isSnippet) {
            var detail = item.detail || '';
            var sigMatch = detail.match(/^func\(([^)]*)\)/);
            if (sigMatch) {
                var params = sigMatch[1].trim();
                if (params === '') {
                    insertText = insertText + '()';
                } else {
                    var paramList = params.split(',').map(function(p) { return p.trim(); });
                    var placeholders = paramList.map(function(p, i) { return '${' + (i + 1) + ':' + p + '}'; });
                    insertText = insertText + '(' + placeholders.join(', ') + ')$0';
                    isSnippet = true;
                }
            } else if (detail.indexOf('func(') === 0 || detail.indexOf('func (') === 0) {
                insertText = insertText + '()';
            }
        }

        var s = {
            label: item.label,
            kind: kind,
            detail: item.detail || '',
            documentation: item.documentation || '',
            insertText: insertText,
            sortText: item.sortText || item.label,
            filterText: item.filterText || item.label
        };
        if (isSnippet) {
            s.insertTextRules = Snippet.InsertAsSnippet;
        }
        if (item.deprecated) {
            s.tags = [monacoRef.languages.CompletionItemTag.Deprecated];
        }
        if (item.additionalTextEdits && item.additionalTextEdits.length > 0) {
            s.additionalTextEdits = item.additionalTextEdits.map(function (edit) {
                return {
                    range: plsRangeToMonaco(edit.range),
                    text: edit.newText
                };
            });
        }
        return s;
    });

    resolve({suggestions: suggestions});
});

// === Signature Help Response Handler ===
ipc.on('gopls-signatureHelp-response', function (reqID, resultJSON) {
    if (!pendingSignatureHelp || pendingSignatureHelp.reqID !== reqID) return;

    var resolve = pendingSignatureHelp.resolve;
    pendingSignatureHelp = null;

    if (!resultJSON) { resolve(null); return; }

    var data = null;
    try { data = JSON.parse(resultJSON); } catch (e) { resolve(null); return; }

    if (!data.signatures || data.signatures.length === 0) { resolve(null); return; }

    var signatures = data.signatures.map(function (sig) {
        return {
            label: sig.label,
            documentation: sig.documentation || '',
            parameters: (sig.parameters || []).map(function (p) {
                return { label: p.label, documentation: p.documentation || '' };
            })
        };
    });

    resolve({
        value: {
            signatures: signatures,
            activeSignature: data.activeSignature >= 0 ? data.activeSignature : 0,
            activeParameter: data.activeParameter >= 0 ? data.activeParameter : 0
        },
        dispose: function () {}
    });
});

// === Code Action Response Handler ===
ipc.on('gopls-codeAction-response', function (reqID, resultJSON) {
    var actions = [];
    try { actions = JSON.parse(resultJSON); } catch (e) {}
    if (!Array.isArray(actions)) actions = [];

    var isAutoOrganize = (!pendingCodeAction || pendingCodeAction.reqID !== reqID);

    if (isAutoOrganize) {
        actions.forEach(function (action) {
            if (action.kind === 'source.organizeImports' && action.edit && action.edit.changes) {
                applyWorkspaceEdit(action.edit);
            }
        });
        if (pendingOrganizeSave && pendingOrganizeSave.reqID === reqID) {
            clearTimeout(saveAfterOrganizeTimeout);
            saveAfterOrganizeTimeout = setTimeout(function () {
                var fp = pendingOrganizeSave.filePath;
                pendingOrganizeSave = null;
                if (fp && files.has(fp)) {
                    var model = getModelByFilePath(fp);
                    if (model) {
                        ipc.emit('save-file', [{file: fp, content: model.getValue()}], function (result) {
                            if (result === 'ok') {
                                files.get(fp).isDirty = false;
                            }
                        });
                    }
                }
            }, 100);
        }
        return;
    }

    var resolve = pendingCodeAction.resolve;
    pendingCodeAction = null;

    var codeActions = actions.map(function (action) {
        var ca = {
            title: action.title,
            kind: action.kind || 'quickfix',
            diagnostics: [],
            edit: { edits: [] }
        };

        if (action.edit && action.edit.changes) {
            var allEdits = [];
            for (var filePath in action.edit.changes) {
                if (!action.edit.changes.hasOwnProperty(filePath)) continue;
                var edits = action.edit.changes[filePath];
                var model = getModelByFilePath(filePath);
                if (!model) continue;
                for (var i = 0; i < edits.length; i++) {
                    var edit = edits[i];
                    allEdits.push({
                        resource: model.uri,
                        edit: {
                            range: plsRangeToMonaco(edit.range),
                            text: edit.newText
                        }
                    });
                }
            }
            ca.edit.edits = allEdits;
        }

        return ca;
    });

    resolve({actions: codeActions, dispose: function(){}});
});

// === Diagnostics Handler ===
ipc.on('gopls-diagnostics', function (filePath, diagnosticsJSON) {
    if (!monacoRef) return;
    var model = getModelByFilePath(filePath);
    if (!model) return;

    var diags = [];
    try { diags = JSON.parse(diagnosticsJSON); } catch (e) { return; }
    if (!Array.isArray(diags)) return;

    currentDiagnostics.set(filePath, diags);

    var markers = diags.map(function (d) {
        return {
            severity: diagSeverityToMonaco(d.severity),
            message: d.message,
            startLineNumber: (d.range && d.range.start ? d.range.start.line + 1 : 1),
            startColumn: (d.range && d.range.start ? d.range.start.character + 1 : 1),
            endLineNumber: (d.range && d.range.end ? d.range.end.line + 1 : 1),
            endColumn: (d.range && d.range.end ? d.range.end.character + 1 : 1)
        };
    });

    monacoRef.editor.setModelMarkers(model, 'gopls', markers);

    var hasUnusedImport = diags.some(function (d) {
        return d.message && d.message.indexOf('imported and not used') >= 0;
    });
    if (hasUnusedImport && filePath === currentFilePath) {
        autoOrganizeImports(filePath, model);
    }
});

// Request organize imports (used both for auto and before-save)
function requestOrganizeImports(filePath, model, onSave) {
    clearTimeout(organizeImportsTimeout);
    organizeImportsTimeout = setTimeout(function () {
        codeActionRequestID++;
        var reqID = codeActionRequestID;
        var kinds = 'source.organizeImports';
        var diags = currentDiagnostics.get(filePath) || [];
        var fullRange = model.getFullModelRange();

        if (onSave) {
            pendingOrganizeSave = {filePath: filePath, reqID: reqID};
        }

        ipc.emit('gopls-codeAction', [{
            requestID: reqID, file: filePath,
            startLine: fullRange.startLineNumber - 1, startChar: fullRange.startColumn - 1,
            endLine: fullRange.endLineNumber - 1, endChar: fullRange.endColumn - 1,
            kinds: kinds, diagnostics: diags
        }]);
    }, onSave ? 0 : 800);
}

function autoOrganizeImports(filePath, model) {
    requestOrganizeImports(filePath, model, false);
}

function applyWorkspaceEdit(workspaceEdit) {
    if (!workspaceEdit || !workspaceEdit.changes) return;
    applyingEdit = true;
    try {
        for (var filePath in workspaceEdit.changes) {
            if (!workspaceEdit.changes.hasOwnProperty(filePath)) continue;
            var edits = workspaceEdit.changes[filePath];
            var model = getModelByFilePath(filePath);
            if (!model) continue;
            var monacoEdits = plsEditsToMonaco(edits);
            model.applyEdits(monacoEdits);
        }
    } finally {
        applyingEdit = false;
        for (var filePath in workspaceEdit.changes) {
            if (!workspaceEdit.changes.hasOwnProperty(filePath)) continue;
            if (files.has(filePath)) {
                var info = files.get(filePath);
                info.version++;
                ipc.emit('gopls-didChange', [{
                    file: filePath,
                    content: info.model.getValue(),
                    version: info.version
                }]);
            }
        }
    }
}
