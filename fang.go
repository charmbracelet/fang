// Package fang provides styling for cobra commands.
package fang

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
	mango "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

const shaLen = 7

// ErrorHandler handles an error, printing them to the given [io.Writer].
type ErrorHandler = func(w io.Writer, styles Styles, err error)

// ColorSchemeFunc gets a [lipgloss.LightDarkFunc] and returns a [ColorScheme].
type ColorSchemeFunc = func(lipgloss.LightDarkFunc) ColorScheme

type settings struct {
	completions bool
	manpages    bool
	version     string
	commit      string
	colorscheme ColorSchemeFunc
	errHandler  ErrorHandler
}

// Option changes fang settings.
type Option func(*settings)

// WithoutCompletions disables completions.
func WithoutCompletions() Option {
	return func(s *settings) {
		s.completions = false
	}
}

// WithoutManpage disables man pages.
func WithoutManpage() Option {
	return func(s *settings) {
		s.manpages = false
	}
}

// WithColorSchemeFunc sets a function that return colorscheme.
func WithColorSchemeFunc(cs ColorSchemeFunc) Option {
	return func(s *settings) {
		s.colorscheme = cs
	}
}

// WithTheme sets the colorscheme.
//
// Deprecated: use [WithColorSchemeFunc] instead.
func WithTheme(theme ColorScheme) Option {
	return func(s *settings) {
		s.colorscheme = func(lipgloss.LightDarkFunc) ColorScheme {
			return theme
		}
	}
}

// WithVersion sets the version.
func WithVersion(version string) Option {
	return func(s *settings) {
		s.version = version
	}
}

// WithCommit sets the commit SHA.
func WithCommit(commit string) Option {
	return func(s *settings) {
		s.commit = commit
	}
}

// WithErrorHandler sets the error handler.
func WithErrorHandler(handler ErrorHandler) Option {
	return func(s *settings) {
		s.errHandler = handler
	}
}

// Execute applies fang to the command and executes it.
func Execute(ctx context.Context, root *cobra.Command, options ...Option) error {
	opts := settings{
		manpages:    true,
		completions: true,
		colorscheme: DefaultColorScheme,
		errHandler:  DefaultErrorHandler,
	}

	for _, option := range options {
		option(&opts)
	}

	helpFunc := func(c *cobra.Command, _ []string) {
		w := colorprofile.NewWriter(c.OutOrStdout(), os.Environ())
		helpFn(c, w, makeStyles(mustColorscheme(opts.colorscheme)))
	}

	root.SilenceUsage = true
	root.SilenceErrors = true
	root.Version = buildVersion(opts)
	root.SetHelpFunc(helpFunc)

	if opts.manpages {
		root.AddCommand(&cobra.Command{
			Use:                   "man",
			Short:                 "Generates manpages",
			SilenceUsage:          true,
			DisableFlagsInUseLine: true,
			Hidden:                true,
			Args:                  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, _ []string) error {
				page, err := mango.NewManPage(1, cmd.Root())
				if err != nil {
					//nolint:wrapcheck
					return err
				}
				_, err = fmt.Fprint(os.Stdout, page.Build(roff.NewDocument()))
				//nolint:wrapcheck
				return err
			},
		})
	}

	if !opts.completions {
		root.CompletionOptions.DisableDefaultCmd = true
	}

	if err := root.ExecuteContext(ctx); err != nil {
		w := colorprofile.NewWriter(root.ErrOrStderr(), os.Environ())
		opts.errHandler(w, makeStyles(mustColorscheme(opts.colorscheme)), err)
		return err //nolint:wrapcheck
	}
	return nil
}

func buildVersion(opts settings) string {
	commit := opts.commit
	version := opts.version
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			version = info.Main.Version
			commit = getKey(info, "vcs.revision")
		} else {
			version = "unknown (built from source)"
		}
	}
	if len(commit) >= shaLen {
		version += " (" + commit[:shaLen] + ")"
	}
	return version
}

func getKey(info *debug.BuildInfo, key string) string {
	if info == nil {
		return ""
	}
	for _, iter := range info.Settings {
		if iter.Key == key {
			return iter.Value
		}
	}
	return ""
}
