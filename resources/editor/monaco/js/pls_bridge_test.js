// pls_bridge.js 的 Node 回归测试（不依赖浏览器/LCL）。
// 运行: node resources/editor/monaco/js/pls_bridge_test.js
// 用 vm 按 index.html 的顺序加载脚本，注入假 DOM / ipc / Monaco / 时钟。

'use strict';

const fs = require('fs');
const path = require('path');
const vm = require('vm');
const assert = require('assert');

const JS_DIR = __dirname;
let passed = 0;
let failed = 0;

function test(name, fn) {
    try {
        fn();
        console.log('PASS  ' + name);
        passed++;
    } catch (e) {
        console.log('FAIL  ' + name + ' :: ' + (e && e.stack ? e.stack : e));
        failed++;
    }
}

function makeEl() {
    return {
        textContent: '',
        style: {},
        _handlers: {},
        classList: {
            _s: new Set(),
            add(c) { this._s.add(c); },
            remove(c) { this._s.delete(c); },
            contains(c) { return this._s.has(c); }
        },
        addEventListener(ev, fn) { this._handlers[ev] = fn; }
    };
}

function makeModel(value, uri, h) {
    return {
        _value: value,
        _applyValue: null,
        uri: uri || '',
        edits: [],
        getValue() { return this._value; },
        setValue(v) { this._value = v; },
        dispose() {},
        applyEdits(edits) {
            this.edits.push(edits);
            if (this._applyValue !== null) this._value = this._applyValue;
            // 模拟 Monaco 的 onDidChangeModelContent
            if (h && h.state.contentChange) h.state.contentChange();
        },
        getFullModelRange() {
            return {startLineNumber: 1, startColumn: 1, endLineNumber: 10, endColumn: 1};
        }
    };
}

// 同步 Promise：让响应处理器调用 resolve 后能立刻断言结果。
function SyncPromise(executor) {
    this.state = 'pending';
    this.value = undefined;
    const self = this;
    executor(function (v) { self.state = 'resolved'; self.value = v; },
        function (e) { self.state = 'rejected'; self.value = e; });
}
SyncPromise.prototype.then = function (fn) {
    if (this.state === 'resolved') return SyncPromise.resolve(fn(this.value));
    return this;
};
SyncPromise.prototype.catch = function (fn) {
    if (this.state === 'rejected') return SyncPromise.resolve(fn(this.value));
    return this;
};
SyncPromise.resolve = function (v) {
    return new SyncPromise(function (res) { res(v); });
};

const SCRIPTS = ['utils.js', 'editor.js', 'pls_bridge.js', 'file_manager.js', 'ipc_handlers.js'];

function harness(opts) {
    opts = opts || {};
    const els = {};
    ['file-path-bar', 'gopls-status-bar', 'conflict-bar', 'conflict-message',
        'conflict-reload', 'conflict-keep', 'toast', 'container'].forEach(function (id) {
        els[id] = makeEl();
    });

    const state = {
        clock: 0,
        timers: [],
        timerSeq: 0,
        listeners: {},
        emitted: [],
        providers: {},
        markers: [],
        revealed: null,
        position: null,
        focused: 0,
        currentModel: null,
        contentChange: null,
        mouseDown: null,
        diskContent: 'DISK CONTENT',
        saveFileResult: 'ok',
        statusResult: opts.statusResult || {ok: true, reason: '', indexing: false}
    };

    const ipc = {
        on(name, fn) { state.listeners[name] = fn; },
        emit(name, args, cb) {
            state.emitted.push({name: name, args: args, cb: cb});
            if (cb && name === 'save-file') cb(state.saveFileResult);
            if (cb && name === 'reload-file-request') {
                cb(JSON.stringify({content: state.diskContent, modTime: 1, language: 'go'}));
            }
            if (cb && name === 'open-file-request') {
                cb(JSON.stringify({content: state.diskContent, modTime: 1, language: 'go', readOnly: false}));
            }
            if (cb && name === 'gopls-status-request') cb(state.statusResult);
        }
    };

    const fakeEditor = {
        setModel() {}, getModel() { return state.currentModel; },
        saveViewState() { return {}; }, restoreViewState() {},
        updateOptions() {},
        revealLineInCenter(n) { state.revealed = n; },
        setPosition(p) { state.position = p; },
        focus() { state.focused++; },
        addCommand() {},
        onDidChangeModelContent(fn) { state.contentChange = fn; },
        onMouseDown(fn) { state.mouseDown = fn; }
    };

    const K = {Text: 0, Method: 1, Function: 2, Snippet: 15};
    const fakeMonaco = {
        KeyMod: {CtrlCmd: 1}, KeyCode: {KeyS: 2},
        MarkerSeverity: {Error: 8, Warning: 4, Info: 2, Hint: 1},
        languages: {
            CompletionItemKind: K,
            CompletionItemInsertTextRule: {InsertAsSnippet: 1},
            CompletionItemTag: {Deprecated: 1},
            registerCompletionItemProvider(lang, p) { state.providers.completion = p; },
            registerSignatureHelpProvider(lang, p) { state.providers.signature = p; },
            registerCodeActionProvider(lang, p) { state.providers.codeAction = p; },
            registerHoverProvider(lang, p) { state.providers.hover = p; }
        },
        editor: {
            create() { return fakeEditor; },
            createModel(value, lang, uri) { return makeModel(value, uri, null); },
            registerLinkOpener(opener) { state.linkOpener = opener; },
            setModelMarkers(model, owner, markers) {
                state.markers.push({model: model, owner: owner, markers: markers});
            }
        },
        Uri: {parse(s) { return s; }}
    };

    function FakeDate() {}
    FakeDate.now = function () { return state.clock; };

    const requireCallbacks = [];
    const ctx = {
        console: console,
        document: {getElementById(id) { return els[id] || null; }},
        ipc: ipc,
        setTimeout(fn, ms) {
            const id = ++state.timerSeq;
            state.timers.push({id: id, at: state.clock + (ms || 0), fn: fn});
            return id;
        },
        clearTimeout(id) {
            const idx = state.timers.findIndex(function (t) { return t.id === id; });
            if (idx >= 0) state.timers.splice(idx, 1);
        },
        // AMD loader 是异步的：等所有脚本加载完再回调。
        require: Object.assign(function (deps, cb) {
            if (cb) requireCallbacks.push(cb);
        }, {config() {}}),
        Promise: SyncPromise,
        Date: FakeDate,
        decodeURIComponent: decodeURIComponent
    };
    vm.createContext(ctx);
    SCRIPTS.forEach(function (f) {
        vm.runInContext(fs.readFileSync(path.join(JS_DIR, f), 'utf8'), ctx, {filename: f});
    });
    requireCallbacks.forEach(function (cb) { cb(fakeMonaco); });

    return {ctx: ctx, els: els, state: state, fakeMonaco: fakeMonaco};
}

function g(h, name) { return vm.runInContext(name, h.ctx); }
function run(h, code) { return vm.runInContext(code, h.ctx); }
function emitted(h, name) {
    return h.state.emitted.filter(function (e) { return e.name === name; });
}
function clearEmitted(h) { h.state.emitted.length = 0; }
function addFile(h, key, opts) {
    opts = opts || {};
    const model = makeModel(opts.value === undefined ? 'package main\n' : opts.value, key, h);
    if (opts.applyValue !== undefined) model._applyValue = opts.applyValue;
    g(h, 'files').set(key, {
        model: model, modTime: 1, isDirty: !!opts.isDirty,
        viewState: null, readOnly: !!opts.readOnly, version: opts.version || 1
    });
    return model;
}
function setCurrent(h, key) {
    run(h, 'currentFilePath = ' + JSON.stringify(key) + ';');
}
function advance(h, ms) {
    h.state.clock += ms;
    for (let guard = 0; guard < 1000; guard++) {
        let due = null, idx = -1;
        for (let i = 0; i < h.state.timers.length; i++) {
            const t = h.state.timers[i];
            if (t.at <= h.state.clock && (due === null || t.at < due.at)) {
                due = t;
                idx = i;
            }
        }
        if (!due) break;
        h.state.timers.splice(idx, 1);
        due.fn();
    }
}
function edit(line, char, text) {
    return {
        range: {start: {line: line, character: char}, end: {line: line, character: char + 1}},
        newText: text
    };
}

// ---- 超时常量与路径归一化 ----

test('超时常量与 Go 侧对齐，前端始终晚于后端', function () {
    const h = harness();
    const vals = Array.from(g(h, '[PLS_GO_TIMEOUT_MS, PLS_GO_INDEXING_TIMEOUT_MS, PLS_REQUEST_TIMEOUT_MS, PLS_INDEXING_TIMEOUT_MS]'));
    assert.deepStrictEqual(vals, [30000, 120000, 35000, 125000]);
    assert.ok(vals[2] > vals[0] && vals[3] > vals[1]);
    assert.strictEqual(run(h, 'plsRequestTimeout()'), 35000);
    h.state.listeners['gopls-indexing'](true);
    assert.strictEqual(run(h, 'plsRequestTimeout()'), 125000);
});

test('normalizePathKey 统一盘符大小写、斜杠与 URI 转义', function () {
    const h = harness();
    const forms = [
        'E:\\app\\中文 目录\\a.go',
        'e:/app/中文 目录/a.go',
        'file:///e:/app/%E4%B8%AD%E6%96%87%20%E7%9B%AE%E5%BD%95/a.go',
        'FILE:///e:/app/%E4%B8%AD%E6%96%87%20%E7%9B%AE%E5%BD%95/a.go'
    ];
    const keys = forms.map(function (f) {
        return run(h, 'normalizePathKey(' + JSON.stringify(f) + ')');
    });
    assert.ok(keys.every(function (k) { return k === keys[0]; }), keys.join(' | '));
    assert.strictEqual(keys[0], 'e:/app/中文 目录/a.go');
});

test('getModelByFilePath 按归一化路径命中已打开 model', function () {
    const h = harness();
    const key = 'E:\\app\\中文 目录\\a.go';
    const model = addFile(h, key, {value: 'A'});
    assert.strictEqual(run(h, 'getModelByFilePath(' + JSON.stringify('e:/app/中文 目录/a.go') + ')'), model);
    assert.strictEqual(run(h, 'getModelByFilePath(' + JSON.stringify('file:///e:/app/%E4%B8%AD%E6%96%87%20%E7%9B%AE%E5%BD%95/a.go') + ')'), model);
    assert.strictEqual(run(h, 'getModelByFilePath("other.go")'), null);
});

// ---- 补全 ----

test('补全请求参数与响应映射（kind/detail/snippet/additionalTextEdits）', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.completion.provideCompletionItems(
        model, {lineNumber: 2, column: 3}, null, {triggerKind: 2, triggerCharacter: '.'});
    const req = emitted(h, 'gopls-completion')[0].args[0];
    assert.strictEqual(req.file, 'a.go');
    assert.strictEqual(req.line, 1);
    assert.strictEqual(req.column, 2);
    assert.strictEqual(req.triggerKind, 2);
    assert.strictEqual(req.triggerChar, '.');

    h.state.listeners['gopls-completion-response'](req.requestID, JSON.stringify([{
        label: 'fmt.Println', kind: 3, detail: 'func(a ...interface{})',
        documentation: 'doc', insertTextFormat: 2, deprecated: true,
        additionalTextEdits: [edit(0, 0, 'import "fmt"\n')]
    }]));
    assert.strictEqual(p.state, 'resolved');
    const s = p.value.suggestions[0];
    assert.strictEqual(s.kind, h.fakeMonaco.languages.CompletionItemKind.Function);
    assert.strictEqual(s.detail, 'func(a ...interface{})');
    assert.strictEqual(s.documentation, 'doc');
    assert.strictEqual(s.insertText, 'fmt.Println');
    assert.strictEqual(s.insertTextRules, 1);
    assert.strictEqual(s.tags[0], 1);
    assert.strictEqual(s.additionalTextEdits.length, 1);
    assert.strictEqual(s.additionalTextEdits[0].range.startLineNumber, 1);
    assert.strictEqual(s.additionalTextEdits[0].text, 'import "fmt"\n');
    assert.strictEqual(g(h, 'pendingCompletion').size, 0);
});

test('补全并发：两个请求各自 resolve，不被单变量覆盖', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const pos = {lineNumber: 1, column: 1};
    const p1 = h.state.providers.completion.provideCompletionItems(model, pos, null, null);
    const p2 = h.state.providers.completion.provideCompletionItems(model, pos, null, null);
    const reqs = emitted(h, 'gopls-completion');
    assert.strictEqual(reqs.length, 2);
    assert.notStrictEqual(reqs[0].args[0].requestID, reqs[1].args[0].requestID);
    assert.strictEqual(g(h, 'pendingCompletion').size, 2);

    h.state.listeners['gopls-completion-response'](reqs[1].args[0].requestID, JSON.stringify([{label: 'B', kind: 1}]));
    h.state.listeners['gopls-completion-response'](reqs[0].args[0].requestID, JSON.stringify([{label: 'A', kind: 1}]));
    assert.strictEqual(p1.value.suggestions[0].label, 'A');
    assert.strictEqual(p2.value.suggestions[0].label, 'B');
    assert.strictEqual(g(h, 'pendingCompletion').size, 0);
});

test('补全非法响应回空建议', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.completion.provideCompletionItems(model, {lineNumber: 1, column: 1}, null, null);
    const reqID = emitted(h, 'gopls-completion')[0].args[0].requestID;
    h.state.listeners['gopls-completion-response'](reqID, 'not-json');
    assert.strictEqual(p.value.suggestions.length, 0);
});

// ---- 签名帮助 ----

test('签名帮助响应映射与空响应', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.signature.provideSignatureHelp(model, {lineNumber: 3, column: 5}, null, null);
    const req = emitted(h, 'gopls-signatureHelp')[0].args[0];
    assert.strictEqual(req.line, 2);
    assert.strictEqual(req.column, 4);
    h.state.listeners['gopls-signatureHelp-response'](req.requestID, JSON.stringify({
        signatures: [{label: 'f(x int)', documentation: 'doc', parameters: [{label: 'x int', documentation: 'p'}]}],
        activeSignature: 0, activeParameter: -1
    }));
    assert.strictEqual(p.state, 'resolved');
    assert.strictEqual(p.value.value.signatures[0].label, 'f(x int)');
    assert.strictEqual(p.value.value.signatures[0].parameters[0].label, 'x int');
    assert.strictEqual(p.value.value.activeParameter, 0);

    const p2 = h.state.providers.signature.provideSignatureHelp(model, {lineNumber: 1, column: 1}, null, null);
    const req2 = emitted(h, 'gopls-signatureHelp')[1].args[0].requestID;
    h.state.listeners['gopls-signatureHelp-response'](req2, '');
    assert.strictEqual(p2.value, null);
});

test('签名帮助并发不串扰', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const pos = {lineNumber: 1, column: 1};
    const p1 = h.state.providers.signature.provideSignatureHelp(model, pos, null, null);
    const p2 = h.state.providers.signature.provideSignatureHelp(model, pos, null, null);
    const reqs = emitted(h, 'gopls-signatureHelp');
    h.state.listeners['gopls-signatureHelp-response'](reqs[0].args[0].requestID,
        JSON.stringify({signatures: [{label: 'A'}]}));
    h.state.listeners['gopls-signatureHelp-response'](reqs[1].args[0].requestID,
        JSON.stringify({signatures: [{label: 'B'}]}));
    assert.strictEqual(p1.value.value.signatures[0].label, 'A');
    assert.strictEqual(p2.value.value.signatures[0].label, 'B');
});

// ---- Hover（鼠标悬停提示） ----

test('hover provider 请求参数与响应映射（Markdown/代码块/range）', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.hover.provideHover(model, {lineNumber: 3, column: 5}, null);
    const req = emitted(h, 'gopls-hover')[0].args[0];
    assert.strictEqual(req.file, 'a.go');
    assert.strictEqual(req.line, 2);
    assert.strictEqual(req.column, 4);

    h.state.listeners['gopls-hover-response'](req.requestID, JSON.stringify({
        contents: [
            {value: '```go\nfunc Selftest()\n```'},
            {value: 'func Selftest()', language: 'go'}
        ],
        range: {start: {line: 2, character: 4}, end: {line: 2, character: 16}}
    }));
    assert.strictEqual(p.state, 'resolved');
    const hover = p.value;
    assert.strictEqual(hover.contents.length, 2);
    assert.strictEqual(hover.contents[0].value, '```go\nfunc Selftest()\n```');
    assert.ok(hover.contents[1].value.indexOf('```go\n') === 0);
    assert.ok(hover.contents[1].value.indexOf('func Selftest()') !== -1);
    assert.strictEqual(hover.range.startLineNumber, 3);
    assert.strictEqual(hover.range.startColumn, 5);
    assert.strictEqual(hover.range.endLineNumber, 3);
    assert.strictEqual(hover.range.endColumn, 17);
    assert.strictEqual(g(h, 'pendingHover').size, 0);
});

test('hover 空/null/非法响应返回 null', function () {
    const cases = ['null', '', 'not-json', '{"contents":[]}', '{}'];
    cases.forEach(function (payload) {
        const h = harness();
        const model = addFile(h, 'a.go', {});
        const p = h.state.providers.hover.provideHover(model, {lineNumber: 1, column: 1}, null);
        const reqID = emitted(h, 'gopls-hover')[0].args[0].requestID;
        h.state.listeners['gopls-hover-response'](reqID, payload);
        assert.strictEqual(p.value, null, 'payload=' + payload);
        assert.strictEqual(g(h, 'pendingHover').size, 0);
    });
});

test('hover 并发不串扰', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p1 = h.state.providers.hover.provideHover(model, {lineNumber: 1, column: 1}, null);
    const p2 = h.state.providers.hover.provideHover(model, {lineNumber: 2, column: 1}, null);
    const reqs = emitted(h, 'gopls-hover');
    assert.strictEqual(reqs.length, 2);
    assert.strictEqual(g(h, 'pendingHover').size, 2);
    h.state.listeners['gopls-hover-response'](reqs[1].args[0].requestID, JSON.stringify({contents: [{value: 'B'}]}));
    h.state.listeners['gopls-hover-response'](reqs[0].args[0].requestID, JSON.stringify({contents: [{value: 'A'}]}));
    assert.strictEqual(p1.value.contents[0].value, 'A');
    assert.strictEqual(p2.value.contents[0].value, 'B');
    assert.strictEqual(g(h, 'pendingHover').size, 0);
});

test('hover 超时清空 pending 并返回 null', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.hover.provideHover(model, {lineNumber: 1, column: 1}, null);
    advance(h, 35000);
    assert.strictEqual(p.state, 'resolved');
    assert.strictEqual(p.value, null);
    assert.strictEqual(g(h, 'pendingHover').size, 0);
});

test('hover 未打开文件时不发请求', function () {
    const h = harness();
    const stray = makeModel('package x\n', 'stray.go', null);
    const r = h.state.providers.hover.provideHover(stray, {lineNumber: 1, column: 1}, null);
    assert.strictEqual(r, null);
    assert.strictEqual(emitted(h, 'gopls-hover').length, 0);
});

// ---- CodeAction ----

test('CodeAction 响应映射，中文/空格路径按归一化命中 model', function () {
    const h = harness();
    const key = 'E:\\app\\中文 目录\\a.go';
    const model = addFile(h, key, {});
    g(h, 'currentDiagnostics').set(key, []);
    const p = h.state.providers.codeAction.provideCodeActions(model,
        {startLineNumber: 1, startColumn: 1, endLineNumber: 1, endColumn: 1}, {});
    const req = emitted(h, 'gopls-codeAction')[0].args[0];
    assert.strictEqual(req.kinds, 'quickfix,source.organizeImports');

    const changes = {};
    changes['e:/app/中文 目录/a.go'] = [edit(2, 3, 'X')];
    h.state.listeners['gopls-codeAction-response'](req.requestID, JSON.stringify([
        {title: 'fix', kind: 'quickfix', isPreferred: true, edit: {changes: changes}}
    ]));
    assert.strictEqual(p.state, 'resolved');
    const ca = p.value.actions[0];
    assert.strictEqual(ca.title, 'fix');
    assert.strictEqual(ca.isPreferred, true);
    assert.strictEqual(ca.edit.edits.length, 1);
    assert.strictEqual(ca.edit.edits[0].resource, key);
    assert.strictEqual(ca.edit.edits[0].textEdit.range.startLineNumber, 3);
});

test('CodeAction 并发不串扰', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const range = {startLineNumber: 1, startColumn: 1, endLineNumber: 1, endColumn: 1};
    const p1 = h.state.providers.codeAction.provideCodeActions(model, range, {});
    const p2 = h.state.providers.codeAction.provideCodeActions(model, range, {});
    const reqs = emitted(h, 'gopls-codeAction');
    h.state.listeners['gopls-codeAction-response'](reqs[0].args[0].requestID,
        JSON.stringify([{title: 'A'}]));
    h.state.listeners['gopls-codeAction-response'](reqs[1].args[0].requestID,
        JSON.stringify([{title: 'B'}]));
    assert.strictEqual(p1.value.actions[0].title, 'A');
    assert.strictEqual(p2.value.actions[0].title, 'B');
});

// ---- 诊断 ----

test('诊断按归一化路径命中、写入 markers 与真实 key 的 currentDiagnostics', function () {
    const h = harness();
    const key = 'E:\\app\\中文 目录\\a.go';
    const model = addFile(h, key, {});
    h.state.listeners['gopls-diagnostics']('e:/app/中文 目录/a.go', JSON.stringify([
        {severity: 1, message: 'err', range: {start: {line: 2, character: 3}, end: {line: 2, character: 7}}}
    ]));
    assert.strictEqual(h.state.markers.length, 1);
    assert.strictEqual(h.state.markers[0].model, model);
    assert.strictEqual(h.state.markers[0].owner, 'gopls');
    assert.strictEqual(h.state.markers[0].markers[0].startLineNumber, 3);
    assert.strictEqual(h.state.markers[0].markers[0].endColumn, 8);
    assert.ok(g(h, 'currentDiagnostics').has(key));
});

test('诊断指向未打开文件时不产生 markers', function () {
    const h = harness();
    h.state.listeners['gopls-diagnostics']('not-open.go', JSON.stringify([
        {severity: 1, message: 'err', range: {start: {line: 0, character: 0}, end: {line: 0, character: 1}}}
    ]));
    assert.strictEqual(h.state.markers.length, 0);
    assert.strictEqual(g(h, 'currentDiagnostics').size, 0);
});

test('unused import 诊断触发自动整理导入', function () {
    const h = harness();
    const key = 'a.go';
    addFile(h, key, {});
    setCurrent(h, key);
    h.state.listeners['gopls-diagnostics'](key, JSON.stringify([
        {severity: 1, message: '"fmt" imported and not used', range: {start: {line: 0, character: 0}, end: {line: 0, character: 1}}}
    ]));
    advance(h, 800);
    const reqs = emitted(h, 'gopls-codeAction');
    assert.ok(reqs.some(function (e) {
        return e.args[0].kinds === 'source.organizeImports';
    }));
});

// ---- 格式化 / 保存 / applyingEdit ----

test('格式化 edits 倒序应用，applyingEdit 抑制重复 didChange 后显式同步', function () {
    const h = harness();
    const key = 'a.go';
    const model = addFile(h, key, {value: 'BEFORE', applyValue: 'FORMATTED', version: 5});
    h.state.currentModel = model;
    setCurrent(h, key);
    clearEmitted(h);

    run(h, 'manualSave();');
    const organizeReqID = emitted(h, 'gopls-codeAction')[0].args[0].requestID;
    const orgChanges = {};
    orgChanges[key] = [edit(0, 0, 'I1'), edit(1, 0, 'I2')];
    h.state.listeners['gopls-codeAction-response'](organizeReqID, JSON.stringify([
        {title: 'organize', kind: 'source.organizeImports', edit: {changes: orgChanges}}
    ]));

    // applyWorkspaceEdit 内部倒序：第二个编辑排在最前。
    assert.strictEqual(model.edits[0][0].text, 'I2');
    assert.strictEqual(model.edits[0][1].text, 'I1');
    assert.strictEqual(g(h, 'applyingEdit'), false);

    const fmtReq = emitted(h, 'gopls-formatting')[0].args[0];
    h.state.listeners['gopls-formatting-response'](fmtReq.requestID, JSON.stringify([
        edit(0, 0, 'A'), edit(1, 0, 'B')
    ]));

    assert.strictEqual(model.edits[1][0].text, 'B');
    assert.strictEqual(model.edits[1][1].text, 'A');
    assert.strictEqual(g(h, 'applyingEdit'), false);

    // applyingEdit 抑制了 onDidChangeModelContent，但每次编辑都要显式同步一次。
    const changes = emitted(h, 'gopls-didChange');
    assert.strictEqual(changes.length, 2);
    assert.strictEqual(changes[0].args[0].content, 'FORMATTED');
    assert.strictEqual(changes[0].args[0].version, 6);
    assert.strictEqual(changes[1].args[0].version, 7);

    const save = emitted(h, 'save-file')[0].args[0];
    assert.strictEqual(save.file, key);
    assert.strictEqual(save.content, 'FORMATTED');
    assert.strictEqual(g(h, 'saveInProgress'), false);
});

test('applyWorkspaceEdit 编辑倒序且只同步一次 didChange', function () {
    const h = harness();
    const key = 'a.go';
    const model = addFile(h, key, {value: 'X', applyValue: 'Y'});
    h.state.currentModel = model;
    const changes = {};
    changes[key] = [edit(0, 0, 'A'), edit(1, 0, 'B')];
    run(h, 'applyWorkspaceEdit')({changes: changes});
    assert.strictEqual(model.edits[0][0].text, 'B');
    assert.strictEqual(emitted(h, 'gopls-didChange').length, 1);
    assert.strictEqual(emitted(h, 'gopls-didChange')[0].args[0].content, 'Y');
});

test('格式化响应按 reqID 分发，并发不覆盖', function () {
    const h = harness();
    addFile(h, 'a.go', {});
    const seen = [];
    run(h, 'requestFormatting')('a.go', function (e) { seen.push('a:' + e.length); });
    run(h, 'requestFormatting')('a.go', function (e) { seen.push('b:' + e.length); });
    const reqs = emitted(h, 'gopls-formatting');
    h.state.listeners['gopls-formatting-response'](reqs[0].args[0].requestID, JSON.stringify([edit(0, 0, 'A')]));
    h.state.listeners['gopls-formatting-response'](reqs[1].args[0].requestID, '[]');
    assert.deepStrictEqual(seen, ['a:1', 'b:0']);
    assert.strictEqual(g(h, 'pendingFormatting').size, 0);
});

test('重载文件保持 version 单调递增', function () {
    const h = harness();
    const key = 'a.go';
    addFile(h, key, {value: 'OLD', version: 7});
    h.state.diskContent = 'NEW';
    run(h, 'reloadFileFromDisk(' + JSON.stringify(key) + ', "")');
    const info = g(h, 'files').get(key);
    assert.strictEqual(info.version, 8);
    assert.strictEqual(info.model.getValue(), 'NEW');
    const open = emitted(h, 'gopls-didOpen')[0].args[0];
    assert.strictEqual(open.version, 8);
    assert.strictEqual(open.content, 'NEW');
});

// ---- 超时与索引状态 ----

test('非索引期间 35s 超时清空 completion pending', function () {
    const h = harness();
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.completion.provideCompletionItems(model, {lineNumber: 1, column: 1}, null, null);
    advance(h, 35000);
    assert.strictEqual(p.state, 'resolved');
    assert.strictEqual(p.value.suggestions.length, 0);
    assert.strictEqual(g(h, 'pendingCompletion').size, 0);
});

test('首次索引期间补全超时延长，不再在默认超时丢结果', function () {
    const h = harness();
    h.state.listeners['gopls-indexing'](true);
    const model = addFile(h, 'a.go', {});
    const p = h.state.providers.completion.provideCompletionItems(model, {lineNumber: 1, column: 1}, null, null);
    advance(h, 15000);
    assert.strictEqual(p.state, 'pending');
    assert.strictEqual(g(h, 'pendingCompletion').size, 1);
    // 索引结束后仍能收到迟到的结果（不被提前清空）。
    const reqID = emitted(h, 'gopls-completion')[0].args[0].requestID;
    h.state.listeners['gopls-completion-response'](reqID, JSON.stringify([{label: 'LATE', kind: 1}]));
    assert.strictEqual(p.value.suggestions[0].label, 'LATE');
});

test('gopls 状态条区分启动 / 首次索引 / 恢复', function () {
    const h = harness();
    const bar = h.els['gopls-status-bar'];
    h.state.listeners['gopls-status'](false, 'gopls 正在启动');
    assert.ok(bar.classList.contains('show'));
    assert.ok(bar.textContent.indexOf('正在启动') !== -1);
    h.state.listeners['gopls-status'](true, '');
    assert.ok(!bar.classList.contains('show'));
    h.state.listeners['gopls-indexing'](true);
    assert.ok(bar.classList.contains('show'));
    assert.ok(bar.textContent.indexOf('索引') !== -1);
    h.state.listeners['gopls-indexing'](false);
    assert.ok(!bar.classList.contains('show'));
});

test('状态查询响应中的 indexing 会同步到前端', function () {
    const h = harness({statusResult: {ok: true, reason: '', indexing: true}});
    assert.strictEqual(g(h, 'plsIndexing'), true);
    assert.ok(h.els['gopls-status-bar'].textContent.indexOf('索引') !== -1);
});

test('保存流程响应丢失时不会永久卡住 saveInProgress', function () {
    const h = harness();
    const key = 'a.go';
    addFile(h, key, {});
    setCurrent(h, key);
    run(h, 'manualSave();');
    assert.strictEqual(g(h, 'saveInProgress'), true);
    advance(h, 40000);
    assert.strictEqual(g(h, 'saveInProgress'), false);
    assert.strictEqual(g(h, 'pendingOrganizeSave').size, 0);
});

// ---- 定义跳转 ----

test('定义跳转：归一化后同文件直接定位，跨文件交给 Go', function () {
    const h = harness();
    const key = 'E:\\app\\a.go';
    const model = addFile(h, key, {});
    h.state.currentModel = model;
    h.state.mouseDown({event: {ctrlKey: true, leftButton: true}, target: {position: {lineNumber: 2, column: 3}}});
    const req = emitted(h, 'gopls-definition')[0].args[0];
    assert.strictEqual(req.line, 1);
    assert.strictEqual(req.column, 2);

    h.state.listeners['gopls-definition-response'](req.requestID, JSON.stringify({
        file: 'e:/app/a.go', range: {start: {line: 5, character: 2}, end: {line: 5, character: 6}}
    }));
    assert.strictEqual(h.state.revealed, 6);
    assert.strictEqual(h.state.position.lineNumber, 6);
    assert.strictEqual(h.state.focused, 1);
    assert.strictEqual(g(h, 'pendingDefinition').size, 0);

    h.state.mouseDown({event: {ctrlKey: true, leftButton: true}, target: {position: {lineNumber: 1, column: 1}}});
    const req2 = emitted(h, 'gopls-definition')[1].args[0];
    h.state.listeners['gopls-definition-response'](req2.requestID, JSON.stringify({
        file: 'E:\\app\\b.go', range: {start: {line: 0, character: 0}, end: {line: 0, character: 1}}
    }));
    assert.strictEqual(emitted(h, 'go-to-definition').length, 1);
});

// ---- hover 定义链接 ----

test('hover 的 pkg.go.dev 链接改写为定义跳转，其它外链保持原样', function () {
    const h = harness();
    const key = 'E:\\app\\a.go';
    const model = addFile(h, key, {});
    h.state.currentModel = model;

    const promise = h.state.providers.hover.provideHover(model, {lineNumber: 6, column: 11}, null);
    const req = emitted(h, 'gopls-hover')[0].args[0];
    assert.strictEqual(req.line, 5);
    assert.strictEqual(req.column, 10);

    const md = '```go\npackage http\n```\n\n' +
        '[`http` on pkg.go.dev](https://pkg.go.dev/net/http)\n\n' +
        '[spec](https://go.dev/ref/spec)';
    h.state.listeners['gopls-hover-response'](req.requestID, JSON.stringify({contents: [{value: md}]}));

    const value = promise.value.contents[0].value;
    assert.ok(value.indexOf('https://pkg.go.dev') === -1);
    assert.ok(value.indexOf('[转到定义](file:///designer-goto-definition?file=') !== -1);
    assert.ok(value.indexOf('[spec](https://go.dev/ref/spec)') !== -1);
});

test('hover 定义链接优先使用 hover range 起点', function () {
    const h = harness();
    const key = 'E:\\app\\a.go';
    const model = addFile(h, key, {});
    h.state.currentModel = model;

    const promise = h.state.providers.hover.provideHover(model, {lineNumber: 9, column: 5}, null);
    const req = emitted(h, 'gopls-hover')[0].args[0];
    h.state.listeners['gopls-hover-response'](req.requestID, JSON.stringify({
        contents: [{value: '[`x` on pkg.go.dev](https://pkg.go.dev/x)'}],
        range: {start: {line: 2, character: 3}, end: {line: 2, character: 4}}
    }));

    const value = promise.value.contents[0].value;
    const link = /\((file:\/\/\/designer-goto-definition\?[^)]+)\)/.exec(value);
    assert.ok(link, value);
    const query = decodeURIComponent(link[1].split('?')[1]);
    assert.strictEqual(h.state.linkOpener.open({scheme: 'file', path: '/designer-goto-definition', query: query}), true);

    const def = emitted(h, 'gopls-definition')[0].args[0];
    assert.strictEqual(def.file, key);
    assert.strictEqual(def.line, 2);
    assert.strictEqual(def.column, 3);
});

test('hover 定义链接复用 gopls-definition，不触发外部打开', function () {
    const h = harness();
    const key = 'E:\\app\\a.go';
    addFile(h, key, {});

    assert.ok(h.state.linkOpener, 'registerLinkOpener 已注册');
    assert.strictEqual(h.state.linkOpener.open({scheme: 'https', path: '/net/http', query: ''}), false);
    assert.strictEqual(h.state.linkOpener.open({
        scheme: 'file', path: '/designer-goto-definition',
        query: 'file=' + encodeURIComponent(key) + '&line=4&column=7'
    }), true);

    const req = emitted(h, 'gopls-definition')[0].args[0];
    assert.strictEqual(req.file, key);
    assert.strictEqual(req.line, 4);
    assert.strictEqual(req.column, 7);
});

test('定义链接路径包含 & = % 时参数不串位', function () {
    const h = harness();
    const key = 'E:\\app\\a&b=c%25.go';
    addFile(h, key, {});

    const uri = run(h, 'definitionLinkUri(' + JSON.stringify(key) + ', 8, 3)');
    const query = decodeURIComponent(uri.split('?')[1]);
    assert.strictEqual(h.state.linkOpener.open({scheme: 'file', path: '/designer-goto-definition', query: query}), true);

    const req = emitted(h, 'gopls-definition')[0].args[0];
    assert.strictEqual(req.file, key);
    assert.strictEqual(req.line, 8);
    assert.strictEqual(req.column, 3);
});

console.log('\n' + passed + ' passed, ' + failed + ' failed');
process.exit(failed === 0 ? 0 : 1);
