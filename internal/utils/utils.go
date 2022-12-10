package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func MustGetUserHomeDir() string {
	var (
		home string
		err  error
	)

	if home, err = os.UserHomeDir(); err != nil {
		panic(fmt.Sprintf("failed to get home directory, err=%+v", err))
	}

	if _, err = os.Stat(home); err != nil {
		panic(fmt.Sprintf("failed to access home dir, err=%+v", err))
	}

	return home
}

func GetUserHomeDir() (string, error) {
	var (
		home string
		err  error
	)

	if home, err = os.UserHomeDir(); err != nil {
		return "", errors.Wrap(err, "failed to get home directory")
	}

	if _, err = os.Stat(home); err != nil {
		return "", errors.Wrap(err, "failed to access home directory")
	}

	return home, nil
}

func GetGVMWorkDir() (string, error) {
	var (
		home string
		err  error
	)

	if home, err = GetUserHomeDir(); err != nil {
		return "", err
	}

	return filepath.Join(home, ".gvm"), nil
}
