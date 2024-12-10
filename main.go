package main

import (
	"flag"
	"log"

	tgClient "github.com/goget-milk/telegram-bot/clients/telegram"
	eventconsumer "github.com/goget-milk/telegram-bot/consumer/event-consumer"
	"github.com/goget-milk/telegram-bot/events/telegram"
	"github.com/goget-milk/telegram-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
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
