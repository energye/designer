# 编辑器接口定义评审与 MVP 草案

## 结论

原定义的分层方向是合理的：`Adapter` 屏蔽具体编辑器实现，`Service` 管理文件和状态，`Bridge` 负责设计器模型同步，`GenerationService` 负责模型生成代码。这个方向适合后续扩展。

但初期实现常规编辑功能时，原定义偏复杂，主要问题是：

- `CodeEditorAdapter` 和 `CodeEditorService` 都暴露了内容、格式化、诊断、命令等能力，容易出现状态归属不清。
- `CodeDesignerBridge` 和 `CodeGenerationService` 属于设计器同步和代码生成能力，不应该阻塞第一阶段编辑器落地。
- `Completion`、`Hover`、`GotoDefinition`、`Rename`、`Diff`、`Commands` 等高级能力需要更完整协议，MVP 阶段可以先不定义。
- `Position` / `Range` 没有说明行列从 0 还是 1 开始，也没有说明 `End` 是否包含，后续很容易和 Monaco、LSP、Go 后端产生偏差。
- 保存、版本号、dirty 状态需要有唯一事实来源。建议由 `CodeEditorService` 管，`Adapter` 只负责和具体编辑器通信。

第一阶段建议只实现：

- 挂载 / 销毁编辑器
- 打开 / 关闭 / 切换文件
- 读取 / 写入内容
- 管理 active file、dirty、readonly、version
- 保存 / 全部保存
- 基础选区、定位、聚焦、布局
- 基础诊断显示
- 可选格式化

`Bridge`、代码生成、反解析、智能提示、重命名、跳转、diff、命令注册建议放到第二阶段。

## 推荐依赖关系

```text
CodeEditorService
        |
        v
CodeEditorAdapter
```

第一阶段只保留这两层：

- `CodeEditorAdapter`：具体编辑器驱动，例如 Monaco WebView、普通 textarea、远程编辑器。
- `CodeEditorService`：文件状态中心，负责 files、content、dirty、save、diagnostics。

后续再增加：

```text
CodeDesignerBridge -> CodeEditorService
CodeDesignerBridge -> CodeGenerationService
```

## MVP 接口定义

```go
package codeeditor

import "context"

// Language 表示编辑器可识别的语言类型。
// 初期只放常用语言；后续可以补充 sql、yaml、markdown 等。
type Language string

const (
	LanguageJavaScript Language = "javascript"
	LanguageTypeScript Language = "typescript"
	LanguageJSON       Language = "json"
	LanguageCSS        Language = "css"
	LanguageHTML       Language = "html"
	LanguageVue        Language = "vue"
	LanguageGo         Language = "go"
	LanguageText       Language = "text"
)

// Position 表示编辑器里的光标位置。
// 约定：Line 和 Column 都从 1 开始，和 Monaco 的默认坐标保持一致。
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Range 表示一段代码范围。
// 约定：Start 包含，End 不包含，便于表达替换范围。
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// File 是编辑器里的文件模型。
// Content 是当前内容快照；Dirty、Version 由 CodeEditorService 维护。
type File struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Path     string         `json:"path,omitempty"`
	Language Language       `json:"language"`
	Content  string         `json:"content"`
	Readonly bool           `json:"readonly,omitempty"`
	Dirty    bool           `json:"dirty,omitempty"`
	Version  int64          `json:"version,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
}

// FilePatch 用于局部更新文件信息。
// 使用指针是为了区分“未修改”和“修改为空值”。
type FilePatch struct {
	Name     *string        `json:"name,omitempty"`
	Path     *string        `json:"path,omitempty"`
	Language *Language      `json:"language,omitempty"`
	Content  *string        `json:"content,omitempty"`
	Readonly *bool          `json:"readonly,omitempty"`
	Dirty    *bool          `json:"dirty,omitempty"`
	Version  *int64         `json:"version,omitempty"`
	Meta     map[string]any `json:"meta,omitempty"`
}

type DiagnosticSeverity string

const (
	DiagnosticError   DiagnosticSeverity = "error"
	DiagnosticWarning DiagnosticSeverity = "warning"
	DiagnosticInfo    DiagnosticSeverity = "info"
	DiagnosticHint    DiagnosticSeverity = "hint"
)

// Diagnostic 表示文件里的错误、警告或提示。
// Source 可写 eslint、typescript、designer、generator 等来源。
type Diagnostic struct {
	FileID   string             `json:"fileId"`
	Message  string             `json:"message"`
	Severity DiagnosticSeverity `json:"severity"`
	Range    Range              `json:"range"`
	Source   string             `json:"source,omitempty"`
	Code     string             `json:"code,omitempty"`
}

// EditorOptions 是编辑器展示选项。
// 这些配置只影响编辑体验，不作为业务状态来源。
type EditorOptions struct {
	Theme       string `json:"theme,omitempty"`
	Readonly    bool   `json:"readonly,omitempty"`
	FontSize    int    `json:"fontSize,omitempty"`
	TabSize     int    `json:"tabSize,omitempty"`
	WordWrap    bool   `json:"wordWrap,omitempty"`
	Minimap     bool   `json:"minimap,omitempty"`
	LineNumbers bool   `json:"lineNumbers,omitempty"`
}

// EditorCapabilities 用于描述具体编辑器驱动已实现的能力。
// Service 调用可选能力前应先判断，或者由 Adapter 返回 ErrUnsupported。
type EditorCapabilities struct {
	MultiFile   bool `json:"multiFile"`
	Diagnostics bool `json:"diagnostics"`
	Formatting  bool `json:"formatting"`
	Selection   bool `json:"selection"`
	Theme       bool `json:"theme"`
}

// MountRequest 是挂载编辑器时的初始参数。
type MountRequest struct {
	ContainerID  string        `json:"containerId,omitempty"`
	Files        []File        `json:"files,omitempty"`
	ActiveFileID string        `json:"activeFileId,omitempty"`
	Options      EditorOptions `json:"options,omitempty"`
}

// CodeEditorAdapter 是具体编辑器实现的驱动接口。
// 它只负责和 Monaco、WebView、textarea 等实际编辑器通信，不保存业务状态。
type CodeEditorAdapter interface {
	// ID 返回驱动唯一标识，例如 "monaco-webview"。
	ID() string

	// Name 返回展示名称，例如 "Monaco Editor"。
	Name() string

	// Capabilities 返回当前驱动支持的能力。
	Capabilities() EditorCapabilities

	// Mount 在指定容器里创建编辑器实例，并加载初始文件。
	Mount(ctx context.Context, req MountRequest) error

	// Dispose 销毁编辑器实例，释放 WebView、事件监听等资源。
	Dispose(ctx context.Context) error

	// OpenFile 在编辑器中打开文件。
	// 如果不支持多文件，可用新文件替换当前文件。
	OpenFile(ctx context.Context, file File) error

	// CloseFile 关闭指定文件。
	CloseFile(ctx context.Context, fileID string) error

	// SwitchFile 切换当前激活文件。
	SwitchFile(ctx context.Context, fileID string) error

	// ActiveFileID 返回编辑器当前激活的文件 ID。
	ActiveFileID(ctx context.Context) (string, error)

	// GetValue 从编辑器读取指定文件的当前内容。
	GetValue(ctx context.Context, fileID string) (string, error)

	// SetValue 把内容写入编辑器。
	// 业务层应避免把用户正在输入的内容反复覆盖。
	SetValue(ctx context.Context, fileID string, value string) error

	// SetOptions 更新编辑器显示选项。
	SetOptions(ctx context.Context, options EditorOptions) error

	// Focus 让编辑器获得焦点。
	Focus(ctx context.Context) error

	// Layout 通知编辑器重新计算尺寸。
	Layout(ctx context.Context) error

	// GetSelection 获取当前选区；不支持选区时可返回 nil。
	GetSelection(ctx context.Context) (*Range, error)

	// SetSelection 设置当前选区。
	SetSelection(ctx context.Context, r Range) error

	// RevealRange 滚动到指定代码范围。
	RevealRange(ctx context.Context, r Range) error

	// Format 格式化指定文件。
	// 这是可选能力；不支持时返回 ErrUnsupported。
	Format(ctx context.Context, fileID string) error

	// SetDiagnostics 在编辑器中显示诊断信息。
	SetDiagnostics(ctx context.Context, diagnostics []Diagnostic) error

	// ClearDiagnostics 清除指定文件的诊断信息。
	ClearDiagnostics(ctx context.Context, fileID string) error
}

// ServiceOptions 是 CodeEditorService 的行为配置。
type ServiceOptions struct {
	AutoSave         bool `json:"autoSave,omitempty"`
	AutoSaveDelayMS  int  `json:"autoSaveDelayMs,omitempty"`
	ValidateOnChange bool `json:"validateOnChange,omitempty"`
	FormatOnSave     bool `json:"formatOnSave,omitempty"`
}

type ChangeSource string

const (
	ChangeSourceUser   ChangeSource = "user"
	ChangeSourceAPI    ChangeSource = "api"
	ChangeSourceFormat ChangeSource = "format"
)

// SaveResult 表示保存结果。
// Value 是保存时的最终内容，可用于写入磁盘、数据库或远端接口。
type SaveResult struct {
	FileID  string `json:"fileId"`
	OK      bool   `json:"ok"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message,omitempty"`
}

// CodeEditorService 是编辑器文件状态中心。
// 它依赖 Adapter，但不绑定具体编辑器实现。
type CodeEditorService interface {
	// Init 初始化服务并挂载底层编辑器。
	Init(ctx context.Context, adapter CodeEditorAdapter, options ServiceOptions) error

	// Dispose 释放服务和底层编辑器资源。
	Dispose(ctx context.Context) error

	// Adapter 返回当前使用的编辑器驱动。
	Adapter() CodeEditorAdapter

	// SetFiles 批量设置文件，并指定激活文件。
	SetFiles(ctx context.Context, files []File, activeFileID string) error

	// Files 返回当前文件列表。
	Files(ctx context.Context) ([]File, error)

	// File 返回指定文件。
	File(ctx context.Context, fileID string) (*File, error)

	// AddFile 添加文件并打开到编辑器。
	AddFile(ctx context.Context, file File) error

	// UpdateFile 更新文件元信息或内容。
	UpdateFile(ctx context.Context, fileID string, patch FilePatch) error

	// RemoveFile 移除文件并关闭编辑器里的文件。
	RemoveFile(ctx context.Context, fileID string) error

	// ActiveFile 返回当前激活文件。
	ActiveFile(ctx context.Context) (*File, error)

	// SetActiveFile 切换当前激活文件。
	SetActiveFile(ctx context.Context, fileID string) error

	// Content 读取文件内容。
	// Service 应优先从 Adapter 读取最新编辑内容，再同步到内部状态。
	Content(ctx context.Context, fileID string) (string, error)

	// SetContent 设置文件内容。
	// source 用于避免编辑器变更和程序同步之间形成循环触发。
	SetContent(ctx context.Context, fileID string, content string, source ChangeSource) error

	// IsDirty 判断文件是否有未保存变更。
	IsDirty(ctx context.Context, fileID string) (bool, error)

	// MarkSaved 把文件标记为已保存。
	MarkSaved(ctx context.Context, fileID string) error

	// Save 保存指定文件。
	Save(ctx context.Context, fileID string) (SaveResult, error)

	// SaveAll 保存所有 dirty 文件。
	SaveAll(ctx context.Context) ([]SaveResult, error)

	// Format 格式化指定文件。
	Format(ctx context.Context, fileID string) error

	// SetDiagnostics 设置诊断信息，并同步给 Adapter 显示。
	SetDiagnostics(ctx context.Context, diagnostics []Diagnostic) error

	// Diagnostics 返回指定文件的诊断信息。
	Diagnostics(ctx context.Context, fileID string) ([]Diagnostic, error)

	// ClearDiagnostics 清除指定文件的诊断信息。
	ClearDiagnostics(ctx context.Context, fileID string) error
}
```

## 示例注释

### 打开并编辑文件

```go
func ExampleOpenFile(ctx context.Context, svc CodeEditorService) error {
	file := File{
		ID:       "home.vue",
		Name:     "Home.vue",
		Path:     "src/views/Home.vue",
		Language: LanguageVue,
		Content:  "<template><div>Hello</div></template>",
	}

	// 添加文件时，Service 负责保存文件状态，Adapter 负责把文件显示到编辑器。
	if err := svc.AddFile(ctx, file); err != nil {
		return err
	}

	// 设置为当前文件后，Adapter 需要切换编辑器 tab 或替换当前内容。
	return svc.SetActiveFile(ctx, file.ID)
}
```

### 用户输入后同步内容

```go
func ExampleUserChange(ctx context.Context, svc CodeEditorService, fileID string, value string) error {
	// 用户输入产生的内容变更应标记为 dirty。
	// Service 可以在这里递增 Version，并根据配置触发自动保存或诊断。
	return svc.SetContent(ctx, fileID, value, ChangeSourceUser)
}
```

### 保存文件

```go
func ExampleSave(ctx context.Context, svc CodeEditorService, fileID string) error {
	result, err := svc.Save(ctx, fileID)
	if err != nil {
		return err
	}
	if !result.OK {
		return fmt.Errorf("save failed: %s", result.Message)
	}

	// 保存成功后，Service 应将文件标记为非 dirty。
	return svc.MarkSaved(ctx, fileID)
}
```

### 显示诊断信息

```go
func ExampleDiagnostics(ctx context.Context, svc CodeEditorService) error {
	diagnostics := []Diagnostic{
		{
			FileID:   "home.vue",
			Message:  "缺少结束标签",
			Severity: DiagnosticError,
			Range: Range{
				Start: Position{Line: 1, Column: 11},
				End:   Position{Line: 1, Column: 16},
			},
			Source: "designer",
		},
	}

	// Service 保存诊断状态，并调用 Adapter 在编辑器中展示红线、标记等。
	return svc.SetDiagnostics(ctx, diagnostics)
}
```

## 暂缓到第二阶段的接口

以下能力建议先不要进入 MVP 接口，等常规编辑闭环跑通后再加：

- `CodeDesignerBridge`
- `CodeGenerationService`
- `RegisterCommand` / `ExecuteCommand`
- `Completion` / `Hover` / `GotoDefinition` / `Rename`
- `Diff`
- 代码反解析为设计器模型
- 复杂增量 patch 协议

这样第一阶段可以先把编辑器作为“可靠的多文件文本编辑器”跑起来，后面再扩展设计器同步和智能能力。
