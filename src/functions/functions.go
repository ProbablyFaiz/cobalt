package functions

import "fmt"

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
	var result int
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

func Div(args []interface{}) (int, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("div: expected 2 arguments, got %d", len(args))
	}

	var result int
	for i, arg := range args {
		switch arg.(type) {
		case int:
			result /= arg.(int)
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
			result %= arg.(int)
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
