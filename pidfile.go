package pidfile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	ErrFileInvalid    = errors.New("pidfile has invalid contents")
	ErrProcessRunning = errors.New("process is running")
)

type PIDFile struct {
	path string
}

func New(p string) (*PIDFile, error) {
	if err := checkPidExists(p); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(p), os.FileMode(0755)); err != nil {
		return nil, err
	}

	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(p, []byte(pid), os.FileMode(0644)); err != nil {
		return nil, err
	}

	return &PIDFile{p}, nil
}

func (f PIDFile) Remove() error {
	return os.Remove(f.path)
}

func checkPidExists(p string) error {
	if pidByte, err := ioutil.ReadFile(p); err == nil {
		pidStr := strings.TrimSpace(string(pidByte))
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return ErrFileInvalid
		}

		if isProcessRunning(pid) {
			return ErrProcessRunning
		}
	}

	return nil
}

func isProcessRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}

	return true
}
