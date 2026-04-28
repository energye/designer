package docs

import (
	"testing"
)

func TestGetClassDesc(t *testing.T) {
	tests := []struct {
		name      string
		className string
		wantEmpty bool
	}{
		{
			name:      "existing class TAction",
			className: "TAction",
			wantEmpty: false,
		},
		{
			name:      "existing class TApplication",
			className: "TApplication",
			wantEmpty: false,
		},
		{
			name:      "non-existing class",
			className: "NonExistentClass",
			wantEmpty: true,
		},
		{
			name:      "empty class name",
			className: "",
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetClassDesc(tt.className)
			if tt.wantEmpty && desc != "" {
				t.Errorf("GetClassDesc(%s) = %v, want empty", tt.className, desc)
			}
			if !tt.wantEmpty && desc == "" {
				t.Errorf("GetClassDesc(%s) returned empty, want non-empty description", tt.className)
			}
		})
	}
}

func TestHasClass(t *testing.T) {
	tests := []struct {
		name      string
		className string
		want      bool
	}{
		{
			name:      "existing class TAction",
			className: "TAction",
			want:      true,
		},
		{
			name:      "existing class TApplication",
			className: "TApplication",
			want:      true,
		},
		{
			name:      "non-existing class",
			className: "NonExistentClass",
			want:      false,
		},
		{
			name:      "empty class name",
			className: "",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasClass(tt.className)
			if got != tt.want {
				t.Errorf("HasClass(%s) = %v, want %v", tt.className, got, tt.want)
			}
		})
	}
}

func TestHasProperty(t *testing.T) {
	tests := []struct {
		name         string
		className    string
		propertyName string
		want         bool
	}{
		{
			name:         "existing property in TAction",
			className:    "TAction",
			propertyName: "Caption",
			want:         true,
		},
		{
			name:         "existing property in TApplication",
			className:    "TApplication",
			propertyName: "Title",
			want:         true,
		},
		{
			name:         "non-existing property",
			className:    "TAction",
			propertyName: "NonExistentProperty",
			want:         false,
		},
		{
			name:         "non-existing class",
			className:    "NonExistentClass",
			propertyName: "SomeProperty",
			want:         false,
		},
		{
			name:         "empty property name",
			className:    "TAction",
			propertyName: "",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasProperty(tt.className, tt.propertyName)
			if got != tt.want {
				t.Errorf("HasProperty(%s, %s) = %v, want %v", tt.className, tt.propertyName, got, tt.want)
			}
		})
	}
}

func TestGetPropertyDescWithReference(t *testing.T) {
	desc := GetPropertyDesc("TAction", "ActionList")
	if desc == "" {
		t.Error("GetPropertyDesc for referenced property TAction.ActionList returned empty")
	}
}

func TestConcurrentAccess(t *testing.T) {
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			_ = GetClassDesc("TAction")
			_ = GetPropertyDesc("TAction", "Caption")
			_ = HasClass("TApplication")
			_ = HasProperty("TAction", "Enabled")
			_ = GetClassProperties("TActionList")
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
