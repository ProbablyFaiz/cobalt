package main

import (
	"pasado/parser"
)

func main() {
	// Prints "Hello, World!" 10 times
	for i := 0; i < 3; i++ {
		//println("Hello, World!")
		parser.Parse(`Hello, world!`)
	}
}
