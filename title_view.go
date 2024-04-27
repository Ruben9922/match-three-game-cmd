package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type titleView struct{}

type titleViewKeyMap struct {
	Quit           key.Binding
	ToggleGameType key.Binding
	Start          key.Binding
}

var titleViewKeys = titleViewKeyMap{
	Quit: sharedKeys.Quit,
	ToggleGameType: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "change game type"),
	),
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "start"),
	),
}

func (k titleViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Start, k.ToggleGameType, k.Quit}
}

func (k titleViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.ToggleGameType, k.Quit},
	}
}

func (tv titleView) draw(m model) string {
	const titlePart1 = "  __  __       _       _       _____ _                   \n |  \\/  | __ _| |_ ___| |__   |_   _| |__  _ __ ___  ___ \n | |\\/| |/ _` | __/ __| '_ \\    | | | '_ \\| '__/ _ \\/ _ \\\n | |  | | (_| | || (__| | | |   | | | | | | | |  __/  __/\n |_|  |_|\\__,_|\\__\\___|_| |_|   |_| |_| |_|_|  \\___|\\___|"
	const titlePart2 = "   ____                      \n  / ___| __ _ _ __ ___   ___ \n | |  _ / _` | '_ ` _ \\ / _ \\\n | |_| | (_| | | | | | |  __/\n  \\____|\\__,_|_| |_| |_|\\___|"
	const text = "\n Press any key to start..."

	radioButtons := drawRadioButtons([]gameType{Endless, LimitedMoves}, m.options.gameType, "Game type", "T")
	helpView := m.help.View(titleViewKeys)
	return lipgloss.JoinVertical(lipgloss.Left, titlePart1, titlePart2, text, radioButtons, "", helpView)
}

func (tv titleView) update(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, titleViewKeys.Quit):
		return showQuitConfirmationView(m)

	case key.Matches(msg, titleViewKeys.ToggleGameType):
		m.options.gameType = toggleGameType(m.options.gameType)
	case key.Matches(msg, titleViewKeys.Start):
		//m.view = RefreshGridView
		refreshGrid(&m.grid, m.rand, &m.score, false)

		// todo: make this code not duplicated
		m.potentialMatch = make([]vector2d, 0)
		for len(m.potentialMatch) == 0 {
			// Check if there are any possible matches; if no possible matches then create a new grid
			m.potentialMatch = findPotentialMatch(m.grid)
			if len(m.potentialMatch) == 0 {
				m.grid = newGrid(m.rand)
				refreshGrid(&m.grid, m.rand, &m.score, true)
			}
		}

		m.view = selectFirstPointView{}
		m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
		m.point2 = emptyVector2d

		//return m, tickCmd()
	}

	return m, nil
}
