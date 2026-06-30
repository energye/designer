package clui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type simpleUI struct {
	opt    Options
	state  *state
	reader *bufio.Reader
}

func (u *simpleUI) Info(args ...any) {
	u.finishProgressLine()
	fmt.Fprintln(u.opt.Out, args...)
}
func (u *simpleUI) Warn(args ...any) {
	u.finishProgressLine()
	fmt.Fprintln(u.opt.Err, append([]any{"[WARN]"}, args...)...)
}
func (u *simpleUI) Error(args ...any) {
	u.finishProgressLine()
	fmt.Fprintln(u.opt.Err, append([]any{"[ERROR]"}, args...)...)
}
func (u *simpleUI) Success(args ...any) {
	u.finishProgressLine()
	fmt.Fprintln(u.opt.Out, args...)
}

func (u *simpleUI) Input(title, def string) (string, error) {
	for {
		if def == "" {
			fmt.Fprintf(u.opt.Out, "%s: ", title)
		} else {
			fmt.Fprintf(u.opt.Out, "%s [%s]: ", title, def)
		}
		text, err := u.reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		text = strings.TrimSpace(text)
		if text != "" {
			return text, nil
		}
		if def != "" {
			return def, nil
		}
	}
}

func (u *simpleUI) Select(title string, items []string, def int) (int, error) {
	if len(items) == 0 {
		return -1, errors.New("clui: select has no items")
	}
	for {
		fmt.Fprintln(u.opt.Out, title+":")
		for i, item := range items {
			fmt.Fprintf(u.opt.Out, "  %d. %s\n", i+1, item)
		}
		fmt.Fprint(u.opt.Out, "Select number")
		if def >= 0 && def < len(items) {
			fmt.Fprintf(u.opt.Out, " [%d, 0 to cancel]", def+1)
		} else {
			fmt.Fprint(u.opt.Out, " (0 to cancel)")
		}
		fmt.Fprint(u.opt.Out, ": ")
		text, err := u.reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return -1, err
		}
		text = strings.TrimSpace(text)
		if text == "0" || (text == "" && (def < 0 || def >= len(items))) {
			return -1, errors.New("clui: selection canceled")
		}
		if text == "" && def >= 0 && def < len(items) {
			return def, nil
		}
		idx, err := strconv.Atoi(text)
		if err == nil && idx > 0 && idx <= len(items) {
			return idx - 1, nil
		}
	}
}

func (u *simpleUI) Confirm(title string, def bool) (bool, error) {
	if u.opt.AssumeYes {
		return true, nil
	}
	for {
		if def {
			fmt.Fprintf(u.opt.Out, "%s [Y/n]: ", title)
		} else {
			fmt.Fprintf(u.opt.Out, "%s [y/N]: ", title)
		}
		text, err := u.reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}
		text = strings.ToLower(strings.TrimSpace(text))
		switch text {
		case "":
			return def, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
	}
}

func (u *simpleUI) Progress(title string, total int64) Progress {
	return &lineProgress{
		out:      u.opt.Out,
		title:    title,
		total:    total,
		ansi:     useANSIProgress(u.opt.Out),
		shared:   u.state,
		progress: startProgress(u.state),
	}
}

func (u *simpleUI) finishProgressLine() {
	finishProgressLine(u.opt.Out, u.state)
}
