package mid

import (
	"context"
	"net/http"

	"github.com/neogan74/svc3/app/fondation/web"
	"github.com/neogan74/svc3/business/sys/validate"
	"go.uber.org/zap"
)

// Error handles comming out of the call chain. it detects normal application errors
// which are used to repsond to the client in a uniform way. Unexpected errro (status >= 500) are logged
func Errors(log *zap.SugaredLogger) web.Midleware {
	// This is a actual middleware fuction to the executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// if the context is missing this value, request the service to be shutdown gracefully
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}
			// run the next handler and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {

				// Log the error.
				log.Errorw("ERROR", "traceid", v.TraceID, "ERROR", err)

				//Build out the error respose.
				var er validate.ErrorResponse
				var status int
				switch act := validate.Cause(err).(type) {
				case validate.FieldErrors:
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: act.Error(),
					}
					status = http.StatusBadRequest
				case *validate.RequestError:
					er = validate.ErrorResponse{
						Error: act.Error(),
					}
					status = act.Status
				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}
				// Respond with the error back to the client
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutown err we need to return it
				// back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}

			}

			return nil
		}
		return h
	}
	return m
}
