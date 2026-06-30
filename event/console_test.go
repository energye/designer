package event

import "testing"

func TestConsoleWriterReceivesRawMessage(t *testing.T) {
	var gotLevel ConsoleLevel
	var gotMessage string
	SetConsoleWriter(func(level ConsoleLevel, message string) {
		gotLevel = level
		gotMessage = message
	})
	defer SetConsoleWriter(nil)

	ConsoleWriteError("build failed")

	if gotLevel != ConsoleLevelError {
		t.Fatalf("level=%v want %v", gotLevel, ConsoleLevelError)
	}
	if gotMessage != "build failed" {
		t.Fatalf("message=%q want build failed", gotMessage)
	}
}
