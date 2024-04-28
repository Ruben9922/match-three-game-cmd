package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type selectFirstPointViewKeyMap struct {
	sharedKeyMap
	Select     key.Binding
	ToggleHint key.Binding
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
}

var selectFirstPointViewKeys = selectFirstPointViewKeyMap{
	sharedKeyMap: sharedKeys,
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	ToggleHint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "show hint"),
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

func (k selectFirstPointViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k selectFirstPointViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.ToggleHint},
		{k.Help, k.Quit},
	}
}

type selectFirstPointViewHintKeyMap struct {
	Quit       key.Binding
	ToggleHint key.Binding
}

var selectFirstPointViewHintKeys = selectFirstPointViewHintKeyMap{
	Quit: sharedKeys.Quit,
	ToggleHint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "hide hint"),
	),
}

func (k selectFirstPointViewHintKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleHint, k.Quit}
}

func (k selectFirstPointViewHintKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleHint, k.Quit},
	}
}

type selectFirstPointView struct{}

func (s selectFirstPointView) update(msg tea.KeyMsg, m model) (tea.Model, tea.Cmd) {
	if m.showHint {
		switch {
		case key.Matches(msg, selectFirstPointViewHintKeys.Quit):
			return showQuitConfirmationView(m)
		case key.Matches(msg, selectFirstPointViewHintKeys.ToggleHint):
			m.showHint = false
		}
		return m, nil
	}

	switch {
	case key.Matches(msg, selectFirstPointViewKeys.Quit):
		return showQuitConfirmationView(m)
	case key.Matches(msg, selectFirstPointViewKeys.Help):
		return toggleHelp(m)

	case key.Matches(msg, selectFirstPointViewKeys.ToggleHint):
		m.showHint = true

	case key.Matches(msg, selectFirstPointViewKeys.Select):
		m.view = selectSecondPointView{}
		m.point2 = getInitialPoint2(m.point1)

	case key.Matches(msg, selectFirstPointViewKeys.Up):
		m.point1.y--
		m.point1.y = (m.point1.y + gridHeight) % gridHeight // Clamp y coordinate between 0 and gridHeight - 1
	case key.Matches(msg, selectFirstPointViewKeys.Down):
		m.point1.y++
		m.point1.y = (m.point1.y + gridHeight) % gridHeight // Clamp y coordinate between 0 and gridHeight - 1
	case key.Matches(msg, selectFirstPointViewKeys.Left):
		m.point1.x--
		m.point1.x = (m.point1.x + gridWidth) % gridWidth // Clamp x coordinate between 0 and gridWidth - 1
	case key.Matches(msg, selectFirstPointViewKeys.Right):
		m.point1.x++
		m.point1.x = (m.point1.x + gridWidth) % gridWidth // Clamp x coordinate between 0 and gridWidth - 1
	}
	return m, nil
}

func getInitialPoint2(point1 vector2d) vector2d {
	if point1.y == 0 {
		if point1.x == gridWidth-1 {
			return vector2d{
				x: point1.x - 1,
				y: point1.y,
			}
		} else {
			return vector2d{
				x: point1.x + 1,
				y: point1.y,
			}
		}
	} else {
		return vector2d{
			x: point1.x,
			y: point1.y - 1,
		}
	}
}

func (s selectFirstPointView) draw(m model) string {
	const text = "Select two points to swap (selecting point 1)..."
	var keys help.KeyMap
	if m.showHint {
		keys = selectFirstPointViewHintKeys
	} else {
		keys = selectFirstPointViewKeys
	}
	helpView := m.help.View(keys)
	selectFirstPointText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)

	var selectedPoints []vector2d
	if m.showHint {
		selectedPoints = m.potentialMatch
	} else {
		selectedPoints = []vector2d{m.point1}
	}
	gridText := createGrid(m, selectedPoints)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		gridText,
		lipgloss.NewStyle().MarginLeft(3).Render(selectFirstPointText),
	)
}
