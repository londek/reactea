package main

import (
	"github.com/londek/reactea"
	"github.com/londek/reactea/example/dynamicRoutes/app"
)

func main() {
	program := reactea.NewProgram(app.New())

	if err := program.Start(); err != nil {
		panic(err)
	}
}
