package gvm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Manager struct {
	// GVM Home directory. Defaults to $HOME/.gvm
	Home string

	// GOOS settings. Defaults to current OS.
	GOOS string

	// GOARCH setting. Defaults to the current architecture.
	GOARCH string

	// Golang binary store URL. Used to download listing and go binaries.
	// Defaults to https://storage.googleapis.com/golang
	GoStorageHome string

	cacheDir    string
	versionsDir string
}

func (m *Manager) Init() error {
	if m.Home == "" {
		home, err := homeDir()
		if err != nil {
			return err
		}

		m.Home = filepath.Join(home, ".gvm")
	}

	if m.GoStorageHome == "" {
		m.GoStorageHome = "https://storage.googleapis.com/golang"
	}

	if m.GOOS == "" {
		m.GOOS = runtime.GOOS
	}
	if m.GOARCH == "" {
		switch runtime.GOARCH {
		default:
			m.GOARCH = runtime.GOARCH
		case "arm":
			// The only binary releases are for ARM v6.
			m.GOARCH = "armv6l"
		}
	}

	m.cacheDir = filepath.Join(m.Home, "cache")
	m.versionsDir = filepath.Join(m.Home, "versions")
	return m.ensureDirStruct()
}

func (m *Manager) ensureDirStruct() error {
	for _, dir := range []string{m.cacheDir, m.versionsDir} {
		if err := os.MkdirAll(dir, os.ModeDir|0755); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Available() ([]*GoVersion, error) {
	return m.AvailableBinaries()
}

func (m *Manager) Remove(version *GoVersion) error {
	dir := m.VersionGoROOT(version)

	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("Version %v not installed\n", version)
	}

	if !fi.IsDir() {
		return fmt.Errorf("Path %v is no directory", dir)
	}

	return os.RemoveAll(dir)
}

// Installed returns all installed go version (unparsed version numbers)
func (m *Manager) Installed() ([]string, error) {
	files, err := ioutil.ReadDir(m.versionsDir)
	if err != nil {
		return nil, err
	}

	versionSuffix := fmt.Sprintf(".%v.%v", m.GOOS, m.GOARCH)

	var list []string
	for _, fi := range files {
		name := fi.Name()
		name = strings.TrimSuffix(name, versionSuffix)
		name = strings.TrimPrefix(name, "go")
		list = append(list, name)
	}
	return list, nil
}

// HasVersion checks if a given go version is installed
func (m *Manager) HasVersion(version *GoVersion) (bool, error) {
	return existsDir(m.VersionGoROOT(version))
}

// VersionGoROOT returns the GOROOT path for a go version. VersionGoROOT does
// not check if the version is installed.
func (m *Manager) VersionGoROOT(version *GoVersion) string {
	return filepath.Join(m.versionsDir, m.versionDir(version))
}

func (m *Manager) versionDir(version *GoVersion) string {
	return fmt.Sprintf("go%v.%v.%v", version, m.GOOS, m.GOARCH)
}

func (m *Manager) Install(version *GoVersion) (string, error) {
	has, err := m.HasVersion(version)
	if err != nil {
		return "", err
	}
	if has {
		return m.VersionGoROOT(version), nil
	}

	return m.installBinary(version)
}
