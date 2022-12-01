package watchdog

import "github.com/google/uuid"

type TaskIDGenerator interface {
	// Generate generates a new task id.
	Generate() string
}

type uuidTaskIDGenerator struct{}

func (g *uuidTaskIDGenerator) Generate() string {
	return uuid.New().String()
}
