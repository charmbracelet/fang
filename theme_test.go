package fang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTitleFirstWord verifies that the first word of a string is capitalized
// while preserving all whitespace, including newlines and tabs. This is used
// to style error messages and flag descriptions in help output.
func TestTitleFirstWord(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple string",
			input:    "hello world",
			expected: "Hello world",
		},
		{
			name:     "leading whitespace preserved",
			input:    "  hello world",
			expected: "  Hello world",
		},
		{
			name:     "newline preserved",
			input:    "hello\nworld",
			expected: "Hello\nworld",
		},
		{
			name:     "tab preserved",
			input:    "hello\tworld",
			expected: "Hello\tworld",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "   ",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "already capitalized",
			input:    "Hello world",
			expected: "Hello world",
		},
		{
			name:     "error prefix gets capitalized",
			input:    "error: something happened",
			expected: "Error: something happened",
		},
		{
			name:     "Error prefix preserved",
			input:    "Error: something happened",
			expected: "Error: something happened",
		},
		{
			name:     "multiline error message with error prefix",
			input:    "error:\ndetails here\nmore info",
			expected: "Error:\ndetails here\nmore info",
		},
		{
			name:     "multiline error message with Error prefix",
			input:    "Error:\ndetails here\nmore info",
			expected: "Error:\ndetails here\nmore info",
		},
		{
			name:     "mixed whitespace",
			input:    "\t\n  hello\tworld",
			expected: "\t\n  Hello\tworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := titleFirstWord(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
