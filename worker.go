package watchdog

type Worker interface {
	// Run runs the worker.
	Run(payload Payload) ([]byte, error)
}
