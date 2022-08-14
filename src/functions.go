package src

import "fmt"

func ExecuteFn(fnName string, args []interface{}) (interface{}, error) {
	switch fnName {
	case "concat":
		return Concat(args)
	case "add":
		return Add(args)
	default:
		return nil, fmt.Errorf("execute: unknown function %s", fnName)
	}
}

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
