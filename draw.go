package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	"sort"
	"strings"
	"unicode/utf8"
)

func drawTitleScreen(s tcell.Screen) {
	s.Clear()

	screenWidth, screenHeight := s.Size()

	const titlePart1 = "  __  __       _       _       _____ _                   \n |  \\/  | __ _| |_ ___| |__   |_   _| |__  _ __ ___  ___ \n | |\\/| |/ _` | __/ __| '_ \\    | | | '_ \\| '__/ _ \\/ _ \\\n | |  | | (_| | || (__| | | |   | | | | | | | |  __/  __/\n |_|  |_|\\__,_|\\__\\___|_| |_|   |_| |_| |_|_|  \\___|\\___|"
	const titlePart2 = "   ____                      \n  / ___| __ _ _ __ ___   ___ \n | |  _ / _` | '_ ` _ \\ / _ \\\n | |_| | (_| | | | | | |  __/\n  \\____|\\__,_|_| |_| |_|\\___|"
	const text = "\n Press any key to start..."

	drawText(s, 0, 0, screenWidth-1, screenHeight-1, defaultStyle, strings.Join([]string{
		titlePart1,
		titlePart2,
		text,
	}, "\n"))
}

func draw(s tcell.Screen, g grid, selectedPoints []vector2d, text string, controls []control, score int) {
	s.Clear()

	// Draw grid
	drawGrid(s, g, selectedPoints)

	// Draw text
	screenWidth, screenHeight := s.Size()
	const textOffsetX = (gridWidth * 2) + 3
	drawText(s, textOffsetX, 0, screenWidth-1, 5, defaultStyle, text)

	// Draw controls
	drawControls(s, controls, textOffsetX, 6)

	drawText(s, 0, gridHeight+1, screenWidth-1, screenHeight-1, defaultStyle,
		fmt.Sprintf("Score: %s", humanize.Comma(int64(score))))

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
