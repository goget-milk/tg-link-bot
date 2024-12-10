package telegram

import (
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

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// sart: /start: hi + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case SartCmd:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessages(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExist {
		// return p.tg.SendMessages(chatID, msgAllreadyExists)
		return sendMsg(msgAllreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessages(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessages(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessages(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessages(chatID, msgHello)
}

// Замыкание
func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessages(chatID, msg)
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
