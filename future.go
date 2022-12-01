package watchdog

import (
	"context"
	"fmt"
)

type Future struct {
	taskId string
	engine *Engine
}

func (f *Future) TaskID() string {
	return f.taskId
}

func (f *Future) Get(ctx context.Context, v any) error {
	data, err := f.engine.ResultBackend().Result(ctx, f.TaskID())
	if err != nil {
		return err
	}
	return f.engine.marshaller.Unmarshal(data, v)
}

func (f *Future) Err(ctx context.Context) error {
	return f.engine.ResultBackend().Error(ctx, f.TaskID())
}

func (f *Future) Forget(ctx context.Context) error {
	return f.engine.ResultBackend().Forget(ctx, f.TaskID())
}

func (f *Future) Cancel(ctx context.Context) error {
	state, err := f.State(ctx)
	if err != nil {
		return err
	}
	switch state {
	case StatePending:
		return f.engine.ResultBackend().MarkState(ctx, f.TaskID(), StateCancelled)
	default:
		return fmt.Errorf("cannot cancel task %s in state %s", f.TaskID(), state)
	}
}

func (f *Future) State(ctx context.Context) (State, error) {
	return f.engine.ResultBackend().State(ctx, f.TaskID())
}
