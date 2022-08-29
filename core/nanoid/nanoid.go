package nanoid

import (
	"github.com/jaevor/go-nanoid"
)

var (
	generator func() string = nil
)

func Next(length int) string {
	gen, _ := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", length) // ignore err since it only checks length

	if generator == nil {
		generator = gen
	}

	return generator()
}
