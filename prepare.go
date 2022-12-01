package watchdog

type Preparer interface {
	Prepare(engine *Engine)
}
