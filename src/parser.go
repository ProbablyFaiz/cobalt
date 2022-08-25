package src

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type PFormula struct {
	Argument     *PArgument    `Eq @@`
	ValueLiteral *PBareLiteral `| @@`
}

type PArgument struct {
	ValueLiteral *PArgLiteral   `@@`
	FunctionCall *PFunctionCall `| @@`
}

type PBareLiteral struct {
	IntLiteral *int `@Int`
	//StringLiteral *string `| @BareString`
}

type PArgLiteral struct {
	IntLiteral    *int    `@Int`
	StringLiteral *string `| @String`
}

type PFunctionCall struct {
	Name *string      `@Ident`
	Args []*PArgument `LPar ( @@ Sep )* @@? RPar`
}

func Parse(input string) (FormulaNode, error) {
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

	// TODO: I think this should not happen on every single parse.
	parser, err := participle.Build[PFormula](
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

func (formula *PFormula) toAst() FormulaNode {
	if formula.Argument != nil {
		return formula.Argument.toAst()
	}
	return formula.ValueLiteral.toAst()
}

func (argument *PArgument) toAst() FormulaNode {
	if argument.FunctionCall != nil {
		return argument.FunctionCall.toAst()
	}
	return argument.ValueLiteral.toAst()
}

func (literal *PBareLiteral) toAst() FormulaNode {
	return &LiteralNode{Value: *literal.IntLiteral}
}

func (literal *PArgLiteral) toAst() FormulaNode {
	if literal.StringLiteral != nil {
		return &LiteralNode{Value: *literal.StringLiteral}
	}
	return &LiteralNode{Value: *literal.IntLiteral}
}

func (call *PFunctionCall) toAst() FormulaNode {
	newArgs := make([]FormulaNode, len(call.Args))
	for i, arg := range call.Args {
		newArgs[i] = arg.toAst()
	}
	return &FunctionNode{Name: strings.ToUpper(*call.Name), Args: newArgs}
}

func astToString(node FormulaNode) string {
	switch node := node.(type) {
	// Cases: literal node string, literal node int, function node, reference node
	case *LiteralNode:
		// Check if it's an int
		if _, ok := node.Value.(int); ok {
			return fmt.Sprintf("%d", node.Value)
		}
		return fmt.Sprintf(`"%s"`, node.Value)
	case *FunctionNode:
		return fmt.Sprintf("%s(%s)", node.Name, strings.Join(Map(node.Args, astToString), ", "))
	}
	return ""
}
