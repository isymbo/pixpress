package util

import (
	"os"
	"os/exec"
	"path/filepath"

	log "gopkg.in/clog.v1"
)

// execPath returns the executable path, which includes executable filename.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Error(2, "Error to determine executable path", err)
	}

	return filepath.Abs(file)
}

// ExecDir returns the executable directory path
func ExecDir() string {
	execpath, _ := execPath()

	return filepath.Dir(execpath)
}
