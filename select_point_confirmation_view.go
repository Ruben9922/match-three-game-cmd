package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize/english"
)

type selectPointConfirmationViewKeyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}

var selectPointConfirmationViewKeys = selectPointConfirmationViewKeyMap{
	Quit: sharedKeys.Quit,
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "continue"),
	),
}

func (s selectPointConfirmationViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Confirm, s.Quit}
}

func (s selectPointConfirmationViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Confirm, s.Quit},
	}
}

type selectPointConfirmationView struct{}

func (s selectPointConfirmationView) update(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, selectPointConfirmationViewKeys.Quit):
		return showQuitConfirmationView(m)
	case key.Matches(msg, selectPointConfirmationViewKeys.Confirm):
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
	default:
		return m, nil
	}
}

func (s selectPointConfirmationView) draw(m model) string {
	matches := findMatches(m.grid)
	var text string
	var selectedPoints []vector2d
	if len(matches) != 0 {
		matchesScore := computeScore(matches)
		symbol1 := m.grid[m.point1.y][m.point1.x]
		symbol2 := m.grid[m.point2.y][m.point2.x]
		swappedText := fmt.Sprintf("Swapped %s (%d, %d) and %s (%d, %d).",
			m.symbolSet.formatSymbol(symbol1), m.point1.x, m.point1.y, m.symbolSet.formatSymbol(symbol2), m.point2.x,
			m.point2.y)
		matchText := fmt.Sprintf("%s formed!", english.PluralWord(len(matches), "Match", ""))
		pointsGainedText := fmt.Sprintf("+%d points!", matchesScore)
		text = lipgloss.JoinVertical(lipgloss.Left, swappedText, "", matchText, pointsGainedText)

		selectedPoints = convertMatchesToPoints(matches)
	} else {
		text = "Not swapping as swap would not result in a match.\nPlease try again."
		selectedPoints = []vector2d{m.point1, m.point2}
	}
	helpView := m.help.View(selectPointConfirmationViewKeys)
	selectPointConfirmationText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)

	gridText := createGrid(m, selectedPoints)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		gridText,
		lipgloss.NewStyle().MarginLeft(3).Render(selectPointConfirmationText),
	)
}
