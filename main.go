package main

import (
	"fmt"
	"github.com/dustin/go-humanize/english"
	"github.com/gdamore/tcell"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"
	"unicode"
	"unicode/utf8"
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

var emptyVector2d = vector2d{x: -1, y: -1}

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

type control struct {
	key         string
	description string
}

func draw(s tcell.Screen, g grid, selectedPoints []vector2d, text string, controls []control) {
	s.Clear()

	// Draw grid
	drawGrid(s, g, selectedPoints)

	// Draw text
	screenWidth, _ := s.Size()
	const textOffsetX = (gridWidth * 2) + 3
	drawText(s, textOffsetX, 0, screenWidth-1, 5, defaultStyle, text)

	// Draw controls
	drawControls(s, controls, textOffsetX, 6)

	s.Show()
}

func drawControls(s tcell.Screen, controls []control, offsetX int, offsetY int) {
	screenWidth, screenHeight := s.Size()

	// TODO: Maybe extract into separate function or use lo.MinBy
	keyLengths := make([]int, 0, len(controls))
	for _, c := range controls {
		keyLengths = append(keyLengths, utf8.RuneCountInString(c.key))
	}
	sort.Ints(keyLengths)
	maxKeyLength := keyLengths[len(keyLengths)-1]

	drawText(s, offsetX, offsetY, screenWidth-1, offsetY+1, defaultStyle, "Controls:")
	for i, c := range controls {
		y1 := offsetY + i + 1
		y2 := offsetY + i + 2
		if y1 < screenHeight {
			drawText(s, offsetX, y1, offsetX+maxKeyLength, y2, defaultStyle, c.key)
			drawText(s, offsetX+maxKeyLength+2, y1, screenWidth-1, y2, defaultStyle, c.description)
		}
	}
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

	const text = "Refreshing grid..."
	controls := []control{{key: "<Any key>", description: "Skip"}}
	draw(s, *g, []vector2d{}, text, controls)

	for {
		matches := findMatches(*g)
		if len(matches) == 0 {
			break
		}

		points := convertMatchesToPoints(matches)

		// Shifting algorithm assumes points are unique; duplicate points will cause strange behaviour
		points = removeDuplicatePoints(points)

		// Set points in matches to empty
		for _, p := range points {
			g[p.y][p.x] = emptySymbol
		}

		waitForKeyPressOrTimeout()
		draw(s, *g, []vector2d{}, text, controls)

		// Shift symbols down and insert random symbol at top of column
		// Instead of manually updating points list, could maybe just search through grid for empty points
		// Assumes points are unique - duplicate points will cause strange behaviour
		// Want to shift lower points first - hence sorting such that lower points (points with higher y) come first
		sort.Slice(points, func(i, j int) bool {
			if points[i].x == points[j].x {
				return points[i].y > points[j].y
			}
			return points[i].x < points[j].x
		})
		for len(points) > 0 {
			p := points[0]
			points = points[1:]

			for y := p.y; y > 0; y-- {
				g[y][p.x] = g[y-1][p.x]
			}
			g[0][p.x] = getRandomSymbol()

			// Shift down remaining points in same column to account for shifting of corresponding empty points in grid
			// p1.y++ for each point p1 in this column (with same x)
			for i := 0; i < len(points); i++ {
				p1 := &points[i]
				// If point is in the same column and (strictly) above current point
				if p1.x == p.x && p1.y < p.y {
					p1.y++
				}
			}

			waitForKeyPressOrTimeout()
			draw(s, *g, []vector2d{}, text, controls)
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

func removeDuplicatePoints(points []vector2d) []vector2d {
	pointsMap := make(map[vector2d]bool, len(points))
	updatedPoints := make([]vector2d, 0, len(points))
	for _, p := range points {
		if _, present := pointsMap[p]; !present {
			pointsMap[p] = true
			updatedPoints = append(updatedPoints, p)
		}
	}
	return updatedPoints
}

func swapPoints(s tcell.Screen, g *grid, potentialMatch []vector2d) {
	point1 := vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
	point2 := emptyVector2d
	for point2 == emptyVector2d {
		point1 = selectFirstPoint(s, *g, potentialMatch, point1)
		point2 = selectSecondPoint(s, *g, point1)
	}

	gUpdated := *g
	gUpdated[point1.y][point1.x], gUpdated[point2.y][point2.x] =
		gUpdated[point2.y][point2.x], gUpdated[point1.y][point1.x]
	matches := findMatches(gUpdated)
	controls := []control{{key: "<Any key>", description: "Continue"}}
	if len(matches) != 0 {
		*g = gUpdated
		text := fmt.Sprintf("Swapped %c (%d, %d) and %c (%d, %d); %s formed",
			g[point1.y][point1.x], point1.x, point1.y, g[point2.y][point2.x], point2.x, point2.y,
			english.PluralWord(len(matches), "match", ""))
		draw(s, *g, convertMatchesToPoints(matches), text, controls)
	} else {
		text := "Not swapping as swap would not result in a match; please try again"
		draw(s, *g, []vector2d{point1, point2}, text, controls)
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

func selectFirstPoint(s tcell.Screen, g grid, potentialMatch []vector2d, point1Initial vector2d) vector2d {
	point1 := point1Initial

	generateText := func() string {
		return fmt.Sprintf(
			"Select two points to swap (selecting point 1)...\n\nCurrent selection: %c (%d, %d)",
			g[point1.y][point1.x], point1.x, point1.y)
	}
	controls := []control{
		{key: "← ↑ → ↓", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "H", description: "Show hint"},
	}
	hintControls := []control{
		{key: "<Any key>", description: "Hide hint"},
	}

	draw(s, g, []vector2d{point1}, generateText(), controls)
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
			draw(s, g, potentialMatch, generateText()+"\nShowing hint", hintControls)
		} else {
			draw(s, g, []vector2d{point1}, generateText(), controls)
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
			"Select two points to swap (selecting point 2)...\n\n"+
				"Point 1: %c (%d, %d)\n"+
				"Current selection: %c (%d, %d)",
			g[point1.y][point1.x], point1.x, point1.y, g[point2.y][point2.x], point2.x, point2.y)
	}
	controls := []control{
		{key: "← ↑ → ↓", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "Escape", description: "Cancel selection"},
	}

	draw(s, g, []vector2d{point1, point2}, generateText(), controls)
	selected := false
	point2Updated := point2
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
				point2Updated = point1
				point2Updated.y--
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'S' {
				//point2 = vector2d{
				//	x: point1.x,
				//	y: point1.y + 1,
				//}
				point2Updated = point1
				point2Updated.y++
			} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'A' {
				//point2 = vector2d{
				//	x: point1.x - 1,
				//	y: point1.y,
				//}
				point2Updated = point1
				point2Updated.x--
			} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'D' {
				//point2 = vector2d{
				//	x: point1.x + 1,
				//	y: point1.y,
				//}
				point2Updated = point1
				point2Updated.x++
			} else if ev.Key() == tcell.KeyEnter {
				selected = true
			} else if ev.Key() == tcell.KeyEscape {
				return emptyVector2d
			}
		}

		if isPointInsideGrid(point2Updated) {
			point2 = point2Updated
		}

		draw(s, g, []vector2d{point1, point2}, generateText(), controls)
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
// TODO: Somehow remove overlapping matches - noticed single match of 4 is counting as two matches
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
	horizontalFilters := make([][]vector2d, 0, (minMatchLength*2)+2)
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
		horizontalFilters = append(horizontalFilters, filter)

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
		horizontalFilters = append(horizontalFilters, filter)
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
	horizontalFilters = append(horizontalFilters, filter)

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
	horizontalFilters = append(horizontalFilters, filter)

	verticalFilters := make([][]vector2d, 0, len(horizontalFilters))
	// Copy horizontal filters but flip the x and y values
	for _, f := range horizontalFilters {
		fVertical := make([]vector2d, 0, len(f))
		for _, p := range f {
			fVertical = append(fVertical, vector2d{x: p.y, y: p.x})
		}
		verticalFilters = append(verticalFilters, fVertical)
	}

	return append(horizontalFilters, verticalFilters...)
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

	for y := gridHeight - 1; y >= 0; y-- {
		for x := 0; x < gridWidth; x++ {
			for _, f := range filters {
				// Don't need to compute size; could just check all filter's points are within grid
				filterSize := computeObjectSize(f)

				// Check filter would be inside the grid when positioned at current x,y coords
				if x >= gridWidth-filterSize.x+1 || y < filterSize.y-1 {
					continue
				}

				sameSymbol := true
				origin := vector2d{x: x, y: y}
				reference := f[0]
				referenceGridCoords := vector2d{x: origin.x + reference.x, y: origin.y - reference.y}
				fGridCoords := make([]vector2d, 0, len(f))
				for _, p := range f {
					pGridCoords := vector2d{x: origin.x + p.x, y: origin.y - p.y}
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
