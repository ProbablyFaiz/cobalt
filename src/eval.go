package src

import (
	"fmt"
	"pasado/src/functions"
)

func (cell *Cell) GetOrComputeValue() (interface{}, error) {
	ss := cell.Sheet.Spreadsheet
	if ss.DirtySet.Contains(cell.Uuid) {
		res, err := (*cell.Formula).eval(&EvalContext{Cell: cell})
		cell.Value, cell.Error = res, err
		ss.DirtySet.Remove(cell.Uuid)
	}
	return cell.Value, nil
}

func (ln *LiteralNode) eval(_ *EvalContext) (interface{}, error) {
	return ln.Value, nil
}

func (rn *ReferenceNode) eval(ctx *EvalContext) (interface{}, error) {
	return ctx.Cell.Sheet.Spreadsheet.CellMap[rn.ResolvedUuid].GetOrComputeValue()
}

func (rn *RangeNode) eval(ctx *EvalContext) (interface{}, error) {
	// First, evaluate the range's dirty parent cells.
	ss := ctx.Cell.Sheet.Spreadsheet
	for cellId := range ss.RangeDirtyParents[rn.ResolvedUuid].Iter() {
		_, _ = ss.CellMap[cellId].GetOrComputeValue()
		ss.RangeDirtyParents[rn.ResolvedUuid].Remove(cellId)
	}

	// Gets the range of cells in the sheet
	sheet := rn.End.Sheet
	startRow, startCol := rn.Start.Row, rn.Start.Col
	endRow, endCol := rn.End.Row, rn.End.Col

	return sheet.GetRange(startRow, startCol, endRow, endCol)
}

func (fn *FunctionNode) eval(ctx *EvalContext) (interface{}, error) {
	args := make([]interface{}, len(fn.Args))
	for i, arg := range fn.Args {
		res, err := arg.eval(ctx)
		if err != nil {
			return nil, err
		}
		args[i] = res
	}
	return ExecuteFn(fn.Name, args)
}

func (_ *NilNode) eval(ctx *EvalContext) (interface{}, error) {
	return nil, nil
}

func ExecuteFn(fnName string, args []interface{}) (interface{}, error) {
	switch fnName {
	case "CONCAT":
		return functions.Concat(args)
	case "ADD":
		return functions.Add(args)
	case "+":
		return functions.Add(args)
	case "SUB":
		return functions.Sub(args)
	case "-":
		return functions.Sub(args)
	case "MUL":
		return functions.Mul(args)
	case "*":
		return functions.Mul(args)
	case "DIV":
		return functions.Div(args)
	case "/":
		return functions.Div(args)
	case "MOD":
		return functions.Mod(args)
	case "%":
		return functions.Mod(args)
	case "POW":
		return functions.Pow(args)
	case "IF":
		return functions.If(args)
	default:
		return nil, fmt.Errorf("execute: unknown function %s", fnName)
	}
}
