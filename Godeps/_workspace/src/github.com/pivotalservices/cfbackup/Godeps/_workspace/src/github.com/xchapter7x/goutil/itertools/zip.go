package itertools

func ZipLongest(fillValue string, args ...interface{}) (out chan []interface{}) {
	out = make(chan []interface{}, GetIterBuffer())

	go func() {
		defer close(out)
		var argsSlice [][]interface{}
		maxSliceLength := 0

		for _, arg := range args {
			var argSlice []interface{}

			for p := range Iterate(arg) {

				if currentIndexGreaterThanMaxLength(p, maxSliceLength) {
					maxSliceLength++
				}
				argSlice = append(argSlice, p.Second)
			}
			argsSlice = append(argsSlice, argSlice)
		}

		for i := 0; i < maxSliceLength; i++ {
			var row []interface{}

			for _, a := range argsSlice {

				if balancedSliceLength(a, i) {
					row = append(row, a[i])

				} else {
					row = append(row, fillValue)
				}
			}
			out <- row
		}
	}()
	return
}

func Zip(fillValue string, args ...interface{}) (out chan []interface{}) {
	out = make(chan []interface{}, GetIterBuffer())

	go func() {
		defer close(out)
		var argsSlice [][]interface{}
		maxSliceLength := 0

		for _, arg := range args {
			var argSlice []interface{}

			for p := range Iterate(arg) {

				if currentIndexGreaterThanMaxLength(p, maxSliceLength) {
					maxSliceLength++
				}
				argSlice = append(argSlice, p.Second)
			}
			argsSlice = append(argsSlice, argSlice)
		}

		for i := 0; i < maxSliceLength; i++ {
			var row []interface{}
			unevenSlices := false

			for _, a := range argsSlice {

				if balancedSliceLength(a, i) {
					row = append(row, a[i])

				} else {
					unevenSlices = true
					break
				}
			}

			if unevenSlices {
				break
			} else {
				out <- row
			}
		}
	}()
	return
}

func balancedSliceLength(a []interface{}, i int) bool {
	return len(a)-1 >= i
}

func currentIndexGreaterThanMaxLength(p Pair, maxSliceLength int) bool {
	return (p.First.(int) + 1) > maxSliceLength
}
