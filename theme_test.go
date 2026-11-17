package fang

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

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

// TestStyleUsage preserves positional arguments in both full usage and command lists.
func TestStyleUsage(t *testing.T) {
	for _, tt := range []struct {
		name, use, want string
		child, flag     bool
	}{
		{name: "optional before required", use: "run [SERIAL] FILE", want: "run [SERIAL] FILE"},
		{name: "optional after required", use: "run FILE [SERIAL]", want: "run FILE [SERIAL]"},
		{name: "multiple optional arguments", use: "run [FIRST] FILE [SECOND] TAIL", want: "run [FIRST] FILE [SECOND] TAIL"},
		{name: "optional containing spaces", use: "run [COLOR NAME] FILE", want: "run [COLOR NAME] FILE"},
		{name: "required only", use: "run FILE", want: "run FILE"},
		{name: "implicit commands", use: "tool", child: true, want: "tool [command]"},
		{name: "implicit flags", use: "run [SERIAL] FILE", flag: true, want: "run [SERIAL] FILE [--flags]"},
		{name: "explicit flags", use: "run [SERIAL] FILE [flags]", flag: true, want: "run [SERIAL] FILE [--flags]"},
		{name: "explicit double dash flags", use: "run [SERIAL] FILE [--flags]", flag: true, want: "run [SERIAL] FILE [--flags]"},
		{name: "canonical with infix optional", use: "tool [args] [FIRST] FILE [SECOND] TAIL", child: true, flag: true, want: "tool [command] [args] [FIRST] FILE [SECOND] TAIL [--flags]"},
		{name: "canonical placeholders", use: "tool [args] [something-else]", child: true, flag: true, want: "tool [command] [args] [something-else] [--flags]"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := &cobra.Command{Use: tt.use}
			if tt.child {
				c.AddCommand(&cobra.Command{Use: "child", Run: func(*cobra.Command, []string) {}})
			}
			if tt.flag {
				c.Flags().Bool("verbose", false, "verbose output")
			}
			for _, complete := range []bool{false, true} {
				require.Equal(t, tt.want, styleUsage(c, Program{}, complete), "complete=%v", complete)
			}
		})
	}

	t.Run("argument styles and parent path", func(t *testing.T) {
		styles := Program{
			Name:           lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")),
			Command:        lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")),
			DimmedArgument: lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff")),
		}
		root := &cobra.Command{Use: "tool"}
		c := &cobra.Command{Use: "run [FIRST] FILE [SECOND] TAIL"}
		root.AddCommand(c)
		arguments := styles.DimmedArgument.Render(" [FIRST]") +
			styles.Command.Render(" FILE") +
			styles.DimmedArgument.Render(" [SECOND]") +
			styles.Command.Render(" TAIL")
		require.Equal(t, styles.Command.Render("run")+arguments, styleUsage(c, styles, false))
		require.Equal(t, styles.Name.Render("tool")+styles.Command.Render(" run")+arguments, styleUsage(c, styles, true))
	})
}
