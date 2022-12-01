package watchdog

import (
	"context"
	"fmt"
)

type Engine struct {
	marshaller    Marshaller
	resultBackend ResultBackend
	broker        Broker
	idGenerator   TaskIDGenerator
	receivers     map[string]Worker
	errHandler    func(error)
}

func (e *Engine) ResultBackend() ResultBackend {
	return e.resultBackend
}

func (e *Engine) Broker() Broker {
	return e.broker
}

func NewEngine(broker Broker, resultBackend ResultBackend) *Engine {
	engine := &Engine{
		marshaller:    &JSONMarshaller{},
		resultBackend: resultBackend,
		broker:        broker,
		idGenerator:   &uuidTaskIDGenerator{},
	}
	broker.Prepare(engine)
	resultBackend.Prepare(engine)
	return engine
}

func (e *Engine) Publish(ctx context.Context, key string, data []byte) (*Future, error) {
	taskID := e.idGenerator.Generate()
	payload := Payload{
		TaskID:     taskID,
		RoutingKey: key,
		Data:       data,
	}
	if err := e.broker.Publish(ctx, payload); err != nil {
		return nil, err
	}
	return &Future{
		taskId: taskID,
		engine: e,
	}, nil
}

func (e *Engine) Register(key string, worker Worker) {
	if e.receivers == nil {
		e.receivers = make(map[string]Worker)
	}
	e.receivers[key] = worker
}

func (e *Engine) ListenAndServe() error {
	ctx := context.Background()
	ch, err := e.broker.Subscribe(ctx)
	if err != nil {
		return err
	}
	for payload := range ch {
		state, err := e.resultBackend.State(ctx, payload.TaskID)
		if err != nil {
			return err
		}
		if state != StatePending {
			continue
		}
		if err := e.resultBackend.MarkState(ctx, payload.TaskID, StateRunning); err != nil {
			return err
		}
		go func(payload Payload) {
			workers := e.receivers[payload.RoutingKey]
			result, err := workers.Run(payload)
			if err != nil {
				_ = e.resultBackend.SetError(ctx, payload.TaskID, err)
				return
			}
			data, err := e.marshaller.Marshal(result)
			if err != nil {
				_ = e.resultBackend.SetError(ctx, payload.TaskID, err)
				return
			}
			_ = e.resultBackend.SetResult(ctx, payload.TaskID, data)

			if err = e.resultBackend.MarkState(ctx, payload.TaskID, StateDone); err != nil {
				fmt.Println(err)
			}
		}(payload)
	}
	<-ctx.Done()
	return nil
}
