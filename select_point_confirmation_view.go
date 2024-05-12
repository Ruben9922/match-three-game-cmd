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

func (s selectPointConfirmationView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, selectPointConfirmationViewKeys.Quit):
			return showQuitConfirmationView(m)
		case key.Matches(msg, selectPointConfirmationViewKeys.Confirm):
			matches := findMatches(m.grid)
			if len(matches) == 0 {
				return returnToSelectFirstPointView(m)
			} else {
				return showRefreshGridView(m)
			}
		}
	}

	return m, nil
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

		selectedPoints = flatten(matches)
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
