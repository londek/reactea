package modal

type ModalResult[T any] struct {
	Return T
	Err    error
}

// Allows for value, err := result.Get() syntactic sugar
// instead of manually destructuring fields
func (result *ModalResult[T]) Get() (T, error) {
	return result.Return, result.Err
}

func Ok[T any](ret T) ModalResult[T] {
	return ModalResult[T]{ret, nil}
}

func Error[T any](err error) ModalResult[T] {
	var ret T
	return ModalResult[T]{ret, err}
}
