package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize/english"
)

type selectPointConfirmationViewKeyMap struct {
	sharedKeyMap
	Select key.Binding
}

var selectPointConfirmationViewKeys = selectPointConfirmationViewKeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "continue"),
	),
}

type selectPointConfirmationView struct{}

func (s selectPointConfirmationView) update(msg tea.KeyMsg, m model) (model, tea.Cmd) {
	// Refresh grid
	//m.view = RefreshGridView

	refreshGrid(&m.grid, m.rand, &m.score, true)

	// todo: use nil everywhere instead of empty slice
	m.potentialMatch = make([]vector2d, 0)
	for len(m.potentialMatch) == 0 {
		// Check if there are any possible matches; if no possible matches then create a new grid
		m.potentialMatch = findPotentialMatch(m.grid)
		if len(m.potentialMatch) == 0 {
			m.grid = newGrid(m.rand)
			refreshGrid(&m.grid, m.rand, &m.score, true)
		}
	}

	if m.options.gameType != LimitedMoves || m.remainingMoveCount > 0 {
		m.view = selectFirstPointView{}
		m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
		m.point2 = emptyVector2d
	} else {
		m.view = gameOverView{}
		m.point1 = emptyVector2d
		m.point2 = emptyVector2d
	}

	return m, nil
}

func (s selectPointConfirmationView) draw(m model) string {
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