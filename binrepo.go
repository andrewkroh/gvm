package gvm

import (
	"fmt"
	"os"

	"github.com/andrewkroh/gvm/common"
)

func (m *Manager) installBinary(version *GoVersion) (string, error) {
	// Fetch releases to find the correct file
	releases, err := m.fetchGoReleases()
	if err != nil {
		return "", fmt.Errorf("failed to fetch releases: %w", err)
	}

	// Find the release for this version
	versionStr := fmt.Sprintf("go%v", version)
	var targetRelease *GoRelease
	for i := range releases {
		if releases[i].Version == versionStr {
			targetRelease = &releases[i]
			break
		}
	}

	if targetRelease == nil {
		return "", common.ErrNotFound
	}

	// Find the archive file for this OS/arch
	file := targetRelease.findArchiveFile(m.GOOS, m.GOARCH)
	if file == nil {
		return "", common.ErrNotFound
	}

	godir := m.versionDir(version)

	tmp, err := os.MkdirTemp("", godir)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)

	// Construct download URL using the filename from the API
	goURL := constructDownloadURL(file.Filename)
	path, err := common.DownloadFile(goURL, tmp, m.HTTPTimeout, common.DefaultRetryParams)
	if err != nil {
		return "", fmt.Errorf("failed downloading from %v: %w", goURL, err)
	}

	return extractTo(m.VersionGoROOT(version), path)
}

func (m *Manager) AvailableBinaries() ([]*GoVersion, error) {
	releases, err := m.fetchGoReleases()
	if err != nil {
		return nil, err
	}

	// Use a map to store unique versions
	versionSet := make(map[string]*GoVersion)

	for _, release := range releases {
		// Find the archive file for this OS/arch combination
		file := release.findArchiveFile(m.GOOS, m.GOARCH)
		if file == nil {
			continue
		}

		// Parse the version (remove "go" prefix)
		versionStr := release.Version
		if len(versionStr) > 2 && versionStr[:2] == "go" {
			versionStr = versionStr[2:]
		}

		ver, err := ParseVersion(versionStr)
		if err != nil {
			continue
		}

		// Store in map to ensure uniqueness
		versionSet[ver.String()] = ver
	}

	// Convert map to slice
	list := make([]*GoVersion, 0, len(versionSet))
	for _, ver := range versionSet {
		list = append(list, ver)
	}

	sortVersions(list)
	return list, nil
}
