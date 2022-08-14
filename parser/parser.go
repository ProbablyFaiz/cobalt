package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"pasado/base"
	"strings"
)

type Formula struct {
	Argument     *Argument    `Eq @@`
	ValueLiteral *BareLiteral `| @@`
}

type Argument struct {
	ValueLiteral *ArgLiteral   `@@`
	FunctionCall *FunctionCall `| @@`
}

type BareLiteral struct {
	IntLiteral *int `@Int`
	//StringLiteral *string `| @BareString`
}

type ArgLiteral struct {
	IntLiteral    *int    `@Int`
	StringLiteral *string `| @String`
}

type FunctionCall struct {
	Name *string     `@Ident`
	Args []*Argument `LPar ( @@ Sep )* @@? RPar`
}

func Parse(input string) (base.FormulaNode, error) {
	langLexer := lexer.MustSimple([]lexer.SimpleRule{
		{"Eq", `=`},
		{"Sep", `,`},
		{"LPar", `\(`},
		{"RPar", `\)`},
		{"Ident", `[a-zA-Z_]\w*`},
		{"Int", `[-+]?\d+`},
		{"String", `"(\\"|[^"])*"`},
		//{"BareString", `.*`},
		{"Whitespace", `\s+`},
	})

	parser, err := participle.Build[Formula](
		participle.Lexer(langLexer),
		participle.CaseInsensitive("Ident"),
		participle.Elide("Whitespace"),
		participle.Unquote("String"))
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%#v\n", langLexer.Symbols())
	//tokens, err := parser.Lex("", strings.NewReader(input))
	//if err != nil {
	//	panic(err)
	//}
	//for _, token := range tokens {
	//	fmt.Printf("%#v\n", token.Type)
	//}
	formula, err := parser.Parse("", strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	return formula.toAst(), nil
}

func (formula *Formula) toAst() base.FormulaNode {
	if formula.Argument != nil {
		return formula.Argument.toAst()
	}
	return formula.ValueLiteral.toAst()
}

func (argument *Argument) toAst() base.FormulaNode {
	if argument.FunctionCall != nil {
		return argument.FunctionCall.toAst()
	}
	return argument.ValueLiteral.toAst()
}

func (literal *BareLiteral) toAst() base.FormulaNode {
	return &base.LiteralNode{Value: *literal.IntLiteral}
}

func (literal *ArgLiteral) toAst() base.FormulaNode {
	if literal.StringLiteral != nil {
		return &base.LiteralNode{Value: *literal.StringLiteral}
	}
	return &base.LiteralNode{Value: *literal.IntLiteral}
}

func (call *FunctionCall) toAst() base.FormulaNode {
	newArgs := make([]base.FormulaNode, len(call.Args))
	for i, arg := range call.Args {
		newArgs[i] = arg.toAst()
	}
	return &base.FunctionNode{Name: strings.ToLower(*call.Name), Args: newArgs}
}
