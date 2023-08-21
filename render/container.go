package render

type container struct {
	*Context
}

func Container() container {
	return container{&Context{}}
}

func (container) Value(string) {} // forbidden, no-op
