package main

import (
	"github.com/janqii/pusher/server"
	"log"
)

func main() {
	cfg, err := server.NewPusherConfig()
	if err != nil {
		log.Fatalf("parse config error, %v", err)
	}

	if err = server.Startable(cfg); err != nil {
		log.Fatalf("server startable error, %v", err)
	}
}
