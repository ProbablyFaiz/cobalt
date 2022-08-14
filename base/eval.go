package base

func (cell Cell) GetValue() interface{} {
	if cell.Dirty {
		// TODO: parse cell.RawContent
		cell.Value = cell.Formula.Eval(&EvalContext{Cell: &cell})
		cell.Dirty = false
	}
	return cell.Value
}

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
	return sheetRef.Cells[rn.Row][rn.Col].GetValue()
}

func (fn FunctionNode) Eval(ctx *EvalContext) interface{} {
	// TODO
	return 0
}
