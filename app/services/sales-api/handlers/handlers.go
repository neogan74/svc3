package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/neogan74/svc3/app/fondation/web"
	"github.com/neogan74/svc3/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/neogan74/svc3/app/services/sales-api/handlers/v1/testgrp"
	"github.com/neogan74/svc3/business/web/mid"
	"go.uber.org/zap"
)

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	chg := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}

	mux.HandleFunc("/debug/readiness", chg.Readiness)
	mux.HandleFunc("/debug/liveness", chg.Liveness)

	return mux
}

// ApiMuxConfig type
type ApiMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

func ApiMux(cfg ApiMuxConfig) *web.App {
	// Construct the web.App which holds all routes as well as common Midleware
	app := web.NewApp(
		cfg.Shutdown,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
	)

	// Load the routes for the different versions of the API.
	v1(app, cfg)

	return app
}

// v1 Binds all version 1 routes
func v1(app *web.App, cfg ApiMuxConfig) {
	const version = "v1"

	thg := testgrp.Handlers{
		Log: cfg.Log,
	}

	app.Handle(http.MethodGet, version, "/test", thg.Test)

}
