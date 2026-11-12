// 前端编辑器逻辑的 Node 回归测试（不依赖浏览器/LCL）。
// 运行: node resources/editor/monaco/js/editor_logic_test.js
// 用 vm 加载与 index.html 相同的脚本，注入假的 DOM / ipc / Monaco 实例。

'use strict';

const fs = require('fs');
const path = require('path');
const vm = require('vm');
const assert = require('assert');

const JS_DIR = __dirname;
let passed = 0;
let failed = 0;

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

function makeModel(value) {
    return {
        _value: value,
        getValue() { return this._value; },
        setValue(v) { this._value = v; },
        dispose() {},
        applyEdits() {},
        getFullModelRange() { return {}; }
    };
}

function harness() {
    const els = {};
    ['file-path-bar', 'gopls-status-bar', 'conflict-bar', 'conflict-message',
        'conflict-reload', 'conflict-keep', 'toast', 'container'].forEach(function (id) {
        els[id] = makeEl();
    });

    const state = {
        listeners: {},
        emitted: [],
        diskContent: 'DISK CONTENT',
        saveFileResult: 'ok',
        reloadAsObject: false,
        reloadReadOnly: false,
        reloadRaw: null,
        goplsStatusResult: {ok: false, reason: '初始未启动'}
    };
    const ipc = {
        on(name, fn) { state.listeners[name] = fn; },
        emit(name, args, cb) {
            state.emitted.push({name: name, args: args});
            if (cb && name === 'save-file') cb(state.saveFileResult);
            if (cb && name === 'reload-file-request') {
                if (state.reloadRaw !== null) {
                    cb(state.reloadRaw);
                    return;
                }
                const payload = {
                    file: Array.isArray(args) ? args[0] : '',
                    content: state.diskContent, modTime: 1, language: 'go',
                    readOnly: state.reloadReadOnly
                };
                cb(state.reloadAsObject ? payload : JSON.stringify(payload));
            }
            if (cb && name === 'open-file-request') {
                cb(JSON.stringify({content: state.diskContent, modTime: 1,
                    language: 'go', readOnly: false}));
            }
            if (cb && name === 'gopls-status-request') cb(state.goplsStatusResult);
        }
    };

    const files = new Map();
    const fakeEditor = {
        setModel() {}, getModel() { return null; },
        saveViewState() { return {}; }, restoreViewState() {},
        updateOptions() {}, revealLineInCenter() {}, setPosition() {}, focus() {},
        addCommand() {}, onDidChangeModelContent() {}, onMouseDown() {}
    };
    const fakeMonaco = {
        editor: {
            createModel(value) { return makeModel(value); }
        },
        Uri: {parse() { return ''; }}
    };

    const ctx = {
        console: console,
        document: {getElementById(id) { return els[id] || null; }},
        ipc: ipc,
        // 不真正调度定时器，避免测试进程等待且便于断言 toast 状态
        setTimeout() { return 0; },
        clearTimeout() {},
        require: Object.assign(function () {}, {config() {}}),
        JSON: JSON, Math: Math, Date: Date, Set: Set, Map: Map,
        __h: {files: files, fakeEditor: fakeEditor, fakeMonaco: fakeMonaco}
    };
    vm.createContext(ctx);

    // 与 index.html 相同的加载顺序；pls_bridge.js 的顶层注册不需要参与这些断言
    ['utils.js', 'editor.js', 'file_manager.js', 'ipc_handlers.js'].forEach(function (f) {
        vm.runInContext(fs.readFileSync(path.join(JS_DIR, f), 'utf8'), ctx, {filename: f});
    });
    // editor.js 顶层的 let 绑定会遮蔽注入对象，这里重新绑定
    vm.runInContext([
        'editor = __h.fakeEditor;',
        'monacoRef = __h.fakeMonaco;',
        'files = __h.files;',
        'currentFilePath = "";',
        'applyingEdit = false;'
    ].join('\n'), ctx);

    return {ctx: ctx, els: els, state: state, files: files, makeModel: makeModel};
}

function addFile(h, name, opts) {
    const model = h.makeModel(opts.value || 'DISK CONTENT');
    h.files.set(name, {
        model: model,
        modTime: 1,
        isDirty: !!opts.isDirty,
        viewState: null,
        readOnly: !!opts.readOnly,
        version: opts.version || 1
    });
    return model;
}

function test(name, fn) {
    try {
        fn();
        console.log('PASS  ' + name);
        passed++;
    } catch (e) {
        console.log('FAIL  ' + name + ' :: ' + (e && e.message));
        failed++;
    }
}

// ---- 外部变更冲突 ----

test('dirty + 外部变更: 不自动重载', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'LOCAL UNSAVED', isDirty: true});
    h.state.listeners['file-conflict-detected']('a.go');
    assert.strictEqual(m.getValue(), 'LOCAL UNSAVED');
    assert.strictEqual(h.files.get('a.go').isDirty, true);
    assert.ok(h.els['conflict-bar'].classList.contains('show'));
});

test('冲突选择「保留我的修改」: 保留内容与脏状态', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'LOCAL UNSAVED', isDirty: true});
    h.state.listeners['file-conflict-detected']('a.go');
    h.els['conflict-keep']._handlers.click();
    assert.strictEqual(m.getValue(), 'LOCAL UNSAVED');
    assert.strictEqual(h.files.get('a.go').isDirty, true);
    assert.ok(!h.els['conflict-bar'].classList.contains('show'));
});

test('冲突选择「重新加载」: 读磁盘并清脏状态', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'LOCAL UNSAVED', isDirty: true});
    h.state.diskContent = 'NEW FROM DISK';
    h.state.listeners['file-conflict-detected']('a.go');
    h.els['conflict-reload']._handlers.click();
    assert.strictEqual(m.getValue(), 'NEW FROM DISK');
    assert.strictEqual(h.files.get('a.go').isDirty, false);
    assert.ok(h.state.emitted.some(function (e) { return e.name === 'gopls-didOpen'; }));
});

test('file-changed-externally 且无脏状态: 自动重载', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'OLD', isDirty: false});
    h.state.diskContent = 'NEW';
    h.state.listeners['file-changed-externally']('a.go');
    assert.strictEqual(m.getValue(), 'NEW');
});

test('file-changed-externally 但状态不同步(仍 dirty): 弹冲突而非静默跳过', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'LOCAL', isDirty: true});
    h.state.diskContent = 'NEW';
    h.state.listeners['file-changed-externally']('a.go');
    assert.strictEqual(m.getValue(), 'LOCAL');
    assert.ok(h.els['conflict-bar'].classList.contains('show'));
});

// ---- 保存门禁联动 ----

test('save-all-request: 只保存 dirty 且可写的文件，内容取自当前 model', function () {
    const h = harness();
    addFile(h, 'a.go', {value: 'A-CONTENT', isDirty: true, readOnly: false});
    addFile(h, 'b.go', {value: 'B-CONTENT', isDirty: false, readOnly: false});
    addFile(h, 'ro.go', {value: 'RO-CONTENT', isDirty: true, readOnly: true});
    h.state.listeners['save-all-request'](1);
    const saves = h.state.emitted.filter(function (e) { return e.name === 'save-file'; });
    assert.strictEqual(saves.length, 1);
    assert.strictEqual(saves[0].args[0].file, 'a.go');
    assert.strictEqual(saves[0].args[0].content, 'A-CONTENT');
    assert.strictEqual(saves[0].args[0].readOnly, false);
    assert.strictEqual(h.files.get('a.go').isDirty, false);
});

test('save-all-request 保存失败: 保留脏状态并给出提示', function () {
    const h = harness();
    addFile(h, 'a.go', {value: 'A', isDirty: true});
    h.state.saveFileResult = 'error: file is read-only';
    h.state.listeners['save-all-request'](1);
    assert.strictEqual(h.files.get('a.go').isDirty, true);
    assert.ok(h.els.toast.textContent.indexOf('自动保存失败') !== -1);
});

test('save-all-request 成功后清除该文件挂起的冲突提示', function () {
    const h = harness();
    addFile(h, 'a.go', {value: 'A', isDirty: true});
    h.state.listeners['file-conflict-detected']('a.go');
    assert.ok(h.els['conflict-bar'].classList.contains('show'));
    h.state.listeners['save-all-request'](1);
    assert.ok(!h.els['conflict-bar'].classList.contains('show'));
});

test('关闭文件时移除该文件挂起的冲突提示', function () {
    const h = harness();
    addFile(h, 'a.go', {value: 'A', isDirty: true});
    h.state.listeners['file-conflict-detected']('a.go');
    h.state.listeners['close-file']('a.go');
    assert.ok(!h.els['conflict-bar'].classList.contains('show'));
});

// ---- 重载响应兼容（真机发现 WebView2 回传对象而非 JSON 字符串） ----

test('reload-file-request 返回对象时仍能重载且 version 单调递增', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'LOCAL UNSAVED', isDirty: true, version: 7});
    h.state.reloadAsObject = true;
    h.state.diskContent = 'NEW FROM DISK';
    vm.runInContext('reloadFileFromDisk("a.go", "")', h.ctx);
    assert.strictEqual(m.getValue(), 'NEW FROM DISK');
    assert.strictEqual(h.files.get('a.go').isDirty, false);
    assert.strictEqual(h.files.get('a.go').version, 8);
    assert.ok(h.state.emitted.some(function (e) {
        return e.name === 'gopls-didOpen' && e.args[0].version === 8;
    }));
});

test('reload-file-request 返回非法 JSON 时不改动 model', function () {
    const h = harness();
    const m = addFile(h, 'a.go', {value: 'LOCAL', isDirty: true});
    h.state.reloadRaw = 'not-json';
    vm.runInContext('reloadFileFromDisk("a.go", "")', h.ctx);
    assert.strictEqual(m.getValue(), 'LOCAL');
    assert.strictEqual(h.files.get('a.go').isDirty, true);
    assert.strictEqual(vm.runInContext('applyingEdit', h.ctx), false);
});

test('重载后同步服务端只读标记与路径栏', function () {
    const h = harness();
    addFile(h, 'a.go', {value: 'OLD'});
    h.state.reloadAsObject = true;
    h.state.reloadReadOnly = true;
    h.state.diskContent = 'NEW';
    vm.runInContext('reloadFileFromDisk("a.go", "")', h.ctx);
    assert.strictEqual(h.files.get('a.go').readOnly, true);
    assert.ok(h.els['file-path-bar'].textContent.indexOf('[Read Only]') !== -1);
});

// ---- 只读保护 ----

test('autoSave 对只读文件直接返回', function () {
    const h = harness();
    addFile(h, 'ro.go', {value: 'RO', isDirty: false, readOnly: true});
    vm.runInContext('currentFilePath = "ro.go";', h.ctx);
    h.state.emitted.length = 0;
    vm.runInContext('autoSave();', h.ctx);
    assert.strictEqual(h.state.emitted.length, 0);
});

test('manualSave 对只读文件直接返回', function () {
    const h = harness();
    addFile(h, 'ro.go', {value: 'RO', isDirty: false, readOnly: true});
    vm.runInContext('currentFilePath = "ro.go";', h.ctx);
    h.state.emitted.length = 0;
    vm.runInContext('manualSave();', h.ctx);
    assert.strictEqual(h.state.emitted.length, 0);
});

// ---- gopls 状态提示 ----

test('gopls 不可用时提示具体原因与修复方式', function () {
    const h = harness();
    h.state.listeners['gopls-status'](false, 'gopls start: exec: "gopls": not found');
    const bar = h.els['gopls-status-bar'];
    assert.ok(bar.classList.contains('show'));
    assert.ok(bar.textContent.indexOf('not found') !== -1);
    assert.ok(bar.textContent.indexOf('DESIGNER_GOPLS_PATH') !== -1);
});

test('gopls 恢复后隐藏提示', function () {
    const h = harness();
    h.state.listeners['gopls-status'](false, '启动失败');
    h.state.listeners['gopls-status'](true, '');
    assert.ok(!h.els['gopls-status-bar'].classList.contains('show'));
    assert.strictEqual(h.els['gopls-status-bar'].textContent, '');
});

test('页面脚本就绪时主动同步一次 gopls 状态', function () {
    const h = harness();
    assert.strictEqual(h.state.emitted.filter(function (e) {
        return e.name === 'gopls-status-request';
    }).length, 1);
    assert.ok(h.els['gopls-status-bar'].textContent.indexOf('初始未启动') !== -1);
});

console.log('\n' + passed + ' passed, ' + failed + ' failed');
process.exit(failed === 0 ? 0 : 1);
