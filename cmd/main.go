package main

import (
	"aggromeration/internal"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	application := internal.New()

	go func() {
		if err := application.Run(); err != nil {
			log.Fatal(ctx, fmt.Sprintf("launch application: %v", err))
		}
	}()
	<-ctx.Done()

	log.Println("Gracefully shutting down process...")
	if err := application.Stop(); err != nil {
		log.Fatal(ctx, fmt.Sprintf("stop application: %v", err))
	}
}
