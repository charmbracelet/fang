package fang_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestAdditionalHelpTopics(t *testing.T) {
	mkroot := func() *cobra.Command {
		root := &cobra.Command{Use: "topics", Short: "Commands and reference documentation"}
		root.AddGroup(&cobra.Group{ID: "work", Title: "Work commands"})
		root.AddCommand(
			&cobra.Command{Use: "guide", Short: "Read the guide"},
			&cobra.Command{Use: "formats", Short: "Supported formats", GroupID: "work"},
			&cobra.Command{Use: "hidden-topic", Hidden: true},
			&cobra.Command{
				Use: "run", Short: "Run the application", Run: func(*cobra.Command, []string) {},
			},
			&cobra.Command{
				Use: "work", Short: "Run a task", GroupID: "work",
				RunE: func(*cobra.Command, []string) error { return nil },
			},
		)
		reference := &cobra.Command{Use: "reference", Short: "Reference topics"}
		reference.AddCommand(&cobra.Command{Use: "details", Short: "More details"})
		root.AddCommand(reference)
		remote := &cobra.Command{Use: "remote", Short: "Manage remotes"}
		remote.AddCommand(&cobra.Command{
			Use: "list", Run: func(*cobra.Command, []string) {},
		})
		root.AddCommand(remote)
		return root
	}

	doExercise(t, mkroot, []string{"--help"},
		func(t *testing.T, err error, stdout, stderr bytes.Buffer) {
			require.NoError(t, err, stderr.String())
			commands, topics, found := strings.Cut(stdout.String(), "ADDITIONAL HELP TOPICS")
			require.True(t, found, stdout.String())
			for _, topic := range []string{"guide", "formats", "reference"} {
				require.NotContains(t, commands, "    "+topic)
				require.Contains(t, topics, "    "+topic)
			}
			for _, command := range []string{"run", "work", "remote [command]", "help [command]"} {
				require.Contains(t, commands, "    "+command)
				require.NotContains(t, topics, "    "+command)
			}
			require.Contains(t, commands, "WORK COMMANDS")
			require.NotContains(t, stdout.String(), "hidden-topic")
			golden.RequireEqual(t, stdout.Bytes())
		})
}
