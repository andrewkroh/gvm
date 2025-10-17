package gvm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GoRelease represents a Go release from the go.dev API
type GoRelease struct {
	Version string   `json:"version"`
	Stable  bool     `json:"stable"`
	Files   []GoFile `json:"files"`
}

// GoFile represents a downloadable file for a Go release
type GoFile struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
}

// fetchGoReleases fetches the list of Go releases from the go.dev API
func (m *Manager) fetchGoReleases() ([]GoRelease, error) {
	apiURL := "https://go.dev/dl/?mode=json&include=all"

	client := &http.Client{
		Timeout: m.HTTPTimeout,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Go releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var releases []GoRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	return releases, nil
}

// findArchiveFile finds the archive file for the given OS/arch combination
func (r *GoRelease) findArchiveFile(goos, goarch string) *GoFile {
	archToMatch := goarch

	// Special case: ARM binary releases are only available for ARMv6
	// When GOARCH is "arm", we need to look for "armv6l" files
	if goarch == "arm" {
		archToMatch = "armv6l"
	}

	for i := range r.Files {
		file := &r.Files[i]

		// Skip non-archive files (installers, source)
		if file.Kind != "archive" {
			continue
		}

		// Match OS
		if file.OS != goos {
			continue
		}

		// Match architecture
		if file.Arch != archToMatch {
			continue
		}

		// For Windows, ensure we get .zip files
		if goos == "windows" && !hasExtension(file.Filename, ".zip") {
			continue
		}

		// For non-Windows, ensure we get .tar.gz files
		if goos != "windows" && !hasExtension(file.Filename, ".tar.gz") {
			continue
		}

		return file
	}

	return nil
}

// hasExtension checks if filename has the given extension
func hasExtension(filename, ext string) bool {
	return len(filename) > len(ext) && filename[len(filename)-len(ext):] == ext
}

// constructDownloadURL constructs the download URL for a given filename
func constructDownloadURL(filename string) string {
	return fmt.Sprintf("https://go.dev/dl/%s", filename)
}
