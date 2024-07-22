package utils

import (
	"strings"
	"unicode/utf8"
)

var RomanNumbers = []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII"}
var Weekdays = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

func StringLen(str string) int {
	return utf8.RuneCountInString(str)
}

func StringStartsWithToken(tokens []string, str string) (tokenIndex int) {
	for i := len(tokens) - 1; i >= 0; i-- {
		if strings.HasPrefix(str, tokens[i]) {
			return i
		}
	}
	return -1
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func Must(err error, context ...string) {
	if err != nil {
		if context == nil {
			panic(err.Error())
		} else {
			c := " | " + strings.Join(context, " ; ")
			panic(err.Error() + c)
		}
	}
}
