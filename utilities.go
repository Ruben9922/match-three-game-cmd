package main

import (
	"github.com/gdamore/tcell"
	"math/rand"
)

func getRandomSymbol(r *rand.Rand) rune {
	index := r.Intn(len(symbols))
	return symbols[index]
}

func isPointInsideGrid(p vector2d) bool {
	return p.x >= 0 && p.x < gridWidth && p.y >= 0 && p.y < gridHeight
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func waitForKeyPress(s tcell.Screen) {
	for {
		ev := s.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		}
	}
}
