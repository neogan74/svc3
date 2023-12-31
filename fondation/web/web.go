package web

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Midleware
}

func NewApp(shutdown chan os.Signal, mw ...Midleware) *App {

	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Midleware) {

	//First wrap handler specific middleware around this handler
	handler = wrapMiddleware(mw, handler)

	//Add the application's general middleware to the handler chain
	handler = wrapMiddleware(a.mw, handler)

	// The function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request) {

		// Inject code  here

		// PRE CODE PROCESSING HERE
		// Logging  Started
		// Call the wrapped handler functions
		if err := handler(r.Context(), w, r); err != nil {
			//Pull the context from request and
			// use it as a separate parameter.
			ctx := r.Context()

			v := Values{
				// TraceID: span.SpanContext().TraceID().String(),
				TraceID: uuid.New().String(),
				Now:     time.Now(),
			}

			ctx = context.WithValue(ctx, key, &v)

			// Call the wrapped handler functions
			if err := handler(ctx, w, r); err != nil {
				a.SignalShutdown()
				return
			}
			// Logging error - handle it
			// error handling
			return
		}

		// Inject code here
		// Logging ended here
		// POST CODE PROCESSING

	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}

	a.ContextMux.Handle(method, finalPath, h)

}
