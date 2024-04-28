package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectSecondPointViewKeyMap struct {
	sharedKeyMap
	Select key.Binding
	Cancel key.Binding
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
}

var selectSecondPointViewKeys = selectSecondPointViewKeyMap{
	sharedKeyMap: sharedKeys,
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("escape"),
		key.WithHelp("esc", "cancel"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "w"),
		key.WithHelp("↑/w", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "s"),
		key.WithHelp("↓/s", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "a"),
		key.WithHelp("←/a", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "d"),
		key.WithHelp("→/d", "right"),
	),
}

func (s selectSecondPointViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Help, s.Quit}
}

func (s selectSecondPointViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Up, s.Down, s.Left, s.Right},
		{s.Select, s.Cancel},
		{s.Help, s.Quit},
	}
}

type selectSecondPointView struct{}

func (s selectSecondPointView) update(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, selectSecondPointViewKeys.Quit):
		return showQuitConfirmationView(m)
	case key.Matches(msg, selectSecondPointViewKeys.Help):
		return toggleHelp(m)

	case key.Matches(msg, selectSecondPointViewKeys.Select):
		// Swap the points, if it would result in a match
		updatedGrid := m.grid
		updatedGrid[m.point1.y][m.point1.x], updatedGrid[m.point2.y][m.point2.x] =
			updatedGrid[m.point2.y][m.point2.x], updatedGrid[m.point1.y][m.point1.x]
		matches := findMatches(updatedGrid)
		if len(matches) != 0 {
			m.grid = updatedGrid

			if m.options.gameType == LimitedMoves {
				m.remainingMoveCount--
			}
		}

		m.view = selectPointConfirmationView{}
		return m, nil

	case key.Matches(msg, selectSecondPointViewKeys.Cancel):
		m.view = selectFirstPointView{}
		m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
		m.point2 = emptyVector2d
		return m, nil
	}

	var point2Updated vector2d
	switch {
	case key.Matches(msg, selectSecondPointViewKeys.Up):
		point2Updated = vector2d{
			x: m.point1.x,
			y: m.point1.y - 1,
		}
	case key.Matches(msg, selectSecondPointViewKeys.Down):
		point2Updated = vector2d{
			x: m.point1.x,
			y: m.point1.y + 1,
		}
	case key.Matches(msg, selectSecondPointViewKeys.Left):
		point2Updated = vector2d{
			x: m.point1.x - 1,
			y: m.point1.y,
		}
	case key.Matches(msg, selectSecondPointViewKeys.Right):
		point2Updated = vector2d{
			x: m.point1.x + 1,
			y: m.point1.y,
		}
	}
	if isPointInsideGrid(point2Updated) {
		m.point2 = point2Updated
	}

	return m, nil
}

func (s selectSecondPointView) draw(m model) string {
	const text = "Select two points to swap (selecting point 2)..."
	helpView := m.help.View(selectSecondPointViewKeys)
	selectSecondPointText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)

	gridText := createGrid(m, []vector2d{m.point1, m.point2})

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		gridText,
		lipgloss.NewStyle().MarginLeft(3).Render(selectSecondPointText),
	)
}
