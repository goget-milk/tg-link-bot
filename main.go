package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/goget-milk/telegram-bot/clients/telegram"
	eventconsumer "github.com/goget-milk/telegram-bot/consumer/event-consumer"
	"github.com/goget-milk/telegram-bot/events/telegram"
	"github.com/goget-milk/telegram-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	storagePath       = "files_storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// s := files.New(storagePath),
	s, err := sqlite.New(sqliteStoragePath)

	if err != nil {
		log.Fatalf("can't connect to storage: %w", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage : %w", err)
	}

	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("servise is stopped", err)
	}

}

func mustToken() string {

	token := flag.String("tg-bot-token", "", "token for access to telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token

}
