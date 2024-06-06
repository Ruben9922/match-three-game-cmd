package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func showWindowTooSmallView(m model) (tea.Model, tea.Cmd) {
	m.previousView = m.view
	m.view = windowTooSmallView{}
	m.help.ShowAll = false

	return m, nil
}

type windowTooSmallViewKeyMap struct {
	Quit key.Binding
}

var windowTooSmallViewKeys = windowTooSmallViewKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
}

func (w windowTooSmallViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{w.Quit}
}

func (w windowTooSmallViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{w.Quit},
	}
}

type windowTooSmallView struct{}

func (w windowTooSmallView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, windowTooSmallViewKeys.Quit):
			return showQuitConfirmationView(m)
		}
	}

	return m, nil
}

func (w windowTooSmallView) draw(m model) string {
	text := fmt.Sprintf("Window is too small. Please resize the window to at least %dx%d (currently %dx%d).",
		minWindowSize.x, minWindowSize.y, m.windowSize.x, m.windowSize.y)

	m.help.Width = m.windowSize.x - 8
	helpView := m.help.View(windowTooSmallViewKeys)

	return lipgloss.NewStyle().
		Width(m.windowSize.x - 8).
		Render(lipgloss.JoinVertical(lipgloss.Left, text, "", helpView))
}
