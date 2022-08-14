package parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"strings"
)

type Formula struct {
	Argument *Argument `"="@@`
	//Value    *BareLiteral `| @@`
}

type Argument struct {
	Value        *ArgLiteral   `@@`
	FunctionCall *FunctionCall `| @@`
}

type BareLiteral struct {
	IntLiteral    *int    `@Int`
	StringLiteral *string `| @String`
}

type ArgLiteral struct {
	StringLiteral *string `"@String"`
	IntLiteral    *int    `| @Int`
}

type FunctionCall struct {
	Name *string     `@Ident`
	Args []*Argument `"(" @@* ")"`
}

func Parse(input string) {
	parser, err := participle.Build[Formula]()
	if err != nil {
		panic(err)
	}
	_, err = parser.Lex("", strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	formula, err := parser.Parse("", strings.NewReader(input))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", *formula.Argument.Value.IntLiteral)
}
