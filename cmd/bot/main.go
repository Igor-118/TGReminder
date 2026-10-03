package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tg-reminder/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	bigJopa, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer bigJopa.Close()
	err = bigJopa.Run(ctx)
	if err != nil {
		log.Println("отака фигня собачка", err)
	}

}
