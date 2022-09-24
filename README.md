# reactea

[![Go Reference](https://pkg.go.dev/badge/github.com/londek/reactea.svg)](https://pkg.go.dev/github.com/londek/reactea)
![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability-percentage/Londek/reactea)

Rather simple **Bubbletea companion** for **handling hierarchy** and support for **lifting state up.**
It Reactifies Bubbletea philosophy and makes it especially easy to work with in bigger projects.

For me, personally - It's a must in project with multiple pages and component communication

## Installation

`go get -u github.com/londek/reactea`

## Info

Always return `reactea.Destroy` instead of `tea.Quit` in order to follow our convention

The goal is to create components which are

- dimensions-aware
- propful
- easy to lift the state up

Most info is currently in source code so I suggest checking it out

## Reactea Routes API

Routes API allows developers for easy creation of multi-page apps.
Routes are kind of substitute for window.Location inside bubbletea

### reactea.CurrentRoute() Route

Returns **copy** of current route

### reactea.LastRoute() Route

Returns **copy** of last route

### reactea.WasRouteChanged() bool

returns `LastRoute() != CurrentRoute()`

## Router Component

Router Component is basic implementation of how routing could look in your application.
It doesn't support wildcards yet or relative pathing. All data is provided from within props

### router.Props

router.Props is a map of route initializers keyed by routes serialized to strings following format `r1/r2/r3...etc`

What is `RouteInitializer`?

`RouteInitializer` is function that initializes the current route component
