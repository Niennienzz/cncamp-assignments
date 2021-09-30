package main

import (
	"cncamp_a01/api"
	"os"
	"os/signal"
)

func main() {
	srv := api.New()
	go srv.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	srv.Shutdown()
}
