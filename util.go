package gvm

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/andrewkroh/gvm/common"
	"github.com/pkg/errors"
)

func homeDir() (string, error) {
	var homeDir string
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	} else {
		homeDir = os.Getenv("HOME")
	}

	if _, err := os.Stat(homeDir); err != nil {
		return "", errors.Wrap(err, "failed to access home dir")
	}

	return homeDir, nil
}

func extractTo(to, file string) (string, error) {
	tmpDir := to + ".tmp"
	if err := os.Mkdir(tmpDir, 0755); err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	if err := common.Extract(file, tmpDir); err != nil {
		return "", err
	}

	// Move into the final location.
	if err := os.Rename(filepath.Join(tmpDir, "go"), to); err != nil {
		return "", err
	}
	return to, nil
}

func existsDir(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
