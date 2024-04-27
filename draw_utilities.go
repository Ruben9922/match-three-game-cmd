package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"slices"
	"strings"
)

type radioButtonItem interface {
	comparable
	String() string
}

func drawRadioButtons[T radioButtonItem](options []T, selected T, label string, key string) string {
	var builder strings.Builder
	builder.WriteString(label)
	builder.WriteString(": ")
	for i, option := range options {
		if option == selected {
			builder.WriteString(option.String() + " [▪]")
		} else {
			builder.WriteString(option.String() + " [ ]")
		}

		if i != len(options)-1 {
			builder.WriteString(";")
		}

		builder.WriteString(" ")
	}
	builder.WriteString(fmt.Sprintf("(press %s)", strings.ToUpper(key)))

	return builder.String()
}

func createGrid(m model, selectedPoints []vector2d) string {
	var stringBuilder strings.Builder
	for y, row := range m.grid {
		for x, symbol := range row {
			point := vector2d{x: x, y: y}

			var style lipgloss.Style
			if slices.Contains(selectedPoints, point) {
				style = symbolHighlightedColors[symbol]
			} else {
				style = symbolColors[symbol]
			}

			stringBuilder.WriteString(style.Render(string(symbol)))

			if x != len(row)-1 {
				stringBuilder.WriteString(" ")
			}
		}

		if y != len(m.grid)-1 {
			stringBuilder.WriteString("\n")
		}
	}
	gridString := stringBuilder.String()

	scoreString := fmt.Sprintf("Score: %s", humanize.Comma(int64(m.score)))

	var remainingMovesString string
	if m.options.gameType == LimitedMoves {
		remainingMovesString = fmt.Sprintf("Remaining moves: %d", m.remainingMoveCount)
	} else {
		remainingMovesString = ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, gridString, scoreString, remainingMovesString)
}
