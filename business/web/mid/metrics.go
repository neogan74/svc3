package mid

import (
	"context"
	"net/http"

	"github.com/neogan74/svc3/app/fondation/web"
	"github.com/neogan74/svc3/business/sys/metrics"
)

func Metrics() web.Midleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Add the metrics into the context for metric gathering.
			ctx = metrics.Set(ctx)

			// Call next handler
			err := handler(ctx, w, r)

			metrics.AddRequests(ctx)
			metrics.AddGoroutines(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}
			return err
		}
		return h
	}
	return m
}
