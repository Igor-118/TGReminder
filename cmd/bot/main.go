package main

import (
	"context"
	"log"
	"log/slog"
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
		log.Printf("Ошибка при запуске бота %v", err)
		return
	}
	defer bigJopa.Close()
	err = bigJopa.Run(ctx)
	if err != nil {
		slog.Error("отака фигня собачка", "error", err)
	}
}
