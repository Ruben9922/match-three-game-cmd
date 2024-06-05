package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func showRefreshGridView(m model) (tea.Model, tea.Cmd) {
	m.view = newRefreshGridView()
	m.point1 = emptyVector2d
	m.point2 = emptyVector2d
	m.help.ShowAll = false

	return m, tickCmd()
}

type refreshGridViewKeyMap struct {
	EndGame key.Binding
	Skip    key.Binding
}

func newRefreshGridViewKeys() refreshGridViewKeyMap {
	return refreshGridViewKeyMap{
		EndGame: newEndGameKeyBinding(),
		Skip: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "skip"),
		),
	}
}

func (r refreshGridViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{r.Skip, r.EndGame}
}

func (r refreshGridViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{r.Skip, r.EndGame},
	}
}

type refreshGridView struct {
	keys refreshGridViewKeyMap
}

func newRefreshGridView() refreshGridView {
	return refreshGridView{
		keys: newRefreshGridViewKeys(),
	}
}

func (r refreshGridView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var skipped bool
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, r.keys.EndGame):
			return showEndGameConfirmationView(m)
		case key.Matches(msg, r.keys.Skip):
			skipped = true
		default:
			return m, nil
		}
	case tickMsg:
		skipped = false
	default:
		return m, nil
	}

	finished := false
	var scorePointer *int // If hint was shown, don't update the score (both for the player's match and cascading matches)
	if m.hintShown {
		scorePointer = nil
	} else {
		scorePointer = &m.score
	}
	for {
		finished = refreshGrid(&m.grid, m.rand, scorePointer)

		if finished {
			isPlaying := m.options.gameType != LimitedMoves || m.moveCount < moveLimit
			if isPlaying {
				// Check if there is a potential match; if not, then navigate to "no possible moves" view to create a new grid
				potentialMatch := findPotentialMatch(m.grid)
				if len(potentialMatch) == 0 {
					showNoPossibleMovesView(m)
				}

				return showSelectFirstPointView(m)
			} else {
				return showGameOverView(m, "No more moves left.")
			}
		}

		if !skipped {
			return m, tickCmd()
		}
	}
}

func (r refreshGridView) draw(m model) string {
	const text = "Refreshing grid..."
	gridText := drawGrid(m, []vector2d{m.point1, m.point2})
	m.help.Width = m.windowSize.x - lipgloss.Width(gridText) - 8 - 3
	helpView := m.help.View(r.keys)
	refreshGridText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)
	gridLayoutText := drawGridLayout(m, gridText, refreshGridText)

	return gridLayoutText
}
