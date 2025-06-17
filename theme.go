package fang

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

// ColorScheme describes a colorscheme.
type ColorScheme struct {
	Base           color.Color
	Title          color.Color
	Codeblock      color.Color
	Program        color.Color
	DimmedArgument color.Color
	Comment        color.Color
	Flag           color.Color
	Command        color.Color
	QuotedString   color.Color
	Argument       color.Color
	Help           color.Color
	Dash           color.Color
	ErrorHeader    [2]color.Color // 0=fg 1=bg
	ErrorDetails   color.Color
}

// DefaultTheme is the default colorscheme.
func DefaultTheme(isDark bool) ColorScheme {
	c := lipgloss.LightDark(isDark)
	return ColorScheme{
		Base:           c(charmtone.Charcoal, charmtone.Ash),
		Title:          charmtone.Charple,
		Codeblock:      c(charmtone.Salt, lipgloss.Color("#2F2E36")),
		Program:        charmtone.Malibu,
		DimmedArgument: charmtone.Squid,
		Comment:        c(charmtone.Squid, lipgloss.Color("#747282")),
		Flag:           c(lipgloss.Color("#00BC82"), charmtone.Julep),
		Argument:       c(charmtone.Charcoal, charmtone.Ash),
		Command:        c(charmtone.Pony, charmtone.Dolly),
		QuotedString:   c(charmtone.Coral, charmtone.Salmon),
		ErrorHeader: [2]color.Color{
			charmtone.Butter,
			charmtone.Cherry,
		},
	}
}

// Styles represents all the styles used.
type Styles struct {
	Text        lipgloss.Style
	Title       lipgloss.Style
	Span        lipgloss.Style
	Default     lipgloss.Style
	ErrorHeader lipgloss.Style
	ErrorText   lipgloss.Style
	Codeblock   Codeblock
	Program     Program
}

// Codeblock styles.
type Codeblock struct {
	Base    lipgloss.Style
	Program Program
	Text    lipgloss.Style
	Comment lipgloss.Style
}

// Program name, args, flags, styling.
type Program struct {
	Name           lipgloss.Style
	Command        lipgloss.Style
	Flag           lipgloss.Style
	Argument       lipgloss.Style
	DimmedArgument lipgloss.Style
	QuotedString   lipgloss.Style
}

func makeStyles(cs ColorScheme) Styles {
	//nolint:mnd
	return Styles{
		Text: lipgloss.NewStyle().Foreground(cs.Base),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(cs.Title).
			Transform(strings.ToUpper).
			Padding(1, 0).
			Margin(0, 2).
			Width(width()),
		Codeblock: Codeblock{
			Base: lipgloss.NewStyle().
				Background(cs.Codeblock).
				Foreground(cs.Base).
				Margin(0, 2).
				Padding(1, 2),
			Text: lipgloss.NewStyle().
				Background(cs.Codeblock),
			Comment: lipgloss.NewStyle().
				Background(cs.Codeblock).
				Foreground(cs.Comment),
			Program: Program{
				Name: lipgloss.NewStyle().
					Background(cs.Codeblock).
					Foreground(cs.Program),
				Flag: lipgloss.NewStyle().
					PaddingLeft(1).
					Background(cs.Codeblock).
					Foreground(cs.Flag),
				Argument: lipgloss.NewStyle().
					PaddingLeft(1).
					Background(cs.Codeblock).
					Foreground(cs.Argument),
				DimmedArgument: lipgloss.NewStyle().
					Background(cs.Codeblock).
					Foreground(cs.DimmedArgument),
				Command: lipgloss.NewStyle().
					PaddingLeft(1).
					Background(cs.Codeblock).
					Foreground(cs.Command),
				QuotedString: lipgloss.NewStyle().
					PaddingLeft(1).
					Background(cs.Codeblock).
					Foreground(cs.QuotedString),
			},
		},
		Program: Program{
			Name: lipgloss.NewStyle().
				Foreground(cs.Program),
			Argument: lipgloss.NewStyle().
				Foreground(cs.Argument).
				PaddingLeft(1),
			DimmedArgument: lipgloss.NewStyle().
				Foreground(cs.DimmedArgument).
				PaddingLeft(1),
			Flag: lipgloss.NewStyle().
				Foreground(cs.Flag).
				PaddingLeft(1),
			Command: lipgloss.NewStyle().
				Foreground(cs.Command),
			QuotedString: lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(cs.QuotedString),
		},
		Span: lipgloss.NewStyle().
			Background(cs.Codeblock),
		ErrorText: lipgloss.NewStyle().
			MarginLeft(2),
		ErrorHeader: lipgloss.NewStyle().
			Foreground(cs.ErrorHeader[0]).
			Background(cs.ErrorHeader[1]).
			Bold(true).
			Padding(0, 1).
			Margin(1).
			MarginLeft(2).
			SetString("ERROR"),
	}
}
