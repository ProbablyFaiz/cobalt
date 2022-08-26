package src

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type PFormula struct {
	Argument     *PArgumentWithInfix `Eq @@`
	ValueLiteral *PBareLiteral       `| @@`
}

type PArgumentWithInfix struct {
	Left       *PArgument   `@@`
	InfixRight *PInfixRight `@@?`
}

type PInfixRight struct {
	Op    *string             `@InfixOp`
	Right *PArgumentWithInfix `@@`
}

type PArgument struct {
	ValueLiteral    *PArgLiteral        `@@`
	FunctionCall    *PFunctionCall      `| @@`
	Reference       *PReference         `| @@`
	WrappedArgument *PArgumentWithInfix `| LPar @@ RPar`
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
	Name *string               `@Ident`
	Args []*PArgumentWithInfix `LPar @@? ( ArgSep @@ )* RPar`
}

type PReference struct {
	A1    *string `@A1Ref`
	A1End *string `(RangeSep @A1Ref)?`
}

var langLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Eq", `=`},
	{"ArgSep", `,`},
	{"LPar", `\(`},
	{"RPar", `\)`},
	{"RangeSep", `:`},
	{"A1Ref", `[A-Z]+[0-9]+`},
	{"Ident", `[a-zA-Z_]\w*`},
	{"Int", `[-+]?\d+`},
	{"InfixOp", `[-+*/]`},
	{"String", `"(\\"|[^"])*"`},
	//{"BareString", `^[^=].*`},
	{"Whitespace", `\s+`},
})

var parser, parseBuildError = participle.Build[PFormula](
	participle.Lexer(langLexer),
	participle.CaseInsensitive("Ident"),
	participle.Elide("Whitespace"),
	participle.Unquote("String"))

func Parse(input string) (FormulaNode, error) {
	if parseBuildError != nil {
		panic(parseBuildError)
	}

	if input == "" {
		return &NilNode{}, nil
	}

	//fmt.Printf("%#v\n", langLexer.Symbols())
	//// Create a map where keys are the symbols and values are the names of the symbols
	//symbolNames := make(map[lexer.TokenType]string)
	//for name, symbol := range langLexer.Symbols() {
	//	symbolNames[symbol] = name
	//}
	//
	//tokens, err := parser.Lex("", strings.NewReader(input))
	//if err != nil {
	//	panic(err)
	//}
	//for _, token := range tokens {
	//	fmt.Printf("%#v\n", symbolNames[token.Type])
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
	if formula.ValueLiteral != nil {
		return formula.ValueLiteral.toAst()
	}
	panic("Impossible state in PFormula.toAst()")
}

func (argument *PArgument) toAst() FormulaNode {
	if argument.FunctionCall != nil {
		return argument.FunctionCall.toAst()
	}
	if argument.Reference != nil {
		return argument.Reference.toAst()
	}
	if argument.ValueLiteral != nil {
		return argument.ValueLiteral.toAst()
	}
	if argument.WrappedArgument != nil {
		return argument.WrappedArgument.toAst()
	}
	panic("Impossible state in PArgument.toAst()")
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

func (call *PArgumentWithInfix) toAst() FormulaNode {
	if call.InfixRight == nil {
		return call.Left.toAst()
	}
	return &FunctionNode{Name: *call.InfixRight.Op, Args: []FormulaNode{call.Left.toAst(), call.InfixRight.Right.toAst()}}
}

func (reference *PReference) toAst() FormulaNode {
	row, col, err := parseA1Notation(*reference.A1)
	if err != nil {
		// This should never happen, as the parser should have already validated the A1 notation
		panic(err)
	}
	startNode := &ReferenceNode{Row: row, Col: col}
	if reference.A1End == nil {
		return startNode
	}
	endRow, endCol, err := parseA1Notation(*reference.A1End)
	if err != nil {
		// See above comment.
		panic(err)
	}
	endNode := &ReferenceNode{Row: endRow, Col: endCol}
	return &RangeNode{From: startNode, To: endNode}
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
