package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"slices"
	"unicode/utf8"
)

type titleView struct{}

type titleViewKeyMap struct {
	Quit            key.Binding
	ToggleGameType  key.Binding
	ToggleSymbolSet key.Binding
	Start           key.Binding
}

var titleViewKeys = titleViewKeyMap{
	Quit: sharedKeys.Quit,
	ToggleGameType: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "change game type"),
	),
	ToggleSymbolSet: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "change symbol set"),
	),
	Start: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "start"),
	),
}

var gameTypes = []gameType{Endless, LimitedMoves}
var symbolSets = []symbolSet{newEmojiSymbolSet(), newShapeSymbolSet(), newLetterSymbolSet(), newNumberSymbolSet()}

func (k titleViewKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Start, k.ToggleGameType, k.ToggleSymbolSet, k.Quit}
}

func (k titleViewKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.ToggleGameType, k.ToggleSymbolSet, k.Quit},
	}
}

func (tv titleView) draw(m model) string {
	const titlePart1 = "  __  __       _       _       _____ _                   \n |  \\/  | __ _| |_ ___| |__   |_   _| |__  _ __ ___  ___ \n | |\\/| |/ _` | __/ __| '_ \\    | | | '_ \\| '__/ _ \\/ _ \\\n | |  | | (_| | || (__| | | |   | | | | | | | |  __/  __/\n |_|  |_|\\__,_|\\__\\___|_| |_|   |_| |_| |_|_|  \\___|\\___|"
	const titlePart2 = "   ____                      \n  / ___| __ _ _ __ ___   ___ \n | |  _ / _` | '_ ` _ \\ / _ \\\n | |_| | (_| | | | | | |  __/\n  \\____|\\__,_|_| |_| |_|\\___|"
	const text = "Press enter key to start..."

	gameTypeRadioButtons := drawRadioButtons(gameTypes, m.options.gameType, "Game type", titleViewKeys.ToggleGameType)
	symbolSetRadioButtons := drawRadioButtons(symbolSets, m.symbolSet, "Symbol set", titleViewKeys.ToggleSymbolSet)
	helpView := m.help.View(titleViewKeys)
	return lipgloss.JoinVertical(lipgloss.Center,
		titlePart1,
		lipgloss.JoinHorizontal(lipgloss.Bottom,
			lipgloss.NewStyle().MarginLeft(utf8.RuneCountInString(version)+2).MarginRight(2).Render(titlePart2),
			version,
		),
		"",
		text,
		"",
		gameTypeRadioButtons,
		symbolSetRadioButtons,
		"",
		helpView,
	)
}

func getNextElement[T comparable](slice []T, element T) T {
	index := slices.Index(slice, element)
	if index == -1 {
		return slice[0]
	}

	return slice[(index+1)%len(slice)]
}

func (tv titleView) update(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, titleViewKeys.Quit):
			return showQuitConfirmationView(m)

		case key.Matches(msg, titleViewKeys.ToggleGameType):
			m.options.gameType = getNextElement(gameTypes, m.options.gameType)
		case key.Matches(msg, titleViewKeys.ToggleSymbolSet):
			m.symbolSet = getNextElement(symbolSets, m.symbolSet)
		case key.Matches(msg, titleViewKeys.Start):
			m.grid = newGridWithMatchesRemoved(m.rand)
			ensurePotentialMatch(&m.grid, m.rand)

			// todo: don't need to navigate to refresh grid view
			m.view = refreshGridView{}
			m.point1 = emptyVector2d
			m.point2 = emptyVector2d

			return m, tickCmd()
		}
	}

	return m, nil
}
