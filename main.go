package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"
	"unicode"
)

const gridHeight int = 10
const gridWidth int = 10
const minMatchLength int = 3

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

const emptySymbol rune = ' '

type grid [gridHeight][gridWidth]rune

func newGrid() (g grid) {
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			g[i][j] = getRandomSymbol()
		}
	}
	return
}

// no longer needed
//func (g grid) String() string {
//	var rowStrings [gridHeight]string
//	for i, row := range g {
//		rowStrings[i] = strings.Join(row, " ")
//	}
//	return strings.Join(rowStrings[:], "\n")
//}

var defaultStyle = tcell.StyleDefault.Background(tcell.ColorDefault).Foreground(tcell.ColorDefault)

func main() {
	// TODO: Consider using something like bubbletea instead
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	s.SetStyle(defaultStyle)

	// Initialise random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialise game
	g := newGrid()

	refreshGrid(s, &g)

	//fixme: for testing purposes only
	//g[0][9] = 'A'

	//g[1][9] = 'A'
	//g[2][8] = 'A'
	//g[3][9] = 'A'
	//g[3][8] = 'A'
	//g[3][7] = 'A'
	//g[3][6] = 'A'

	for {
		// todo: use nil everywhere instead of empty slice
		potentialMatch := make([]vector2d, 0)
		for len(potentialMatch) == 0 {
			// Check if there are any possible matches; if no possible matches then create a new grid
			potentialMatch = findPotentialMatch(g)
			if len(potentialMatch) == 0 {
				g := newGrid()
				refreshGrid(s, &g)
			}
		}

		// todo: fix initial extra key press
		swapPoints(s, &g, potentialMatch)

		refreshGrid(s, &g)
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

func draw(s tcell.Screen, g grid, selectedPoints []vector2d, text string) {
	s.Clear()

	// Draw grid
	drawGrid(s, g, selectedPoints)

	// Draw text
	screenWidth, screenHeight := s.Size()
	drawText(s, (gridWidth*2)+3, 0, screenWidth-1, screenHeight-1, defaultStyle, text)

	s.Show()
}

func refreshGrid(s tcell.Screen, g *grid) {
	ticker := time.NewTicker(150 * time.Millisecond)
	skipped := false
	skippedChannel := make(chan bool)

	go func() {
		for {
			ev := s.PollEvent()
			switch ev.(type) {
			case *tcell.EventKey:
				skippedChannel <- true
			}
		}
	}()

	waitForKeyPressOrTimeout := func() {
		if !skipped {
			select {
			case skipped = <-skippedChannel:
				ticker.Stop()
			case <-ticker.C:
			}
		}
	}

	const skipHint string = "Press any key to skip"
	draw(s, *g, []vector2d{}, skipHint)

	for {
		matches := findMatches(*g)
		if len(matches) == 0 {
			break
		}

		points := convertMatchesToPoints(matches)

		// Set points in matches to empty
		for _, p := range points {
			g[p.y][p.x] = emptySymbol
		}

		waitForKeyPressOrTimeout()
		draw(s, *g, []vector2d{}, skipHint)

		// Shift symbols down and insert random symbol at top of column
		sort.Slice(points, func(i, j int) bool {
			if points[i].y == points[j].y {
				return points[i].x < points[j].x
			}
			return points[i].y < points[j].y
		})
		for _, p := range points {
			for y := p.y; y > 0; y-- {
				g[y][p.x] = g[y-1][p.x]
			}
			g[0][p.x] = getRandomSymbol()

			waitForKeyPressOrTimeout()
			draw(s, *g, []vector2d{}, skipHint)
		}
	}

	//ticker := time.NewTicker(500 * time.Millisecond)
	//skipped := make(chan bool)
	//
	//go func() {
	//	for {
	//		select {
	//		case <-skipped:
	//			return
	//		case <-ticker.C:
	//			grid := newGrid()
	//			drawGrid(s, grid)
	//			s.Show()
	//		}
	//
	//	}
	//}()
	//
	//ticker.Stop()
	//skipped <- true
}

func convertMatchesToPoints(matches []match) []vector2d {
	// Calculating actual capacity would require looping through `matches`, so just making a guess
	points := make([]vector2d, 0, len(matches)*5)
	for _, m := range matches {
		for i := 0; i < m.length; i++ {
			point := vector2d{
				x: m.position.x + (i * m.direction.x),
				y: m.position.y + (i * m.direction.y),
			}
			points = append(points, point)
		}
	}
	return points
}

func swapPoints(s tcell.Screen, g *grid, potentialMatch []vector2d) {
	point1 := selectFirstPoint(s, *g, potentialMatch)
	point2 := selectSecondPoint(s, *g, point1)

	gUpdated := *g
	gUpdated[point1.y][point1.x], gUpdated[point2.y][point2.x] =
		gUpdated[point2.y][point2.x], gUpdated[point1.y][point1.x]
	matches := findMatches(gUpdated)
	if len(matches) != 0 {
		*g = gUpdated
		text := fmt.Sprintf("Swapped %c (%d, %d) and %c (%d, %d); match formed\nPress any key to continue",
			g[point1.y][point1.x], point1.x, point1.y, g[point2.y][point2.x], point2.x, point2.y)
		draw(s, *g, convertMatchesToPoints(matches), text)
	} else {
		text := "Not swapping as swap would not result in a match; please try again\nPress any key to continue"
		draw(s, *g, []vector2d{point1, point2}, text)
	}

	keyPressed := false
	for !keyPressed {
		ev := s.PollEvent()
		switch ev.(type) {
		//case *tcell.EventResize:
		//	s.Sync()
		case *tcell.EventKey:
			keyPressed = true
		}
	}
}

func selectFirstPoint(s tcell.Screen, g grid, potentialMatch []vector2d) vector2d {
	// Initialise point 1 to centre of grid
	point1 := vector2d{x: gridWidth / 2, y: gridHeight / 2}

	generateText := func() string {
		return fmt.Sprintf(
			"Select two points to swap (selecting point 1)...\n"+
				"Press arrow keys (← ↑ → ↓) to move selection; press enter to continue; press H to toggle hint\n\n"+
				"Current selection: %c (%d, %d)",
			g[point1.y][point1.x], point1.x, point1.y)
	}

	draw(s, g, []vector2d{point1}, generateText())
	selected := false
	showPotentialMatch := false
	for !selected {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		//case *tcell.EventResize:
		//	s.Sync()
		case *tcell.EventKey:
			if showPotentialMatch {
				showPotentialMatch = false
			} else {
				switch ev.Key() {
				case tcell.KeyUp:
					point1.y--
				case tcell.KeyDown:
					point1.y++
				case tcell.KeyLeft:
					point1.x--
				case tcell.KeyRight:
					point1.x++
				case tcell.KeyEnter:
					selected = true
				}

				point1.x = (point1.x + gridWidth) % gridWidth
				point1.y = (point1.y + gridHeight) % gridHeight

				if unicode.ToLower(ev.Rune()) == 'h' {
					// todo: score no points if hint shown (?)
					showPotentialMatch = true
				}
			}
		}

		if showPotentialMatch {
			draw(s, g, potentialMatch, generateText()+"\nShowing hint")
		} else {
			draw(s, g, []vector2d{point1}, generateText())
		}
	}

	return point1
}

func selectSecondPoint(s tcell.Screen, g grid, point1 vector2d) vector2d {
	point2 := point1
	if point1.y == 0 {
		if point1.x == gridWidth-1 {
			point2.x--
		} else {
			point2.x++
		}
	} else {
		point2.y--
	}

	generateText := func() string {
		return fmt.Sprintf(
			"Select two points to swap (selecting point 2)...\n"+
				"Press arrow keys (← ↑ → ↓) to move selection; press enter to continue\n\n"+
				"Point 1: %c (%d, %d)\n"+
				"Current selection: %c (%d, %d)",
			g[point1.y][point1.x], point1.x, point1.y, g[point2.y][point2.x], point2.x, point2.y)
	}

	draw(s, g, []vector2d{point1, point2}, generateText())
	selected := false
	for !selected {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		//case *tcell.EventResize:
		//	s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyUp || ev.Rune() == 'W' {
				//point2 = vector2d{
				//	x: point1.x,
				//	y: point1.y - 1,
				//}
				point2 = point1
				point2.y--
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'S' {
				//point2 = vector2d{
				//	x: point1.x,
				//	y: point1.y + 1,
				//}
				point2 = point1
				point2.y++
			} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'A' {
				//point2 = vector2d{
				//	x: point1.x - 1,
				//	y: point1.y,
				//}
				point2 = point1
				point2.x--
			} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'D' {
				//point2 = vector2d{
				//	x: point1.x + 1,
				//	y: point1.y,
				//}
				point2 = point1
				point2.x++
			} else if ev.Key() == tcell.KeyEnter {
				selected = true
			}
		}

		point2.x = (point2.x + gridWidth) % gridWidth
		point2.y = (point2.y + gridHeight) % gridHeight

		draw(s, g, []vector2d{point1, point2}, generateText())
	}

	return point2
}

type vector2d struct {
	x, y int
}

//var emptyVector2d = vector2d{x: -1, y: -1}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// Maybe remove this and just use slice of points instead (?)
type match struct {
	position  vector2d
	direction vector2d
	length    int
}

var emptyMatch = newMatch(vector2d{x: -1, y: -1}, vector2d{}, 0)

func newMatch(position, direction vector2d, length int) match {
	return match{
		position:  position,
		direction: direction,
		length:    length,
	}
}

func isPointInsideGrid(p vector2d) bool {
	return p.x >= 0 && p.x < gridWidth && p.y >= 0 && p.y < gridHeight
}

// Match algorithm works in a way that scores the player the most points
// * Prefers longer matches (checks for all possible matches and chooses the longest one) - to score the player more points
// * Prefers lower down matches - idea is that this would score the player more points as more pieces falling means potentially more "automatic" matches
// * "Maximal munch" behaviour - matches will be as long as possible; matches can be longer than the minimum match length
// TODO: If match lengths equal, then prefer matches lower in grid
// TODO: Return slice of matches so multiple matches are removed in one go
func findMatches(g grid) []match {
	directions := []vector2d{
		{x: 1, y: 0},
		{x: 0, y: 1},
	}
	matches := make([]match, 0, 10)
	for _, d := range directions {
		offset := vector2d{
			x: max((d.x*minMatchLength)-1, 0),
			y: max((d.y*minMatchLength)-1, 0),
		}

		d.y = -d.y

		for i := gridHeight - 1; i >= offset.y; i-- {
			for j := 0; j < gridWidth-offset.x; j++ {
				matchLength := 0
				originPoint := vector2d{x: j, y: i}
				for {
					currentPoint := vector2d{
						x: j + (matchLength * d.x),
						y: i + (matchLength * d.y),
					}

					if !isPointInsideGrid(currentPoint) {
						break
					}

					isSameSymbol := g[originPoint.y][originPoint.x] == g[currentPoint.y][currentPoint.x]
					if !isSameSymbol {
						break
					}

					matchLength++
				}

				if matchLength >= minMatchLength {
					matches = append(matches, newMatch(originPoint, d, matchLength))
				}
			}
		}
	}

	return matches
}

func generatePotentialMatchFilters() [][]vector2d {
	filters := make([][]vector2d, 0, (minMatchLength*2)+2)
	for i := 0; i < minMatchLength; i++ {
		// Filters of the form:
		// X   |  X  |   X
		//  XX | X X | XX
		filter := make([]vector2d, 0, 3)
		for j := 0; j < minMatchLength; j++ {
			if j == i {
				filter = append(filter, vector2d{x: j, y: 0})
			} else {
				filter = append(filter, vector2d{x: j, y: 1})
			}
		}
		filters = append(filters, filter)

		// Filters of the form:
		//  XX | X X | XX
		// X   |  X  |   X
		filter = make([]vector2d, 0, 3)
		for j := 0; j < minMatchLength; j++ {
			if j == i {
				filter = append(filter, vector2d{x: j, y: 1})
			} else {
				filter = append(filter, vector2d{x: j, y: 0})
			}
		}
		filters = append(filters, filter)
	}

	// Filter of the form:
	// X XX
	filter := make([]vector2d, 0, 3)
	for j := 0; j < minMatchLength; j++ {
		if j == 0 {
			filter = append(filter, vector2d{x: 0, y: 0})
		} else {
			filter = append(filter, vector2d{x: j + 1, y: 0})
		}
	}
	filters = append(filters, filter)

	// Filter of the form:
	// XX X
	filter = make([]vector2d, 0, 3)
	for j := 0; j < minMatchLength; j++ {
		if j == minMatchLength-1 {
			filter = append(filter, vector2d{x: minMatchLength, y: 0})
		} else {
			filter = append(filter, vector2d{x: j, y: 0})
		}
	}
	filters = append(filters, filter)

	// todo: add rotated versions of filters

	return filters
}

func computeObjectSize(object []vector2d) vector2d {
	xs := make([]int, 0, len(object))
	for _, p := range object {
		xs = append(xs, p.x)
	}
	sort.Ints(xs)

	xMin := xs[0]
	xMax := xs[len(xs)-1]
	xSize := (xMax - xMin) + 1

	ys := make([]int, 0, len(object))
	for _, p := range object {
		ys = append(ys, p.y)
	}
	sort.Ints(ys)

	yMin := ys[0]
	yMax := ys[len(ys)-1]
	ySize := (yMax - yMin) + 1

	return vector2d{x: xSize, y: ySize}
}

// May want to revise this to allow potential matches longer than minimum match length
// Add text warning that it may not be the optimal match
func findPotentialMatch(g grid) []vector2d {
	filters := generatePotentialMatchFilters()

	// todo put filters loop on inside
	// todo go from bottom to top
	for _, f := range filters {
		filterSize := computeObjectSize(f)

		for x := 0; x < gridWidth-filterSize.x+1; x++ {
			for y := 0; y < gridHeight-filterSize.y+1; y++ {
				sameSymbol := true
				origin := vector2d{x: x, y: y}
				reference := f[0]
				referenceGridCoords := vector2d{x: origin.x + reference.x, y: origin.y + reference.y}
				fGridCoords := make([]vector2d, 0, len(f))
				for _, p := range f {
					pGridCoords := vector2d{x: origin.x + p.x, y: origin.y + p.y}
					if g[pGridCoords.y][pGridCoords.x] != g[referenceGridCoords.y][referenceGridCoords.x] {
						sameSymbol = false
						break
					}

					fGridCoords = append(fGridCoords, pGridCoords)
				}
				if sameSymbol {
					return fGridCoords
				}
			}
		}
	}

	return []vector2d{}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 || r == '\n' {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func drawGrid(s tcell.Screen, g grid, selectedPoints []vector2d) {
	for i, row := range g {
		for j, symbol := range row {
			isSelected := false
			for _, selectedPoint := range selectedPoints {
				if i == selectedPoint.y && j == selectedPoint.x {
					isSelected = true
					break
				}
			}

			var style tcell.Style
			if isSelected {
				style = symbolHighlightedColors[symbol]
			} else {
				style = symbolColors[symbol]
			}

			s.SetContent(j*2, i, symbol, nil, style)
		}
	}
}

func getRandomSymbol() rune {
	index := rand.Intn(len(symbols))
	return symbols[index]
}
