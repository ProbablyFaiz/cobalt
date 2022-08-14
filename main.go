package main

import (
	"fmt"
	"pasado/parser"
)

func main() {
	//parser.Parse(`Hello, World!`)
	p1, _ := parser.Parse(`5`)
	fmt.Printf("%#v\n", p1)
	p2, _ := parser.Parse(`="Hello, World!"`)
	fmt.Printf("%#v\n", p2)
	p3, _ := parser.Parse(`=5`)
	fmt.Printf("%#v\n", p3)
	p4, _ := parser.Parse(`=Add(1, 2 ,3)`)
	fmt.Printf("%#v\n", p4)
	p5, _ := parser.Parse(`=Concat("Hello", ",", "World!")`)
	fmt.Printf("%#v\n", p5)
}
