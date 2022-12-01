package watchdog

import "sync"

type Scheduler struct {
	maxWorkers     int
	mu             sync.Mutex
	currentWorkers int
	ch             chan Payload
}
