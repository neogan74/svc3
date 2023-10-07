package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/dimfeld/httptreemux/v5"
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

func ApiMux(cfg ApiMuxConfig) http.Handler {
	mux := httptreemux.NewContextMux()

	thg := testgrp.Handlers{
		Log: cfg.Log,
	}

	mux.Handle(http.MethodGet, "/test", thg.Test)

	return mux
}
