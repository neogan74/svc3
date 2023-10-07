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

			traceID := "000000000000000000"
			statusCode := http.StatusOK
			now := time.Now()
			// Logging here
			log.Infow("request started", "traceid", traceID, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)
			time.Sleep(time.Second)
			if err != nil {
				return err
			}

			// Logging here
			log.Infow("request completed", "traceid", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statusCode", statusCode, "since", time.Since(now))

			return nil
		}

		return h
	}

	return m
}
