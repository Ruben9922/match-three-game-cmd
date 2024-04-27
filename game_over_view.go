package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameOverViewKeyMap struct {
	Confirm key.Binding
}

var gameOverViewKeys = gameOverViewKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "quit"),
	),
}

func (s gameOverViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Confirm}
}

func (s gameOverViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Confirm},
	}
}

type gameOverView struct{}

func (g gameOverView) update(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	if key.Matches(msg, gameOverViewKeys.Confirm) {
		return m, tea.Quit
	}

	return m, nil
}

func (g gameOverView) draw(m model) string {
	const text = "Game over!\n\nNo more moves left."
	helpView := m.help.View(gameOverViewKeys)
	return lipgloss.JoinVertical(lipgloss.Left, text, helpView)
}
