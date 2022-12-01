package watchdog

type State int

const (
	// StatePending is the default state of a broker.
	StatePending State = iota
	// StateRunning is the state of a broker that is running.
	StateRunning
	// StateCancelled is the state of a broker that is stopped.
	StateCancelled
	// StateError is the state of a broker that is in error.
	StateError
	// StateDone is the state of a broker that is done.
	StateDone
)

func (s State) String() string {
	switch s {
	case StatePending:
		return "pending"
	case StateRunning:
		return "running"
	case StateCancelled:
		return "cancelled"
	case StateError:
		return "error"
	case StateDone:
		return "done"
	default:
		return "unknown"
	}
}
