package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type confirmationViewKeyMap struct {
	Confirm key.Binding
	Cancel  key.Binding
}

var confirmationViewKeys = confirmationViewKeyMap{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

func (c confirmationViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{c.Confirm, c.Cancel}
}

func (c confirmationViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{c.Confirm, c.Cancel},
	}
}

type confirmationView struct {
	text          string
	keys          confirmationViewKeyMap
	confirmAction func(m model) (tea.Model, tea.Cmd)
}

func (c confirmationView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keys.Confirm):
			return c.confirmAction(m)
		case key.Matches(msg, c.keys.Cancel):
			m.view = m.previousView
			return m, nil
		}
	}

	return m, nil
}

func (c confirmationView) draw(m model) string {
	helpView := m.help.View(c.keys)
	return lipgloss.JoinVertical(lipgloss.Left, c.text, "", helpView)
}

func newQuitConfirmationView() quitConfirmationView {
	const text = "Are you sure you want to quit?"
	confirmAction := func(m model) (tea.Model, tea.Cmd) {
		return m, tea.Quit
	}
	q := quitConfirmationView{
		confirmationView: confirmationView{
			text:          text,
			keys:          confirmationViewKeys,
			confirmAction: confirmAction,
		},
	}

	const confirmKeyDescription = "quit"
	q.confirmationView.keys.Confirm.SetHelp(q.confirmationView.keys.Confirm.Help().Key, confirmKeyDescription)

	return q
}

func showQuitConfirmationView(m model) (tea.Model, tea.Cmd) {
	m.previousView = m.view
	m.view = newQuitConfirmationView()

	return m, nil
}

type quitConfirmationView struct {
	confirmationView
}

func newEndGameConfirmationView() endGameConfirmationView {
	const text = "Are you sure you want to end the game?\n\nAny game progress will be lost."
	confirmAction := func(m model) (tea.Model, tea.Cmd) {
		return showGameOverView(m, "You ended the game.")
	}
	q := endGameConfirmationView{
		confirmationView: confirmationView{
			text:          text,
			keys:          confirmationViewKeys,
			confirmAction: confirmAction,
		},
	}

	const confirmKeyDescription = "end game"
	q.confirmationView.keys.Confirm.SetHelp(q.confirmationView.keys.Confirm.Help().Key, confirmKeyDescription)

	return q
}

func showEndGameConfirmationView(m model) (tea.Model, tea.Cmd) {
	m.previousView = m.view
	m.view = newEndGameConfirmationView()

	return m, nil
}

type endGameConfirmationView struct {
	confirmationView
}
