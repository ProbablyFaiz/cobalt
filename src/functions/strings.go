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
