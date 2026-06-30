package clui

import (
	"bytes"
	"strings"
	"testing"
)

func TestSimpleInputDefault(t *testing.T) {
	ui := NewSimple(Options{In: strings.NewReader("\n"), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})
	got, err := ui.Input("Name", "demo")
	if err != nil {
		t.Fatal(err)
	}
	if got != "demo" {
		t.Fatalf("Input got %q want demo", got)
	}
}

func TestSimpleSelect(t *testing.T) {
	ui := NewSimple(Options{In: strings.NewReader("2\n"), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})
	got, err := ui.Select("Version", []string{"109", "127"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("Select got %d want 1", got)
	}
}

func TestSimpleSequentialInputAndSelect(t *testing.T) {
	ui := NewSimple(Options{In: strings.NewReader("demo\n2\n"), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})
	name, err := ui.Input("Name", "")
	if err != nil {
		t.Fatal(err)
	}
	if name != "demo" {
		t.Fatalf("Input got %q want demo", name)
	}
	got, err := ui.Select("Version", []string{"109", "127"}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Fatalf("Select got %d want 1", got)
	}
}

func TestSimpleSelectCancel(t *testing.T) {
	ui := NewSimple(Options{In: strings.NewReader("0\n"), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})
	if _, err := ui.Select("Version", []string{"109", "127"}, -1); err == nil {
		t.Fatal("expected cancel error")
	}
}

func TestSimpleConfirmDefault(t *testing.T) {
	ui := NewSimple(Options{In: strings.NewReader("\n"), Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})
	got, err := ui.Confirm("Overwrite?", true)
	if err != nil {
		t.Fatal(err)
	}
	if !got {
		t.Fatal("Confirm got false want true")
	}
}

func TestLineProgress(t *testing.T) {
	var out bytes.Buffer
	progress := NewSimple(Options{Out: &out}).Progress("Downloading", 100)
	progress.Set(50)
	progress.Finish()
	got := out.String()
	if !strings.Contains(got, "50%") || !strings.Contains(got, "100%") {
		t.Fatalf("progress output missing percentages: %q", got)
	}
	if !strings.Contains(got, "Downloading\n") {
		t.Fatalf("progress output missing title line: %q", got)
	}
}

func TestProgressSink(t *testing.T) {
	var out bytes.Buffer
	sink := NewProgressSink(NewSimple(Options{Out: &out}))
	sink.Update("Preparing", 0, 0)
	sink.Update("Downloading", 50, 100)
	sink.Update("Extracting", 1, 2)
	sink.Finish()
	got := out.String()
	for _, want := range []string{"Preparing", "Downloading", "Extracting", "50%", "100%"} {
		if !strings.Contains(got, want) {
			t.Fatalf("progress sink output missing %q: %q", want, got)
		}
	}
}

func TestBubbleProgress(t *testing.T) {
	var out bytes.Buffer
	progress := NewBubble(Options{Out: &out}, &state{}).Progress("Downloading", 100)
	progress.Set(50)
	progress.Finish()
	got := out.String()
	if !strings.Contains(got, "Downloading") || !strings.Contains(got, "100%") {
		t.Fatalf("progress output missing expected content: %q", got)
	}
	if strings.Contains(got, "[##########----------]") {
		t.Fatalf("bubble progress should not use simple hash bar: %q", got)
	}
}
