package main

import (
	"math/rand"
	"slices"
)

func findEmptyPoints(g grid) []vector2d {
	emptyPoints := make([]vector2d, 0, gridWidth*gridHeight)
	for y := 0; y < gridHeight; y++ {
		for x := 0; x < gridWidth; x++ {
			if g[y][x] == emptySymbol {
				emptyPoints = append(emptyPoints, vector2d{x: x, y: y})
			}
		}
	}
	return emptyPoints
}

// todo: do initial refresh without animation and scoring
func refreshGrid(g *grid, r *rand.Rand, score *int, isScoring bool) bool {
	emptyPoints := findEmptyPoints(*g)
	if len(emptyPoints) == 0 {
		matches := findMatches(*g)
		if len(matches) == 0 {
			return true
		}

		if isScoring {
			matchesScore := computeScore(matches)
			*score += matchesScore
		}

		points := convertMatchesToPoints(matches)

		// Set points in matches to empty
		for _, p := range points {
			g[p.y][p.x] = emptySymbol
		}

		return false
	}

	// Shift symbols down and insert random symbol at top of column
	shiftPoint(g, r)

	return false
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

func shiftPoint(g *grid, r *rand.Rand) {
	emptyPoints := findEmptyPoints(*g)
	m := make(map[int][]int, gridWidth)
	for _, p := range emptyPoints {
		if m[p.x] == nil {
			m[p.x] = make([]int, 0, gridHeight)
		}

		m[p.x] = append(m[p.x], p.y)
	}

	for x, ys := range m {
		// Want to shift lower points first - hence getting the lowest point (point with highest y value)
		maxY := slices.Max(ys)

		for y := maxY; y > 0; y-- {
			g[y][x] = g[y-1][x]
		}
		g[0][x] = r.Intn(symbolCount)
	}
}
