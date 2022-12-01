package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
	"watchdog"
)

type Result struct {
	Name string `json:"name"`
}

type worker struct{}

func (w *worker) Run(payload watchdog.Payload) (any, error) {
	time.Sleep(time.Second)
	return Result{Name: payload.TaskID}, nil
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	broker := watchdog.NewRedisBroker(client, "queue")

	resultBackend := watchdog.NewRedisResultBackend(client)

	engine := watchdog.NewEngine(broker, resultBackend)

	go func() {
		for {
			result, err := engine.Publish(context.Background(), "key", []byte("hello world"))
			if err != nil {
				fmt.Println(err)
			}
			if err := result.Cancel(context.Background()); err != nil {
				fmt.Println(err)
				return
			}
			var loop = true
			for loop {
				state, err := result.State(context.Background())
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(state)
				switch state {
				case watchdog.StateDone:
					var data Result
					fmt.Println(result.Get(context.Background(), &data))
					fmt.Println(data)
					loop = false
					break
				case watchdog.StateCancelled:
					loop = false
					break
				default:
					time.Sleep(time.Second * 1)
				}
			}

			time.Sleep(time.Second)
		}
	}()
	engine.Register("key", &worker{})

	fmt.Println(engine.ListenAndServe())
}
