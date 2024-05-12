package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func showQuitConfirmationView(m model) (tea.Model, tea.Cmd) {
	m.previousView = m.view
	m.view = quitConfirmationView{}
	return m, nil
}

type quitConfirmationViewKeyMap struct {
	Quit   key.Binding
	Cancel key.Binding
}

var quitConfirmationViewKeys = quitConfirmationViewKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "quit"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

func (q quitConfirmationViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{q.Quit, q.Cancel}
}

func (q quitConfirmationViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{q.Quit, q.Cancel},
	}
}

type quitConfirmationView struct{}

func (q quitConfirmationView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, quitConfirmationViewKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, quitConfirmationViewKeys.Cancel):
			m.view = m.previousView
			return m, nil
		}
	}

	return m, nil
}

func (q quitConfirmationView) draw(m model) string {
	const text = "Are you sure you want to quit?\n\nAny game progress will be lost."
	helpView := m.help.View(quitConfirmationViewKeys)
	return lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)
}
