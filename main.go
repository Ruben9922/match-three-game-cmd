package main

import (
	"github.com/gdamore/tcell"
	"log"
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

const gridHeight int = 10
const gridWidth int = 10
const minMatchLength int = 3
const scorePerMatchedSymbol = 40

const emptySymbol rune = ' '

var symbols = []rune{'A', 'B', 'C', 'D', 'E', 'F'}
var symbolColors = map[rune]tcell.Style{
	symbols[0]: tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorWhite),
	symbols[1]: tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDarkCyan),
	symbols[2]: tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDarkMagenta),
	symbols[3]: tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorGreen),
	symbols[4]: tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorRed),
	symbols[5]: tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorYellow),
}
var symbolHighlightedColors = map[rune]tcell.Style{
	symbols[0]: tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorDefault),
	symbols[1]: tcell.StyleDefault.Background(tcell.ColorDarkCyan).Foreground(tcell.ColorDefault),
	symbols[2]: tcell.StyleDefault.Background(tcell.ColorDarkMagenta).Foreground(tcell.ColorDefault),
	symbols[3]: tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorDefault),
	symbols[4]: tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorDefault),
	symbols[5]: tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorDefault),
}

var defaultStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)

func main() {
	// TODO: Consider using something like bubbletea instead
	// TODO: Quit on Escape / "q" key
	// TODO: Try using emojis instead of letters (maybe make this optional)
	// TODO: Add different game modes - e.g. endless, timed, limited number of moves
	// TODO: Check resizing
	// TODO: Reorder code
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	s.SetStyle(defaultStyle)

	// Display title screen
	drawTitleScreen(s)
	waitForKeyPress(s)

	// Initialise random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialise game
	g := newGrid(r)
	score := 0

	refreshGrid(s, &g, r, &score, false)

	for {
		// todo: use nil everywhere instead of empty slice
		potentialMatch := make([]vector2d, 0)
		for len(potentialMatch) == 0 {
			// Check if there are any possible matches; if no possible matches then create a new grid
			potentialMatch = findPotentialMatch(g)
			if len(potentialMatch) == 0 {
				g = newGrid(r)
				refreshGrid(s, &g, r, &score, false)
			}
		}

		// todo: fix initial extra key press
		swapPoints(s, &g, potentialMatch, score)

		refreshGrid(s, &g, r, &score, true)
	}

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}

	s.Clear()
}
