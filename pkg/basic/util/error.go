package util

import "github.com/pkg/errors"

func PanicIfError(err error) {
	if err != nil {
		panic(errors.WithStack(err))
	}
	return
}
