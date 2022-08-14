package src

func (cell Cell) GetValue() (interface{}, error) {
	if cell.Dirty {
		formula, err := Parse(cell.RawContent)
		if err != nil {
			return nil, err
		}
		cell.Formula = formula
		res, err := cell.Formula.Eval(&EvalContext{Cell: &cell})
		if err != nil {
			return nil, err
		}
		cell.Value = res
		cell.Dirty = false
	}
	return cell.Value, nil
}

func (ln LiteralNode) Eval(ctx *EvalContext) (interface{}, error) {
	return ln.Value, nil
}

func (rn ReferenceNode) Eval(ctx *EvalContext) (interface{}, error) {
	var sheetRef *Sheet
	if rn.Sheet != nil {
		sheetRef = rn.Sheet
	} else {
		sheetRef = ctx.Cell.Sheet
	}
	return sheetRef.Cells[rn.Row][rn.Col].GetValue()
}

func (fn FunctionNode) Eval(ctx *EvalContext) (interface{}, error) {
	args := make([]interface{}, len(fn.Args))
	for i, arg := range fn.Args {
		res, err := arg.Eval(ctx)
		if err != nil {
			return nil, err
		}
		args[i] = res
	}
	return ExecuteFn(fn.Name, args)
}
