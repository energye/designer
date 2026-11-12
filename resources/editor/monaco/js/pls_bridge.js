// PLS provider registration and response handlers

let completionRequestID = 0;
let pendingCompletion = new Map();       // reqID -> {resolve}
let signatureHelpRequestID = 0;
let pendingSignatureHelp = new Map();   // reqID -> {resolve}
let codeActionRequestID = 0;
let pendingCodeAction = new Map();      // reqID -> {resolve}
let hoverRequestID = 0;
let pendingHover = new Map();           // reqID -> {resolve, file, line, column}

// hover 中的 pkg.go.dev 文档链接改写为定义跳转，避免点击后打开系统浏览器。
// 必须使用 Monaco markdown 允许的协议（file/http/https/mailto），自定义协议会被剥成纯文本。
const DEFINITION_LINK_SCHEME = 'file';
const DEFINITION_LINK_PATH = '/designer-goto-definition';

function definitionLinkUri(filePath, line, column) {
    // 双重编码：Monaco 解析 Uri 时会先解码一次 query，避免路径中的 & = % 干扰参数拆分。
    return DEFINITION_LINK_SCHEME + '://' + DEFINITION_LINK_PATH +
        '?file=' + encodeURIComponent(encodeURIComponent(filePath)) +
        '&line=' + line + '&column=' + column;
}

function replaceHoverDefinitionLinks(markdown, filePath, line, column) {
    if (!markdown || markdown.indexOf('pkg.go.dev') < 0) return markdown;
    let uri = definitionLinkUri(filePath, line, column);
    return markdown.replace(/\[[^\]]*\]\(https?:\/\/pkg\.go\.dev\/[^)\s]*\)/g, '[转到定义](' + uri + ')');
}

function parseDefinitionLinkQuery(query) {
    let params = {};
    (query || '').split('&').forEach(function (pair) {
        let idx = pair.indexOf('=');
        if (idx < 0) return;
        try {
            params[decodeURIComponent(pair.slice(0, idx))] = decodeURIComponent(pair.slice(idx + 1));
        } catch (e) {
        }
    });
    return {
        file: params.file || '',
        line: parseInt(params.line || '0', 10) || 0,
        column: parseInt(params.column || '0', 10) || 0
    };
}

function registerDefinitionLinkOpener() {
    if (!monacoRef || !monacoRef.editor || typeof monacoRef.editor.registerLinkOpener !== 'function') {
        return;
    }
    monacoRef.editor.registerLinkOpener({
        open: function (resource) {
            if (resource.scheme !== DEFINITION_LINK_SCHEME || resource.path !== DEFINITION_LINK_PATH) {
                return false;
            }
            let target = parseDefinitionLinkQuery(resource.query);
            if (target.file) {
                requestDefinition(target.file, target.line, target.column);
            }
            return true;
        }
    });
}

function registerPLSProviders() {
    registerDefinitionLinkOpener();
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
                pendingCompletion.set(reqID, {resolve: resolve});
                schedulePLSTimeout(function () {
                    return pendingCompletion.has(reqID);
                }, function () {
                    pendingCompletion.delete(reqID);
                    resolve({suggestions: []});
                });

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
                pendingSignatureHelp.set(reqID, {resolve: resolve});
                schedulePLSTimeout(function () {
                    return pendingSignatureHelp.has(reqID);
                }, function () {
                    pendingSignatureHelp.delete(reqID);
                    resolve(null);
                });

                ipc.emit('gopls-signatureHelp', [{
                    requestID: reqID, file: filePath, line: line, column: column
                }]);
            });
        }
    });

    // Hover provider：鼠标停在代码上时显示 gopls 的类型/文档信息。
    monacoRef.languages.registerHoverProvider('go', {
        provideHover: function (model, position, token) {
            let filePath = getFilePathByModel(model);
            if (!filePath) return null;

            hoverRequestID++;
            let reqID = hoverRequestID;
            return new Promise(function (resolve) {
                pendingHover.set(reqID, {
                    resolve: resolve, file: filePath,
                    line: position.lineNumber - 1, column: position.column - 1
                });
                schedulePLSTimeout(function () {
                    return pendingHover.has(reqID);
                }, function () {
                    pendingHover.delete(reqID);
                    resolve(null);
                });

                ipc.emit('gopls-hover', [{
                    requestID: reqID, file: filePath,
                    line: position.lineNumber - 1, column: position.column - 1
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
                pendingCodeAction.set(reqID, {resolve: resolve});
                schedulePLSTimeout(function () {
                    return pendingCodeAction.has(reqID);
                }, function () {
                    pendingCodeAction.delete(reqID);
                    resolve({
                        actions: [], dispose: function () {
                        }
                    });
                });

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
    let pending = pendingCompletion.get(reqID);
    if (!pending) return;
    pendingCompletion.delete(reqID);

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
                    // 占位符语法：${1:name}，末尾 $0 表示最终光标位置。
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
    pending.resolve({suggestions: suggestions});
});

// === Signature Help Response Handler ===
ipc.on('gopls-signatureHelp-response', function (reqID, resultJSON) {
    let pending = pendingSignatureHelp.get(reqID);
    if (!pending) return;
    pendingSignatureHelp.delete(reqID);

    if (!resultJSON) {
        pending.resolve(null);
        return;
    }

    let data = null;
    try {
        data = JSON.parse(resultJSON);
    } catch (e) {
        pending.resolve(null);
        return;
    }

    if (!data.signatures || data.signatures.length === 0) {
        pending.resolve(null);
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

    pending.resolve({
        value: {
            signatures: signatures,
            activeSignature: data.activeSignature >= 0 ? data.activeSignature : 0,
            activeParameter: data.activeParameter >= 0 ? data.activeParameter : 0
        },
        dispose: function () {
        }
    });
});

// === Hover Response Handler ===
ipc.on('gopls-hover-response', function (reqID, resultJSON) {
    let pending = pendingHover.get(reqID);
    if (!pending) return;
    pendingHover.delete(reqID);

    if (!resultJSON || resultJSON === 'null') {
        pending.resolve(null);
        return;
    }
    let data;
    try {
        data = JSON.parse(resultJSON);
    } catch (e) {
        pending.resolve(null);
        return;
    }
    if (!data || !Array.isArray(data.contents) || data.contents.length === 0) {
        pending.resolve(null);
        return;
    }

    // Monaco 的 contents 是 IMarkdownString[]；LSP MarkedString 带 language
    // 时包成 fenced code，保证签名代码块高亮。
    // 定义链接优先用 hover range 起点（符号起点），鼠标恰好停在词边界时更稳。
    let defLine = pending.line;
    let defColumn = pending.column;
    if (data.range && data.range.start) {
        if (typeof data.range.start.line === 'number') defLine = data.range.start.line;
        if (typeof data.range.start.character === 'number') defColumn = data.range.start.character;
    }
    let contents = data.contents.map(function (c) {
        let value = c.value || '';
        if (c.language) value = '```' + c.language + '\n' + value + '\n```';
        if (pending.file) {
            value = replaceHoverDefinitionLinks(value, pending.file, defLine, defColumn);
        }
        return {value: value};
    });
    let hover = {contents: contents};
    if (data.range && data.range.start && data.range.end) {
        hover.range = {
            startLineNumber: data.range.start.line + 1,
            startColumn: data.range.start.character + 1,
            endLineNumber: data.range.end.line + 1,
            endColumn: data.range.end.character + 1
        };
    }
    pending.resolve(hover);
});

// === Code Action Response Handler ===
ipc.on('gopls-codeAction-response', function (reqID, resultJSON) {
    let actions = [];
    try {
        actions = JSON.parse(resultJSON);
    } catch (e) {
    }
    if (!Array.isArray(actions)) actions = [];

    // 保存流程（手动 Ctrl+S / 自动保存）发出的 organizeImports 请求。
    let saveReq = pendingOrganizeSave.get(reqID);
    if (saveReq) {
        pendingOrganizeSave.delete(reqID);
        applyOrganizeImports(actions);
        if (saveReq.phase === 'organize') {
            manualSaveAfterOrganize(saveReq.filePath);
        } else {
            autoSaveAfterOrganize(saveReq.filePath);
        }
        return;
    }

    // 诊断触发的自动整理导入（或已超时的 provider 请求）：只应用编辑。
    if (!pendingCodeAction.has(reqID)) {
        applyOrganizeImports(actions);
        return;
    }

    let resolve = pendingCodeAction.get(reqID).resolve;
    pendingCodeAction.delete(reqID);

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

function applyOrganizeImports(actions) {
    actions.forEach(function (action) {
        if (action.kind === 'source.organizeImports' && action.edit && action.edit.changes) {
            applyWorkspaceEdit(action.edit);
        }
    });
}

// 手动保存：整理导入完成 → 格式化 → 落盘。
function manualSaveAfterOrganize(filePath) {
    if (!files.has(filePath)) {
        saveInProgress = false;
        return;
    }
    let afterOrganizeModel = getModelByFilePath(filePath);
    if (!afterOrganizeModel) {
        saveInProgress = false;
        return;
    }

    requestFormatting(filePath, function (formatEdits) {
        if (!files.has(filePath)) {
            saveInProgress = false;
            return;
        }
        let currentModel = getModelByFilePath(filePath);
        if (!currentModel) {
            saveInProgress = false;
            return;
        }

        if (formatEdits.length > 0) {
            applyingEdit = true;
            try {
                // 倒序应用，避免前面的编辑改变后面编辑的坐标。
                let monacoEdits = formatEdits.map(function (edit) {
                    return {range: plsRangeToMonaco(edit.range), text: edit.newText};
                }).reverse();
                currentModel.applyEdits(monacoEdits);
            } finally {
                applyingEdit = false;
            }
            // applyingEdit 抑制了 onDidChangeModelContent，必须显式把格式化后的
            // 文本同步给 gopls，否则它仍持有格式化前的内容。
            syncModelChangeToGopls(filePath, currentModel);
        }
        saveModelToDisk(filePath, function (ok) {
            if (ok) showNotification('已保存: ' + getFileName(filePath));
            saveInProgress = false;
        });
    });
}

// 把当前 model 内容写盘；成功时清脏并撤掉冲突提示。
function saveModelToDisk(filePath, onComplete) {
    if (!files.has(filePath)) {
        if (onComplete) onComplete(false);
        return;
    }
    let model = getModelByFilePath(filePath);
    if (!model) {
        if (onComplete) onComplete(false);
        return;
    }
    let saveInfo = files.get(filePath);
    ipc.emit('save-file', [{
        file: filePath,
        content: model.getValue(),
        readOnly: !!(saveInfo && saveInfo.readOnly)
    }], function (result) {
        let ok = result === 'ok' && files.has(filePath);
        if (ok) {
            files.get(filePath).isDirty = false;
            clearFileConflict(filePath);
        } else {
            showNotification('保存失败: ' + (result || '未知错误'));
        }
        if (onComplete) onComplete(ok);
    });
}

// 自动保存：整理导入完成后直接落盘。
function autoSaveAfterOrganize(filePath) {
    if (!filePath) return;
    saveModelToDisk(filePath, null);
}

// === Formatting Response Handler ===
ipc.on('gopls-formatting-response', function (reqID, resultJSON) {
    let pending = pendingFormatting.get(reqID);
    if (!pending) return;
    pendingFormatting.delete(reqID);

    let edits = [];
    try {
        edits = JSON.parse(resultJSON);
    } catch (e) {
    }
    if (!Array.isArray(edits)) edits = [];

    pending.callback(edits);
});

// === Definition Response Handler (async) ===
ipc.on('gopls-definition-response', function (reqID, resultJSON) {
    let pending = pendingDefinition.get(reqID);
    if (!pending) return;
    pendingDefinition.delete(reqID);

    if (!resultJSON || resultJSON === 'null') return;

    let data = null;
    try {
        data = JSON.parse(resultJSON);
    } catch (e) {
        return;
    }
    if (!data.file) return;

    let targetPath = data.file;
    let targetLine = data.range.start.line;
    let targetCol = data.range.start.character;

    // 用请求发起时的源文件，而不是当前 editor model（可能已被切换）。
    // 归一化比较避免中文/空格/盘符大小写导致同文件被当成跨文件跳转。
    if (normalizePathKey(targetPath) === normalizePathKey(pending.sourceFile)) {
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
    let key = resolveFilePath(filePath);
    if (key === null) return;
    let model = files.get(key).model;
    if (!model) return;

    let diags = [];
    try {
        diags = JSON.parse(diagnosticsJSON);
    } catch (e) {
        return;
    }
    if (!Array.isArray(diags)) return;

    // 用 files 的真实 key 存诊断，CodeAction 请求才能取到。
    currentDiagnostics.set(key, diags);

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
    if (hasUnusedImport && key === currentFilePath) {
        autoOrganizeImports(key, model);
    }
});

// 手动编辑被 applyingEdit 抑制时，显式把 model 快照同步给 gopls。
// version 保持单调递增；Go 侧派发器会再取 max 统一分配。
function syncModelChangeToGopls(filePath, model) {
    let info = files.get(filePath);
    if (!info) return;
    info.version++;
    ipc.emit('gopls-didChange', [{
        file: filePath,
        content: model.getValue(),
        version: info.version
    }]);
}

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
            pendingOrganizeSave.set(reqID, {filePath: filePath});
            schedulePLSTimeout(function () {
                return pendingOrganizeSave.has(reqID);
            }, function () {
                pendingOrganizeSave.delete(reqID);
                showNotification('整理导入超时，已跳过自动保存');
            });
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
    let changed = [];
    applyingEdit = true;
    try {
        for (let filePath in workspaceEdit.changes) {
            if (!workspaceEdit.changes.hasOwnProperty(filePath)) continue;
            let model = getModelByFilePath(filePath);
            if (!model) continue;
            model.applyEdits(plsEditsToMonaco(workspaceEdit.changes[filePath]));
            let key = resolveFilePath(filePath);
            changed.push(key === null ? filePath : key);
        }
    } finally {
        applyingEdit = false;
        // applyingEdit 结束后才发通知，保证内容快照是编辑后的最终状态。
        changed.forEach(function (key) {
            if (files.has(key)) syncModelChangeToGopls(key, files.get(key).model);
        });
    }
}
