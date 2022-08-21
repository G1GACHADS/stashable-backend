package nanoid

import (
	"github.com/jaevor/go-nanoid"
)

var (
	generator func() string = nil
)

func Next() string {
	gen, _ := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 21) // ignore err since it only checks length

	if generator == nil {
		generator = gen
	}

	return generator()
}
