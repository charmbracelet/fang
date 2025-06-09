package fang

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	minSpace = 10
	shortPad = 2
)

var width = sync.OnceValue(func() int {
	if s := os.Getenv("__FANG_TEST_WIDTH"); s != "" {
		w, _ := strconv.Atoi(s)
		return w
	}
	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return 80
	}
	return min(w, 80)
})

func helpFn(c *cobra.Command, w *colorprofile.Writer, styles Styles) {
	writeLongShort(w, styles, cmp.Or(c.Long, c.Short))
	firstUse := use(c, styles)
	_, _ = fmt.Fprintln(w, firstUse)
	if lines := usage(c, styles, lipgloss.Width(firstUse)); len(lines) > 0 {
		_, _ = fmt.Fprintln(w, styles.Title.Render("examples\n"))
		_, _ = fmt.Fprintln(
			w,
			styles.Codeblock.Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					lines[1:]...,
				),
			),
		)
	}

	cmds, cmdKeys := evalCmds(c, styles.nobg())
	flags, flagKeys := evalFlags(c, styles.nobg())
	space := calculateSpace(cmdKeys, flagKeys)

	if len(cmds) > 0 {
		_, _ = fmt.Fprintln(w, styles.Title.Render("commands\n"))
		for _, k := range cmdKeys {
			_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				k,
				strings.Repeat(" ", space-lipgloss.Width(k)),
				cmds[k],
			))
		}
	}

	if len(flags) > 0 {
		_, _ = fmt.Fprintln(w, styles.Title.Render("flags\n"))
		for _, k := range flagKeys {
			_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
				lipgloss.Left,
				k,
				strings.Repeat(" ", space-lipgloss.Width(k)),
				flags[k],
			))
		}
	}

	_, _ = fmt.Fprintln(w)
}

func writeError(w *colorprofile.Writer, styles Styles, err error) {
	_, _ = fmt.Fprintln(w, styles.ErrorHeader.String())
	_, _ = fmt.Fprintln(w, styles.ErrorDetails.Render(titleFirstWord(err.Error()+".")))
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.ErrorDetails.Render("Try"),
		styles.ErrorDetailsFlag.Render("--help"),
		styles.ErrorDetails.UnsetMargins().PaddingLeft(1).Render("for usage."),
	))
	_, _ = fmt.Fprintln(w)
}

func writeLongShort(w *colorprofile.Writer, styles Styles, longShort string) {
	if longShort == "" {
		return
	}
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, styles.Help.Width(width()).PaddingLeft(shortPad).Render(longShort))
	_, _ = fmt.Fprintln(w, styles.Title.Render("usage"))
	_, _ = fmt.Fprintln(w)
}

var otherArgsRe = regexp.MustCompile(`(\[.*\])`)

// use stylized use line for a given command.
func use(c *cobra.Command, styles Styles) string {
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
		styles.Program.
			UnsetBackground().
			PaddingLeft(4).
			Render(u),
	}
	if hasCommands {
		useLine = append(
			useLine,
			styles.Argument.
				UnsetBackground().
				Render("[command]"),
		)
	}
	if hasArgs {
		useLine = append(
			useLine,
			styles.Argument.
				UnsetBackground().
				Render("[args]"),
		)
	}
	for _, arg := range otherArgs {
		useLine = append(
			useLine,
			styles.Argument.
				UnsetBackground().
				Render(arg),
		)
	}
	if hasFlags {
		useLine = append(
			useLine,
			styles.Flag.
				UnsetBackground().
				Render("[--flags]"),
		)
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, useLine...)
}

// usage for a given command.
// will print both the cmd.Use and cmd.Example bits.
func usage(c *cobra.Command, styles Styles, minSize int) []string {
	usage := []string{}
	size := minSize
	examples := strings.Split(c.Example, "\n")
	for i, line := range examples {
		s := evalExample(c, line, i != len(examples)-1, styles)
		size = max(size, lipgloss.Width(s))
		usage = append(usage, s)
	}

	for i, use := range usage {
		if usize := lipgloss.Width(use); size > usize {
			usage[i] = lipgloss.JoinHorizontal(
				lipgloss.Left,
				use,
				styles.Span.Render(strings.Repeat(" ", size-usize)),
			)
		}
	}
	return usage
}

func evalExample(c *cobra.Command, line string, last bool, styles Styles) string {
	line = strings.TrimSpace(line)
	if line == "" && !last {
		return ""
	}

	if strings.HasPrefix(line, "# ") {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.Comment.Render(line),
		)
	}

	args := strings.Fields(line)
	var nextIsFlag bool
	var isQuotedString bool
	for i, arg := range args {
		if i == 0 {
			args[i] = styles.Program.Render(arg)
			continue
		}
		if nextIsFlag {
			args[i] = styles.Flag.Render(arg)
			nextIsFlag = false
			continue
		}
		if strings.HasPrefix(arg, `"`) {
			isQuotedString = true
		}
		if isQuotedString {
			args[i] = styles.QuotedString.Render(arg)
			continue
		}
		if strings.HasSuffix(arg, `"`) {
			isQuotedString = false
			continue
		}
		var dashes string
		if strings.HasPrefix(arg, "-") {
			dashes = "-"
		}
		if strings.HasPrefix(arg, "--") {
			dashes = "--"
		}
		// handle a flag
		if dashes != "" {
			name, value, ok := strings.Cut(arg, "=")
			name = strings.TrimPrefix(name, dashes)
			// it is --flag=value
			if ok {
				args[i] = lipgloss.JoinHorizontal(
					lipgloss.Left,
					styles.Dash.Render(dashes),
					styles.Flag.UnsetPadding().Render(name),
					styles.Comment.UnsetPadding().Render("="),
					styles.Flag.UnsetPadding().Render(value),
				)
				continue
			}
			// it is either --bool-flag or --flag value
			args[i] = lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.Dash.Render(dashes),
				styles.Flag.UnsetPadding().Render(name),
			)
			// if the flag is not a bool flag, next arg continues current flag
			nextIsFlag = !isFlagBool(c, name)
			continue
		}
		args[i] = styles.Argument.Render(arg)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		args...,
	)
}

func evalFlags(c *cobra.Command, styles Styles) (map[string]string, []string) {
	const shortPad = 4
	const noShortPad = shortPad + 3
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
				styles.Dash.PaddingLeft(noShortPad).Render("--"),
				styles.Flag.UnsetPadding().Render(f.Name),
			)
		} else {
			parts = append(
				parts,
				styles.Dash.PaddingLeft(shortPad).Render("-"),
				styles.Flag.UnsetPadding().Render(f.Shorthand),
				styles.Dash.Render("--"),
				styles.Flag.UnsetPadding().Render(f.Name),
			)
		}
		key := lipgloss.JoinHorizontal(lipgloss.Left, parts...)
		help := styles.Help.Render(titleFirstWord(f.Usage))
		if f.DefValue != "" && f.DefValue != "false" && f.DefValue != "0" && f.DefValue != "[]" {
			help = lipgloss.JoinHorizontal(
				lipgloss.Left,
				help,
				styles.Default.PaddingLeft(1).Render("("+f.DefValue+")"),
			)
		}
		flags[key] = help
		keys = append(keys, key)
	})
	return flags, keys
}

func evalCmds(c *cobra.Command, styles Styles) (map[string]string, []string) {
	padStyle := lipgloss.NewStyle().PaddingLeft(0) //nolint:mnd
	keys := []string{}
	cmds := map[string]string{}
	for _, sc := range c.Commands() {
		if sc.Hidden {
			continue
		}
		key := padStyle.Render(use(sc, styles))
		help := styles.Help.Render(titleFirstWord(sc.Short))
		cmds[key] = help
		keys = append(keys, key)
	}
	return cmds, keys
}

func calculateSpace(k1, k2 []string) int {
	const spaceBetween = 2
	space := minSpace
	for _, k := range append(k1, k2...) {
		space = max(space, lipgloss.Width(k)+spaceBetween)
	}
	return space
}

func isFlagBool(c *cobra.Command, name string) bool {
	flag := c.Flags().Lookup(name)
	if flag == nil && len(name) == 1 {
		flag = c.Flags().ShorthandLookup(name)
	}
	if flag == nil {
		return false
	}
	return flag.Value.Type() == "bool"
}

func titleFirstWord(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}
	words[0] = cases.Title(language.AmericanEnglish).String(words[0])
	return strings.Join(words, " ")
}
