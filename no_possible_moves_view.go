package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type noPossibleMovesViewKeyMap struct {
	Quit    key.Binding
	Confirm key.Binding
}

var noPossibleMovesViewKeys = noPossibleMovesViewKeyMap{
	Quit: sharedKeys.Quit,
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "continue"),
	),
}

func (s noPossibleMovesViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Confirm, s.Quit}
}

func (s noPossibleMovesViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Confirm, s.Quit},
	}
}

type noPossibleMovesView struct{}

func (n noPossibleMovesView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, noPossibleMovesViewKeys.Quit):
			return showQuitConfirmationView(m)
		case key.Matches(msg, noPossibleMovesViewKeys.Confirm):
			// todo: use nil everywhere instead of empty slice
			// Check if there is a potential match; if not, then create a new grid
			m.potentialMatch = findPotentialMatch(m.grid)
			for len(m.potentialMatch) == 0 {
				// Check if there are any possible matches; if no possible matches then create a new grid
				m.grid = newGrid(m.rand)

				// Refresh the grid to remove all matches
				finished := false
				for !finished {
					finished = refreshGrid(&m.grid, m.rand, &m.score, false)
				}

				m.potentialMatch = findPotentialMatch(m.grid)
			}

			m.view = refreshGridView{}
			m.point1 = emptyVector2d
			m.point2 = emptyVector2d

			return m, tickCmd()
		}
	}

	return m, nil
}

func (n noPossibleMovesView) draw(m model) string {
	text := fmt.Sprintf("No more possible moves\n\nPress %s to generate a new grid...",
		lipgloss.NewStyle().Bold(true).Render(noPossibleMovesViewKeys.Confirm.Help().Key))
	helpView := m.help.View(noPossibleMovesViewKeys)
	noMorePossibleMovesText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)

	gridText := createGrid(m, []vector2d{m.point1, m.point2})

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		gridText,
		lipgloss.NewStyle().MarginLeft(3).Render(noMorePossibleMovesText),
	)
}
