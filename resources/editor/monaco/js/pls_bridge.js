// PLS provider registration and response handlers

let completionRequestID = 0;
let pendingCompletion = null;
let signatureHelpRequestID = 0;
let pendingSignatureHelp = null;
let codeActionRequestID = 0;
let pendingCodeAction = null;

function registerPLSProviders() {
    // Completion provider
    monacoRef.languages.registerCompletionItemProvider('go', {
        triggerCharacters: ['.', '('],
        provideCompletionItems: function (model, position, token, context) {
            let filePath = getFilePathByModel(model);
            if (!filePath) return {suggestions: []};

            let line = position.lineNumber - 1;
            let column = position.column - 1;
            completionRequestID++;
            let reqID = completionRequestID;

            let triggerKind = 1;
            let triggerChar = '';
            if (context && context.triggerKind === 2) {
                triggerKind = 2;
                triggerChar = context.triggerCharacter || '';
            } else if (context && context.triggerKind === 3) {
                triggerKind = 3;
                triggerChar = context.triggerCharacter || '';
            }

            return new Promise(function (resolve) {
                pendingCompletion = {reqID: reqID, resolve: resolve};
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
            let filePath = getFilePathByModel(model);
            if (!filePath) return null;

            let line = position.lineNumber - 1;
            let column = position.column - 1;
            signatureHelpRequestID++;
            let reqID = signatureHelpRequestID;

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
            let filePath = getFilePathByModel(model);
            if (!filePath) return {
                actions: [], dispose: function () {
                }
            };

            codeActionRequestID++;
            let reqID = codeActionRequestID;

            let kinds = 'quickfix,source.organizeImports';
            let diags = currentDiagnostics.get(filePath) || [];
            let rangeDiags = diags.filter(function (d) {
                let dStartLine = (d.range && d.range.start ? d.range.start.line + 1 : 1);
                let dEndLine = (d.range && d.range.end ? d.range.end.line + 1 : 1);
                return dEndLine >= range.startLineNumber && dStartLine <= range.endLineNumber;
            });

            return new Promise(function (resolve) {
                pendingCodeAction = {reqID: reqID, resolve: resolve};
                setTimeout(function () {
                    if (pendingCodeAction && pendingCodeAction.reqID === reqID) {
                        pendingCodeAction = null;
                        resolve({
                            actions: [], dispose: function () {
                            }
                        });
                    }
                }, 3000);

                ipc.emit('gopls-codeAction', [{
                    requestID: reqID, file: filePath,
                    startLine: range.startLineNumber - 1, startChar: range.startColumn - 1,
                    endLine: range.endLineNumber - 1, endChar: range.endColumn - 1,
                    kinds: kinds, diagnostics: rangeDiags
                }]);
            }).catch(function () {
                // Prevent "Uncaught (in promise) Canceled" errors
                return {
                    actions: [], dispose: function () {
                    }
                };
            });
        }
    });
}

// === Completion Response Handler ===
ipc.on('gopls-completion-response', function (reqID, resultJSON) {
    if (!pendingCompletion || pendingCompletion.reqID !== reqID) return;

    let resolve = pendingCompletion.resolve;
    pendingCompletion = null;

    let items = [];
    try {
        let parsed = JSON.parse(resultJSON);
        if (Array.isArray(parsed)) {
            items = parsed;
        }
    } catch (e) {
    }

    let K = monacoRef.languages.CompletionItemKind;
    let Snippet = monacoRef.languages.CompletionItemInsertTextRule;
    let suggestions = items.map(function (item) {
        let insertText = item.insertText || item.label;
        let kind = plsKindToMonaco(item.kind, K);
        let isSnippet = item.insertTextFormat === 2;

        // Enhance function/method completions with parentheses and parameter placeholders
        if ((kind === K.Function || kind === K.Method) && !isSnippet) {
            let detail = item.detail || '';
            let sigMatch = detail.match(/^func\(([^)]*)\)/);
            if (sigMatch) {
                let params = sigMatch[1].trim();
                if (params === '') {
                    insertText = insertText + '()';
                } else {
                    let paramList = params.split(',').map(function (p) {
                        return p.trim();
                    });
                    let placeholders = paramList.map(function (p, i) {
                        return '${' + (i + 1) + ':' + p + '}';
                    });
                    insertText = insertText + '(' + placeholders.join(', ') + ')$0';
                    isSnippet = true;
                }
            } else if (detail.indexOf('func(') === 0 || detail.indexOf('func (') === 0) {
                insertText = insertText + '()';
            }
        }

        let s = {
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

    let resolve = pendingSignatureHelp.resolve;
    pendingSignatureHelp = null;

    if (!resultJSON) {
        resolve(null);
        return;
    }

    let data = null;
    try {
        data = JSON.parse(resultJSON);
    } catch (e) {
        resolve(null);
        return;
    }

    if (!data.signatures || data.signatures.length === 0) {
        resolve(null);
        return;
    }

    let signatures = data.signatures.map(function (sig) {
        return {
            label: sig.label,
            documentation: sig.documentation || '',
            parameters: (sig.parameters || []).map(function (p) {
                return {label: p.label, documentation: p.documentation || ''};
            })
        };
    });

    resolve({
        value: {
            signatures: signatures,
            activeSignature: data.activeSignature >= 0 ? data.activeSignature : 0,
            activeParameter: data.activeParameter >= 0 ? data.activeParameter : 0
        },
        dispose: function () {
        }
    });
});

// === Code Action Response Handler ===
ipc.on('gopls-codeAction-response', function (reqID, resultJSON) {
    let actions = [];
    try {
        actions = JSON.parse(resultJSON);
    } catch (e) {
    }
    if (!Array.isArray(actions)) actions = [];

    let isManualSave = (pendingOrganizeSave && pendingOrganizeSave.reqID === reqID && pendingOrganizeSave.phase === 'organize');
    let isAutoOrganize = (!pendingCodeAction || pendingCodeAction.reqID !== reqID) && !isManualSave;

    if (isAutoOrganize || isManualSave) {
        // Apply organize imports edits
        actions.forEach(function (action) {
            if (action.kind === 'source.organizeImports' && action.edit && action.edit.changes) {
                applyWorkspaceEdit(action.edit);
            }
        });

        if (isManualSave) {
            // Manual save flow: organize imports done, now format then save
            let fp = pendingOrganizeSave.filePath;
            pendingOrganizeSave = null;

            if (!files.has(fp)) {
                saveInProgress = false;
                return;
            }
            let afterOrganizeModel = getModelByFilePath(fp);
            if (!afterOrganizeModel) {
                saveInProgress = false;
                return;
            }

            requestFormatting(fp, function (formatEdits) {
                if (!files.has(fp)) {
                    saveInProgress = false;
                    return;
                }
                let currentModel = getModelByFilePath(fp);
                if (!currentModel) {
                    saveInProgress = false;
                    return;
                }

                if (formatEdits.length > 0) {
                    applyingEdit = true;
                    try {
                        let monacoEdits = formatEdits.map(function (edit) {
                            return {range: plsRangeToMonaco(edit.range), text: edit.newText};
                        }).reverse();
                        currentModel.applyEdits(monacoEdits);
                    } finally {
                        applyingEdit = false;
                    }
                }
                // Save the file to disk
                ipc.emit('save-file', [{file: fp, content: currentModel.getValue()}], function (result) {
                    if (result === 'ok' && files.has(fp)) {
                        files.get(fp).isDirty = false;
                    }
                    saveInProgress = false;
                });
                showNotification('已保存: ' + getFileName(fp));
            });
        } else if (pendingOrganizeSave && pendingOrganizeSave.reqID === reqID) {
            // Auto-save path (from autoSave function, no phase set)
            let fp = pendingOrganizeSave.filePath;
            pendingOrganizeSave = null;
            if (fp && files.has(fp)) {
                let model = getModelByFilePath(fp);
                if (model) {
                    ipc.emit('save-file', [{file: fp, content: model.getValue()}], function (result) {
                        if (result === 'ok') {
                            files.get(fp).isDirty = false;
                        }
                    });
                }
            }
        }
        return;
    }

    let resolve = pendingCodeAction.resolve;
    pendingCodeAction = null;

    let codeActions = actions.map(function (action) {
        let ca = {
            title: action.title,
            kind: action.kind || 'quickfix',
            diagnostics: [],
            isPreferred: !!action.isPreferred
        };

        if (action.edit && action.edit.changes) {
            let allEdits = [];
            for (let filePath in action.edit.changes) {
                if (!action.edit.changes.hasOwnProperty(filePath)) continue;
                let edits = action.edit.changes[filePath];
                let model = getModelByFilePath(filePath);
                if (!model) continue;
                for (let i = 0; i < edits.length; i++) {
                    let edit = edits[i];
                    allEdits.push({
                        resource: model.uri,
                        versionId: undefined,
                        textEdit: {
                            range: plsRangeToMonaco(edit.range),
                            text: edit.newText
                        }
                    });
                }
            }
            if (allEdits.length > 0) {
                ca.edit = {edits: allEdits};
            }
        }

        return ca;
    });

    resolve({
        actions: codeActions, dispose: function () {
        }
    });
});

// === Formatting Response Handler ===
ipc.on('gopls-formatting-response', function (reqID, resultJSON) {
    if (!pendingFormatting || pendingFormatting.reqID !== reqID) return;

    let callback = pendingFormatting.callback;
    pendingFormatting = null;

    let edits = [];
    try {
        edits = JSON.parse(resultJSON);
    } catch (e) {
    }
    if (!Array.isArray(edits)) edits = [];

    callback(edits);
});

// === Definition Response Handler (async) ===
ipc.on('gopls-definition-response', function (reqID, resultJSON) {
    if (!pendingDefinition || pendingDefinition.reqID !== reqID) return;
    pendingDefinition = null;

    if (!resultJSON || resultJSON === 'null') return;

    let data = null;
    try {
        data = JSON.parse(resultJSON);
    } catch (e) {
        return;
    }
    if (!data.file) return;

    let sourceFile = getFilePathByModel(editor.getModel());
    let targetPath = data.file;
    let targetLine = data.range.start.line;
    let targetCol = data.range.start.character;

    if (targetPath === sourceFile) {
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

// === Diagnostics Handler ===
ipc.on('gopls-diagnostics', function (filePath, diagnosticsJSON) {
    if (!monacoRef) return;
    let model = getModelByFilePath(filePath);
    if (!model) return;

    let diags = [];
    try {
        diags = JSON.parse(diagnosticsJSON);
    } catch (e) {
        return;
    }
    if (!Array.isArray(diags)) return;

    currentDiagnostics.set(filePath, diags);

    let markers = diags.map(function (d) {
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

    let hasUnusedImport = diags.some(function (d) {
        return d.message && d.message.indexOf('imported and not used') >= 0;
    });
    if (hasUnusedImport && filePath === currentFilePath) {
        autoOrganizeImports(filePath, model);
    }
});

// Request organize imports (used for auto-diagnostics only, NOT during manualSave)
function requestOrganizeImports(filePath, model, onSave) {
    clearTimeout(organizeImportsTimeout);
    organizeImportsTimeout = setTimeout(function () {
        codeActionRequestID++;
        let reqID = codeActionRequestID;
        let kinds = 'source.organizeImports';
        let diags = currentDiagnostics.get(filePath) || [];
        let fullRange = model.getFullModelRange();

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
        for (let filePath in workspaceEdit.changes) {
            if (!workspaceEdit.changes.hasOwnProperty(filePath)) continue;
            let edits = workspaceEdit.changes[filePath];
            let model = getModelByFilePath(filePath);
            if (!model) continue;
            let monacoEdits = plsEditsToMonaco(edits);
            model.applyEdits(monacoEdits);
        }
    } finally {
        applyingEdit = false;
        for (let filePath in workspaceEdit.changes) {
            if (!workspaceEdit.changes.hasOwnProperty(filePath)) continue;
            if (files.has(filePath)) {
                let info = files.get(filePath);
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
