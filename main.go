package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ardanlabs/conf"
	"go.uber.org/automaxprocs/maxprocs"
)

var build = "develop"

func main() {
	if _, err := maxprocs.Set(); err != nil {
		fmt.Printf("maxprocs: %w\n", err)
		os.Exit(1)
	}
	//automaxprocs.New()
	_ = conf.Field{}

	g := runtime.GOMAXPROCS(1)
	log.Printf("starting server LEO build [%s] CPU[%d]  ", build, g)
	defer log.Println("Service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Stopping service")
}
