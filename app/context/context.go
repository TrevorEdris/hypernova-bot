package context

import (
	"context"
	"errors"
)

type (
	ContextKey string
)

var (
	// SomeCustomKey is a key that can be used to store a value in a context.Context.
	SomeCustomKey ContextKey = "a_key_to_store_in_ctx"
)

func IsCanceledError(err error) bool {
	return errors.Is(err, context.Canceled)
}
