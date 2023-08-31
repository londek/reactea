package modal

type ModalResult[T any] struct {
	Ok     bool
	Return T
	Err    error
}

type ModalResultKind uint8

const (
	OkResult ModalResultKind = iota
	BadResult
	ErrResult
)

func Ok[T any](ret T) ModalResult[T] {
	return ModalResult[T]{true, ret, nil}
}

func Fail[T any]() ModalResult[T] {
	var ret T
	return ModalResult[T]{false, ret, nil}
}

func Error[T any](err error) ModalResult[T] {
	var ret T
	return ModalResult[T]{false, ret, err}
}
