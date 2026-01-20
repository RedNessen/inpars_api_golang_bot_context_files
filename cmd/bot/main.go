package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RedNessen/inpars-telegram-bot/internal/config"
	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
	"github.com/RedNessen/inpars-telegram-bot/internal/monitor"
	"github.com/RedNessen/inpars-telegram-bot/internal/storage"
	"github.com/RedNessen/inpars-telegram-bot/internal/telegram"
)

func main() {
	log.Println("Starting InPars Telegram Bot...")

	// Загрузка конфигурации
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Configuration loaded successfully")

	// Подключение к базе данных
	store, err := storage.NewPostgresStorage(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()
	log.Println("Database connection established")

	// Создание клиента InPars API
	inparsClient := inpars.NewClient(cfg.InParsToken)
	log.Println("InPars API client initialized")

	// Создание Telegram бота
	bot, err := telegram.NewBot(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}
	log.Println("Telegram bot initialized")

	// Создание монитора с хранилищем
	mon := monitor.NewMonitor(inparsClient, bot, store, cfg)
	log.Println("Monitor initialized")

	// Запуск бота в отдельной горутине
	go func() {
		if err := bot.Start(); err != nil {
			log.Fatalf("Bot stopped with error: %v", err)
		}
	}()

	// Запуск монитора в отдельной горутине
	go func() {
		if err := mon.Start(); err != nil {
			log.Fatalf("Monitor stopped with error: %v", err)
		}
	}()

	log.Println("Bot and monitor are running. Press Ctrl+C to stop.")

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	log.Println(mon.GetStatus())

	// Закрываем соединение с БД
	if err := store.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Goodbye!")
}
