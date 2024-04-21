package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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

type quitConfirmationView struct{}

func (q quitConfirmationView) update(msg tea.KeyMsg, m model) (model, tea.Cmd) {
	switch {
	case key.Matches(msg, quitConfirmationViewKeys.Quit):
		return m, tea.Quit
	case key.Matches(msg, quitConfirmationViewKeys.Cancel):
		m.view = m.previousView
		return m, nil
	}

	return m, nil
}

func (q quitConfirmationView) draw(m model) string {
	const text = "Are you sure you want to quit?\n\nAny game progress will be lost."
	controls := []control{
		{key: "enter", description: "quit"},
		{key: "any other key", description: "cancel"},
	}
	controlsString := controlsToString(controls)

	return lipgloss.JoinVertical(lipgloss.Left, text, controlsString)
}
