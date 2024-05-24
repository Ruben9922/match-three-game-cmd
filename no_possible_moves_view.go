package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type noPossibleMovesViewKeyMap struct {
	EndGame key.Binding
	Confirm key.Binding
}

var noPossibleMovesViewKeys = noPossibleMovesViewKeyMap{
	EndGame: sharedKeys.EndGame,
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "continue"),
	),
}

func (s noPossibleMovesViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Confirm, s.EndGame}
}

func (s noPossibleMovesViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Confirm, s.EndGame},
	}
}

type noPossibleMovesView struct{}

func (n noPossibleMovesView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, noPossibleMovesViewKeys.EndGame):
			return showEndGameConfirmationView(m)
		case key.Matches(msg, noPossibleMovesViewKeys.Confirm):
			ensurePotentialMatch(&m.grid, m.rand)

			return showSelectFirstPointView(m)
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
