# reactea

[![Go Reference](https://pkg.go.dev/badge/github.com/londek/reactea.svg)](https://pkg.go.dev/github.com/londek/reactea)
![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability-percentage/Londek/reactea)

Rather simple **Bubbletea companion** for **handling hierarchy** and support for **lifting state up.**
It Reactifies Bubbletea philosophy and makes it especially easy to work with in bigger projects.

For me, personally - It's a must in project with multiple pages and component communication

## Installation

`go get github.com/londek/reactea`

## Info

Always return `reactea.Destroy` instead of `tea.Quit` in order to follow our convention

The goal is to create components which are

- dimensions-aware
- propful
- easy to lift the state up

Most info is currently in source code so I suggest checking it out
