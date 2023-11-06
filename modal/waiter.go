package modal

type waiter chan struct{}

func (w waiter) Wait() {
	w <- struct{}{}
}

func (w waiter) Signal() {
	select {
	case <-w:
		return
	default:
		return
	}
}
