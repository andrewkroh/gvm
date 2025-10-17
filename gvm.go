// Package gvm provides a Go Version Manager that allows you to install and manage
// multiple versions of Go. It supports downloading pre-built binaries from the
// official Go downloads API as well as building from source.
package gvm

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/andrewkroh/gvm/common"
)

type AvailableVersion struct {
	Version *GoVersion
	Source  bool // Available to install from source.
	Binary  bool // Available to download as a binary.
}

func (av AvailableVersion) String() string {
	switch {
	case av.Source && av.Binary:
		return fmt.Sprintf("%v\t(source, binary)", av.Version)
	case av.Source:
		return fmt.Sprintf("%v\t(source)", av.Version)
	case av.Binary:
		return fmt.Sprintf("%v\t(binary)", av.Version)
	default:
		return av.Version.String()
	}
}

type Manager struct {
	// GVM Home directory. Defaults to $HOME/.gvm
	Home string

	// GOOS settings. Defaults to current OS.
	GOOS string

	// GOARCH setting. Defaults to the current architecture.
	GOARCH string

	// GoStorageHome is the base URL for the Go downloads API.
	// Defaults to https://go.dev/dl
	GoStorageHome string

	// GoSourceURL configres the update git repository to download and update local
	// source checkouts from.
	// Defaults to https://go.googlesource.com/go
	GoSourceURL string

	HTTPTimeout time.Duration

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
		m.GoStorageHome = "https://go.dev/dl"
	}

	if m.GoSourceURL == "" {
		m.GoSourceURL = "https://go.googlesource.com/go"
	}

	if m.HTTPTimeout == 0 {
		m.HTTPTimeout = 3 * time.Minute
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
		if err := os.MkdirAll(dir, os.ModeDir|0o755); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) Available() ([]AvailableVersion, error) {
	if !m.hasSrcCache() {
		versions, err := m.AvailableBinaries()
		if err != nil {
			return nil, err
		}

		available := make([]AvailableVersion, 0, len(versions))
		for _, ver := range versions {
			available = append(available, AvailableVersion{Version: ver, Binary: true})
		}
		return available, nil
	}

	src, err := m.AvailableSource()
	if err != nil {
		return nil, err
	}

	versionSet := make(map[string]*AvailableVersion, len(src))
	for _, ver := range src {
		versionSet[ver.String()] = &AvailableVersion{Version: ver, Source: true}
	}

	toSlice := func() []AvailableVersion {
		available := make([]AvailableVersion, 0, len(versionSet))
		for _, ver := range versionSet {
			available = append(available, *ver)
		}
		sort.Slice(available, func(i, j int) bool {
			return available[i].Version.LessThan(available[j].Version)
		})
		return available
	}

	bin, err := m.AvailableBinaries()
	if err != nil {
		// Return the source versions if we cannot get binary info.
		m.Logger.WithError(err).Info("Failed to list available binary versions.")
		return toSlice(), nil
	}

	// Merge source and binary versions.
	for _, ver := range bin {
		if avail, found := versionSet[ver.String()]; found {
			avail.Binary = true
			continue
		}
		versionSet[ver.String()] = &AvailableVersion{Version: ver, Binary: true}
	}

	return toSlice(), nil
}

func (m *Manager) Remove(version *GoVersion) error {
	dir := m.VersionGoROOT(version)

	fi, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("version %q not installed", version)
	}

	if !fi.IsDir() {
		return fmt.Errorf("path %q is not a directory", dir)
	}

	return os.RemoveAll(dir)
}

// Installed returns all installed go version
func (m *Manager) Installed() ([]*GoVersion, error) {
	files, err := os.ReadDir(m.versionsDir)
	if err != nil {
		return nil, err
	}

	versionSuffix := fmt.Sprintf(".%v.%v", m.GOOS, m.GOARCH)

	list := make([]*GoVersion, 0, len(files))
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

	if tryBinary := !version.IsTip(); tryBinary {
		dir, err := m.installBinary(version)
		if err == nil {
			return dir, nil
		}
		// Only continue to installing from source if the server confirms 404.
		if !errors.Is(err, common.ErrNotFound) {
			return "", err
		}
		m.Logger.Debug("Binary release not found on server. Trying to install from source.")
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
