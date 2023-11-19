package main

import (
	"fmt"
	"github.com/dustin/go-humanize/english"
	"github.com/gdamore/tcell"
	"unicode"
)

func swapPoints(s tcell.Screen, g *grid, potentialMatch []vector2d, score int) {
	point1 := vector2d{x: gridWidth / 2, y: gridHeight / 2} // Initialise point 1 to centre of grid
	point2 := emptyVector2d
	for point2 == emptyVector2d {
		point1 = selectFirstPoint(s, *g, potentialMatch, point1, score)
		point2 = selectSecondPoint(s, *g, point1, score)
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
		draw(s, *g, convertMatchesToPoints(matches), text, controls, score)
	} else {
		text := "Not swapping as swap would not result in a match; please try again"
		draw(s, *g, []vector2d{point1, point2}, text, controls, score)
	}

	waitForKeyPress(s)
}

func selectFirstPoint(s tcell.Screen, g grid, potentialMatch []vector2d, point1Initial vector2d, score int) vector2d {
	point1 := point1Initial

	generateText := func() string {
		return fmt.Sprintf(
			"Select two points to swap (selecting point 1)...\n\nCurrent selection: %c (%d, %d)",
			g[point1.y][point1.x], point1.x, point1.y)
	}
	controls := []control{
		{key: "← ↑ → ↓ / WASD", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "H", description: "Show hint"},
	}
	hintControls := []control{
		{key: "<Any key>", description: "Hide hint"},
	}

	draw(s, g, []vector2d{point1}, generateText(), controls, score)
	selected := false
	showPotentialMatch := false
	for !selected {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if showPotentialMatch {
				showPotentialMatch = false
			} else {
				if ev.Key() == tcell.KeyUp || unicode.ToLower(ev.Rune()) == 'w' {
					point1.y--
				} else if ev.Key() == tcell.KeyDown || unicode.ToLower(ev.Rune()) == 's' {
					point1.y++
				} else if ev.Key() == tcell.KeyLeft || unicode.ToLower(ev.Rune()) == 'a' {
					point1.x--
				} else if ev.Key() == tcell.KeyRight || unicode.ToLower(ev.Rune()) == 'd' {
					point1.x++
				} else if ev.Key() == tcell.KeyEnter {
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
			draw(s, g, potentialMatch, generateText()+"\nShowing hint", hintControls, score)
		} else {
			draw(s, g, []vector2d{point1}, generateText(), controls, score)
		}
	}

	return point1
}

func selectSecondPoint(s tcell.Screen, g grid, point1 vector2d, score int) vector2d {
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
		{key: "← ↑ → ↓ / WASD", description: "Move selection"},
		{key: "Enter", description: "Select"},
		{key: "Escape", description: "Cancel selection"},
	}

	draw(s, g, []vector2d{point1, point2}, generateText(), controls, score)
	selected := false
	point2Updated := point2
	for !selected {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyUp || unicode.ToLower(ev.Rune()) == 'w' {
				point2Updated = vector2d{
					x: point1.x,
					y: point1.y - 1,
				}
			} else if ev.Key() == tcell.KeyDown || unicode.ToLower(ev.Rune()) == 's' {
				point2Updated = vector2d{
					x: point1.x,
					y: point1.y + 1,
				}
			} else if ev.Key() == tcell.KeyLeft || unicode.ToLower(ev.Rune()) == 'a' {
				point2Updated = vector2d{
					x: point1.x - 1,
					y: point1.y,
				}
			} else if ev.Key() == tcell.KeyRight || unicode.ToLower(ev.Rune()) == 'd' {
				point2Updated = vector2d{
					x: point1.x + 1,
					y: point1.y,
				}
			} else if ev.Key() == tcell.KeyEnter {
				selected = true
			} else if ev.Key() == tcell.KeyEscape {
				return emptyVector2d
			}
		}

		if isPointInsideGrid(point2Updated) {
			point2 = point2Updated
		}

		draw(s, g, []vector2d{point1, point2}, generateText(), controls, score)
	}

	return point2
}

func waitForKeyPress(s tcell.Screen) {
	keyPressed := false
	for !keyPressed {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			keyPressed = true
		}
	}
}
