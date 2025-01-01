package selectbox

import (
	"fmt"
	"testing"

	"github.com/trysourcetool/sourcetool-go/internal/selectbox"
)

func TestOptions(t *testing.T) {
	opts := &selectbox.Options{}
	options := []string{"Option 1", "Option 2", "Option 3"}

	option := Options(options...)
	option(opts)

	if len(opts.Options) != len(options) {
		t.Errorf("Options length = %v, want %v", len(opts.Options), len(options))
	}

	for i, opt := range opts.Options {
		if opt != options[i] {
			t.Errorf("Option[%d] = %v, want %v", i, opt, options[i])
		}
	}
}

func TestPlaceholder(t *testing.T) {
	opts := &selectbox.Options{}
	placeholder := "Select an option"

	option := Placeholder(placeholder)
	option(opts)

	if opts.Placeholder != placeholder {
		t.Errorf("Placeholder = %v, want %v", opts.Placeholder, placeholder)
	}
}

func TestDefaultValue(t *testing.T) {
	opts := &selectbox.Options{}
	defaultValue := "Option 1"

	option := DefaultValue(defaultValue)
	option(opts)

	if opts.DefaultValue == nil {
		t.Fatal("DefaultValue is nil")
	}
	if *opts.DefaultValue != defaultValue {
		t.Errorf("DefaultValue = %v, want %v", *opts.DefaultValue, defaultValue)
	}
}

func TestRequired(t *testing.T) {
	tests := []struct {
		name     string
		required bool
	}{
		{"Required true", true},
		{"Required false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &selectbox.Options{}
			option := Required(tt.required)
			option(opts)

			if opts.Required != tt.required {
				t.Errorf("Required = %v, want %v", opts.Required, tt.required)
			}
		})
	}
}

func TestDisabled(t *testing.T) {
	tests := []struct {
		name     string
		disabled bool
	}{
		{"Disabled true", true},
		{"Disabled false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &selectbox.Options{}
			option := Disabled(tt.disabled)
			option(opts)

			if opts.Disabled != tt.disabled {
				t.Errorf("Disabled = %v, want %v", opts.Disabled, tt.disabled)
			}
		})
	}
}

func TestFormatFunc(t *testing.T) {
	opts := &selectbox.Options{}
	formatFunc := func(value string, index int) string {
		return fmt.Sprintf("%d. %s", index+1, value)
	}

	option := FormatFunc(formatFunc)
	option(opts)

	if opts.FormatFunc == nil {
		t.Fatal("FormatFunc is nil")
	}

	testValue := "Test"
	testIndex := 0
	expected := "1. Test"
	result := opts.FormatFunc(testValue, testIndex)

	if result != expected {
		t.Errorf("FormatFunc result = %v, want %v", result, expected)
	}
}

func TestMultipleOptions(t *testing.T) {
	opts := &selectbox.Options{}
	options := []string{"Option 1", "Option 2"}
	defaultValue := "Option 1"
	placeholder := "Select an option"
	required := true
	disabled := false

	Options(options...)(opts)
	DefaultValue(defaultValue)(opts)
	Placeholder(placeholder)(opts)
	Required(required)(opts)
	Disabled(disabled)(opts)

	if len(opts.Options) != len(options) {
		t.Errorf("Options length = %v, want %v", len(opts.Options), len(options))
	}

	if opts.DefaultValue == nil {
		t.Fatal("DefaultValue is nil")
	}
	if *opts.DefaultValue != defaultValue {
		t.Errorf("DefaultValue = %v, want %v", *opts.DefaultValue, defaultValue)
	}

	if opts.Placeholder != placeholder {
		t.Errorf("Placeholder = %v, want %v", opts.Placeholder, placeholder)
	}

	if opts.Required != required {
		t.Errorf("Required = %v, want %v", opts.Required, required)
	}

	if opts.Disabled != disabled {
		t.Errorf("Disabled = %v, want %v", opts.Disabled, disabled)
	}
}
