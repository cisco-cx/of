package v2

import (
	"log"

	"github.com/fsnotify/fsnotify"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
)

// Implements Watcher interface
type Notifier struct {
	l       *logger.Logger
	path    string
	w       *fsnotify.Watcher
	Changed chan string
	Errored chan error
}

// Init path Watcher.
func NewPath(path string, l *logger.Logger) (*Notifier, error) {
	fs := Notifier{}
	fs.l = l
	var err error
	fs.w, err = fsnotify.NewWatcher()
	if err != nil {
		fs.l.WithError(err).Errorf("Failed to init watcher.")
		return nil, err
	}

	fs.Changed = make(chan string)
	fs.Errored = make(chan error)
	return &fs, nil
}

// Watch for change in given path. Path can be a file or directory.
func (fs *Notifier) Watch() error {
	go func() {
		for {
			select {
			case event, ok := <-fs.w.Events:
				// Channel closed.
				if !ok {
					return
				}
				fs.l.Tracef("event: %+v", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fs.l.Debugf("modified file: %s", event.Name)
					fs.Changed <- event.Name
				}
			case err, ok := <-fs.w.Errors:
				// Channel closed.
				if !ok {
					return
				}
				fs.l.WithError(err).Errorf("Error occured while watching, %s", fs.path)
				fs.Errored <- err
			}
		}
	}()
	err := fs.w.Add(fs.path)
	if err != nil {
		log.Fatal(err)
		fs.l.WithError(err).Errorf("Failed to watch %s", fs.path)
		return err
	}
	return nil
}

// Stop watching given path.
func (fs *Notifier) Unwatch() error {
	fs.w.Remove(fs.path)
	close(fs.Changed)
	close(fs.Errored)
	return fs.w.Close()
}
