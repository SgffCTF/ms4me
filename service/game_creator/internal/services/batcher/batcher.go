package batcher

import (
	"context"
	"log/slog"
	gameclient "ms4me/game_creator/pkg/http/game"
	"sync"
	"time"
)

const (
	processInterval = 200 * time.Millisecond
	maxBatchSize    = 50
)

type Batcher struct {
	log        *slog.Logger
	gameClient *gameclient.GameClient
	batchPool  *sync.Pool
	newEvents  chan gameclient.Event
	shutdown   chan struct{}
	wg         sync.WaitGroup
}

func New(log *slog.Logger, gameClient *gameclient.GameClient) *Batcher {
	return &Batcher{
		log:        log,
		gameClient: gameClient,
		newEvents:  make(chan gameclient.Event),
		shutdown:   make(chan struct{}),
		batchPool: &sync.Pool{
			New: func() any {
				sl := make([]gameclient.Event, 0, maxBatchSize)
				return &sl
			},
		},
	}
}

func (b *Batcher) processBatch(batch []gameclient.Event) {
	const op = "services.batch.processBatch"
	log := b.log.With(slog.String("op", op), slog.Int("batch_size", len(batch)))

	b.wg.Add(1)

	go func() {
		defer b.wg.Done()
		defer func() {
			batch = batch[:0]
			b.batchPool.Put(&batch)
		}()

		log.Debug("processing batch")
		err := b.gameClient.LoadEvents(batch)
		if err != nil {
			b.log.Debug("error processing batch", slog.String("op", op))
			return
		}
		log.Debug("batch processed successfully")
	}()
}

func (b *Batcher) Start() {
	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		ticker := time.NewTicker(processInterval)
		defer ticker.Stop()

		batch := *b.batchPool.Get().(*[]gameclient.Event)

		for {
			select {
			case <-ticker.C:
				if len(batch) != 0 {
					b.log.Debug("flushing batch", slog.Int("batch_len", len(batch)))
					b.processBatch(batch)
					batch = *b.batchPool.Get().(*[]gameclient.Event)
				}
			case event := <-b.newEvents:
				batch = append(batch, event)

				b.log.Debug("new item added to batch", slog.Int("batch_len", len(batch)))

				if len(batch) == cap(batch) {
					b.log.Debug("flushing batch", slog.Int("batch_len", len(batch)))
					b.processBatch(batch)
					ticker.Reset(processInterval)
					batch = *b.batchPool.Get().(*[]gameclient.Event)
				}
			case <-b.shutdown:
				if len(batch) != 0 {
					b.log.Debug("flushing remaining batch", slog.Int("batch_len", len(batch)))
					b.processBatch(batch)
				}
				return
			}
		}
	}()
}

func (b *Batcher) AddEvents(ctx context.Context, event gameclient.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.newEvents <- event:
		return nil
	}
}

func (b *Batcher) Shutdown() {
	close(b.shutdown)
	b.wg.Wait()
}
