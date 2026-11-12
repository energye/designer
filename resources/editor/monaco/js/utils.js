// PLS kind mapping, diagnostic severity and path utilities

// 与 Go 侧 gopls.DefaultRequestTimeout / IndexingRequestTimeout 成对约束：
// 前端必须比 Go 侧晚超时，否则结果回来时 pending 已被清空。
// gopls 首次索引（go.work 覆盖大工作区时实测约 45s）期间使用长超时。
const PLS_GO_TIMEOUT_MS = 30000;
const PLS_GO_INDEXING_TIMEOUT_MS = 120000;
const PLS_REQUEST_TIMEOUT_MS = 35000;
const PLS_INDEXING_TIMEOUT_MS = 125000;
const PLS_TICK_MS = 15000;

function plsRequestTimeout() {
    return plsIndexing ? PLS_INDEXING_TIMEOUT_MS : PLS_REQUEST_TIMEOUT_MS;
}

// 到点后若 gopls 仍在建立索引，继续等待，而不是把迟到结果丢掉。
function schedulePLSTimeout(isPending, onTimeout) {
    let startedAt = Date.now();
    let tick = function () {
        if (!isPending()) return;
        if (Date.now() - startedAt < plsRequestTimeout()) {
            setTimeout(tick, PLS_TICK_MS);
            return;
        }
        onTimeout();
    };
    setTimeout(tick, PLS_TICK_MS);
}

// 与 Go 侧 CanonicalPathKey 对应：忽略盘符大小写、斜杠方向与 file URI 转义。
function normalizePathKey(p) {
    let s = String(p == null ? '' : p);
    if (s.indexOf('%') !== -1 && typeof decodeURIComponent === 'function') {
        try {
            s = decodeURIComponent(s);
        } catch (e) {
        }
    }
    if (s.slice(0, 7).toLowerCase() === 'file://') {
        s = s.slice(7);
        if (s.charAt(0) !== '/') {
            let idx = s.indexOf('/');
            if (idx >= 0) s = s.substring(idx);
        }
    }
    s = s.replace(/\\/g, '/');
    let uriDrive = /^\/([A-Za-z]:\/)/.exec(s);
    if (uriDrive) s = s.substring(1);
    let drive = /^([A-Za-z]):\//.exec(s);
    if (drive) s = drive[1].toLowerCase() + s.substring(1);
    while (s.indexOf('//') !== -1) s = s.replace(/\/\//g, '/');
    return s;
}

// 找到 files 里真实的 key；gopls 回传路径的盘符大小写/斜杠可能与打开时不同。
function resolveFilePath(filePath) {
    if (files.has(filePath)) return filePath;
    let target = normalizePathKey(filePath);
    let iter = files.keys();
    let pair = iter.next();
    while (!pair.done) {
        if (normalizePathKey(pair.value) === target) return pair.value;
        pair = iter.next();
    }
    return null;
}

function plsKindToMonaco(kind, K) {
    switch (kind) {
        case 1:
            return K.Text;
        case 2:
            return K.Method;
        case 3:
            return K.Function;
        case 4:
            return K.Constructor;
        case 5:
            return K.Field;
        case 6:
            return K.Variable;
        case 7:
            return K.Class;
        case 8:
            return K.Interface;
        case 9:
            return K.Module;
        case 10:
            return K.Property;
        case 11:
            return K.Unit;
        case 12:
            return K.Value;
        case 13:
            return K.Enum;
        case 14:
            return K.Keyword;
        case 15:
            return K.Snippet;
        case 16:
            return K.Color;
        case 17:
            return K.File;
        case 18:
            return K.Reference;
        case 19:
            return K.Folder;
        case 20:
            return K.EnumMember;
        case 21:
            return K.Constant;
        case 22:
            return K.Struct;
        case 23:
            return K.Event;
        case 24:
            return K.Operator;
        case 25:
            return K.TypeParameter;
        default:
            return K.Text;
    }
}

function diagSeverityToMonaco(severity) {
    let S = monacoRef.MarkerSeverity;
    switch (severity) {
        case 1:
            return S.Error;
        case 2:
            return S.Warning;
        case 3:
            return S.Info;
        case 4:
            return S.Hint;
        default:
            return S.Error;
    }
}

function getFilePathByModel(model) {
    let iter = files.entries();
    let pair = iter.next();
    while (!pair.done) {
        let entry = pair.value;
        if (entry[1].model === model) return entry[0];
        pair = iter.next();
    }
    return null;
}

function getModelByFilePath(filePath) {
    let key = resolveFilePath(filePath);
    return key === null ? null : files.get(key).model;
}

function getFileName(filePath) {
    return filePath.replace(/\\/g, '/').split('/').pop();
}

function showNotification(message) {
    let toast = document.getElementById('toast');
    toast.textContent = message;
    toast.classList.add('show');
    setTimeout(function () {
        toast.classList.remove('show');
    }, 3000);
}

// Convert PLS range to Monaco range (0-based to 1-based)
function plsRangeToMonaco(range) {
    return {
        startLineNumber: range.start.line + 1,
        startColumn: range.start.character + 1,
        endLineNumber: range.end.line + 1,
        endColumn: range.end.character + 1
    };
}

// Convert PLS text edits to Monaco edits (reversed for safe application)
function plsEditsToMonaco(edits) {
    return edits.slice().reverse().map(function (edit) {
        return {
            range: plsRangeToMonaco(edit.range),
            text: edit.newText
        };
    });
}
