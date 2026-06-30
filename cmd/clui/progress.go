package clui

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func startProgress(state *state) *progressState {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.progress = &progressState{active: true}
	return state.progress
}

func finishProgressLine(out io.Writer, state *state) {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.progress != nil && state.progress.active {
		fmt.Fprintln(out)
		state.progress.active = false
	}
}

type lineProgress struct {
	out       io.Writer
	title     string
	message   string
	total     int64
	current   int64
	started   bool
	ansi      bool
	lastTitle string
	done      bool
	state     *progressState
}

func (p *lineProgress) Set(current int64) {
	p.Update(current, "")
}

func (p *lineProgress) Message(message string) {
	p.Update(-1, message)
}

func (p *lineProgress) Update(current int64, message string) {
	if current >= 0 {
		p.current = current
	}
	if message != "" {
		p.message = message
	}
	p.render(false)
}

func (p *lineProgress) Finish() {
	p.render(true)
	if !p.done {
		fmt.Fprintln(p.out)
		p.deactivate()
		p.done = true
	}
}

func (p *lineProgress) render(finish bool) {
	if p.done {
		return
	}
	title := p.progressTitle()
	p.start(title)
	current := p.current
	if current < 0 {
		current = 0
	}
	if p.total <= 0 {
		p.writeProgress(title, title, 80)
		return
	}
	if current > p.total || finish {
		current = p.total
	}
	percent := current * 100 / p.total
	width := int(percent / 5)
	bar := strings.Repeat("#", width) + strings.Repeat("-", 20-width)
	text := fmt.Sprintf("[%s] %d%%", bar, percent)
	if !p.ansi && p.started && p.lastTitle != title {
		text += " " + title
	}
	p.writeProgress(title, text, 100)
}

func (p *lineProgress) start(title string) {
	if p.started {
		return
	}
	fmt.Fprintln(p.out, title)
	p.started = true
	p.lastTitle = title
}

func (p *lineProgress) progressTitle() string {
	if p.message != "" {
		return p.message
	}
	return p.title
}

func (p *lineProgress) writeProgress(title, text string, width int) {
	if p.ansi && p.lastTitle != title {
		fmt.Fprintf(p.out, "\r\033[1A\033[2K\r%s\r\n\033[2K\r", trimLine(title, width))
		p.lastTitle = title
	}
	p.writeLine(text, width)
}

func (p *lineProgress) writeLine(text string, width int) {
	fmt.Fprintf(p.out, "\r%s\r%s", strings.Repeat(" ", width), trimLine(text, width))
}

func (p *lineProgress) deactivate() {
	if p.state != nil {
		p.state.active = false
	}
}

func trimLine(value string, max int) string {
	if len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

func useANSIProgress(out io.Writer) bool {
	if out != os.Stdout && out != os.Stderr {
		return false
	}
	file, ok := out.(*os.File)
	if !ok {
		return false
	}
	stat, err := file.Stat()
	if err != nil || stat.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	if runtime.GOOS == "windows" {
		return false
	}
	return strings.ToLower(os.Getenv("TERM")) != "dumb"
}

type bubbleProgress struct {
	program *tea.Program
	total   int64
	current int64
	message string
	doneCh  chan struct{}
	done    bool
	state   *progressState
}

func (p *bubbleProgress) Set(current int64) {
	p.Update(current, "")
}

func (p *bubbleProgress) Message(message string) {
	p.Update(-1, message)
}

func (p *bubbleProgress) Update(current int64, message string) {
	if p.done {
		return
	}
	if current >= 0 {
		p.current = current
	}
	if message != "" {
		p.message = message
	}
	p.send(p.current, false)
}

func (p *bubbleProgress) Finish() {
	if p.done {
		return
	}
	p.done = true
	p.send(p.total, true)
	select {
	case <-p.doneCh:
	case <-time.After(2 * time.Second):
	}
	p.deactivate()
}

func (p *bubbleProgress) send(current int64, done bool) {
	percent := 0.0
	if done {
		percent = 1
	} else if p.total > 0 {
		if current < 0 {
			current = 0
		}
		if current > p.total {
			current = p.total
		}
		percent = float64(current) / float64(p.total)
	}
	p.program.Send(progressUpdateMsg{percent: percent, message: p.message, done: done})
}

func (p *bubbleProgress) deactivate() {
	if p.state != nil {
		p.state.active = false
	}
}
