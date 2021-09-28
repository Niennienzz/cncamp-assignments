package main

import (
	"cncamp_a01/api"
	"cncamp_a01/config"
)

func main() {
	cfg := config.Parse()
	srv := api.New(cfg)
	srv.Run()
}
