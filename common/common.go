package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("package", "common")

// ErrNotFound is returned when the download fails due to HTTP 404 Not Found.
var ErrNotFound = errors.New("not found")

func DownloadFile(url, destinationDir string) (string, error) {
	log.WithField("url", url).Debug("Downloading file")
	var name string
	var err error
	var retry bool
	for a := 1; a <= 3; a++ {
		name, err, retry = downloadFile(url, destinationDir)
		if err != nil && retry {
			log.WithError(err).Debugf("Download attempt %d failed", a)
			continue
		}
		break
	}
	return name, err
}

func downloadFile(url, destinationDir string) (path string, err error, retryable bool) {
	client := http.Client{
		Timeout: time.Duration(3 * time.Minute),
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "http get failed"), true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return "", ErrNotFound, false
		}
		return "", errors.Errorf("download failed with http status %v", resp.StatusCode), true
	}

	name := filepath.Join(destinationDir, filepath.Base(url))
	f, err := os.Create(name)
	if err != nil {
		return "", errors.Wrap(err, "failed to create output file"), false
	}

	numBytes, err := io.Copy(f, resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to write file to disk"), true
	}
	log.WithFields(logrus.Fields{"file": name, "size_bytes": numBytes}).Debug("Download complete")

	return name, nil, false
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
