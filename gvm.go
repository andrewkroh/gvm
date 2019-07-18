package gvm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
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

	// GoSourceURL configres the update git repository to download and update local
	// source checkouts from.
	// Defaults to https://go.googlesource.com/go
	GoSourceURL string

	Logger logrus.FieldLogger

	cacheDir    string
	versionsDir string
	logsDir     string
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

	if m.GoSourceURL == "" {
		m.GoSourceURL = "https://go.googlesource.com/go"
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

	if m.Logger == nil {
		m.Logger = logrus.StandardLogger()
	}

	m.cacheDir = filepath.Join(m.Home, "cache")
	m.versionsDir = filepath.Join(m.Home, "versions")
	m.logsDir = filepath.Join(m.Home, "logs")
	return m.ensureDirStruct()
}

func (m *Manager) UpdateCache() error {
	return m.updateSrcCache()
}

func (m *Manager) ensureDirStruct() error {
	for _, dir := range []string{m.cacheDir, m.versionsDir, m.logsDir} {
		if err := os.MkdirAll(dir, os.ModeDir|0755); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Available() ([]*GoVersion, []bool, error) {
	if !m.hasSrcCache() {
		versions, err := m.AvailableBinaries()
		if err != nil {
			return nil, nil, err
		}

		hasBin := make([]bool, len(versions))
		for i := range hasBin {
			hasBin[i] = true
		}
		return versions, hasBin, nil
	}

	src, err := m.AvailableSource()
	if err != nil {
		return nil, nil, err
	}

	hasBin := make([]bool, len(src))
	bin, err := m.AvailableBinaries()
	if err != nil {
		return src, hasBin, nil
	}

	for i, ver := range src {
		hasBin[i] = findVersion(ver, bin) >= 0
	}

	return src, hasBin, nil
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

// Installed returns all installed go version
func (m *Manager) Installed() ([]*GoVersion, error) {
	files, err := ioutil.ReadDir(m.versionsDir)
	if err != nil {
		return nil, err
	}

	versionSuffix := fmt.Sprintf(".%v.%v", m.GOOS, m.GOARCH)

	var list []*GoVersion
	for _, fi := range files {
		name := fi.Name()
		name = strings.TrimSuffix(name, versionSuffix)
		name = strings.TrimPrefix(name, "go")

		v, err := ParseVersion(name)
		if err != nil {
			continue
		}
		list = append(list, v)
	}

	sortVersions(list)
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

func (m *Manager) Build(version *GoVersion) (string, error) {
	if version.IsTip() {
		return m.ensureUpToDateTip()
	}

	has, err := m.HasVersion(version)
	if err != nil {
		return "", err
	}
	if has {
		return m.VersionGoROOT(version), nil
	}

	return m.installSrc(version)
}

func (m *Manager) Install(version *GoVersion) (string, error) {
	if version.IsTip() {
		return m.ensureUpToDateTip()
	}

	has, err := m.HasVersion(version)
	if err != nil {
		return "", err
	}
	if has {
		return m.VersionGoROOT(version), nil
	}

	tryBinary := !version.IsTip()
	if tryBinary {
		dir, err := m.installBinary(version)
		if err == nil {
			return dir, err
		}
	}

	return m.installSrc(version)
}

func (m *Manager) ensureUpToDateTip() (string, error) {
	version, _ := ParseVersion("tip")

	has, err := m.HasVersion(version)
	if err != nil {
		return "", err
	}

	// no updates since last build -> return installed version
	if has {
		updates, err := m.tryRefreshSrcCache()
		if err != nil {
			return "", err
		}

		if !updates {
			return m.VersionGoROOT(version), nil
		}
	}

	// new updates in cache -> rebuild
	return m.installSrc(version)
}
