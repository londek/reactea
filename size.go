package reactea

// Currently used as utility

type Size int

func (size *Size) Take(n Size) Size {
	if *size < Size(n) {
		taken := *size
		*size = 0
		return taken
	}

	*size -= Size(n)

	return n
}

func (size *Size) TakeUnsafe(n Size) Size {
	*size -= Size(n)
	return n
}

func (size *Size) Sub(n Size) *Size {
	if *size < Size(n) {
		*size = 0
		return size
	}

	*size -= Size(n)
	return size
}

func (size *Size) SubUnsafe(n Size) *Size {
	*size -= Size(n)
	return size
}

func (size *Size) Remaining() Size {
	return *size
}
