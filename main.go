package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/opoccomaxao/shiki-github-graph/pkg/server"
)

func main() {
	_ = godotenv.Load()

	config := server.Config{}

	err := env.Parse(&config)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	server, err := server.New(config)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	ctx, cancelFn := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancelFn()

	go func() {
		<-ctx.Done()

		err = server.Close()
		if err != nil {
			log.Fatalf("%+v", err)
		}
	}()

	err = server.Serve(ctx)
	if err != nil {
		log.Printf("%+v", err)
	}

	<-ctx.Done()
}
