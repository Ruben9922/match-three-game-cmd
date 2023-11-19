package main

import "sort"

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
