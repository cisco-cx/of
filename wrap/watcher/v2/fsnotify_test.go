package v2_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	of "github.com/cisco-cx/of/pkg/v2"
	logger "github.com/cisco-cx/of/wrap/logrus/v2"
	watcher "github.com/cisco-cx/of/wrap/watcher/v2"
)

// Enforce interface implementation.
func TestFSInterface(t *testing.T) {
	var _ of.Watcher = &watcher.Notifier{}
}

// Test file watch.
func TestFileWatch(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "watcher.*")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	l := logger.New()

	// Watching file
	fs, err := watcher.NewPath(tmpFile.Name(), l)
	require.NoError(t, err)
	go func() {
		fileName := <-fs.Changed
		require.Equal(t, tmpFile.Name(), fileName)
		err := fs.Unwatch()
		require.NoError(t, err)
	}()

	_, err = tmpFile.Write([]byte("Testing watcher."))
	require.NoError(t, err)
}

// Test directory watch.
func TestDirWatch(t *testing.T) {
	dir, err := ioutil.TempDir("", "watcher_dir")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	tmpFile, err := ioutil.TempFile(dir, "watcher.*")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	l := logger.New()

	// Watching dir
	fs, err := watcher.NewPath(dir, l)
	require.NoError(t, err)
	go func() {
		fileName := <-fs.Changed
		require.Equal(t, tmpFile.Name(), fileName)
		err := fs.Unwatch()
		require.NoError(t, err)
	}()

	_, err = tmpFile.Write([]byte("Testing watcher."))
	require.NoError(t, err)
}
