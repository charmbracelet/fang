package fang

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// Theme describes a colorscheme.
type Theme struct {
	Codeblock    color.Color
	Program      color.Color
	Title        color.Color
	Comment      color.Color
	Command      color.Color
	QuotedString color.Color
	Argument     color.Color
	Flag         color.Color
	Help         color.Color
	Default      color.Color
	Dash         color.Color
	ErrorHeader  [2]color.Color // 0=fg 1=bg
	ErrorDetails [2]color.Color // 0=fg 1=flag
}

// DefaultThemeDark is the default colorscheme.
var DefaultThemeDark = Theme{
	Codeblock:    lipgloss.Color("#2F2E36"),
	Title:        charmtone.Charple,
	Comment:      lipgloss.Color("#747282"),
	Flag:         charmtone.Ash,
	Argument:     charmtone.Ash,
	Help:         charmtone.Ash,
	Dash:         charmtone.Smoke,
	Default:      lipgloss.Color("#747282"),
	Program:      charmtone.Blush,
	Command:      charmtone.Ash,
	QuotedString: charmtone.Julep,
	ErrorHeader: [2]color.Color{
		charmtone.Butter,
		charmtone.Cherry,
	},
	ErrorDetails: [2]color.Color{
		charmtone.Ash,
		charmtone.Julep,
	},
}

// DefaultThemeDark is the default colorscheme.
var DefaultThemeLight = Theme{
	Codeblock:    charmtone.Salt,
	Title:        charmtone.Charple,
	Comment:      charmtone.Squid,
	Flag:         charmtone.Charcoal,
	Argument:     charmtone.Charcoal,
	Help:         charmtone.Charcoal,
	Dash:         charmtone.Charcoal,
	Default:      charmtone.Squid,
	Program:      charmtone.Dolly,
	Command:      charmtone.Charcoal,
	QuotedString: lipgloss.Color("#00BC82"),
	ErrorHeader: [2]color.Color{
		charmtone.Butter,
		charmtone.Cherry,
	},
	ErrorDetails: [2]color.Color{
		charmtone.Charcoal,
		lipgloss.Color("#00BC82"),
	},
}

// Styles represents all the styles used.
type Styles struct {
	Codeblock        lipgloss.Style
	Program          lipgloss.Style
	Command          lipgloss.Style
	Comment          lipgloss.Style
	Argument         lipgloss.Style
	QuotedString     lipgloss.Style
	Flag             lipgloss.Style
	Title            lipgloss.Style
	Span             lipgloss.Style
	Dash             lipgloss.Style
	Help             lipgloss.Style
	Default          lipgloss.Style
	ErrorHeader      lipgloss.Style
	ErrorDetails     lipgloss.Style
	ErrorDetailsFlag lipgloss.Style
}

func makeStyles(theme Theme) Styles {
	//nolint:mnd
	return Styles{
		QuotedString: lipgloss.NewStyle().
			PaddingLeft(1).
			Background(theme.Codeblock).
			Foreground(theme.QuotedString),
		Codeblock: lipgloss.NewStyle().
			Background(theme.Codeblock).
			MarginLeft(2).
			Padding(1, 3, 0, 1),
		Program: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Program).
			PaddingLeft(1),
		Command: lipgloss.NewStyle().
			Foreground(theme.Command),
		Comment: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Comment).
			PaddingLeft(1),
		Argument: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Argument).
			PaddingLeft(1),
		Flag: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Flag).
			PaddingLeft(1),
		Dash: lipgloss.NewStyle().
			Background(theme.Codeblock).
			Foreground(theme.Dash).
			PaddingLeft(1),
		Span: lipgloss.NewStyle().
			Background(theme.Codeblock),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Title).
			Transform(strings.ToUpper).
			Margin(1, 0, 0, 2),
		Help: lipgloss.NewStyle().
			Foreground(theme.Help),
		Default: lipgloss.NewStyle().
			Foreground(theme.Default),
		ErrorHeader: lipgloss.NewStyle().
			Foreground(theme.ErrorHeader[0]).
			Background(theme.ErrorHeader[1]).
			Bold(true).
			Padding(0, 1).
			Margin(1).
			MarginLeft(2).
			SetString("ERROR"),
		ErrorDetails: lipgloss.NewStyle().
			Foreground(theme.ErrorDetails[0]).
			MarginLeft(2),
		ErrorDetailsFlag: lipgloss.NewStyle().
			Foreground(theme.ErrorDetails[1]).
			PaddingLeft(1),
	}
}

func (s Styles) nobg() Styles {
	return Styles{
		Codeblock: s.Codeblock.UnsetBackground(),
		Program:   s.Program.UnsetBackground(),
		Comment:   s.Comment.UnsetBackground(),
		Argument:  s.Argument.UnsetBackground(),
		Flag:      s.Flag.UnsetBackground(),
		Dash:      s.Dash.UnsetBackground(),
		Span:      s.Span.UnsetBackground(),
		Help:      s.Help,
	}
}
