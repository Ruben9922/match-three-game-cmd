package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell"
	"sort"
	"strings"
	"unicode"
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
		drawRadioButtons([]gameType{Endless, LimitedMoves}, options.gameType, "Game type", "T"),
	}, "\n"))

	s.Show()
}

type radioButtonItem interface {
	comparable
	String() string
}

func drawRadioButtons[T radioButtonItem](options []T, selected T, label string, key string) string {
	var builder strings.Builder
	builder.WriteString(label)
	builder.WriteString(": ")
	for i, option := range options {
		if option == selected {
			builder.WriteString(option.String() + " [▪]")
		} else {
			builder.WriteString(option.String() + " [ ]")
		}

		if i != len(options)-1 {
			builder.WriteString(";")
		}

		builder.WriteString(" ")
	}
	builder.WriteString(fmt.Sprintf("(press %s)", strings.ToUpper(key)))

	return builder.String()
}

func updateTitleScreen(s tcell.Screen) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if unicode.ToLower(ev.Rune()) == 't' {
				options.gameType = toggleGameType(options.gameType)
				drawTitleScreen(s)
			} else {
				return
			}
		}
	}
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

	drawText(s, 0, gridHeight+1, screenWidth-1, gridHeight+1, defaultStyle,
		fmt.Sprintf("Score: %s", humanize.Comma(int64(score))))

	drawText(s, 0, gridHeight+3, screenWidth-1, gridHeight+3, defaultStyle,
		fmt.Sprintf("Game type: %s", options.gameType))

	if options.gameType == LimitedMoves {
		drawText(s, 0, gridHeight+4, screenWidth-1, screenHeight-1, defaultStyle,
			fmt.Sprintf("Remaining moves: %d", remainingMoveCount))
	}

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

func drawGameOverScreen(s tcell.Screen, g grid, score int) {
	const text = "Game over!\n\nNo more moves left."
	controls := []control{{key: "<Any key>", description: "Exit"}}

	draw(s, g, []vector2d{}, text, controls, score)
}

func drawQuitConfirmationScreen(s tcell.Screen, g grid, score int) {
	const text = "Are you sure you want to quit?\n\nAny game progress will be lost."
	controls := []control{
		{key: "Enter", description: "Quit"},
		{key: "<Any other key>", description: "Cancel"},
	}

	draw(s, g, []vector2d{}, text, controls, score)
}

func updateQuitConfirmationScreen(s tcell.Screen) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEnter {
				quit(s)
			}
			return
		}
	}
}
