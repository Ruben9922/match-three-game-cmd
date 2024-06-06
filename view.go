package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type view interface {
	update(msg tea.Msg, m model) (tea.Model, tea.Cmd)
	draw(m model) string
}

func showPreviousView(m model) (tea.Model, tea.Cmd) {
	if m.previousView == nil {
		return m, nil
	}

	m.view = m.previousView
	m.help.ShowAll = false

	return m, nil
}

func showModal(m model, v view) (tea.Model, tea.Cmd) {
	// Only store the current view if it's not a modal
	// I.e. always return back to the last *non-modal* view
	// I.e. don't show a modal on top of another modal
	// Essentially avoids weird scenarios where a modal (e.g. quit confirmation view) is shown on top of another modal
	// (e.g. window too small view); modals store each other as the previous view, so it's impossible to escape from the
	// modals
	switch m.view.(type) {
	case windowTooSmallView, endGameConfirmationView, quitConfirmationView:
	default:
		m.previousView = m.view
	}

	m.view = v
	m.help.ShowAll = false

	return m, nil
}
