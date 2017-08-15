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
		// Download the package to a temp dir.
		tmp, err := ioutil.TempDir("", filepath.Base(goDir))
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmp)

		file, err := downloadGo(version, runtime.GOOS, runtime.GOARCH, tmp)
		if err != nil {
			return "", err
		}

		// Make the ~/.gvm/versions directory.
		if err := os.MkdirAll(filepath.Dir(goDir), 0755); err != nil {
			if !os.IsExist(err) {
				return "", err
			}
		}

		// Extract to a temp dir on the same volume as the destination to avoid
		// potential EXDEV (invalid cross-device link) errors when renaming.
		goDirTmp := goDir + ".tmp"
		if err = os.Mkdir(goDirTmp, 0755); err != nil {
			return "", err
		}
		defer os.RemoveAll(goDirTmp)

		err = common.Extract(file, goDirTmp)
		if err != nil {
			return "", err
		}

		// Move into the final location.
		if err = os.Rename(filepath.Join(goDirTmp, "go"), goDir); err != nil {
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
