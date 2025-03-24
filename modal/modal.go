package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type Modal[T any] struct {
	ch chan<- ModalResult[T]
	c  *Controller
}

//lint:ignore U1000 This function is used, but through interface
func (modal *Modal[T]) initModal(resultChan chan<- ModalResult[T], controller *Controller) {
	modal.ch = resultChan
	modal.c = controller
}

func (modal *Modal[T]) Return(result ModalResult[T]) tea.Cmd {
	modal.ch <- result
	modal.c.w.Wait()
	return reactea.Rerender
}

func (modal *Modal[T]) Ok(result T) tea.Cmd {
	return modal.Return(Ok(result))
}

func (modal *Modal[T]) Error(err error) tea.Cmd {
	return modal.Return(Error[T](err))
}

func Show[T any](c *Controller, modal ModalComponent[T]) ModalResult[T] {
	c.w.Signal()

	resultChan := make(chan ModalResult[T])

	modal.initModal(resultChan, c)

	c.cond.L.Lock()

	c.modal = modal
	c.initCmd = modal.Init()

	c.cond.Broadcast()
	c.cond.L.Unlock()

	result := <-resultChan

	c.shouldDestruct = true

	return result
}

func Get[T any](c *Controller, modal ModalComponent[T]) T {
	c.w.Signal()

	resultChan := make(chan ModalResult[T])

	modal.initModal(resultChan, c)

	c.cond.L.Lock()

	c.modal = modal
	c.initCmd = modal.Init()

	c.cond.Broadcast()
	c.cond.L.Unlock()

	result := <-resultChan

	c.shouldDestruct = true

	return result.Return
}
