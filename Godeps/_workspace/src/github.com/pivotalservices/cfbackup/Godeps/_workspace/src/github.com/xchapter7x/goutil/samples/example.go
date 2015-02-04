package main

import (
	"container/list"
	"container/ring"
	"fmt"
	"github.com/xchapter7x/goutil/itertools"
)

func f(i int, v string) string {
	fmt.Println(i, v)
	return v
}

func mf(i, v string) string {
	fmt.Println(i, v)
	return v
}

func t(i interface{}) {
	for x := range itertools.Iterate(i) {
		fmt.Println(x)
	}
}

func main() {
	s := []string{"asdf", "asdfasdf", "geeeg", "gggggggg"}
	m := map[string]string{"a": "asdf", "b": "asdfasdf", "c": "geeeg", "d": "gggggggg"}
	itertools.Each(&s, f)
	itertools.Each(&m, mf)
	fmt.Println("\n\nbegin concurrent map\n\n")
	itertools.CEach(s, f)
	itertools.CEach(m, mf)

	fmt.Println("\n\nFilter Sample\n\n")

	f := itertools.Filter(s, func(i, v interface{}) bool {
		il := map[int]int{1: 1, 2: 2}
		_, ok := il[i.(int)]
		return ok
	})

	for i := range f {
		fmt.Println(i)
	}

	fmt.Println("\n\nConcurrent Filter Sample\n\n")

	fC := itertools.CFilter(s, func(i, v interface{}) bool {
		il := map[int]int{1: 1, 2: 2}
		_, ok := il[i.(int)]
		return ok
	})

	for i := range fC {
		fmt.Println(i)
	}

	fmt.Println("\n\nlets iterate a string")

	t("this is a test")

	fmt.Println("\n\nlets iterate a list")

	l := list.New()
	l.PushFront(1)
	l.PushFront(2)
	l.PushFront(3)
	l.PushFront(4)
	l.PushFront(5)
	l.PushFront(6)
	t(l)

	fmt.Println("\n\nlets iterate a ring")

	r := ring.New(10)
	z := 100
	r.Value = z
	for p := r.Next(); p != r; p = p.Next() {
		z -= 10
		p.Value = z
	}
	t(r)

    fmt.Println("test the ziplong")
    for z := range itertools.ZipLongest("-", "abcdefghijk", "ABCde", "aBCDefg") {
        fmt.Println(z)
    }

    fmt.Println("test the zip")
    for z := range itertools.Zip("-", "abcdefghijk", "ABCde", "aBCDefg") {
        fmt.Println(z)
    }
}
