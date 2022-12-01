package watchdog

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"time"
)

type Broker interface {
	Preparer
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
				data, err := r.client.BRPop(ctx, 5*time.Second, r.queue).Result()
				if errors.Is(err, redis.Nil) {
					continue
				}
				if err != nil {
					return
				}
				var payload Payload
				if err = r.Engine.marshaller.Unmarshal([]byte(data[1]), &payload); err != nil {
					return
				}
				ch <- payload
			}
		}
	}()
	return ch, nil
}
