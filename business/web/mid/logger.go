package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/neogan74/svc3/app/fondation/web"
	"go.uber.org/zap"
)

//

// Logger
func Logger(log *zap.SugaredLogger) web.Midleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// traceID := "000000000000000000"
			// statusCode := http.StatusOK
			// now := time.Now()
			v, err := web.GetValues(ctx)
			if err != nil {
				return err
			}

			// Logging here
			log.Infow("request started", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			err = handler(ctx, w, r)
			time.Sleep(time.Second)
			if err != nil {
				return err
			}

			// Logging here
			log.Infow("request completed", "traceid", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statusCode", v.StatusCode, "since", time.Since(v.Now))

			return nil
		}

		return h
	}

	return m
}
