package functions

import (
	"fmt"
)

// TODO: It may make sense to implement a basic type-checker so that
//  we don't have to validate the types of arguments in every function.

func Concat(args []interface{}) (string, error) {
	var result string
	for i, arg := range args {
		switch arg.(type) {
		case string:
			result += arg.(string)
		default:
			return "", fmt.Errorf("concat: argument %d is a %T, not a string", i, arg)
		}

	}
	return result, nil
}

func Add(args []interface{}) (int, error) {
	var result int
	for i, arg := range args {
		switch arg.(type) {
		case int:
			result += arg.(int)
		default:
			return 0, fmt.Errorf("add: argument %d is a %T, not an int", i, arg)
		}

	}
	return result, nil
}

func Sub(args []interface{}) (int, error) {
	// Must have at least 1 argument
	if len(args) < 1 {
		return 0, fmt.Errorf("sub: expected at least 1 argument, got %d", len(args))
	}

	var result int
	for i, arg := range args {
		switch arg.(type) {
		case int:
			if i == 0 {
				result = arg.(int)
			} else {
				result -= arg.(int)
			}
		default:
			return 0, fmt.Errorf("sub: argument %d is a %T, not an int", i, arg)
		}

	}
	return result, nil
}

func Mul(args []interface{}) (int, error) {
	var result int = 1
	for i, arg := range args {
		switch arg.(type) {
		case int:
			result *= arg.(int)
		default:
			return 0, fmt.Errorf("mul: argument %d is a %T, not an int", i, arg)
		}

	}
	return result, nil
}

func Div(args []interface{}) (float64, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("div: expected 2 arguments, got %d", len(args))
	}

	var result float64
	for i, arg := range args {
		switch arg.(type) {
		case int:
			if i == 0 {
				result = arg.(float64)
			} else {
				result /= arg.(float64)
			}
		default:
			return 0, fmt.Errorf("div: argument %d is a %T, not an int", i, arg)
		}
	}
	return result, nil
}

func Mod(args []interface{}) (int, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("mod: expected 2 arguments, got %d", len(args))
	}

	var result int
	for i, arg := range args {
		switch arg.(type) {
		case int:
			if i == 0 {
				result = arg.(int)
			} else {
				result %= arg.(int)
			}
		default:
			return 0, fmt.Errorf("mod: argument %d is a %T, not an int", i, arg)
		}
	}
	return result, nil
}

func Pow(args []interface{}) (int, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("pow: expected 2 arguments, got %d", len(args))
	}

	var result int
	for i, arg := range args {
		switch arg.(type) {
		case int:
			result ^= arg.(int)
		default:
			return 0, fmt.Errorf("pow: argument %d is a %T, not an int", i, arg)
		}
	}
	return result, nil
}

func If(args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return 0, fmt.Errorf("if: expected 3 arguments, got %d", len(args))
	}

	var truthy bool
	switch args[0].(type) {
	case bool:
		truthy = args[0].(bool)
	case int:
		truthy = args[0].(int) != 0
	case string:
		truthy = args[0].(string) != ""
	case nil:
		truthy = false
	default:
		return 0, fmt.Errorf("if: argument 0 is a %T, not a bool, int, string, or nil", args[0])
	}

	if truthy {
		return args[1], nil
	} else {
		return args[2], nil
	}
}

func Not(args []interface{}) (bool, error) {
	if len(args) != 1 {
		return false, fmt.Errorf("not: expected 1 argument, got %d", len(args))
	}

	switch args[0].(type) {
	case bool:
		return !args[0].(bool), nil
	case int:
		return args[0].(int) == 0, nil
	case string:
		return args[0].(string) == "", nil
	case nil:
		return true, nil
	default:
		return false, fmt.Errorf("not: argument 0 is a %T, not a bool, int, string, or nil", args[0])
	}
}

func And(args []interface{}) (bool, error) {
	if len(args) < 1 {
		return false, fmt.Errorf("and: expected at least 1 argument, got %d", len(args))
	}

	var result bool = true
	for i, arg := range args {
		switch arg.(type) {
		case bool:
			result = result && arg.(bool)
		case int:
			result = result && arg.(int) != 0
		case string:
			result = result && arg.(string) != ""
		case nil:
			result = false
		default:
			return false, fmt.Errorf("and: argument %d is a %T, not a bool, int, string, or nil", i, arg)
		}
	}
	return result, nil
}

func Or(args []interface{}) (bool, error) {
	if len(args) < 1 {
		return false, fmt.Errorf("or: expected at least 1 argument, got %d", len(args))
	}

	var result bool = false
	for i, arg := range args {
		switch arg.(type) {
		case bool:
			result = result || arg.(bool)
		case int:
			result = result || arg.(int) != 0
		case string:
			result = result || arg.(string) != ""
		case nil:
			result = false
		default:
			return false, fmt.Errorf("or: argument %d is a %T, not a bool, int, string, or nil", i, arg)
		}
	}
	return result, nil
}

func Eq(args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("eq: expected 2 arguments, got %d", len(args))
	}

	return args[0] == args[1], nil
}

func Ne(args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("ne: expected 2 arguments, got %d", len(args))
	}

	return args[0] != args[1], nil
}

func Gt(args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("gt: expected 2 arguments, got %d", len(args))
	}

	switch args[0].(type) {
	case int:
		return args[0].(int) > args[1].(int), nil
	case string:
		return args[0].(string) > args[1].(string), nil
	default:
		return false, fmt.Errorf("gt: argument 0 is a %T, not an int or string", args[0])
	}
}

func Lt(args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("lt: expected 2 arguments, got %d", len(args))
	}

	switch args[0].(type) {
	case int:
		return args[0].(int) < args[1].(int), nil
	case string:
		return args[0].(string) < args[1].(string), nil
	default:
		return false, fmt.Errorf("lt: argument 0 is a %T, not an int or string", args[0])
	}
}

func Gte(args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("gte: expected 2 arguments, got %d", len(args))
	}

	switch args[0].(type) {
	case int:
		return args[0].(int) >= args[1].(int), nil
	case string:
		return args[0].(string) >= args[1].(string), nil
	default:
		return false, fmt.Errorf("gte: argument 0 is a %T, not an int or string", args[0])
	}
}

func Lte(args []interface{}) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("lte: expected 2 arguments, got %d", len(args))
	}

	switch args[0].(type) {
	case int:
		return args[0].(int) <= args[1].(int), nil
	case string:
		return args[0].(string) <= args[1].(string), nil
	default:
		return false, fmt.Errorf("lte: argument 0 is a %T, not an int or string", args[0])
	}
}

func Sum(args []interface{}) (int, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("sum: expected at least 1 argument, got %d", len(args))
	}

	var result int
	for i, arg := range args {
		switch arg.(type) {
		case [][]interface{}:
			for _, row := range arg.([][]interface{}) {
				for _, cell := range row {
					// If the cell is not an int, skip it.
					val, ok := cell.(int)
					if !ok {
						continue
					}
					result += val
				}
			}
		case int:
			result += arg.(int)
		default:
			return 0, fmt.Errorf("sum: argument %d is a %T, not an int or range", i, arg)
		}
	}
	return result, nil
}

func Count(args []interface{}) (int, error) {
	// Counts over all arguments.
	if len(args) < 1 {
		return 0, fmt.Errorf("count: expected at least 1 argument, got %d", len(args))
	}

	var result int
	for i, arg := range args {
		switch arg.(type) {
		case [][]interface{}:
			for _, row := range arg.([][]interface{}) {
				for _, cell := range row {
					_, ok := cell.(int)
					if !ok {
						continue
					}
					result++
				}
			}
		case int:
			result++
		default:
			return 0, fmt.Errorf("count: argument %d is a %T, not an int or range", i, arg)
		}
	}
	return result, nil
}

func Average(args []interface{}) (float64, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("average: expected at least 1 argument, got %d", len(args))
	}

	// Use sum, count, and div to calculate the average.
	sum, err := Sum(args)
	if err != nil {
		return 0, fmt.Errorf("average: %s", err)
	}
	count, err := Count(args)
	if err != nil {
		return 0, fmt.Errorf("average: %s", err)
	}
	res, err := Div([]interface{}{sum, count})
	if err != nil {
		return 0, fmt.Errorf("average: %s", err)
	}
	return res, nil
}
