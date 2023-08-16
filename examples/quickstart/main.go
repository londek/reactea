package main

import (
	"github.com/londek/reactea"
	"github.com/londek/reactea/example/quickstart/app"
)

func main() {
	// reactea.NewProgram initializes program with
	// "translation layer", so Reactea components work
	program := reactea.NewProgram(app.New())

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
