package fang

import (
	"cmp"
	"fmt"
	"io"
	"maps"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	minSpace = 10
	shortPad = 2
	longPad  = 4
)

var width = sync.OnceValue(func() int {
	if s := os.Getenv("__FANG_TEST_WIDTH"); s != "" {
		w, _ := strconv.Atoi(s)
		return w
	}
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 120
	}
	return min(w, 120)
})

func helpFn(c *cobra.Command, w *colorprofile.Writer, styles Styles) {
	writeLongShort(w, styles, cmp.Or(c.Long, c.Short))
	usage := styleUsage(c, styles.Codeblock.Program, true)
	examples := styleExamples(c, styles)

	padding := styles.Codeblock.Base.GetHorizontalPadding()
	blockWidth := lipgloss.Width(usage)
	for _, ex := range examples {
		blockWidth = max(blockWidth, lipgloss.Width(ex))
	}
	blockWidth = min(width()-padding, blockWidth+padding)

	styles.Codeblock.Base = styles.Codeblock.Base.Width(blockWidth)

	_, _ = fmt.Fprintln(w, styles.Title.Render("usage"))
	_, _ = fmt.Fprintln(w, styles.Codeblock.Base.Render(usage))
	if len(examples) > 0 {
		cw := styles.Codeblock.Base.GetWidth() - styles.Codeblock.Base.GetHorizontalPadding()
		_, _ = fmt.Fprintln(w, styles.Title.Render("examples"))
		for i, example := range examples {
			if lipgloss.Width(example) > cw {
				examples[i] = ansi.Truncate(example, cw, "…")
			}
		}
		_, _ = fmt.Fprintln(w, styles.Codeblock.Base.Render(strings.Join(examples, "\n")))
	}

	cmds, cmdKeys := evalCmds(c, styles)
	flags, flagKeys := evalFlags(c, styles)
	space := calculateSpace(cmdKeys, flagKeys)

	groups := getCmdGroupNames(c)
	// render default group first
	renderCommandGroup(w, styles, space, "commands", cmds[""])
	delete(cmds, "")
	groupIDs := slices.Collect(maps.Keys(groups))
	slices.Sort(groupIDs)
	for _, v := range groupIDs {
		renderCommandGroup(w, styles, space, groups[v], cmds[v])
	}

	if len(flags) > 0 {
		_, _ = fmt.Fprintln(w, styles.Title.Render("flags"))
		for _, k := range flagKeys {
			_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.NewStyle().PaddingLeft(longPad).Render(k),
				strings.Repeat(" ", space-lipgloss.Width(k)),
				flags[k],
			))
		}
	}

	_, _ = fmt.Fprintln(w)
}

// DefaultErrorHandler is the default [ErrorHandler] implementation.
func DefaultErrorHandler(w io.Writer, styles Styles, err error) {
	_, _ = fmt.Fprintln(w, styles.ErrorHeader.String())
	_, _ = fmt.Fprintln(w, styles.ErrorText.Render(err.Error()+"."))
	_, _ = fmt.Fprintln(w)
	if isUsageError(err) {
		_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.ErrorText.UnsetWidth().Render("Try"),
			styles.Program.Flag.Render("--help"),
			styles.ErrorText.UnsetWidth().UnsetMargins().UnsetTransform().PaddingLeft(1).Render("for usage."),
		))
		_, _ = fmt.Fprintln(w)
	}
}

// XXX: this is a hack to detect usage errors.
// See: https://github.com/spf13/cobra/pull/2266
func isUsageError(err error) bool {
	s := err.Error()
	for _, prefix := range []string{
		"flag needs an argument:",
		"unknown flag:",
		"unknown shorthand flag:",
		"unknown command",
		"invalid argument",
	} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func writeLongShort(w *colorprofile.Writer, styles Styles, longShort string) {
	if longShort == "" {
		return
	}
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, styles.Text.Width(width()).PaddingLeft(shortPad).Render(longShort))
}

var otherArgsRe = regexp.MustCompile(`(\[.*\])`)

// styleUsage stylized styleUsage line for a given command.
func styleUsage(c *cobra.Command, styles Program, complete bool) string {
	// XXX: maybe use c.UseLine() here?
	u := c.Use
	hasArgs := strings.Contains(u, "[args]")
	hasFlags := strings.Contains(u, "[flags]") || strings.Contains(u, "[--flags]") || c.HasFlags() || c.HasPersistentFlags() || c.HasAvailableFlags()
	hasCommands := strings.Contains(u, "[command]") || c.HasAvailableSubCommands()
	for _, k := range []string{
		"[args]",
		"[flags]", "[--flags]",
		"[command]",
	} {
		u = strings.ReplaceAll(u, k, "")
	}

	var otherArgs []string //nolint:prealloc
	for _, arg := range otherArgsRe.FindAllString(u, -1) {
		u = strings.ReplaceAll(u, arg, "")
		otherArgs = append(otherArgs, arg)
	}

	u = strings.TrimSpace(u)

	useLine := []string{
		styles.Name.Render(u),
	}
	if !complete {
		useLine[0] = styles.Command.Render(u)
	}
	if hasCommands {
		useLine = append(
			useLine,
			styles.DimmedArgument.Render("[command]"),
		)
	}
	if hasArgs {
		useLine = append(
			useLine,
			styles.DimmedArgument.Render("[args]"),
		)
	}
	for _, arg := range otherArgs {
		useLine = append(
			useLine,
			styles.DimmedArgument.Render(arg),
		)
	}
	if hasFlags {
		useLine = append(
			useLine,
			styles.DimmedArgument.Render("[--flags]"),
		)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, useLine...)
}

// styleExamples for a given command.
// will print both the cmd.Use and cmd.Example bits.
func styleExamples(c *cobra.Command, styles Styles) []string {
	if c.Example == "" {
		return nil
	}
	usage := []string{}
	examples := strings.Split(c.Example, "\n")
	var indent bool
	for i, line := range examples {
		line = strings.TrimSpace(line)
		if (i == 0 || i == len(examples)-1) && line == "" {
			continue
		}
		s := styleExample(c, line, indent, styles.Codeblock)
		usage = append(usage, s)
		indent = len(line) > 1 && line[len(line)-1] == '\\'
	}

	return usage
}

func styleExample(c *cobra.Command, line string, indent bool, styles Codeblock) string {
	if strings.HasPrefix(line, "# ") {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.Comment.Render(line),
		)
	}

	padding := func() int {
		if !indent {
			return 0
		}
		indent = false
		return 2
	}

	var isQuotedString bool
	var foundProgramName bool
	programName := c.Name()
	args := strings.Fields(line)
	for i, arg := range args {
		isQuoteStart := arg[0] == '"'
		isQuoteEnd := arg[len(arg)-1] == '"'
		isFlagStart := arg[0] == '-'

		if i == len(args)-1 && len(arg) == 1 && arg[0] == '\\' {
			args[i] = styles.Program.DimmedArgument.UnsetPadding().Render(arg)
			continue
		}

		if !foundProgramName { //nolint:nestif
			if isQuotedString {
				args[i] = styles.Program.QuotedString.PaddingRight(1).Render(arg)
				isQuotedString = !isQuoteEnd
				continue
			}
			if left, right, ok := strings.Cut(arg, "="); ok {
				args[i] = styles.Program.Flag.UnsetPadding().PaddingLeft(padding()).Render(left + "=")
				if right[0] == '"' {
					isQuotedString = true
					args[i] += styles.Program.QuotedString.UnsetPadding().Render(right)
					continue
				}
				args[i] += styles.Program.Argument.UnsetPadding().PaddingRight(1).Render(right)
				continue
			}

			if arg == programName {
				args[i] = styles.Program.Name.PaddingLeft(padding()).Render(arg)
				foundProgramName = true
				continue
			}
		}

		if !isQuoteStart && !isFlagStart && isSubCommand(c, arg) {
			args[i] = styles.Program.Command.Render(arg)
			continue
		}
		isQuotedString = isQuotedString || isQuoteStart
		if isQuotedString {
			args[i] = styles.Program.QuotedString.Render(arg)
			isQuotedString = !isQuoteEnd
			continue
		}
		// handle a flag
		if arg[0] == '-' {
			name, value, ok := strings.Cut(arg, "=")
			// it is --flag=value
			if ok {
				args[i] = lipgloss.JoinHorizontal(
					lipgloss.Left,
					styles.Program.Flag.Render(name+"="),
					styles.Program.Argument.UnsetPadding().Render(value),
				)
				continue
			}
			// it is either --bool-flag or --flag value
			args[i] = lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.Program.Flag.Render(name),
			)
			continue
		}
		args[i] = styles.Program.Argument.Render(arg)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		args...,
	)
}

func renderCommandGroup(
	w io.Writer,
	styles Styles,
	space int,
	name string,
	help map[string]string,
) {
	if len(help) == 0 {
		return
	}
	_, _ = fmt.Fprintln(w, styles.Title.Render(name))
	keys := slices.Collect(maps.Keys(help))
	slices.Sort(keys)
	for _, k := range keys {
		_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().PaddingLeft(longPad).Render(k),
			strings.Repeat(" ", space-lipgloss.Width(k)),
			help[k],
		))
	}
}

func evalFlags(c *cobra.Command, styles Styles) (map[string]string, []string) {
	flags := map[string]string{}
	keys := []string{}
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		var parts []string
		if f.Shorthand == "" {
			parts = append(
				parts,
				styles.Program.Flag.UnsetPadding().Render("--"+f.Name),
			)
		} else {
			parts = append(
				parts,
				styles.Program.Flag.UnsetPadding().Render("-"+f.Shorthand),
				styles.Program.Flag.Render("--"+f.Name),
			)
		}
		key := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
		help := styles.FlagDescription.Render(f.Usage)
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			help = lipgloss.JoinHorizontal(
				lipgloss.Left,
				help,
				styles.FlagDefault.Render("("+f.DefValue+")"),
			)
		}
		flags[key] = help
		keys = append(keys, key)
	})
	return flags, keys
}

type commandsHelp map[string]string

func evalCmds(c *cobra.Command, styles Styles) (map[string]commandsHelp, []string) {
	padStyle := lipgloss.NewStyle().PaddingLeft(0) //nolint:mnd
	keys := []string{}
	cmds := map[string]commandsHelp{}
	for groupID, scs := range getCmdGroups(c) {
		group := map[string]string{}
		for _, sc := range scs {
			if sc.Hidden {
				continue
			}
			key := padStyle.Render(styleUsage(sc, styles.Program, false))
			help := styles.FlagDescription.Render(sc.Short)
			group[key] = help
			keys = append(keys, key)
		}
		cmds[groupID] = group
	}
	return cmds, keys
}

func getCmdGroupNames(c *cobra.Command) map[string]string {
	result := map[string]string{}
	for _, g := range c.Groups() {
		result[g.ID] = g.Title
	}
	return result
}

func getCmdGroups(c *cobra.Command) map[string][]*cobra.Command {
	result := map[string][]*cobra.Command{}
	for _, sc := range c.Commands() {
		result[sc.GroupID] = append(result[sc.GroupID], sc)
	}
	return result
}

func calculateSpace(k1, k2 []string) int {
	const spaceBetween = 2
	space := minSpace
	for _, k := range append(k1, k2...) {
		space = max(space, lipgloss.Width(k)+spaceBetween)
	}
	return space
}

func isSubCommand(c *cobra.Command, arg string) bool {
	cmd, _, _ := c.Traverse([]string{arg})
	return cmd != nil && cmd.Name() == arg
}
