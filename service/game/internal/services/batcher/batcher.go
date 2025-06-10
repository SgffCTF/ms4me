package batcher

import (
	"context"
	"log/slog"
	"ms4me/game/internal/models"
	"ms4me/game/internal/storage/redis"
	"sync"
	"time"

	"github.com/jacute/prettylogger"
)

const (
	processInterval = 10 * time.Millisecond
	maxBatchSize    = 50
)

type Batcher struct {
	log       *slog.Logger
	batchPool *sync.Pool
	newEvents chan models.Event
	shutdown  chan struct{}
	redis     *redis.Redis
	wg        sync.WaitGroup
}

func New(log *slog.Logger, redis *redis.Redis) *Batcher {
	return &Batcher{
		log:       log,
		newEvents: make(chan models.Event),
		shutdown:  make(chan struct{}),
		batchPool: &sync.Pool{
			New: func() any {
				sl := make([]models.Event, 0, maxBatchSize)
				return &sl
			},
		},
		redis: redis,
	}
}

func (b *Batcher) processBatch(batch []models.Event) {
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
		err := b.redis.PublishEvents(context.Background(), batch)
		if err != nil {
			log.Error("error processing batch", prettylogger.Err(err))
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

		batch := *b.batchPool.Get().(*[]models.Event)

		for {
			select {
			case <-ticker.C:
				if len(batch) != 0 {
					b.log.Debug("flushing batch", slog.Int("batch_len", len(batch)))
					b.processBatch(batch)
					batch = *b.batchPool.Get().(*[]models.Event)
				}
			case event := <-b.newEvents:
				batch = append(batch, event)

				b.log.Debug("new item added to batch", slog.Int("batch_len", len(batch)))

				if len(batch) == cap(batch) {
					b.log.Debug("flushing batch", slog.Int("batch_len", len(batch)))
					b.processBatch(batch)
					ticker.Reset(processInterval)
					batch = *b.batchPool.Get().(*[]models.Event)
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

func (b *Batcher) AddEvents(ctx context.Context, event models.Event) error {
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
