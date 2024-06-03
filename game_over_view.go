package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func showGameOverView(m model, text string) (tea.Model, tea.Cmd) {
	m.view = gameOverView{text: text}
	m.help.ShowAll = false

	return m, nil
}

type gameOverViewKeyMap struct {
	TitleView key.Binding
	Quit      key.Binding
}

var gameOverViewKeys = gameOverViewKeyMap{
	TitleView: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "title screen"),
	),
	Quit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "quit"),
	),
}

func (s gameOverViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.TitleView, s.Quit}
}

func (s gameOverViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.TitleView, s.Quit},
	}
}

type gameOverView struct {
	text string
}

func (g gameOverView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gameOverViewKeys.TitleView):
			return showTitleView(m)
		case key.Matches(msg, gameOverViewKeys.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (g gameOverView) draw(m model) string {
	text := "Game over!\n\n" + g.text
	helpView := m.help.View(gameOverViewKeys)
	gameOverViewText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)

	gridText := drawGrid(m, []vector2d{})

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		gridText,
		lipgloss.NewStyle().Width(m.windowSize.x-lipgloss.Width(gridText)-8).PaddingLeft(3).Render(gameOverViewText),
	)
}
