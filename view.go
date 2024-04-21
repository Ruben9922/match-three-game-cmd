package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type view interface {
	update(msg tea.KeyMsg, m model) (model, tea.Cmd)
	draw(m model) string
}

//type keyMap struct {
//	Up    key.Binding
//	Down  key.Binding
//	Left  key.Binding
//	Right key.Binding
//	Help  key.Binding
//	Quit  key.Binding
//}
//
//func (k keyMap) ShortHelp() []key.Binding {
//	return []key.Binding{k.Help, k.Quit}
//}
//
//func (k keyMap) FullHelp() [][]key.Binding {
//	return [][]key.Binding{
//		{k.Up, k.Down, k.Left, k.Right},
//		{k.Help, k.Quit},
//	}
//}q
