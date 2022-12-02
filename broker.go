package watchdog

import (
	"context"
	"errors"
	"time"

	"github.com/eatmoreapple/redis"
)

type Broker interface {
	Publish(ctx context.Context, payload Payload) error
	Subscribe(ctx context.Context) (<-chan Payload, error)
}

func NewRedisBroker(client *redis.Client, queue string) *RedisBroker {
	return &RedisBroker{
		client: client,
		queue:  queue,
	}
}

type RedisBroker struct {
	client *redis.Client
	Engine *Engine
	queue  string
}

func (r *RedisBroker) Prepare(engine *Engine) {
	r.Engine = engine
}

func (r *RedisBroker) Publish(ctx context.Context, payload Payload) error {
	return r.client.LPush(ctx, r.queue, payload).Err()
}

func (r *RedisBroker) Subscribe(ctx context.Context) (<-chan Payload, error) {
	ch := make(chan Payload)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, data, err := r.client.BRPop(ctx, []string{r.queue}, 5*time.Second).Result()
				if errors.Is(err, redis.ErrNil) {
					continue
				}
				if err != nil {
					return
				}
				var payload Payload
				if err = r.Engine.marshaller.Unmarshal([]byte(data), &payload); err != nil {
					return
				}
				ch <- payload
			}
		}
	}()
	return ch, nil
}
