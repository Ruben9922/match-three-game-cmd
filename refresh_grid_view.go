package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type refreshGridViewKeyMap struct {
	Quit key.Binding
	Skip key.Binding
}

var refreshGridViewKeys = refreshGridViewKeyMap{
	Quit: sharedKeys.Quit,
	Skip: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "skip"),
	),
}

func (r refreshGridViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{r.Skip, r.Quit}
}

func (r refreshGridViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{r.Skip, r.Quit},
	}
}

type refreshGridView struct{}

func (r refreshGridView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var skipped bool
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, refreshGridViewKeys.Quit):
			return showQuitConfirmationView(m)
		case key.Matches(msg, refreshGridViewKeys.Skip):
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
	for {
		// todo: sort out scoring
		finished = refreshGrid(&m.grid, m.rand, &m.score, true)

		if finished {
			isPlaying := m.options.gameType != LimitedMoves || m.remainingMoveCount > 0
			if isPlaying {
				// Check if there is a potential match; if not, then navigate to "no possible moves" view to create a new grid
				// todo: use nil everywhere instead of empty slice
				m.potentialMatch = findPotentialMatch(m.grid)
				if len(m.potentialMatch) == 0 {
					m.view = noPossibleMovesView{}

					return m, nil
				}

				m.view = selectFirstPointView{}
				m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
			} else {
				m.view = gameOverView{}
			}

			return m, nil
		}

		if !skipped {
			return m, tickCmd()
		}
	}
}

func (r refreshGridView) draw(m model) string {
	const text = "Refreshing grid..."
	helpView := m.help.View(refreshGridViewKeys)
	refreshGridText := lipgloss.JoinVertical(lipgloss.Left, text, "", helpView)

	gridText := createGrid(m, []vector2d{m.point1, m.point2})

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		gridText,
		lipgloss.NewStyle().MarginLeft(3).Render(refreshGridText),
	)
}