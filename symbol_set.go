package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

type symbolSet interface {
	fmt.Stringer
	formatSymbol(symbol int) string
	formatSymbolHighlighted(symbol int) string
}

type plainSymbolSet struct {
	name        string
	symbolRunes [symbolCount]rune
}

func (p plainSymbolSet) String() string {
	return p.name
}

func (p plainSymbolSet) getSymbolRune(symbol int) rune {
	if symbol < 0 || symbol >= symbolCount {
		return ' '
	}

	return p.symbolRunes[symbol]
}

func (p plainSymbolSet) formatSymbol(symbol int) string {
	symbolRune := p.getSymbolRune(symbol)
	return string(symbolRune)
}

func (p plainSymbolSet) formatSymbolHighlighted(symbol int) string {
	symbolRune := p.getSymbolRune(symbol)
	return lipgloss.NewStyle().Background(whiteColor).Render(string(symbolRune))
}

func newEmojiSymbolSet() plainSymbolSet {
	return plainSymbolSet{name: "Emojis", symbolRunes: [symbolCount]rune{'🍏', '🍇', '🍊', '🍋', '🍒', '🍓'}}
}

type colorSymbolSet struct {
	plainSymbolSet
	symbolColors [symbolCount]lipgloss.Color
}

func (c colorSymbolSet) getSymbolColor(symbol int) lipgloss.TerminalColor {
	if symbol < 0 || symbol >= symbolCount {
		return lipgloss.NoColor{}
	}

	return c.symbolColors[symbol]
}

func (c colorSymbolSet) formatSymbol(symbol int) string {
	color := c.getSymbolColor(symbol)
	symbolRune := c.getSymbolRune(symbol)
	return lipgloss.NewStyle().Foreground(color).Render(string(symbolRune))
}

func (c colorSymbolSet) formatSymbolHighlighted(symbol int) string {
	color := c.getSymbolColor(symbol)
	symbolRune := c.getSymbolRune(symbol)
	return lipgloss.NewStyle().Background(color).Render(string(symbolRune))
}

func newColorSymbolSet(name string, symbolRunes [symbolCount]rune) colorSymbolSet {
	return colorSymbolSet{
		plainSymbolSet: plainSymbolSet{name: name, symbolRunes: symbolRunes},
		symbolColors: [symbolCount]lipgloss.Color{
			"205",
			"34",
			"33",
			"220",
			"93",
			"37",
		},
	}
}

func newLetterSymbolSet() colorSymbolSet {
	return newColorSymbolSet("Letters", [symbolCount]rune{'A', 'B', 'C', 'D', 'E', 'F'})
}

func newShapeSymbolSet() colorSymbolSet {
	return newColorSymbolSet("Shapes", [symbolCount]rune{'▲', '■', '●', '★', '◆', '♥'})
}

func newNumberSymbolSet() colorSymbolSet {
	return newColorSymbolSet("Numbers", [symbolCount]rune{'1', '2', '3', '4', '5', '6'})
}
