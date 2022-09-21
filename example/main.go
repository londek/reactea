package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
	"github.com/londek/reactea/example/app"
)

func main() {
	// reactea.Model is a "translation layer" between
	// reactea components and bubbletea models
	program := tea.NewProgram(reactea.New(app.New()))

	if err := program.Start(); err != nil {
		panic(err)
	}
}
