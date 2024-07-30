package utils

import (
	"errors"
	"strings"
)

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

func MustTrue(val bool, errMessage string) {
	if !val {
		panic(errors.New(errMessage))
	}
}
