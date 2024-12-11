package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/goget-milk/telegram-bot/clients/telegram"
	"github.com/goget-milk/telegram-bot/lib/e"
	"github.com/goget-milk/telegram-bot/storage"
)

const (
	RndCmd  = "/rnd"
	HelpCmd = "/help"
	SartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// sart: /start: hi + help

	if isAddCmd(text) {
		return p.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case SartCmd:
		return p.sendHelp(ctx, chatID)
	default:
		return p.tg.SendMessages(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(
	ctx context.Context,
	chatID int,
	pageURL string,
	username string,
) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := NewMessageSender(ctx, chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(ctx, page)
	if err != nil {
		return err
	}

	if isExist {
		// return p.tg.SendMessages(chatID, msgAllreadyExists)
		return sendMsg(msgAllreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessages(ctx, chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessages(ctx, chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessages(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessages(ctx, chatID, msgHello)
}

// Замыкание
func NewMessageSender(ctx context.Context, chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessages(ctx, chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// ya.ru - NO
	// https://ya.ru - YES
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
