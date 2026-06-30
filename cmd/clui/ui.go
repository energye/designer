package clui

import (
	"bufio"
	"io"
	"os"
	"sync"
)

type UI interface {
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Success(args ...any)
	Input(title, def string) (string, error)
	Select(title string, items []string, def int) (int, error)
	Confirm(title string, def bool) (bool, error)
	Progress(title string, total int64) Progress
}

type Progress interface {
	Set(current int64)
	Message(message string)
	Update(current int64, message string)
	Finish()
}

type Options struct {
	In        io.Reader
	Out       io.Writer
	Err       io.Writer
	NoTUI     bool
	AssumeYes bool
}

type state struct {
	mu       sync.Mutex
	progress *progressState
}

type progressState struct {
	active bool
}

func New(options ...Options) UI {
	opt := defaultOptions()
	if len(options) > 0 {
		opt = mergeOptions(opt, options[0])
	}
	shared := &state{}
	if opt.NoTUI || shouldUseSimpleUI() {
		return NewSimple(opt, shared)
	}
	return NewBubble(opt, shared)
}

func NewSimple(opt Options, states ...*state) UI {
	opt = mergeOptions(defaultOptions(), opt)
	var shared *state
	if len(states) > 0 {
		shared = states[0]
	}
	if shared == nil {
		shared = &state{}
	}
	return &simpleUI{opt: opt, state: shared, reader: bufio.NewReader(opt.In)}
}

func defaultOptions() Options {
	return Options{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}
}

func mergeOptions(base, override Options) Options {
	if override.In != nil {
		base.In = override.In
	}
	if override.Out != nil {
		base.Out = override.Out
	}
	if override.Err != nil {
		base.Err = override.Err
	}
	base.NoTUI = override.NoTUI
	base.AssumeYes = override.AssumeYes
	return base
}
