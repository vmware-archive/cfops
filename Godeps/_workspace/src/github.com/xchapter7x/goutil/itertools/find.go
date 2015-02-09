package itertools

func Find(iterator interface{}, compareFunctor func(Pair) bool) (out Pair) {
	for pair := range Iterate(iterator) {

		if compareFunctor(pair) {
			out = pair
			break
		}
	}
	return
}
