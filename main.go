package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	help           help.Model
}

func initialModel(r *rand.Rand) model {
	return model{
		rand:               r,
		grid:               newGrid(r),
		score:              0,
		options:            options{gameType: Endless},
		remainingMoveCount: moveLimit,
		view:               titleView{},
		point1:             emptyVector2d,
		point2:             emptyVector2d,
		//animationQueue:      make([]grid, 0),
		showHint:       false,
		potentialMatch: make([]vector2d, 0),
		help:           help.New(),
	}
}

type sharedKeyMap struct {
	Quit key.Binding
	Help key.Binding
}

var sharedKeys = sharedKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "show/hide controls"),
	),
}

//type tickMsg time.Time

// TODO: Add different game modes - e.g. endless, timed, limited number of moves
// TODO: Check resizing
// todo: fix having to press twice
// todo: change esc key to different key (?)

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
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg:
		return m.view.update(msg, m)

		//case tickMsg:
		//return m, tickCmd()
	}

	return m, nil
}

func toggleHelp(m model) (tea.Model, tea.Cmd) {
	m.help.ShowAll = !m.help.ShowAll
	return m, nil
}

func showQuitConfirmationView(m model) (tea.Model, tea.Cmd) {
	m.previousView = m.view
	m.view = quitConfirmationView{}
	return m, nil
}

func (m model) View() string {
	return m.view.draw(m)
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := tea.NewProgram(initialModel(r))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
