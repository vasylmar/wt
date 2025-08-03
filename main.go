package main

import (
	"log"
	"wt/internals/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.NewCli())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
