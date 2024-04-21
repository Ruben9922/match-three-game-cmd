package main

import (
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

type selectFirstPointView struct{}

var selectFirstPointViewKeys = selectFirstPointViewKeyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
	ToggleHint: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "show/hide hint"),
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

func (s selectFirstPointView) update(msg tea.KeyMsg, m model) (model, tea.Cmd) {
	if key.Matches(msg, selectFirstPointViewKeys.Select) {
		m.view = selectSecondPointView{}
		m.point2 = getInitialPoint2(m.point1)
		return m, nil
	}
	switch {
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

	m.showHint = !m.showHint && key.Matches(msg, selectFirstPointViewKeys.ToggleHint)
	return m, nil
}

func (s selectFirstPointView) draw(m model) string {
	const text = "Select two points to swap (selecting point 1)..."
	controls := []control{
		{key: "← ↑ → ↓ / WASD", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "H", description: "Show hint"},
		{key: "Q", description: "End Game"},
	}
	hintControls := []control{
		{key: "<Any key>", description: "Hide hint"},
	}

	var controlsString string
	var selectedPoints []vector2d
	if m.showHint {
		controlsString = controlsToString(hintControls)
		selectedPoints = m.potentialMatch
	} else {
		controlsString = controlsToString(controls)
		selectedPoints = []vector2d{m.point1}
	}
	selectFirstPointText := lipgloss.JoinVertical(lipgloss.Left, text, controlsString)

	gridText := createGrid(m, selectedPoints)

	return lipgloss.JoinHorizontal(lipgloss.Top, gridText, selectFirstPointText)
}