package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
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

func (p plainSymbolSet) getSymbolRune(symbol int) string {
	emptySymbolRune := strings.Repeat(" ", lipgloss.Width(string(p.symbolRunes[0])))

	if symbol < 0 || symbol >= symbolCount {
		return emptySymbolRune
	}

	return string(p.symbolRunes[symbol])
}

func (p plainSymbolSet) formatSymbol(symbol int) string {
	return p.getSymbolRune(symbol)
}

func (p plainSymbolSet) formatSymbolHighlighted(symbol int) string {
	symbolRune := p.getSymbolRune(symbol)
	return lipgloss.NewStyle().Background(whiteColor).Render(symbolRune)
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
	return lipgloss.NewStyle().Foreground(color).Render(symbolRune)
}

func (c colorSymbolSet) formatSymbolHighlighted(symbol int) string {
	color := c.getSymbolColor(symbol)
	symbolRune := c.getSymbolRune(symbol)
	return lipgloss.NewStyle().Background(color).Render(symbolRune)
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
