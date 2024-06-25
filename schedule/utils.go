package main

import (
	"strings"
	"unicode/utf8"
)

var romanNumbers = []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII"}
var weekdays = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

func slen(str string) int {
	return utf8.RuneCountInString(str)
}

func stringStartsWithAnyOf(of []string, str string) (success bool, of_index int) {
	for i := len(of) - 1; i >= 0; i-- {
		if strings.HasPrefix(str, of[i]) {
			return true, i
		}
	}
	return false, -1
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
