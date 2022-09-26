# reactea

[![Go Reference](https://pkg.go.dev/badge/github.com/londek/reactea.svg)](https://pkg.go.dev/github.com/londek/reactea)
![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability-percentage/Londek/reactea)
[![Go Report Card](https://goreportcard.com/badge/github.com/londek/reactea)](https://goreportcard.com/report/github.com/londek/reactea)

Rather simple **Bubbletea companion** for **handling hierarchy** and support for **lifting state up.**\
It Reactifies Bubbletea philosophy and makes it especially easy to work with in bigger projects.

For me, personally - It's a must in project with multiple pages and component communication

Check our example code [right here!](/example)

## Installation

`go get -u github.com/londek/reactea`

## General info

The goal is to create components which are

- dimensions-aware (especially unify all setSize conventions)
- propful
- easy to lift the state up
- able to communicate with parent without importing it (I spent too many hours solving import cycles hehe)
- easier to code
- all of that without code duplication

The extreme performance is not main goal of this package, because either way Bubbletea\
refresh rate is only 60hz and 50 allocations in entire **runtime** won't really hurt anyone.\
Most info is currently in source code so I suggest checking it out

Always return `reactea.Destroy` instead of `tea.Quit` in order to follow our convention\

## Component lifecycle

![Component lifecycle image](.github/lifecycle-diagram.png)

reactea takes pointer approach for components
making state modifiable in any lifecycle method\
**There are also 2 additional lifecycle methods: [AfterUpdate()](#afterupdate) and [UpdateProps()](#updateprops)**

### AfterUpdate()

AfterUpdate is the only lifecycle method that is not controlled by parent. It's called right after root component finishes Update(). Components should queue itself with `reactea.AfterUpdate(component)` in Update()

### UpdateProps()

UpdateProps() is lifecycle method that derives state from props, It can happen anytime during lifecycle. Usually called by Init()

### Notes

Update **IS NOT** guaranteed to be called on first-run, Init() for most part is, and critical logic should be there

Lifecycle is **(almost, see [AfterUpdate()](#afterupdate)) fully controlled by parent component** making graph above fully theoritical and possibly invalid for third-party components

## Example code

Check our example code [right here!](/example)

## Stateless components

Stateless components are represented by following function types

|   | Renderer | ProplessRenderer | DumbRenderer |
|----------------|----------|------------------|--------------|
| **Properties** | ✅        | ❌                | ❌            |
| **Dimensions** | ✅        | ✅                | ❌            |
| **Type** | `func(TProps, int, int) string`        | `func(int, int) string`                | `func() string`           |

There are many utility functions for transforming stateless into stateful components or for rendering any component without knowing its type (`reactea.RenderAny`, `reactea.RenderPropless`)

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
