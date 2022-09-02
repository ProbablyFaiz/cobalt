package core

import (
	"fmt"
	"strconv"
)

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func getA1Notation(row int, col int) string {
	colLetters := ""
	for i := col; i >= 0; i = i/26 - 1 {
		colLetters = string(rune('A'+i%26)) + colLetters
	}
	return fmt.Sprintf("%s%d", colLetters, row+1)
}

func parseA1Notation(a1 string) (int, int, error) {
	i := 0

	var col int
	for i = 0; i < len(a1); i++ {
		if a1[i] >= 'A' && a1[i] <= 'Z' {
			letterValue := int(a1[i]-'A') + 1
			col = col*26 + letterValue
		} else {
			break
		}
	}
	row, err := strconv.Atoi(a1[i:])
	if err != nil {
		return 0, 0, fmt.Errorf("parseA1Notation: %s", err)
	}
	return row - 1, col - 1, nil
}
