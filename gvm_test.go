package main

import (
	"bytes"
	"testing"

	"fmt"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestParseTravisYML(t *testing.T) {
	ver, err := parseTravisYml("testdata/travis1.yml")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1.7.4", ver)

	ver, err = parseTravisYml("testdata/travis2.yml")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1.7.4", ver)
}

func TestGVMRun(t *testing.T) {
	var cases = []struct {
		Version string
		Format  string
		Cmds    []string
	}{
		{"1.7.4", "bash", []string{"export GOROOT=", "export PATH"}},
		{"1.7.4", "batch", []string{"set GOROOT=", "set PATH"}},
		{"1.7.4", "powershell", []string{"$env:GOROOT = ", "$env:PATH ="}},
		{"1.5.4", "bash", []string{"export GOROOT=", "export PATH", "export GO15VENDOREXPERIMENT=1"}},
		{"1.5.4", "batch", []string{"set GOROOT=", "set PATH=", "set GO15VENDOREXPERIMENT=1"}},
		{"1.5.4", "powershell", []string{"$env:GOROOT = ", "$env:PATH =", "$env:GO15VENDOREXPERIMENT=1"}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v_%v", tc.Version, tc.Format), func(t *testing.T) {
			out := new(bytes.Buffer)
			g := &GVM{
				Version: tc.Version,
				Format:  tc.Format,
				out:     out,
			}

			err := g.Run(nil)
			if err != nil {
				t.Fatal(err)
			}

			output := out.String()
			t.Log(output)
			lines := strings.Split(output, "\n")

			if !assert.Len(t, lines, len(tc.Cmds)+1, "expected %d lines, got [%v]", strings.Join(lines, "|")) {
				return
			}

			for i, line := range lines[:len(lines)-1] {
				assert.Contains(t, line, tc.Cmds[i])
			}
		})
	}

}
