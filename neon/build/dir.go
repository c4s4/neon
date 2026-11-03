package build

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

// parallel counts the groups of threads running. The current directory is
// global to the process, so it can't be changed while threads are running.
var parallel atomic.Int32

// EnterParallel must be called when starting threads, and LeaveParallel when
// they are done. While threads are running, Chdir fails if asked to change
// to another directory.
func EnterParallel() {
	parallel.Add(1)
}

// LeaveParallel must be called when threads started after EnterParallel are
// done.
func LeaveParallel() {
	parallel.Add(-1)
}

// Chdir changes current directory. While threads are running, changing to a
// directory other than the current one fails, as the current directory is
// shared by all threads.
// - dir: the directory to change to
// Return: an error if something went wrong
func Chdir(dir string) error {
	if parallel.Load() > 0 {
		same, err := isCurrentDir(dir)
		if err != nil {
			return err
		}
		if !same {
			return fmt.Errorf("can't change directory to '%s' while threads are running (current directory is shared by all threads)", dir)
		}
		return nil
	}
	return os.Chdir(dir)
}

// isCurrentDir tells if given directory is the current directory
func isCurrentDir(dir string) (bool, error) {
	current, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("getting current directory: %v", err)
	}
	target, err := filepath.Abs(dir)
	if err != nil {
		return false, err
	}
	if resolved, err := filepath.EvalSymlinks(target); err == nil {
		target = resolved
	}
	if resolved, err := filepath.EvalSymlinks(current); err == nil {
		current = resolved
	}
	return target == current, nil
}
