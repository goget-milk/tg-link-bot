package eventconsumer

import (
	"log"
	"time"

	"github.com/goget-milk/telegram-bot/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}

	}
}

/*
1. Потеря событий:

  - ретраи,

  - возвращение в хранилище,

  - фоллбэк,

  - подтверждение для фетчера,

  - фетчер не будет делать внутри себя сдвиг пока не увидит, что мы коректно обработали текущую пачку,

  - или фетчер не будет делать сдвиги, а мы будем передавать оффсет самостоятельно,

2. Обработка всей пачки:

  - останавливаться после первой ошибки,

  - счётчик ошибок, останавливаться при достижении определённого кол-ва,

3. Паралельная обработка:

  - sync.WaitGroup{},
*/
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
