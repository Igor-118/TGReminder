package main

import (
	// "log"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "github.com/joho/godotenv"
	"tg-reminder/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	bigJopa, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//truba := &app.App{}
	oshibka := bigJopa.Run(ctx)
	if err != nil {
		log.Println("отака фигня собачка", oshibka)
	}
	defer bigJopa.Close()
}
