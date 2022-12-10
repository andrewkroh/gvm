package gvm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/andrewkroh/gvm/common"
)

func extractTo(to, file string) (string, error) {
	tmpDir := to + ".tmp"
	if err := os.Mkdir(tmpDir, 0o755); err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	if err := common.Extract(file, tmpDir); err != nil {
		return "", err
	}

	// Move into the final location.
	if err := common.Rename(filepath.Join(tmpDir, "go"), to); err != nil {
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

func writeJSONFile(filename string, value interface{}) error {
	contents, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, contents, 0o644)
}

func readJSONFile(filename string, to interface{}) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, to)
}
