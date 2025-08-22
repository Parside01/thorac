package main

import (
	"go.uber.org/zap"
	"log"
)

func main() {
	config := NewConfig()
	logger := NewLogger(config)
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal("failed to sync logger", err)
		}
	}()

	fsm := &Fsm{}
	n, err := NewNode(config, logger, fsm)
	if err != nil {
		log.Fatal(err)
	}

	server := NewHttpServer(config, logger, n)

	logger.Fatal("", zap.Error(server.Start()))
}
