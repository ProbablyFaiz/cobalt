package src

func (cell *Cell) GetValue() (interface{}, error) {
	ss := cell.Sheet.Spreadsheet
	if ss.DirtySet.Contains(cell.Uuid) {
		res, err := (*cell.Formula).Eval(&EvalContext{Cell: cell})
		if err != nil {
			return nil, err
		}
		cell.Value = res
		ss.DirtySet.Remove(cell.Uuid)
	}
	return cell.Value, nil
}

func (ln *LiteralNode) Eval(ctx *EvalContext) (interface{}, error) {
	return ln.Value, nil
}

func (rn *ReferenceNode) Eval(ctx *EvalContext) (interface{}, error) {
	return ctx.Cell.Sheet.Spreadsheet.CellMap[rn.ResolvedUuid].GetValue()
}

func (fn *FunctionNode) Eval(ctx *EvalContext) (interface{}, error) {
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

func (_ *NilNode) Eval(ctx *EvalContext) (interface{}, error) {
	return nil, nil
}
