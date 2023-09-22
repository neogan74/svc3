package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf"
)

var build = "develop"

func main() {
	_ = conf.Field{}
	log.Println("starting server .... ", build)
	defer log.Println("Service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("Stopping service")
}
