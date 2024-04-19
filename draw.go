package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"
	"github.com/dustin/go-humanize/english"
	"slices"
	"strings"
)

func createTitleView(m model) string {
	const titlePart1 = "  __  __       _       _       _____ _                   \n |  \\/  | __ _| |_ ___| |__   |_   _| |__  _ __ ___  ___ \n | |\\/| |/ _` | __/ __| '_ \\    | | | '_ \\| '__/ _ \\/ _ \\\n | |  | | (_| | || (__| | | |   | | | | | | | |  __/  __/\n |_|  |_|\\__,_|\\__\\___|_| |_|   |_| |_| |_|_|  \\___|\\___|"
	const titlePart2 = "   ____                      \n  / ___| __ _ _ __ ___   ___ \n | |  _ / _` | '_ ` _ \\ / _ \\\n | |_| | (_| | | | | | |  __/\n  \\____|\\__,_|_| |_| |_|\\___|"
	const text = "\n Press any key to start..."

	radioButtons := drawRadioButtons([]gameType{Endless, LimitedMoves}, m.options.gameType, "Game type", "T")

	titleView := lipgloss.JoinVertical(lipgloss.Left, titlePart1, titlePart2, text, radioButtons)

	controls := []control{
		{key: "t", description: "change game type"},
		{key: "any other key", description: "start"},
	}
	controlsView := controlsToString(controls)
	return lipgloss.JoinVertical(lipgloss.Left, titleView, controlsView)
}

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

func controlsToString(controls []control) string {
	controlStrings := make([]string, 0, len(controls))
	for _, c := range controls {
		controlString := fmt.Sprintf("%s: %s", c.key, c.description)
		controlStrings = append(controlStrings, controlString)
	}
	controlsString := strings.Join(controlStrings, " • ")
	return controlsString
}

func createSelectFirstPointView(m model) string {
	const text = "Select two points to swap (selecting point 1)..."
	controls := []control{
		{key: "← ↑ → ↓ / WASD", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "H", description: "Show hint"},
		{key: "Q", description: "End Game"},
	}
	hintControls := []control{
		{key: "<Any key>", description: "Hide hint"},
	}

	var controlsString string
	var selectedPoints []vector2d
	if m.showHint {
		controlsString = controlsToString(hintControls)
		selectedPoints = m.potentialMatch
	} else {
		controlsString = controlsToString(controls)
		selectedPoints = []vector2d{m.point1}
	}
	selectFirstPointText := lipgloss.JoinVertical(lipgloss.Left, text, controlsString)

	gridText := createGrid(m, selectedPoints)

	return lipgloss.JoinHorizontal(lipgloss.Top, gridText, selectFirstPointText)
}

func createSelectSecondPointView(m model) string {
	const text = "Select two points to swap (selecting point 2)..."
	controls := []control{
		{key: "← ↑ → ↓ / WASD", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "Escape", description: "Cancel selection"},
		{key: "Q", description: "End Game"},
	}
	controlsString := controlsToString(controls)
	selectSecondPointText := lipgloss.JoinVertical(lipgloss.Left, text, controlsString)

	gridText := createGrid(m, []vector2d{m.point1, m.point2})

	return lipgloss.JoinHorizontal(lipgloss.Top, gridText, selectSecondPointText)
}

func createSelectPointConfirmationView(m model) string {
	controls := []control{{key: "<Any key>", description: "Continue"}}
	controlsString := controlsToString(controls)

	matches := findMatches(m.grid)
	var text string
	var selectedPoints []vector2d
	if len(matches) != 0 {
		text = fmt.Sprintf("Swapped %c (%d, %d) and %c (%d, %d); %s formed",
			m.grid[m.point1.y][m.point1.x], m.point1.x, m.point1.y, m.grid[m.point2.y][m.point2.x], m.point2.x,
			m.point2.y, english.PluralWord(len(matches), "match", ""))
		selectedPoints = convertMatchesToPoints(matches)
	} else {
		text = "Not swapping as swap would not result in a match; please try again"
		selectedPoints = []vector2d{m.point1, m.point2}
	}
	selectPointConfirmationText := lipgloss.JoinVertical(lipgloss.Left, text, controlsString)

	gridText := createGrid(m, selectedPoints)

	return lipgloss.JoinHorizontal(lipgloss.Top, gridText, selectPointConfirmationText)
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

	controls := []control{{key: "any key", description: "skip"}}
	controlsString := controlsToString(controls)

	scoreString := fmt.Sprintf("Score: %s", humanize.Comma(int64(m.score)))

	var remainingMovesString string
	if m.options.gameType == LimitedMoves {
		remainingMovesString = fmt.Sprintf("Remaining moves: %d", m.remainingMoveCount)
	} else {
		remainingMovesString = ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, gridString, controlsString, scoreString, remainingMovesString)
}

func createGameOverView() string {
	const text = "Game over!\n\nNo more moves left."
	controls := []control{{key: "<Any key>", description: "Exit"}}
	controlsString := controlsToString(controls)
	return lipgloss.JoinVertical(lipgloss.Left, text, controlsString)
}

func createQuitConfirmationView() string {
	const text = "Are you sure you want to quit?\n\nAny game progress will be lost."
	controls := []control{
		{key: "enter", description: "quit"},
		{key: "any other key", description: "cancel"},
	}
	controlsString := controlsToString(controls)

	return lipgloss.JoinVertical(lipgloss.Left, text, controlsString)
}
