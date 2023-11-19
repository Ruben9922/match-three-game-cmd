package main

import (
	"github.com/gdamore/tcell"
	"math/rand"
	"sort"
	"time"
)

// todo: do initial refresh without animation and scoring
func refreshGrid(s tcell.Screen, g *grid, r *rand.Rand, score *int, isScoring bool) {
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
	draw(s, *g, []vector2d{}, text, controls, *score)

	for {
		matches := findMatches(*g)
		if len(matches) == 0 {
			break
		}

		if isScoring {
			matchesScore := computeScore(matches)
			*score += matchesScore
		}

		points := convertMatchesToPoints(matches)

		// Shifting algorithm assumes points are unique; duplicate points will cause strange behaviour
		points = removeDuplicatePoints(points)

		// Set points in matches to empty
		for _, p := range points {
			g[p.y][p.x] = emptySymbol
		}

		waitForKeyPressOrTimeout()
		draw(s, *g, []vector2d{}, text, controls, *score)

		// Shift symbols down and insert random symbol at top of column
		// Instead of manually updating points list, could maybe just search through grid for empty points
		// Assumes points are unique - duplicate points will cause strange behaviour
		// Want to shift lower points first - hence sorting such that lower points (points with higher y) come first
		sortPoints(points)
		for len(points) > 0 {
			shiftPoint(g, &points, r)

			waitForKeyPressOrTimeout()
			draw(s, *g, []vector2d{}, text, controls, *score)
		}
	}
}

// Match algorithm works in a way that scores the player the most points
// * Prefers longer matches (checks for all possible matches and chooses the longest one) - to score the player more points
// * Prefers lower down matches - idea is that this would score the player more points as more pieces falling means potentially more "automatic" matches
// * "Maximal munch" behaviour - matches will be as long as possible; matches can be longer than the minimum match length
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

func computeScore(matches []match) int {
	totalSymbolCount := 0
	for _, m := range matches {
		totalSymbolCount += m.length
	}
	score := totalSymbolCount * scorePerMatchedSymbol
	return score
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

func sortPoints(points []vector2d) {
	sort.Slice(points, func(i, j int) bool {
		if points[i].x == points[j].x {
			return points[i].y > points[j].y
		}
		return points[i].x < points[j].x
	})
}

func shiftPoint(g *grid, points *[]vector2d, r *rand.Rand) {
	updatedPoints := *points

	p := updatedPoints[0]
	updatedPoints = updatedPoints[1:]

	for y := p.y; y > 0; y-- {
		g[y][p.x] = g[y-1][p.x]
	}
	g[0][p.x] = getRandomSymbol(r)

	// Shift down remaining points in same column to account for shifting of corresponding empty points in grid
	// p1.y++ for each point p1 in this column (with same x)
	for i := 0; i < len(updatedPoints); i++ {
		p1 := &updatedPoints[i]
		// If point is in the same column and (strictly) above current point
		if p1.x == p.x && p1.y < p.y {
			p1.y++
		}
	}

	*points = updatedPoints
}
