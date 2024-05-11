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

func newGridWithMatchesRemoved(r *rand.Rand) grid {
	g := newGrid(r)
	removeMatches(&g, r)
	return g
}

func removeMatches(g *grid, r *rand.Rand) {
	finished := false
	for !finished {
		finished = refreshGrid(g, r, nil)
	}
}

func ensurePotentialMatch(g *grid, r *rand.Rand) {
	potentialMatch := findPotentialMatch(*g)
	for len(potentialMatch) == 0 {
		// Check if there are any possible matches; if no possible matches then create a new grid
		*g = newGridWithMatchesRemoved(r)

		potentialMatch = findPotentialMatch(*g)
	}
}

func refreshGrid(g *grid, r *rand.Rand, score *int) bool {
	emptyPoints := findEmptyPoints(*g)
	if len(emptyPoints) == 0 {
		matches := findMatches(*g)
		if len(matches) == 0 {
			return true
		}

		if score != nil {
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

func findMatches(g grid) []match {
	directions := []vector2d{
		{x: 1, y: 0},
		{x: 0, y: 1},
	}
	matches := make([]match, 0, 10)
	for _, d := range directions {
		offset := vector2d{
			x: maxInt((d.x*minMatchLength)-1, 0),
			y: maxInt((d.y*minMatchLength)-1, 0),
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

				m := newMatch(originPoint, d, matchLength)
				if matchLength >= minMatchLength {
					matches = updateMatches(matches, m)
				}
			}
		}
	}

	return matches
}

// This fixes an issue where longer matches (longer than `minMatchLength`) were being counted more than once
func updateMatches(matches []match, newMatch match) []match {
	updatedMatches := make([]match, 0, len(matches))
	for _, existingMatch := range matches {
		newMatchPoints := convertMatchesToPoints([]match{newMatch})
		existingMatchPoints := convertMatchesToPoints([]match{existingMatch})

		// If new match is a subset of any existing match, then don't add it because it's not needed
		if isSubset(newMatchPoints, existingMatchPoints) {
			return matches
		}

		// If any existing match is a subset of the new match, then remove it as it will be replaced by the new match
		// I.e. only keep existing matches which aren't a subset of the new match
		if !isSubset(existingMatchPoints, newMatchPoints) {
			updatedMatches = append(updatedMatches, existingMatch)
		}
	}
	updatedMatches = append(updatedMatches, newMatch)
	return updatedMatches
}

func isSubset[T comparable](possibleSubset, s []T) bool {
	if len(possibleSubset) > len(s) {
		return false
	}

	// Populate map from `s`
	m := make(map[T]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}

	// Check all values in `possibleSubset` are also present in the map `m`
	for _, v := range possibleSubset {
		if _, present := m[v]; !present {
			return false
		}
	}
	return true
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
