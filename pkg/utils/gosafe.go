package utils

import (
	"context"
	"log/slog"
	"runtime/debug"
)

// Go - safe version of 'go func' which recovers panics
func Go(ctx context.Context, fn func(ctx context.Context)) {
	if fn != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("goroutine panic: %s %s", r, debug.Stack())
				}
			}()

			fn(ctx)
		}()
	}
}
