package clui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TestInputModelEnter(t *testing.T) {
	input := textinput.New()
	input.SetValue("demo")
	model := inputModel{title: "Name", input: input}
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(inputModel)
	if cmd == nil {
		t.Fatal("expected enter to return quit command")
	}
	if model.input.Value() != "demo" || model.canceled {
		t.Fatalf("input model value=%q canceled=%v", model.input.Value(), model.canceled)
	}
}

func TestSelectModelEnter(t *testing.T) {
	items := []list.Item{selectItem{index: 0, text: "109"}, selectItem{index: 1, text: "127"}}
	model := newSelectModel("Version", items, 0)
	model.list.Select(1)
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(selectModel)
	if cmd == nil {
		t.Fatal("expected enter to return quit command")
	}
	if model.selected != 1 || model.canceled {
		t.Fatalf("select model selected=%d canceled=%v", model.selected, model.canceled)
	}
}

func TestSelectModelShowsSmallLists(t *testing.T) {
	items := newSelectItems([]string{"LCL", "WV", "CEF"})
	model := newSelectModel("Select UI", items, 0)
	if model.list.Paginator.PerPage < len(items) {
		t.Fatalf("select per page=%d want at least %d", model.list.Paginator.PerPage, len(items))
	}
	view := model.View()
	for _, want := range []string{"1. LCL", "2. WV", "3. CEF", ">"} {
		if !strings.Contains(view, want) {
			t.Fatalf("select view missing %q: %q", want, view)
		}
	}
}

func TestProgressModelUpdate(t *testing.T) {
	model := progressModel{title: "Downloading", progress: progress.New()}
	updated, cmd := model.Update(progressUpdateMsg{percent: 0.5, message: "Downloading"})
	model = updated.(progressModel)
	if cmd == nil {
		t.Fatal("expected progress update command")
	}
	if model.message != "Downloading" {
		t.Fatalf("progress model message=%q", model.message)
	}
	if model.progress.Percent() != 0.5 {
		t.Fatalf("progress model percent=%v", model.progress.Percent())
	}
}

func TestProgressModelViewUsesTwoLines(t *testing.T) {
	model := progressModel{title: "Downloading", message: "Downloading", progress: progress.New()}
	view := strings.TrimSuffix(model.View(), "\n")
	lines := strings.Split(view, "\n")
	if len(lines) != 2 {
		t.Fatalf("progress view lines=%d view=%q", len(lines), view)
	}
	if !strings.Contains(lines[0], "Downloading") {
		t.Fatalf("progress title line missing message: %q", lines[0])
	}
	if !strings.Contains(lines[1], "%") {
		t.Fatalf("progress bar line missing percentage: %q", lines[1])
	}
}
