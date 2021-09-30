package main

import (
	"cncamp_a01/api"
	"os"
	"os/signal"
)

func main() {
	srv := api.New()

	go func() {
		srv.Run()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	srv.Shutdown()
}
