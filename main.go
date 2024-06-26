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

var version = "dev"

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

var whiteColor = lipgloss.AdaptiveColor{
	Light: "8",
	Dark:  "7",
}
var blackColor = lipgloss.AdaptiveColor{
	Light: "15",
	Dark:  "0",
}

var accentColor = lipgloss.AdaptiveColor{
	Light: "12",
	Dark:  "4",
}

var highlightedStyle = lipgloss.NewStyle().Background(whiteColor).Foreground(blackColor).Bold(true)
var secondaryTextStyle = help.New().Styles.ShortDesc

type model struct {
	rand         *rand.Rand
	grid         grid
	score        int
	options      options
	moveCount    int
	view         view
	previousView view
	point1       vector2d
	point2       vector2d
	help         help.Model
	symbolSet    symbolSet
	windowSize   vector2d
	hintShown    bool
}

func initialModel(r *rand.Rand) model {
	return model{
		rand:      r,
		score:     0,
		options:   options{gameType: Endless},
		moveCount: 0,
		view:      titleView{},
		point1:    emptyVector2d,
		point2:    emptyVector2d,
		help:      help.New(),
		symbolSet: newEmojiSymbolSet(),
		hintShown: false,
	}
}

type tickMsg time.Time

// TODO: Add different game modes - e.g. endless, timed, limited number of moves
// TODO: Check resizing
// todo: change esc key to different key (?)

func (m model) Init() tea.Cmd {
	return nil
}

func tickCmd() tea.Cmd {
	return tea.Tick(180*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.windowSize = vector2d{
			x: msg.Width,
			y: msg.Height,
		}
		if !isWindowLargeEnough(m) && m.view != (windowTooSmallView{}) {
			return showWindowTooSmallView(m)
		} else if isWindowLargeEnough(m) && m.view == (windowTooSmallView{}) {
			return showPreviousView(m)
		}
	case tea.KeyMsg, tickMsg:
		return m.view.update(msg, m)
	}

	return m, nil
}

var minWindowSize = vector2d{
	x: 80,
	y: 22,
}

func isWindowLargeEnough(m model) bool {
	return m.windowSize.x >= minWindowSize.x && m.windowSize.y >= minWindowSize.y
}

func newEndGameKeyBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "end game"),
	)
}

func newHelpKeyBinding(m model) key.Binding {
	var helpKeyDescription string
	if m.help.ShowAll {
		helpKeyDescription = "hide controls"
	} else {
		helpKeyDescription = "show controls"
	}

	return key.NewBinding(
		key.WithKeys("?", "/"), // Include "/" ("?" without pressing shift key) for convenience
		key.WithHelp("?", helpKeyDescription),
	)
}

func (m model) View() string {
	titleBar := drawTitleBar(m)
	mainView := lipgloss.PlaceHorizontal(m.windowSize.x, lipgloss.Center,
		lipgloss.NewStyle().Padding(2, 4).Render(m.view.draw(m)))
	return lipgloss.NewStyle().Height(m.windowSize.y).Render(lipgloss.JoinVertical(lipgloss.Left, titleBar, mainView))
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := tea.NewProgram(initialModel(r))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
