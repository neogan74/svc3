package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/neogan74/svc3/app/fondation/web"
	"github.com/neogan74/svc3/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/neogan74/svc3/app/services/sales-api/handlers/v1/testgrp"
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
	app := web.NewApp(cfg.Shutdown)

	thg := testgrp.Handlers{
		Log: cfg.Log,
	}

	app.Handle(http.MethodGet, "/v1/test", thg.Test)

	return app
}
