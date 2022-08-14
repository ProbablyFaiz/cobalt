package base

func (ln LiteralNode) Eval(ctx *EvalContext) interface{} {
	return ln.Value
}

func (rn ReferenceNode) Eval(ctx *EvalContext) interface{} {
	var sheetRef *Sheet
	if rn.Sheet != nil {
		sheetRef = rn.Sheet
	} else {
		sheetRef = ctx.Cell.Sheet
	}
	return sheetRef.Cells[rn.Row][rn.Col].Value
}
