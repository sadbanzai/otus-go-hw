package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrCannotConvertToInt = errors.New("cannot convert to int")
	ErrInvalidString      = errors.New("invalid string")
)

func Unpack(s string) (string, error) {
	var result strings.Builder
	var lastRune rune
	for i, currentRune := range s {
		if (i == 0 || unicode.IsDigit(lastRune)) && unicode.IsDigit(currentRune) {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(currentRune) {
			val, err := strconv.Atoi(string(currentRune))
			if err != nil {
				return "", ErrCannotConvertToInt
			}
			if val == 0 {
				str := result.String()
				str = strings.TrimSuffix(str, string(lastRune))
				result.Reset()
				result.WriteString(str)
			} else {
				result.WriteString(strings.Repeat(string(lastRune), val-1))
			}
		} else {
			result.WriteString(string(currentRune))
		}
		lastRune = currentRune
	}
	return result.String(), nil
}
