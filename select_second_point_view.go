package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func showSelectSecondPointView(m model) (tea.Model, tea.Cmd) {
	m.help.ShowAll = false // Important that this is updated before creating the view
	m.point2 = getInitialPoint2(m.point1)

	s := newSelectSecondPointView(m)
	m.view = &s

	return m, nil
}

type selectSecondPointViewKeyMap struct {
	EndGame key.Binding
	Help    key.Binding
	Select  key.Binding
	Cancel  key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
}

func newSelectSecondPointViewKeys(m model) selectSecondPointViewKeyMap {
	return selectSecondPointViewKeyMap{
		EndGame: newEndGameKeyBinding(),
		Help:    newHelpKeyBinding(m),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "select"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
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
}

func (s selectSecondPointViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{s.Help, s.EndGame}
}

func (s selectSecondPointViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{s.Up, s.Down, s.Left, s.Right, s.Select, s.Cancel},
		{s.Help, s.EndGame},
	}
}

type selectSecondPointView struct {
	keys selectSecondPointViewKeyMap
}

func newSelectSecondPointView(m model) selectSecondPointView {
	return selectSecondPointView{
		keys: newSelectSecondPointViewKeys(m),
	}
}

func (s *selectSecondPointView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.EndGame):
			return showEndGameConfirmationView(m)
		case key.Matches(msg, s.keys.Help):
			return s.toggleHelp(m)

		case key.Matches(msg, s.keys.Select):
			// Swap the points, if it would result in a match
			updatedGrid := m.grid
			updatedGrid[m.point1.y][m.point1.x], updatedGrid[m.point2.y][m.point2.x] =
				updatedGrid[m.point2.y][m.point2.x], updatedGrid[m.point1.y][m.point1.x]
			matches := findMatches(updatedGrid)
			if len(matches) != 0 {
				m.grid = updatedGrid

				m.moveCount++
			}

			return showSelectPointConfirmationView(m)
		case key.Matches(msg, s.keys.Cancel):
			return showSelectFirstPointView(m)
		}

		var point2Updated vector2d
		switch {
		case key.Matches(msg, s.keys.Up):
			point2Updated = vector2d{
				x: m.point1.x,
				y: m.point1.y - 1,
			}
		case key.Matches(msg, s.keys.Down):
			point2Updated = vector2d{
				x: m.point1.x,
				y: m.point1.y + 1,
			}
		case key.Matches(msg, s.keys.Left):
			point2Updated = vector2d{
				x: m.point1.x - 1,
				y: m.point1.y,
			}
		case key.Matches(msg, s.keys.Right):
			point2Updated = vector2d{
				x: m.point1.x + 1,
				y: m.point1.y,
			}
		default:
			return m, nil
		}
		if isPointInsideGrid(point2Updated) {
			m.point2 = point2Updated
		}
	}

	return m, nil
}

// todo: combine the two copies of this function (?)
func (s *selectSecondPointView) toggleHelp(m model) (tea.Model, tea.Cmd) {
	// Toggle between short and full help in help view
	m.help.ShowAll = !m.help.ShowAll

	// Update help key so the description ("show controls"/"hide controls") is updated accordingly
	s.keys.Help = newHelpKeyBinding(m)

	return m, nil
}

func (s *selectSecondPointView) draw(m model) string {
	const text = "Select two points to swap (selecting point 2)..."
	gridText := drawGrid(m, []vector2d{m.point1, m.point2})
	m.help.Width = m.windowSize.x - lipgloss.Width(gridText) - 8 - 3
	helpView := m.help.View(s.keys)
	selectSecondPointText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)
	gridLayoutText := drawGridLayout(m, gridText, selectSecondPointText)

	return gridLayoutText
}
