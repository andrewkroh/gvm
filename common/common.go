package common

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

var log = logrus.WithField("package", "common")

func DownloadFile(url, destinationDir string) (string, error) {
	log.WithField("url", url).Debug("downloading file")

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "http get failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("download failed with http status %v", resp.StatusCode)
	}

	name := filepath.Join(destinationDir, filepath.Base(url))
	f, err := os.Create(name)
	if err != nil {
		return "", errors.Wrap(err, "failed to create output file")
	}

	numBytes, err := io.Copy(f, resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to write file to disk")
	}
	log.WithFields(logrus.Fields{"file": name, "size_bytes": numBytes}).Debug("download complete")

	return name, nil
}

