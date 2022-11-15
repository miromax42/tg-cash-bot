package util

func OrDone[T any](done <-chan struct{}, input <-chan T) <-chan T {
	outStream := make(chan T)

	go func() {
		defer close(outStream)
		for {
			select {
			case <-done:
				return
			case val, ok := <-input:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case outStream <- val:
				}
			}
		}
	}()

	return outStream
}

func RepeatFn[T any](done <-chan struct{}, fn func() T) <-chan T {
	repeatStream := make(chan T)

	go func() {
		defer close(repeatStream)
		for {
			select {
			case <-done:
				return
			case repeatStream <- fn():
			}
		}
	}()

	return repeatStream
}
