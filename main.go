package main

import (
	"fmt"
	"pasado/src"
)

func main() {
	//parser.Parse(`Hello, World!`)
	p1, _ := src.Parse(`5`)
	fmt.Printf("%#v\n", p1)
	p2, _ := src.Parse(`="Hello, World!"`)
	fmt.Printf("%#v\n", p2)
	p3, _ := src.Parse(`=5`)
	fmt.Printf("%#v\n", p3)
	p4, _ := src.Parse(`=Add(1, 2 ,3)`)
	fmt.Printf("%#v\n", p4)
	p5, _ := src.Parse(`=Concat("Hello", ",", "World!")`)
	fmt.Printf("%#v\n", p5)
}
