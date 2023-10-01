package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/neogan74/svc3/app/services/sales-api/handlers"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"
var service = "SALES-API"

func main() {

	// Constructing application logger
	log, err := initLogger(service)
	if err != nil {
		fmt.Println("Cannot init logger", err)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}

	// if _, err := maxprocs.Set(); err != nil {
	// 	fmt.Printf("maxprocs: %w\n", err)
	// 	os.Exit(1)
	// }
	// //automaxprocs.New()
	// _ = conf.Field{}
	// g := runtime.GOMAXPROCS(1)
	// log.Printf("starting server LEO build [%s] CPU[%d]  ", build, g)
	// defer log.Println("Service ended")
	// shutdown := make(chan os.Signal, 1)
	// signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	// <-shutdown
	// log.Println("Stopping service")
}

func run(log *zap.SugaredLogger) error {
	// GOMAXPROCS

	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("starting...", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =================
	// configuration
	// -----------------

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s,mask"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright",
		},
	}

	const prefix = "SALES"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}
	// =====================
	// App start
	log.Infow("starting service...", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)
	// =====================

	// ====================
	// start debug service
	log.Infow("startup", "status", "Debug router started", "host", cfg.Web.DebugHost)

	// THe Debug function returns a mux to lister and serve on for all the debug
	// related endpoints. This include the standard library endpoints.

	// constuct the mux for the build calls.
	debugMux := handlers.DebugStandardLibraryMux()

	// Start the service listening for debug requests.
	// Not concerned with suttling this down with load shedding.

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			log.Errorf("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	// <-shutdown

	// Start API service

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      nil,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// Make a chanel to listen for errors comming from the listener. use a
	// buffered channel so the goroutine can exit if we don't colelct error.
	serverErrors := make(chan error, 1)

	// Start the service listening for api reqeusts
	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()
	// Shutdown

	//blickong main and watching for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for compleation
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		//Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("Could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	logconfig := zap.NewProductionConfig()
	logconfig.OutputPaths = []string{"stdout"}
	logconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logconfig.DisableStacktrace = true
	logconfig.InitialFields = map[string]interface{}{
		"service": service,
	}

	llog, err := logconfig.Build()
	if err != nil {
		return nil, err
	}

	return llog.Sugar(), nil
}
