package modal

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/londek/reactea"
)

type ModalComponent[TReturn any] interface {
	reactea.Component

	initModal(chan<- ModalResult[TReturn], *Controller)
	Return(ModalResult[TReturn]) tea.Cmd
}
