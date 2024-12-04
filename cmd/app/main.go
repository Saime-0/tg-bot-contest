package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"tgBotContest/tg"
)

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN environment variable is empty")
	}

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

	db, err := sqlx.Connect("sqlite3", "file:main.db?cache=shared")
	if err != nil {
		log.Fatalln(err)
	}

	if err := tg.Run(ctx, token, db); err != nil {
		log.Fatal(err)
	}

	log.Println("Application has shut down gracefully.")
}
