package main

import (
	"cncamp_a01/httpserver/api"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	srv := api.New()
	go srv.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	srv.Shutdown()
}
