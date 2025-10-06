package tui

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
)

var (
	subtleColor    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	subtle    = lipgloss.NewStyle().Foreground(subtleColor)
	highlight = lipgloss.NewStyle().Foreground(highlightColor)
	focus     = highlight
	normal    = lipgloss.NewStyle()

	footer = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderForeground(subtleColor)

	flexBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	flexCellHeader = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(subtleColor).
			MarginRight(2)

	line = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(subtleColor)

	flexCell = lipgloss.NewStyle().
			Border(flexBorder, true).
			BorderForeground(highlightColor).
			Padding(0, 1)

	flexCellFunc = func(cell *flexbox.Cell) lipgloss.Style {
		return lipgloss.NewStyle().
			Width(cell.GetWidth()-6).
			Height(cell.GetHeight()-6).
			Padding(1, 1).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
	}
)
