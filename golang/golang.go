package golang

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"github.com/andrewkroh/gvm/common"
)

const (
	downloadBase = "https://storage.googleapis.com/golang"
)

var executableExtension = ""

func init() {
	if runtime.GOOS == "windows" {
		executableExtension = ".exe"
	}
}

// SetupGolang returns the GOROOT for a Go installation.
func SetupGolang(version string) (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", err
	}

	goDir, err := golangDir(home, version, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return "", err
	}

	if !isGoInstalled(goDir) {
		tmp, err := ioutil.TempDir("", filepath.Base(goDir))
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmp)

		file, err := downloadGo(version, runtime.GOOS, runtime.GOARCH, tmp)
		if err != nil {
			return "", err
		}

		if err := os.MkdirAll(goDir, 0755); err != nil {
			if !os.IsExist(err) {
				return "", err
			}
		}

		err = common.Extract(file, goDir)
		if err != nil {
			// Remove on error because the extracted data could be corrupt.
			os.RemoveAll(goDir)
			return "", err
		}
	}

	return goDir, nil
}

// downloadGo downloads the Golang package over HTTPS.
func downloadGo(version, goos, arch, destinationDir string) (string, error) {
	//	Example: https://storage.googleapis.com/golang/go1.7.3.windows-amd64.zip
	extension := "tar.gz"
	if goos == "windows" {
		extension = "zip"
	}

	goURL := fmt.Sprintf("%s/go%v.%v-%v.%v", downloadBase, version, goos, arch, extension)
	return common.DownloadFile(goURL, destinationDir)
}

func golangDir(home, version, goos, goarch string) (string, error) {
	return filepath.Join(home, ".gvm", "versions", fmt.Sprintf("go%s.%s.%s", version, goos, goarch)), nil
}

func isGoInstalled(goDir string) bool {
	// Test if Go exists at the GOPATH.
	_, err := os.Stat(filepath.Join(goDir, "bin", "go"+executableExtension))
	if err != nil {
		return false
	}

	return true
}

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
