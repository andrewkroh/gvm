package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrewkroh/gvm"
)

var (
	go1_6  = gvm.MustParseVersion("1.6")
	go1_9  = gvm.MustParseVersion("1.9")
	go1_16 = gvm.MustParseVersion("1.16")
)

func TestGVMRunUse(t *testing.T) {
	// When testing building from source GOROOT_BOOTSTRAP must be set.
	t.Setenv("GOROOT_BOOTSTRAP", build.Default.GOROOT)

	cases := []struct {
		Version    string
		Format     string
		Cmds       []string
		FromSource bool
	}{
		// Using 1.16+ allows testing on Apple M1.
		{Version: "1.16.15", Format: "bash", Cmds: []string{"export GOROOT=", "export PATH"}},
		{Version: "1.16.15", Format: "batch", Cmds: []string{"set GOROOT=", "set PATH"}},
		{Version: "1.16.15", Format: "powershell", Cmds: []string{"$env:GOROOT = ", "$env:PATH ="}},
		// Check that newer versions which use go.mod can be build from source.
		{Version: "1.16.14", FromSource: true, Format: "bash", Cmds: []string{"export GOROOT=", "export PATH"}},
		// Check that older versions which did not use go.mod can be build from source.
		{Version: "1.10.8", FromSource: true, Format: "bash", Cmds: []string{"export GOROOT=", "export PATH"}},
		// Check that GO15VENDOREXPERIMENT is added for Go 1.5.
		// NOTE: 1.5 requires Go 1.4 for bootstrapping if built from source.
		{Version: "1.5.4", Format: "bash", Cmds: []string{"export GOROOT=", "export PATH", `export GO15VENDOREXPERIMENT="1"`}},
		{Version: "1.5.4", Format: "batch", Cmds: []string{"set GOROOT=", "set PATH=", `set GO15VENDOREXPERIMENT=1`}},
		{Version: "1.5.4", Format: "powershell", Cmds: []string{"$env:GOROOT = ", "$env:PATH =", `$env:GO15VENDOREXPERIMENT = "1"`}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v_%v", tc.Version, tc.Format), func(t *testing.T) {
			ver := gvm.MustParseVersion(tc.Version)

			// Go introduced support for Apple M1 in Go 1.16.
			if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" && ver.LessThan(go1_16) {
				t.Skip("darwin/arm64 only supports Go 1.16 and newer")
			}

			// Go 1.5 needs Go 1.4 to build it from source so avoid that problem
			// by only testing on platforms where Go 1.5 is available as a binary.
			if ver.LessThan(go1_6) && runtime.GOARCH != "amd64" {
				t.Skip("Binary distributions of Go 1.5 are not available.")
			}

			output, err := withStdout(func() {
				manager := &gvm.Manager{}
				if err := manager.Init(); err != nil {
					t.Fatal(err)
				}

				cmd := &useCmd{
					version: tc.Version,
					format:  tc.Format,
					build:   tc.FromSource,
				}
				err := cmd.Run(manager)
				if err != nil {
					t.Fatal(err)
				}
			})
			if err != nil {
				t.Fatal(err)
			}

			output = strings.TrimSpace(output)
			t.Log(output)
			lines := strings.Split(output, "\n")

			if !assert.Lenf(t, lines, len(tc.Cmds), "expected %d lines, got %d [%v]", len(lines), len(tc.Cmds)+1, strings.Join(lines, "|")) {
				return
			}

			var goroot string
			for i, line := range lines[:len(lines)-1] {
				assert.Contains(t, line, tc.Cmds[i])

				if !strings.Contains(line, "PATH") && strings.Contains(line, "GOROOT") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) != 2 {
						t.Fatal("failed to parse GOROOT", line)
					}
					goroot = strings.Trim(parts[1], ` "`)
				}
			}

			// Test that GOROOT/bin/go exists and is the correct version.
			goVersionCmd := exec.Command(filepath.Join(goroot, "bin", "go"), "version")

			// Go versions less than 1.9 require GOROOT to be set.
			if ver.LessThan(go1_9) {
				goVersionCmd.Env = append(goVersionCmd.Env, "GOROOT="+goroot)
			}

			version, err := goVersionCmd.Output()
			if err != nil {
				var exitErr *exec.ExitError
				if errors.As(err, &exitErr) {
					t.Fatalf("failed to run go version: %v\n%s", err, exitErr.Stderr)
				} else {
					t.Fatal("failed to run go version", err)
				}
			}
			assert.Contains(t, string(version), tc.Version)
		})
	}
}

// capture stdout and return captured string
func withStdout(fn func()) (string, error) {
	stdout := os.Stdout

	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}

	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	outC := make(chan string)
	go func() {
		// capture all output
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		r.Close()
		outC <- buf.String()
	}()

	fn()
	w.Close()
	result := <-outC
	return result, err
}
