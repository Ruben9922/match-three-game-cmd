package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"slices"
	"strings"
)

type radioButtonItem interface {
	comparable
	String() string
}

func drawRadioButtons[T radioButtonItem](options []T, selected T, label string, key key.Binding) string {
	var builder strings.Builder
	builder.WriteString(label)
	builder.WriteString(":  ")
	for _, option := range options {
		var style lipgloss.Style
		if option == selected {
			style = highlightedStyle
		} else {
			style = lipgloss.NewStyle()
		}
		builder.WriteString(style.Render(option.String()))

		builder.WriteString("  ")
	}

	keyString := key.Help().Key
	styledKeyString := lipgloss.NewStyle().Inherit(secondaryTextStyle).Bold(true).Render(keyString)
	// Couldn't get styling to work correctly with `fmt.Sprintf`, hence styling each substring separately then
	// concatenating
	keyDescription := secondaryTextStyle.Render("(press ") + styledKeyString + secondaryTextStyle.Render(" to change)")
	builder.WriteString(keyDescription)

	return builder.String()
}

func createGrid(m model, selectedPoints []vector2d) string {
	var stringBuilder strings.Builder
	for y, row := range m.grid {
		for x, symbol := range row {
			point := vector2d{x: x, y: y}

			var formattedSymbol string
			if slices.Contains(selectedPoints, point) {
				formattedSymbol = m.symbolSet.formatSymbolHighlighted(symbol)
			} else {
				formattedSymbol = m.symbolSet.formatSymbol(symbol)
			}

			stringBuilder.WriteString(formattedSymbol)

			if x != len(row)-1 {
				stringBuilder.WriteString(" ")
			}
		}

		if y != len(m.grid)-1 {
			stringBuilder.WriteString("\n")
		}
	}
	border := lipgloss.RoundedBorder()
	gridStyle := lipgloss.NewStyle().
		BorderForeground(accentColor).
		BorderStyle(border).
		Padding(0, 1)

	gridString := gridStyle.Render(stringBuilder.String())

	scoreString := fmt.Sprintf("Score: %s", humanize.Comma(int64(m.score)))

	var remainingMovesString string
	if m.options.gameType == LimitedMoves {
		remainingMovesString = fmt.Sprintf("Remaining moves: %d", m.remainingMoveCount)
	} else {
		remainingMovesString = ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, gridString, "", scoreString, remainingMovesString)
}
