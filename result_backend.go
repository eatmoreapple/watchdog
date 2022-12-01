package watchdog

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
)

type ResultBackend interface {
	Preparer
	// State returns the state of the given task id.
	State(ctx context.Context, taskId string) (State, error)
	// Result returns the result of the given task id.
	Result(ctx context.Context, taskId string) ([]byte, error)
	// Error returns the error of the given task id.
	Error(ctx context.Context, taskId string) error
	// SetError sets the error of the given task id.
	SetError(ctx context.Context, taskId string, err error) error
	// MarkState marks the state of the given task id.
	MarkState(ctx context.Context, taskId string, state State) error
	// Forget forgets the given task id.
	Forget(ctx context.Context, taskId string) error
	// SetResult the result of the given task id.
	SetResult(ctx context.Context, taskId string, result []byte) error
}

type RedisResultBackend struct {
	redisClient *redis.Client
}

func (r RedisResultBackend) Prepare(engine *Engine) {

}

func (r RedisResultBackend) State(ctx context.Context, taskId string) (State, error) {
	key := "watchdog:task:" + taskId + ":state"
	state, err := r.redisClient.Get(ctx, key).Int64()
	if err == redis.Nil {
		return StatePending, nil
	}
	if err != nil {
		return 0, err
	}
	return State(state), nil
}

func (r RedisResultBackend) Result(ctx context.Context, taskId string) ([]byte, error) {
	key := "watchdog:task:" + taskId + ":result"
	return r.redisClient.Get(ctx, key).Bytes()
}

func (r RedisResultBackend) Error(ctx context.Context, taskId string) error {
	key := "watchdog:task:" + taskId + ":error"
	msg, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New(msg)
}

func (r RedisResultBackend) SetError(ctx context.Context, taskId string, err error) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisResultBackend) MarkState(ctx context.Context, taskId string, state State) error {
	key := "watchdog:task:" + taskId + ":state"
	return r.redisClient.Set(ctx, key, int64(state), 0).Err()
}

func (r RedisResultBackend) Forget(ctx context.Context, taskId string) error {
	//TODO implement me
	panic("implement me")
}

func (r RedisResultBackend) SetResult(ctx context.Context, taskId string, result []byte) error {
	key := "watchdog:task:" + taskId + ":result"
	return r.redisClient.Set(ctx, key, result, 0).Err()
}

func NewRedisResultBackend(client *redis.Client) RedisResultBackend {
	return RedisResultBackend{
		redisClient: client,
	}
}
