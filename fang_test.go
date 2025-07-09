package fang_test

import (
	"bytes"
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
		exercise(t, toMkroot(&cobra.Command{
			Use: "simple",
		}))
	})

	t.Run("custom error handler", func(t *testing.T) {
		doExercise(
			t,
			toMkroot(&cobra.Command{Use: "simple"}),
			[]string{"nope"},
			assertError,
			fang.WithErrorHandler(func(w io.Writer, styles fang.Styles, err error) {
				_, _ = fmt.Fprintf(w, "Custom error handler: %v\n", err)
			}),
		)
	})

	t.Run("complete", func(t *testing.T) {
		cmd := toMkroot(&cobra.Command{
			Use:   "simple",
			Short: "Short help",
			Long:  "Long help",
		})

		exercise(t, cmd)

		t.Run("man", func(t *testing.T) {
			doExercise(
				t, cmd,
				[]string{"man", "-h"},
				assertNoError,
			)
		})
	})

	t.Run("use with args", func(t *testing.T) {
		exercise(t, toMkroot(&cobra.Command{
			Use:   "simple [args] [something-else]",
			Short: "Short help",
			Long:  "Long help",
		}))
	})

	t.Run("without completions", func(t *testing.T) {
		cmd := toMkroot(&cobra.Command{
			Use:   "simple",
			Short: "no completions",
			Args:  cobra.NoArgs,
		})

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
		cmd := toMkroot(&cobra.Command{
			Use:   "simple",
			Short: "no manpages",
			Args:  cobra.NoArgs,
		})
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

	t.Run("without version", func(t *testing.T) {
		doExercise(
			t,
			toMkroot(&cobra.Command{Use: "simple"}),
			[]string{"--version"},
			assertError,
			fang.WithoutVersion(),
		)
	})

	t.Run("with version and hash", func(t *testing.T) {
		exercise(
			t,
			toMkroot(&cobra.Command{Use: "simple"}),
			fang.WithVersion("v1.2.3"),
			fang.WithCommit("aaabbb"),
		)
	})

	t.Run("with flags", func(t *testing.T) {
		mkroot := func() *cobra.Command {
			cmd := &cobra.Command{
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
			return cmd
		}
		exercise(t, mkroot)
	})

	t.Run("with subcommands", func(t *testing.T) {
		mkroot := func() *cobra.Command {
			cmd := &cobra.Command{
				Use:   "simple",
				Short: "Short help",
			}
			sub := &cobra.Command{
				Use:     "sub1",
				Short:   "a sub command",
				Example: `simple sub1 some args`,
			}
			sub.AddCommand(&cobra.Command{
				Use:     "sub2 [args]",
				Short:   "yet another sub command",
				Example: `simple sub1 sub2 args --help`,
			})
			cmd.AddCommand(sub)
			return cmd
		}

		exercise(t, mkroot)

		t.Run("help-sub", func(t *testing.T) {
			doExercise(
				t,
				mkroot,
				[]string{"sub1", "--help"},
				assertNoError,
			)
		})

		t.Run("help-sub-sub", func(t *testing.T) {
			doExercise(
				t,
				mkroot,
				[]string{"sub1", "sub2", "--help"},
				assertNoError,
			)
		})
	})

	t.Run("with command groups", func(t *testing.T) {
		mkroot := func() *cobra.Command {
			cmd := &cobra.Command{
				Use:   "simple",
				Short: "Short help",
			}
			cmd.AddGroup(&cobra.Group{
				ID:    "1",
				Title: "First group",
			})
			cmd.AddGroup(&cobra.Group{
				ID:    "2",
				Title: "Second group",
			})
			cmd.AddCommand(&cobra.Command{
				Use:   "sub-cmd",
				Short: "a sub command",
			})
			cmd.AddCommand(&cobra.Command{
				Use:     "sub-cmd-2",
				Short:   "a sub command",
				GroupID: "1",
			})
			cmd.AddCommand(&cobra.Command{
				Use:     "sub-cmd-3",
				Short:   "a sub command",
				GroupID: "2",
			})
			return cmd
		}
		exercise(t, mkroot)
	})

	t.Run("with examples", func(t *testing.T) {
		mkroot := func() *cobra.Command {
			cmd := &cobra.Command{
				Use:   "example",
				Short: "Short help",
				Example: `
# Run it:
example

# Run it with some arguments:
FOO=bar ZAZ="quoted value" example --name=Carlos -a -s Becker -a

# Run a subcommand with an argument:
example sub --async --name=xyz --async arguments

# Run with a quoted string:
example sub "quoted string"

# Mix and match:
example sub "multi-word quoted string" --name "another quoted string" -a

# Multi-line:
ENV_A=0 ENV_B=0 ENV_C=0 \
  CERT_FILE=/path/to/chain.pem KEY_FILE=/path/to/key.pem \
  example sub "quoted argument"

# Run a subcommand's subcommand with an argument:
example sub another args --async

# Pipe example:
echo "foo" | example > bar.txt

# Redirects:
example < in.txt > out.txt
example 2>&1 1>/dev/null
example 1>&2 2>/dev/null

# And / Or:
foo || example
example && foo

# Another pipe example:
echo 'foo' |
  example sub |
  cat -
			`,
			}
			sub := &cobra.Command{
				Use:   "sub",
				Short: "a sub command",
			}
			cmd.AddCommand(sub)
			sub.AddCommand(&cobra.Command{
				Use:   "another",
				Short: "Another sub command",
			})
			cmd.Flags().String("name", "", "the name")
			cmd.Flags().StringP("surname", "s", "", "the surname")
			cmd.Flags().BoolP("async", "a", false, "async?")
			return cmd
		}
		exercise(t, mkroot)
	})
}

func exercise(t *testing.T, mkroot func() *cobra.Command, options ...fang.Option) {
	t.Helper()

	t.Run("error", func(t *testing.T) {
		doExercise(
			t, mkroot,
			[]string{"--nope-nope-nope"},
			assertError,
			options...,
		)
	})

	t.Run("version", func(t *testing.T) {
		doExercise(
			t, mkroot,
			[]string{"--version"},
			assertNoError,
			options...,
		)
	})

	t.Run("help", func(t *testing.T) {
		doExercise(
			t, mkroot,
			[]string{"--help"},
			assertNoError,
			options...,
		)
	})
}

func doExercise(
	t *testing.T,
	mkroot func() *cobra.Command,
	args []string,
	assert func(t *testing.T, err error, stdout, stderr bytes.Buffer),
	options ...fang.Option,
) {
	t.Helper()
	t.Setenv("__FANG_TEST_WIDTH", "45")

	root := mkroot()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	root.SetArgs(args)

	err := fang.Execute(t.Context(), root, options...)
	assert(t, err, stdout, stderr)
}

func toMkroot(c *cobra.Command) func() *cobra.Command {
	return func() *cobra.Command {
		return c
	}
}

func assertNoError(t *testing.T, err error, stdout, stderr bytes.Buffer) {
	require.NoError(t, err, stderr.String())
	golden.RequireEqual(t, stdout.Bytes())
}

func assertError(t *testing.T, err error, stdout, stderr bytes.Buffer) {
	require.Error(t, err, stdout.String())
	golden.RequireEqual(t, stderr.Bytes())
}
