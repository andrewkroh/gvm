package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/andrewkroh/gvm"
	"github.com/stretchr/testify/assert"
)

func TestGVMRunUse(t *testing.T) {
	var cases = []struct {
		Version string
		Format  string
		Cmds    []string
	}{
		{"1.7.4", "bash", []string{"export GOROOT=", "export PATH"}},
		{"1.7.4", "batch", []string{"set GOROOT=", "set PATH"}},
		{"1.7.4", "powershell", []string{"$env:GOROOT = ", "$env:PATH ="}},
		{"1.5.4", "bash", []string{"export GOROOT=", "export PATH", `export GO15VENDOREXPERIMENT="1"`}},
		{"1.5.4", "batch", []string{"set GOROOT=", "set PATH=", `set GO15VENDOREXPERIMENT=1`}},
		{"1.5.4", "powershell", []string{"$env:GOROOT = ", "$env:PATH =", `$env:GO15VENDOREXPERIMENT = "1"`}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v_%v", tc.Version, tc.Format), func(t *testing.T) {

			output, err := withStdout(func() {
				manager := &gvm.Manager{}
				if err := manager.Init(); err != nil {
					t.Fatal(err)
				}

				cmd := &useCmd{
					version: tc.Version,
					format:  tc.Format,
				}
				err := cmd.Run(manager)
				if err != nil {
					t.Fatal(err)
				}
			})
			if err != nil {
				t.Fatal(err)
			}

			t.Log(output)
			lines := strings.Split(output, "\n")

			if !assert.Len(t, lines, len(tc.Cmds)+1, "expected %d lines, got [%v]", strings.Join(lines, "|")) {
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
					goroot = strings.TrimSpace(parts[1])

					if unquotedPath, err := strconv.Unquote(goroot); err == nil {
						goroot = unquotedPath
					}
				}
			}

			// Test that GOROOT/bin/go exists and is the correct version.
			version, err := exec.Command(filepath.Join(goroot, "bin", "go"), "version").Output()
			if err != nil {
				t.Fatal("failed to run go version", err)
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
