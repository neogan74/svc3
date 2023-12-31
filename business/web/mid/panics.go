package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/neogan74/svc3/business/sys/metrics"
	"github.com/neogan74/svc3/fondation/web"
)

func Panics() web.Midleware {

	// This is the actual middleware function to be executed
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {

					trace := debug.Stack()

					err = fmt.Errorf("PANIC [%v] TRACE: [%s]", rec, trace)

					metrics.AddPanics(ctx)
				}
			}()

			// Call the next handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}
		return h
	}
	return m
}
