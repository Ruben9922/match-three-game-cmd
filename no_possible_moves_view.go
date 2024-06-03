package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func showNoPossibleMovesView(m model) (tea.Model, tea.Cmd) {
	m.view = newNoPossibleMovesView()
	m.help.ShowAll = false

	return m, nil
}

type noPossibleMovesViewKeyMap struct {
	EndGame key.Binding
	Confirm key.Binding
}

func newNoPossibleMovesViewKeys() noPossibleMovesViewKeyMap {
	return noPossibleMovesViewKeyMap{
		EndGame: newEndGameKeyBinding(),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "continue"),
		),
	}
}

func (s noPossibleMovesViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Confirm, s.EndGame}
}

func (s noPossibleMovesViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Confirm, s.EndGame},
	}
}

type noPossibleMovesView struct {
	keys noPossibleMovesViewKeyMap
}

func newNoPossibleMovesView() noPossibleMovesView {
	return noPossibleMovesView{
		keys: newNoPossibleMovesViewKeys(),
	}
}

func (n noPossibleMovesView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, n.keys.EndGame):
			return showEndGameConfirmationView(m)
		case key.Matches(msg, n.keys.Confirm):
			ensurePotentialMatch(&m.grid, m.rand)

			return showSelectFirstPointView(m)
		}
	}

	return m, nil
}

func (n noPossibleMovesView) draw(m model) string {
	text := fmt.Sprintf("No more possible moves\n\nPress %s to generate a new grid...",
		lipgloss.NewStyle().Bold(true).Render(n.keys.Confirm.Help().Key))
	gridText := drawGrid(m, []vector2d{m.point1, m.point2})
	m.help.Width = m.windowSize.x - lipgloss.Width(gridText) - 8 - 3
	helpView := m.help.View(n.keys)
	noMorePossibleMovesText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)
	gridLayoutText := drawGridLayout(m, gridText, noMorePossibleMovesText)

	return gridLayoutText
}
