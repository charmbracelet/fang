package main

import (
	"context"
	"os"
	"time"

	"github.com/charmbracelet/serpentine"
	"github.com/spf13/cobra"
)

func main() {
	var cmd *cobra.Command
	cmd = &cobra.Command{
		Use:   "example [args]",
		Short: "An example program!",
		Long: `A pretty silly example program!

It doesn’t really do anything, but that’s the point.`,
		Example: `
# Run it:
example

# Run it with some arguments:
example --name=Carlos -a -s Becker -a

# Run a subcommand with an argument:
example sub --async --foo=xyz --async arguments
		`,
		// Args: cobra.ArbitraryArgs,
		Run: func(*cobra.Command, []string) {
			_ = cmd.Help()
		},
	}
	var foo string
	var bar int
	var baz float64
	var d time.Duration
	cmd.PersistentFlags().StringVarP(&foo, "surname", "s", "doe", "Your surname")
	cmd.Flags().StringVar(&foo, "name", "john", "Your name")
	cmd.Flags().DurationVar(&d, "duration", 0, "Time since your last commit")
	cmd.Flags().IntVar(&bar, "age", 0, "Your age")
	cmd.Flags().Float64Var(&baz, "idk", 0.0, "I don't know")
	cmd.Flags().BoolP("async", "a", false, "Run async")

	_ = cmd.Flags().MarkHidden("age")
	_ = cmd.Flags().MarkHidden("duration")
	_ = cmd.Flags().MarkHidden("idk")

	cmd.AddCommand(&cobra.Command{
		Use:   "sub [command] [flags] [args]",
		Short: "An example subcommand",
	})

	// This is where the magic happens.
	if err := serpentine.Execute(
		context.Background(),
		cmd,
		serpentine.WithoutManpage(),
		serpentine.WithoutCompletions(),
	); err != nil {
		os.Exit(1)
	}
}
