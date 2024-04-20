package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math/rand"
	"os"
	"time"
)

type grid [gridHeight][gridWidth]rune

func newGrid(r *rand.Rand) (g grid) {
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			g[i][j] = getRandomSymbol(r)
		}
	}
	return
}

type vector2d struct {
	x, y int
}

var emptyVector2d = vector2d{x: -1, y: -1}

// Maybe remove this and just use slice of points instead (?)
type match struct {
	position  vector2d
	direction vector2d
	length    int
}

func newMatch(position, direction vector2d, length int) match {
	return match{
		position:  position,
		direction: direction,
		length:    length,
	}
}

type control struct {
	key         string
	description string
}

type gameType int

const (
	Endless gameType = iota
	LimitedMoves
)

func (gt gameType) String() string {
	return [...]string{"Endless", "Limited moves"}[gt]
}

type options struct {
	gameType gameType
}

const gridHeight int = 10
const gridWidth int = 10
const minMatchLength int = 3
const scorePerMatchedSymbol int = 40
const moveLimit int = 20

const emptySymbol rune = ' '

var symbols = []rune{'●', '▲', '■', '◆', '★', '❤'}
var symbolColors = map[rune]lipgloss.Style{
	symbols[0]: lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
	symbols[1]: lipgloss.NewStyle().Foreground(lipgloss.Color("274")),
	symbols[2]: lipgloss.NewStyle().Foreground(lipgloss.Color("279")),
	symbols[3]: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	symbols[4]: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
	symbols[5]: lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
}
var symbolHighlightedColors = map[rune]lipgloss.Style{
	symbols[0]: lipgloss.NewStyle().Background(lipgloss.Color("15")),
	symbols[1]: lipgloss.NewStyle().Background(lipgloss.Color("274")),
	symbols[2]: lipgloss.NewStyle().Background(lipgloss.Color("279")),
	symbols[3]: lipgloss.NewStyle().Background(lipgloss.Color("2")),
	symbols[4]: lipgloss.NewStyle().Background(lipgloss.Color("9")),
	symbols[5]: lipgloss.NewStyle().Background(lipgloss.Color("11")),
}

type view int

const (
	TitleView view = iota
	SelectFirstPointView
	SelectSecondPointView
	SelectPointConfirmationView
	//RefreshGridView
	QuitConfirmationView
	GameOverView
)

type model struct {
	rand               *rand.Rand
	grid               grid
	score              int
	options            options
	remainingMoveCount int
	view               view
	previousView       view
	point1             vector2d
	point2             vector2d
	//animationQueue      []grid
	showHint       bool
	potentialMatch []vector2d
}

func initialModel(r *rand.Rand) model {
	return model{
		rand:               r,
		grid:               newGrid(r),
		score:              0,
		options:            options{gameType: Endless},
		remainingMoveCount: moveLimit,
		view:               TitleView,
		point1:             emptyVector2d,
		point2:             emptyVector2d,
		//animationQueue:      make([]grid, 0),
		showHint:       false,
		potentialMatch: make([]vector2d, 0),
	}
}

//type tickMsg time.Time

// TODO: Add different game modes - e.g. endless, timed, limited number of moves
// TODO: Check resizing
// todo: fix having to press twice

func toggleGameType(gt gameType) gameType {
	if gt == Endless {
		return LimitedMoves
	}

	return Endless
}

func (m model) Init() tea.Cmd {
	return nil
}

//func tickCmd() tea.Cmd {
//	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
//		return tickMsg(t)
//	})
//}

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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.view != TitleView && m.view != QuitConfirmationView && msg.String() == "q" {
			m.previousView = m.view
			m.view = QuitConfirmationView
			return m, nil
		}

		switch m.view {
		case TitleView:
			if msg.String() == "t" {
				m.options.gameType = toggleGameType(m.options.gameType)
			} else {
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

				m.view = SelectFirstPointView
				m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
				m.point2 = emptyVector2d

				//return m, tickCmd()
			}
		case SelectFirstPointView:
			if msg.String() == "enter" {
				m.view = SelectSecondPointView
				m.point2 = getInitialPoint2(m.point1)
				return m, nil
			}
			switch msg.String() {
			case "up", "w":
				m.point1.y--
				m.point1.y = (m.point1.y + gridHeight) % gridHeight // Clamp y coordinate between 0 and gridHeight - 1
			case "down", "s":
				m.point1.y++
				m.point1.y = (m.point1.y + gridHeight) % gridHeight // Clamp y coordinate between 0 and gridHeight - 1
			case "left", "a":
				m.point1.x--
				m.point1.x = (m.point1.x + gridWidth) % gridWidth // Clamp x coordinate between 0 and gridWidth - 1
			case "right", "d":
				m.point1.x++
				m.point1.x = (m.point1.x + gridWidth) % gridWidth // Clamp x coordinate between 0 and gridWidth - 1
			}

			m.showHint = !m.showHint && msg.String() == "h"
		case SelectSecondPointView:
			if msg.String() == "enter" {
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

				m.view = SelectPointConfirmationView
				return m, nil
			}

			if msg.String() == "escape" {
				m.view = SelectFirstPointView
				m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
				m.point2 = emptyVector2d
				return m, nil
			}

			var point2Updated vector2d
			switch msg.String() {
			case "up", "w":
				point2Updated = vector2d{
					x: m.point1.x,
					y: m.point1.y - 1,
				}
			case "down", "s":
				point2Updated = vector2d{
					x: m.point1.x,
					y: m.point1.y + 1,
				}
			case "left", "a":
				point2Updated = vector2d{
					x: m.point1.x - 1,
					y: m.point1.y,
				}
			case "right", "d":
				point2Updated = vector2d{
					x: m.point1.x + 1,
					y: m.point1.y,
				}
			}
			if isPointInsideGrid(point2Updated) {
				m.point2 = point2Updated
			}
		case SelectPointConfirmationView:
			// Refresh grid
			//m.view = RefreshGridView

			refreshGrid(&m.grid, m.rand, &m.score, true)

			// todo: use nil everywhere instead of empty slice
			m.potentialMatch = make([]vector2d, 0)
			for len(m.potentialMatch) == 0 {
				// Check if there are any possible matches; if no possible matches then create a new grid
				m.potentialMatch = findPotentialMatch(m.grid)
				if len(m.potentialMatch) == 0 {
					m.grid = newGrid(m.rand)
					refreshGrid(&m.grid, m.rand, &m.score, true)
				}
			}

			if m.options.gameType != LimitedMoves || m.remainingMoveCount > 0 {
				m.view = SelectFirstPointView
				m.point1 = vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
				m.point2 = emptyVector2d
			} else {
				m.view = GameOverView
				m.point1 = emptyVector2d
				m.point2 = emptyVector2d
			}
		case GameOverView:
			return m, tea.Quit
		case QuitConfirmationView:
			if msg.String() == "enter" {
				return m, tea.Quit
			}
			m.view = m.previousView
		}
		//case tickMsg:
		//return m, tickCmd()
	}

	return m, nil
}

func (m model) View() string {
	switch m.view {
	case TitleView:
		return createTitleView(m)
	case SelectFirstPointView:
		return createSelectFirstPointView(m)
	case SelectSecondPointView:
		return createSelectSecondPointView(m)
	case SelectPointConfirmationView:
		return createSelectPointConfirmationView(m)
	//case RefreshGridView:
	//	return createGridView(m)
	case QuitConfirmationView:
		return createQuitConfirmationView()
	case GameOverView:
		return createGameOverView()
	default:
		return ""
	}
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := tea.NewProgram(initialModel(r))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
