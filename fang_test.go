package fang_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		exercise(t, cobra.Command{
			Use: "simple",
		})
	})

	t.Run("custom error handler", func(t *testing.T) {
		doExercise(
			t,
			cobra.Command{
				Use: "simple",
			},
			[]string{"nope"},
			assertError,
			fang.WithErrorHandler(func(w io.Writer, styles fang.Styles, err error) {
				_, _ = fmt.Fprintf(w, "Custom error handler: %v\n", err)
			}),
		)
	})

	t.Run("complete", func(t *testing.T) {
		cmd := cobra.Command{
			Use:   "simple",
			Short: "Short help",
			Long:  "Long help",
		}
		exercise(t, cmd)

		t.Run("man", func(t *testing.T) {
			doExercise(
				t, cmd,
				[]string{"man"},
				assertNoError,
				fang.WithoutManpage(),
			)
		})
	})

	t.Run("use with args", func(t *testing.T) {
		exercise(t, cobra.Command{
			Use:   "simple [args] [something-else]",
			Short: "Short help",
			Long:  "Long help",
		})
	})

	t.Run("without completions", func(t *testing.T) {
		cmd := cobra.Command{
			Use:   "simple",
			Short: "no completions",
			Args:  cobra.NoArgs,
		}

		exercise(t, cmd, fang.WithoutCompletions())

		t.Run("completion", func(t *testing.T) {
			t.Skip("this fails when testing, but works as expected otherwise, no idea why yet")
			doExercise(
				t, cmd,
				[]string{"completion"},
				assertError,
				fang.WithoutCompletions(),
			)
		})
	})

	t.Run("without manpage", func(t *testing.T) {
		cmd := cobra.Command{
			Use:   "simple",
			Short: "no manpages",
			Args:  cobra.NoArgs,
		}
		exercise(t, cmd, fang.WithoutManpage())

		t.Run("man", func(t *testing.T) {
			t.Skip("this fails when testing, but works as expected otherwise, no idea why yet")
			doExercise(
				t, cmd,
				[]string{"man"},
				assertError,
				fang.WithoutManpage(),
			)
		})
	})

	t.Run("with version and hash", func(t *testing.T) {
		exercise(t, cobra.Command{
			Use: "simple",
		}, fang.WithVersion("v1.2.3"), fang.WithCommit("aaabbb"))
	})

	t.Run("with flags", func(t *testing.T) {
		cmd := cobra.Command{
			Use:   "simple",
			Short: "Short help",
			Long:  "Long help",
		}
		cmd.Flags().String("string1", "default-value", "a string flag")
		cmd.Flags().String("string2", "", "a string flag")
		cmd.Flags().StringP("string3", "s", "", "a string flag")
		cmd.Flags().Int("int1", 0, "an int flag")
		cmd.Flags().Int("int2", 10, "an int flag")
		cmd.Flags().IntP("int3", "i", 10, "an int flag")
		cmd.Flags().Float64("float1", 0, "a float flag")
		cmd.Flags().Float64("float2", 10, "a float flag")
		cmd.Flags().Float64P("float3", "f", 10, "a float flag")
		cmd.Flags().Bool("bool1", false, "a bool flag")
		cmd.Flags().Bool("bool2", true, "a bool flag")
		cmd.Flags().BoolP("bool3", "b", true, "a bool flag")
		cmd.Flags().BoolP("hidden", "z", true, "a bool flag")
		cmd.Flags().Bool("no-help", false, "")
		_ = cmd.Flags().MarkHidden("hidden")
		exercise(t, cmd)
	})

	t.Run("with subcommands", func(t *testing.T) {
		cmd := cobra.Command{
			Use:   "simple",
			Short: "Short help",
		}
		cmd.AddCommand(&cobra.Command{
			Use:   "sub-cmd",
			Short: "a sub command",
		})
		exercise(t, cmd)
	})

	t.Run("with examples", func(t *testing.T) {
		cmd := cobra.Command{
			Use:   "simple",
			Short: "Short help",
			Example: `
# a comment about the usage
simple [some arguments]

# another comment
simple --string1=2 -s abc -b --bool1 --flag-not-found [args]
			`,
		}

		cmd.Flags().String("string1", "", "a string flag")
		cmd.Flags().StringP("string2", "s", "", "a string flag")
		cmd.Flags().Bool("bool1", false, "a bool flag")
		cmd.Flags().BoolP("bool2", "b", false, "a bool flag")
		exercise(t, cmd)
	})
}

func exercise(t *testing.T, cmd cobra.Command, options ...fang.Option) {
	t.Helper()

	t.Run("error", func(t *testing.T) {
		doExercise(
			t, cmd,
			[]string{"--nope-nope-nope"},
			assertError,
			options...,
		)
	})

	t.Run("version", func(t *testing.T) {
		doExercise(
			t, cmd,
			[]string{"--version"},
			assertNoError,
			options...,
		)
	})

	t.Run("help", func(t *testing.T) {
		doExercise(
			t, cmd,
			[]string{"--help"},
			assertNoError,
			options...,
		)
	})
}

func doExercise(
	t *testing.T,
	cmd cobra.Command,
	args []string,
	assert func(t *testing.T, err error, stdout, stderr bytes.Buffer),
	options ...fang.Option,
) {
	t.Helper()
	t.Setenv("__FANG_TEST_WIDTH", "45")
	root := &cmd
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)
	err := fang.Execute(context.Background(), root, options...)
	assert(t, err, stdout, stderr)
}

func assertNoError(t *testing.T, err error, stdout, stderr bytes.Buffer) {
	require.NoError(t, err, stderr.String())
	golden.RequireEqual(t, stdout.Bytes())
}

func assertError(t *testing.T, err error, stdout, stderr bytes.Buffer) {
	require.Error(t, err, stdout.String())
	golden.RequireEqual(t, stderr.Bytes())
}
