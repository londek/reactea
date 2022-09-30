package main

import (
	"github.com/londek/reactea"
	"github.com/londek/reactea/example/app"
)

func main() {
	// reactea.NewProgram initializes program with
	// "translation layer", so Reactea components work
	program := reactea.NewProgram(app.New())

	if err := program.Start(); err != nil {
		panic(err)
	}
}
