package utils

import (
	"strings"
	"unicode/utf8"
)

func StringLen(str string) int {
	return utf8.RuneCountInString(str)
}

func StringStartsWithToken(tokens []string, str string) (tokenIndex int, found bool) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if strings.HasPrefix(str, tokens[i]) {
			return i, true
		}
	}
	return -1, false
}
