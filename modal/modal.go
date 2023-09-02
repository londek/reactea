package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Modal[T any] struct {
	resultChan chan<- ModalResult[T]
}

type ModalComponent[TReturn any] interface {
	reactea.Component

	initModal(chan<- ModalResult[TReturn])
	Return(ModalResult[TReturn]) tea.Cmd
}

//lint:ignore U1000 This function is used, but through interface
func (modal *Modal[T]) initModal(resultChan chan<- ModalResult[T]) {
	modal.resultChan = resultChan
}

// Returns nil tea.Cmd that allows for following syntactic sugar:
// return modal.Return(result)
func (modal *Modal[T]) Return(result ModalResult[T]) tea.Cmd {
	modal.resultChan <- result
	return nil
}

func (modal *Modal[T]) Ok(result T) tea.Cmd {
	return modal.Return(Ok[T](result))
}

func (modal *Modal[T]) Error(err error) tea.Cmd {
	return modal.Return(Error[T](err))
}

func Show[T any](controller *Controller, modal ModalComponent[T]) ModalResult[T] {
	resultChan := make(chan ModalResult[T])

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

	return c, c.Run(f)
}
