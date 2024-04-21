package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type gameOverView struct{}

func (g gameOverView) update(_ tea.KeyMsg, m model) (model, tea.Cmd) {
	return m, tea.Quit
}

func (g gameOverView) draw(model) string {
	const text = "Game over!\n\nNo more moves left."
	controls := []control{{key: "<Any key>", description: "Exit"}}
	controlsString := controlsToString(controls)
	return lipgloss.JoinVertical(lipgloss.Left, text, controlsString)
}
