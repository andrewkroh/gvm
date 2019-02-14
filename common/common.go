package common

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

var log = logrus.WithField("package", "common")

func DownloadFile(url, destinationDir string) (string, error) {
	log.WithField("url", url).Debug("downloading file")
	var name string
	var err error
	var retry bool
	for a := 1; a <= 3; a++ {
		name, err, retry = downloadFile(url, destinationDir)
		if err != nil && retry {
			log.WithError(err).Debugf("Attempt %d failed", a)
			continue
		}
		break
	}
	return name, err
}

func downloadFile(url, destinationDir string) (string, error, bool) {
	client := http.Client{
		Timeout: time.Duration(3 * time.Minute),
	}
	resp, err := client.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "http get failed"), true
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("download failed with http status %v", resp.StatusCode), resp.StatusCode != http.StatusNotFound
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
	log.WithFields(logrus.Fields{"file": name, "size_bytes": numBytes}).Debug("download complete")

	return name, nil, false
}
