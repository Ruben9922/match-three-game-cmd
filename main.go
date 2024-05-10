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

const symbolCount = 6

type grid [gridHeight][gridWidth]int

func newGrid(r *rand.Rand) (g grid) {
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			g[i][j] = r.Intn(symbolCount)
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

const emptySymbol int = -1

const whiteColor = lipgloss.Color("15")
const blackColor = lipgloss.Color("0")
const accentColor = lipgloss.Color("105")

var highlightedStyle = lipgloss.NewStyle().Background(whiteColor).Foreground(blackColor)

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
	showHint           bool
	help               help.Model
	symbolSet          symbolSet
}

func initialModel(r *rand.Rand) model {
	return model{
		rand:               r,
		score:              0,
		options:            options{gameType: Endless},
		remainingMoveCount: moveLimit,
		view:               titleView{},
		point1:             emptyVector2d,
		point2:             emptyVector2d,
		showHint:           false,
		help:               help.New(),
		symbolSet:          newEmojiSymbolSet(),
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
		key.WithKeys("?", "/"), // Include "/" ("?" without pressing shift key) for convenience
		key.WithHelp("?", "show/hide controls"),
	),
}

type tickMsg time.Time

// TODO: Add different game modes - e.g. endless, timed, limited number of moves
// TODO: Check resizing
// todo: change esc key to different key (?)

func (m model) Init() tea.Cmd {
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width

	case tea.KeyMsg, tickMsg:
		return m.view.update(msg, m)
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
	return lipgloss.NewStyle().Padding(2, 4).Render(m.view.draw(m))
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := tea.NewProgram(initialModel(r))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
