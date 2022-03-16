package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("package", "common")

// ErrNotFound is returned when the download fails due to HTTP 404 Not Found.
var ErrNotFound = errors.New("not found")

func DownloadFile(url, destinationDir string, httpTimeout time.Duration) (string, error) {
	log.WithField("url", url).Debug("Downloading file")
	var name string
	var err error
	var retry bool
	for a := 1; a <= 3; a++ {
		name, retry, err = downloadFile(url, destinationDir, httpTimeout)
		if err != nil && retry {
			log.WithError(err).Debugf("Download attempt %d failed", a)
			continue
		}
		break
	}
	return name, err
}

func downloadFile(url, destinationDir string, httpTimeout time.Duration) (path string, retryable bool, err error) {
	client := http.Client{
		Timeout: httpTimeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", true, fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return "", false, ErrNotFound
		}
		return "", true, fmt.Errorf("download failed with http status %v: %w", resp.StatusCode, err)
	}

	name := filepath.Join(destinationDir, filepath.Base(url))
	f, err := os.Create(name)
	if err != nil {
		return "", false, fmt.Errorf("failed to create output file: %w", err)
	}

	numBytes, err := io.Copy(f, resp.Body)
	if err != nil {
		return "", true, fmt.Errorf("failed to write file to disk: %w", err)
	}
	log.WithFields(logrus.Fields{"file": name, "size_bytes": numBytes}).Debug("Download complete")

	return name, false, nil
}

// Rename renames src to dest. If the rename operation fails it will attempt to
// recursively copy the src to dest then delete src.
func Rename(src, dest string) error {
	err := os.Rename(src, dest)
	if err == nil {
		return nil
	}

	// Try copying.
	log := log.WithFields(logrus.Fields{"function": "Rename"})
	log.WithError(err).Debug("Falling back to a recursive copy after the rename operation failed.")
	if err = copy.Copy(src, dest); err != nil {
		return fmt.Errorf("copy/delete operation failed: %w", err)
	}
	if err = os.RemoveAll(src); err != nil {
		log.WithError(err).Warnf("Failed to delete source (%q) after copy operation.", src)
	}
	return nil
}
