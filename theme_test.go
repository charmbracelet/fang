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
		{
			name:     "number first word",
			input:    "1234 world",
			expected: "1234 world",
		},
		{
			name:     "single letter",
			input:    "a cat",
			expected: "A cat",
		},
		{
			name:     "accented character",
			input:    "été chaud",
			expected: "Été chaud",
		},
		{
			name:     "punctuation first word",
			input:    "- item",
			expected: "- item",
		},
		{
			name:     "all caps",
			input:    "ABC",
			expected: "Abc",
		},
		{
			name:     "all caps with second word",
			input:    "ABC def",
			expected: "Abc def",
		},
		{
			name:     "CJK characters",
			input:    "你好 世界",
			expected: "你好 世界",
		},
		{
			name:     "CJK characters with newline",
			input:    "你好\n世界",
			expected: "你好\n世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := titleFirstWord(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
