package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/tg"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN environment variable is empty")
	}

	dbDSN := os.Getenv("MAIN_DATABASE_DSN")
	if dbDSN == "" {
		log.Fatal("MAIN_DATABASE_DSN environment variable is empty")
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем обработку сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Received shutdown signal, shutting down...")
		cancel() // Отменяем контекст
	}()

	db, err := sqlx.Connect("sqlite", dbDSN)
	if err != nil {
		log.Fatalln(err)
	}

	if err = tg.Run(ctx, token, db); err != nil {
		log.Fatal(err)
	}

	log.Println("Application has shut down gracefully.")
}
