package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"
var service = "SALES-API"

func main() {

	// Constructing logger
	logconfig := zap.NewProductionConfig()
	logconfig.OutputPaths = []string{"stdout"}
	logconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logconfig.DisableStacktrace = true
	logconfig.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := logconfig.Build()
	if err != nil {
		fmt.Println("Error constructing logger", err)
		os.Exit(1)
	}

	defer log.Sync()

	if err := run(log.Sugar()); err != nil {
		log.Sugar().Errorw("startup", "ERROR", err)
		os.Exit(1)
	}

	log.Info("We are starting....")

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
	return nil
}
