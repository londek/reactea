package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Modal[T any] struct {
	resultChan chan<- T
}

type ModalComponent[TReturn, TProps any] interface {
	reactea.Component[TProps]

	initModal(chan<- TReturn)
	Return(TReturn)
}

type SomeModalComponent[TReturn any] interface {
	reactea.SomeComponent

	Init() tea.Cmd

	initModal(chan<- TReturn)
	Return(TReturn)
}

func (modal *Modal[T]) initModal(resultChan chan<- T) {
	modal.resultChan = resultChan
}

func (modal *Modal[T]) Return(result T) {
	modal.resultChan <- result
}

func ShowProps[T, U any](controller *Controller, modal ModalComponent[T, U], props U) T {
	resultChan := make(chan T)

	modal.initModal(resultChan)
	controller.show(modal, modal.Init(props))

	result := <-resultChan

	controller.hide()

	return result
}

func Show[T any](controller *Controller, modal SomeModalComponent[T]) T {
	resultChan := make(chan T)

	modal.initModal(resultChan)
	controller.show(modal, modal.Init())

	result := <-resultChan

	controller.hide()

	return result
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Run(f func(*Controller) tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return f(c)
	}
}

func Execute(f func(*Controller) tea.Cmd) (*Controller, tea.Cmd) {
	c := NewController()

	return c, tea.Batch(c.Init(reactea.NoProps{}), c.Run(f))
}
