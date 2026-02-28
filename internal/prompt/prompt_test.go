package prompt

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(input string) *bytes.Buffer {
	SetReader(strings.NewReader(input))
	out := new(bytes.Buffer)
	SetOutput(out)
	return out
}

func TestString(t *testing.T) {
	defer Reset()

	out := setup("hello world\n")
	result := String("Enter value")

	assert.Equal(t, "hello world", result)
	assert.Contains(t, out.String(), "Enter value:")
}

func TestStringEmpty(t *testing.T) {
	defer Reset()

	setup("\n")
	result := String("Enter value")

	assert.Equal(t, "", result)
}

func TestDefault(t *testing.T) {
	defer Reset()

	setup("custom\n")
	result := Default("Value", "default")
	assert.Equal(t, "custom", result)
}

func TestDefaultEmpty(t *testing.T) {
	defer Reset()

	setup("\n")
	result := Default("Value", "default")
	assert.Equal(t, "default", result)
}

func TestSecret(t *testing.T) {
	defer Reset()

	out := setup("secret123\n")
	result := Secret("Password")

	assert.Equal(t, "secret123", result)
	assert.Contains(t, out.String(), "Password:")
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"y\n", true},
		{"yes\n", true},
		{"Y\n", true},
		{"YES\n", true},
		{"n\n", false},
		{"no\n", false},
		{"\n", false},
		{"anything\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			defer Reset()
			setup(tt.input)
			result := Confirm("Continue?")
			assert.Equal(t, tt.expected, result, "input: %q", tt.input)
		})
	}
}

func TestConfirmDefault(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultYes bool
		expected   bool
	}{
		{"default_yes_empty", "\n", true, true},
		{"default_no_empty", "\n", false, false},
		{"override_yes_with_n", "n\n", true, false},
		{"override_no_with_y", "y\n", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()
			setup(tt.input)
			result := ConfirmDefault("Continue?", tt.defaultYes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetReaderAndOutput(t *testing.T) {
	defer Reset()

	in := strings.NewReader("test\n")
	out := new(bytes.Buffer)

	SetReader(in)
	SetOutput(out)

	result := String("Label")
	assert.Equal(t, "test", result)
	assert.Contains(t, out.String(), "Label:")
}

func TestReset(t *testing.T) {
	Reset()
}
