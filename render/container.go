package ruler

type container struct {
	*renderContext
}

func Container() container {
	return container{&renderContext{}}
}

func (container) Value(string) {} // forbidden, no-op
