package itertools

var (
	iterBuffer int = 10
)

func SetIterBuffer(buff int) {
	iterBuffer = buff
}

func GetIterBuffer() int {
	return iterBuffer
}
