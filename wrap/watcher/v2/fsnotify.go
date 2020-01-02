package v2

import (
	"github.com/rjeczalik/notify"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
)

// Implements Watcher interface
type Notifier struct {
	l       *logger.Logger
	path    string
	Changed chan string
	c       chan notify.EventInfo
}

// Init path Watcher.
func NewPath(path string, l *logger.Logger) (*Notifier, error) {
	fs := Notifier{}
	fs.l = l
	fs.path = path
	fs.Changed = make(chan string)
	return &fs, nil
}

// Watch for change in given path. Path can be a file or directory.
func (fs *Notifier) Watch() error {
	fs.c = make(chan notify.EventInfo, 1)
	if err := notify.Watch(fs.path, fs.c, notify.All); err != nil {
		return err
	}
	go func() {
		for {
			fs.l.Debugf("Waiting for event.\n")
			ei := <-fs.c
			fs.l.Debugf("File/Dir modified, %s", ei.Path())
			fs.Changed <- ei.Path()
		}
	}()

	return nil
}

// Stop watching given path.
func (fs *Notifier) Unwatch() error {
	notify.Stop(fs.c)
	return nil
}
