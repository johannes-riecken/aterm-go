package main

func exampleFunc(foo int, bar []int, baz int, quux []int) {
	var foobar int
	foobar, foo, quux = len(bar), len(quux), []int{foo, baz}
}
