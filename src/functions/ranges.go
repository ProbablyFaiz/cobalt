package functions

import "fmt"

func Sum(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("sum: expected at least 1 argument, got %d", len(args))
	}

	var result interface{} = 0
	for i, arg := range args {
		switch arg.(type) {
		case [][]interface{}:
			for _, row := range arg.([][]interface{}) {
				for _, cell := range row {
					switch cell.(type) {
					case int:
						switch result.(type) {
						case int:
							result = result.(int) + cell.(int)
						case float64:
							result = result.(float64) + float64(cell.(int))
						}
					case float64:
						switch result.(type) {
						case int:
							result = float64(result.(int)) + cell.(float64)
						case float64:
							result = result.(float64) + cell.(float64)
						}
					}
				}
			}
		case int:
			switch result.(type) {
			case int:
				result = result.(int) + arg.(int)
			case float64:
				result = result.(float64) + float64(arg.(int))
			}
		case float64:
			switch result.(type) {
			case int:
				result = float64(result.(int)) + arg.(float64)
			case float64:
				result = result.(float64) + arg.(float64)
			}
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
					// If cell is not a float64 or int, skip it.
					_, ok := cell.(float64)
					if !ok {
						_, ok = cell.(int)
						if !ok {
							continue
						}
					}
					result++
				}
			}
		case int:
			result++
		case float64:
			result++
		default:
			return 0, fmt.Errorf("count: argument %d is a %T, not an int or range", i, arg)
		}
	}
	return result, nil
}

func Average(args []interface{}) (interface{}, error) {
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
