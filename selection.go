package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize/english"
)

//func swapPoints(s tcell.Screen, g *grid, potentialMatch []vector2d, score int) bool {
//	point1 := vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
//	point2 := emptyVector2d
//	for point2 == emptyVector2d {
//		point1 = selectFirstPoint(s, *g, potentialMatch, point1, score)
//		point2 = selectSecondPoint(s, *g, point1, score)
//	}
//
//	if len(matches) != 0 {
//	} else {
//	}
//
//	waitForKeyPress(s)
//
//	return len(matches) != 0
//}

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
		controlsString = controlsToString(controls)
		selectedPoints = m.potentialMatch
	} else {
		controlsString = controlsToString(hintControls)
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
