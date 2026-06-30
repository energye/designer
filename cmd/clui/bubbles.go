package clui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	bubbleInputMinWidth    = 20
	bubbleInputExtraWidth  = 8
	bubbleListWidth        = 64
	bubbleListMinHeight    = 6
	bubbleListMaxHeight    = 18
	bubbleListChromeHeight = 4
	bubbleProgressMaxWidth = 80
	bubbleProgressPadding  = 8
	bubbleDonePause        = 750 * time.Millisecond
)

const (
	bubbleProgressColorA = "#7D56F4"
	bubbleProgressColorB = "#04B575"
)

type bubbleUI struct {
	opt   Options
	state *state
}

func NewBubble(opt Options, state *state) UI {
	return &bubbleUI{opt: opt, state: state}
}

func (u *bubbleUI) Info(args ...any) {
	finishProgressLine(u.opt.Out, u.state)
	fmt.Fprintln(u.opt.Out, args...)
}

func (u *bubbleUI) Warn(args ...any) {
	finishProgressLine(u.opt.Out, u.state)
	fmt.Fprintln(u.opt.Err, append([]any{"[WARN]"}, args...)...)
}

func (u *bubbleUI) Error(args ...any) {
	finishProgressLine(u.opt.Out, u.state)
	fmt.Fprintln(u.opt.Err, append([]any{"[ERROR]"}, args...)...)
}

func (u *bubbleUI) Success(args ...any) {
	finishProgressLine(u.opt.Out, u.state)
	fmt.Fprintln(u.opt.Out, args...)
}

func (u *bubbleUI) Input(title, def string) (string, error) {
	model := inputModel{title: title, input: textinput.New()}
	model.input.SetValue(def)
	model.input.Focus()
	model.input.CharLimit = 256
	model.input.Width = maxInt(bubbleInputMinWidth, len(def)+bubbleInputExtraWidth)
	result, err := tea.NewProgram(model, tea.WithInput(u.opt.In), tea.WithOutput(u.opt.Out)).Run()
	if err != nil {
		return NewSimple(u.opt, u.state).Input(title, def)
	}
	typed, ok := result.(inputModel)
	if !ok || typed.canceled {
		return "", fmt.Errorf("clui: input canceled")
	}
	value := strings.TrimSpace(typed.input.Value())
	if value == "" {
		value = def
	}
	return value, nil
}

func (u *bubbleUI) Select(title string, items []string, def int) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("clui: select has no items")
	}
	model := newSelectModel(title, newSelectItems(items), def)
	result, err := tea.NewProgram(model, tea.WithInput(u.opt.In), tea.WithOutput(u.opt.Out)).Run()
	if err != nil {
		return NewSimple(u.opt, u.state).Select(title, items, def)
	}
	typed, ok := result.(selectModel)
	if !ok || typed.canceled {
		return -1, fmt.Errorf("clui: selection canceled")
	}
	return typed.selected, nil
}

func newSelectItems(items []string) []list.Item {
	listItems := make([]list.Item, 0, len(items))
	for i, item := range items {
		listItems = append(listItems, selectItem{index: i, text: fmt.Sprintf("%d. %s", i+1, item)})
	}
	return listItems
}

func newSelectModel(title string, items []list.Item, def int) selectModel {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Border(lipgloss.Border{Left: ">"}, false, false, false, true)
	model := selectModel{
		list: list.New(items, delegate, bubbleListWidth, selectListHeight(len(items))),
	}
	model.list.Title = title
	model.list.SetFilteringEnabled(false)
	model.list.SetShowStatusBar(false)
	model.list.SetShowPagination(false)
	model.list.SetShowHelp(false)
	model.list.SetHeight(selectListHeight(len(items)))
	model.list.Select(clampIndex(def, len(items)))
	return model
}

func selectListHeight(items int) int {
	return minInt(bubbleListMaxHeight, maxInt(bubbleListMinHeight, items+bubbleListChromeHeight))
}

func (u *bubbleUI) Confirm(title string, def bool) (bool, error) {
	if u.opt.AssumeYes {
		return true, nil
	}
	items := []string{"Yes", "No"}
	defIndex := 1
	if def {
		defIndex = 0
	}
	idx, err := u.Select(title, items, defIndex)
	if err != nil {
		return false, err
	}
	return idx == 0, nil
}

func (u *bubbleUI) Progress(title string, total int64) Progress {
	model := progressModel{
		title:    title,
		progress: progress.New(progress.WithScaledGradient(bubbleProgressColorA, bubbleProgressColorB)),
	}
	program := tea.NewProgram(model, tea.WithInput(nil), tea.WithOutput(u.opt.Out))
	p := &bubbleProgress{
		program: program,
		total:   total,
		message: title,
		doneCh:  make(chan struct{}),
		state:   startProgress(u.state),
	}
	go func() {
		_, _ = program.Run()
		close(p.doneCh)
	}()
	return p
}

type inputModel struct {
	title    string
	input    textinput.Model
	canceled bool
}

func (m inputModel) Init() tea.Cmd { return textinput.Blink }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyEsc, tea.KeyCtrlC:
			m.canceled = true
			return m, tea.Quit
		}
	}
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m inputModel) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s\n", m.title, m.input.View(), "(enter submit  esc cancel)")
}

type selectItem struct {
	index int
	text  string
}

func (i selectItem) Title() string       { return i.text }
func (i selectItem) Description() string { return "" }
func (i selectItem) FilterValue() string { return i.text }

type selectModel struct {
	list     list.Model
	selected int
	canceled bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyCtrlC:
			m.canceled = true
			return m, tea.Quit
		case tea.KeyEnter:
			if item, ok := m.list.SelectedItem().(selectItem); ok {
				m.selected = item.index
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectModel) View() string {
	if m.canceled {
		return ""
	}
	return "\n" + m.list.View()
}

type progressModel struct {
	title    string
	message  string
	progress progress.Model
	err      error
}

type progressUpdateMsg struct {
	percent float64
	message string
	done    bool
}

func (m progressModel) Init() tea.Cmd { return nil }

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = minInt(bubbleProgressMaxWidth, maxInt(1, msg.Width-bubbleProgressPadding))
		return m, nil
	case progressUpdateMsg:
		m.message = msg.message
		cmd := m.progress.SetPercent(msg.percent)
		if msg.done {
			return m, tea.Batch(cmd, tea.Sequence(finalPause(), tea.Quit))
		}
		return m, cmd
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m progressModel) View() string {
	if m.err != nil {
		return "Error: " + m.err.Error() + "\n"
	}
	message := m.message
	if message == "" {
		message = m.title
	}
	return "  " + message + "\n  " + m.progress.View() + "\n"
}

func finalPause() tea.Cmd {
	return tea.Tick(bubbleDonePause, func(time.Time) tea.Msg { return nil })
}

func clampIndex(index, length int) int {
	if length <= 0 {
		return 0
	}
	if index < 0 {
		return 0
	}
	if index >= length {
		return length - 1
	}
	return index
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
