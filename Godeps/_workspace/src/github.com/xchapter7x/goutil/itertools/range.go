package itertools

func Range(low, high int) (out chan int) {
	out = make(chan int, GetIterBuffer())

	go func() {
		defer close(out)

		for i := low; i <= high; i++ {
			out <- i
		}
	}()
	return
}
