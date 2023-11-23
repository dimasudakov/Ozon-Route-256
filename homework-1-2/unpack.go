package main

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result string
	var escape bool

	for i, r := range str {
		if escape {
			if !unicode.IsDigit(r) && r != '\\' && r != '"' {
				return "", ErrInvalidString
			}
			result += string(r)
			escape = false
		} else if unicode.IsDigit(r) {
			if len(result) == 0 {
				return "", ErrInvalidString
			}
			if i-1 >= 0 && unicode.IsDigit(rune(str[i-1])) && (i-2 < 0 || str[i-2] != '\\') {
				return "", ErrInvalidString
			}
			cnt := int(r - '0')
			if cnt == 0 {
				result = result[:len(result)-1]
			} else {
				for j := 0; j < cnt-1; j++ {
					result += string(str[i-1])
				}
			}
		} else if r == '\\' {
			escape = true
		} else {
			result += string(r)
		}
	}
	return result, nil

}
