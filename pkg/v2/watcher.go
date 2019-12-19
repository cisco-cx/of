package v2

// Interface to watch for file/dir changes or syskill signals.
type Watcher interface {
	Watch() error
	Unwatch() error
}
