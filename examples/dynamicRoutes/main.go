package main

import (
	"github.com/londek/reactea"
	"github.com/londek/reactea/examples/dynamicRoutes/app"
)

func main() {
	program := reactea.NewProgram(app.New())

	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
