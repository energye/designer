// PLS kind mapping and diagnostic severity utilities

var goplsStatus = 'loading'; // 'loading' | 'ready' | 'unavailable'

function plsKindToMonaco(kind, K) {
    switch (kind) {
        case 1: return K.Text; case 2: return K.Method; case 3: return K.Function;
        case 4: return K.Constructor; case 5: return K.Field; case 6: return K.Variable;
        case 7: return K.Class; case 8: return K.Interface; case 9: return K.Module;
        case 10: return K.Property; case 11: return K.Unit; case 12: return K.Value;
        case 13: return K.Enum; case 14: return K.Keyword; case 15: return K.Snippet;
        case 16: return K.Color; case 17: return K.File; case 18: return K.Reference;
        case 19: return K.Folder; case 20: return K.EnumMember; case 21: return K.Constant;
        case 22: return K.Struct; case 23: return K.Event; case 24: return K.Operator;
        case 25: return K.TypeParameter; default: return K.Text;
    }
}

function diagSeverityToMonaco(severity) {
    var S = monacoRef.MarkerSeverity;
    switch (severity) {
        case 1: return S.Error; case 2: return S.Warning;
        case 3: return S.Info; case 4: return S.Hint;
        default: return S.Error;
    }
}

function getFilePathByModel(model) {
    var iter = files.entries();
    var pair = iter.next();
    while (!pair.done) {
        var entry = pair.value;
        if (entry[1].model === model) return entry[0];
        pair = iter.next();
    }
    return null;
}

function getModelByFilePath(filePath) {
    var info = files.get(filePath);
    return info ? info.model : null;
}

function getFileName(filePath) {
    return filePath.replace(/\\/g, '/').split('/').pop();
}

function showNotification(message) {
    var toast = document.getElementById('toast');
    toast.textContent = message;
    toast.classList.add('show');
    setTimeout(function () { toast.classList.remove('show'); }, 3000);
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
