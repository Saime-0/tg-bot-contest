package main

import (
	"context"
	"log/slog"
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
		slog.Error("TOKEN environment variable is empty")
		os.Exit(1)
	}

	dbDSN := os.Getenv("MAIN_DATABASE_DSN")
	if dbDSN == "" {
		slog.Error("MAIN_DATABASE_DSN environment variable is empty")
		os.Exit(1)
	}

	debug := os.Getenv("DEBUG")
	if debug == "true" || debug == "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем обработку сигналов
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		slog.Info("Received shutdown signal, shutting down...")
		cancel() // Отменяем контекст
	}()

	db, err := sqlx.Connect("sqlite", dbDSN)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err = tg.Run(ctx, token, db); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("Application has shut down gracefully")
}
