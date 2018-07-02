package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/andrewkroh/gvm"
	"github.com/andrewkroh/gvm/cmd/gvm/internal/shellfmt"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func useCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	var (
		version string
		format  string
	)
	cmd.Arg("version", "Go version to install (e.g. 1.10).").StringVar(&version)
	cmd.Flag("format", "Format to use for the shell commands. Options: bash, batch, powershell").
		Short('f').
		Default(shellfmt.DefaultFormat()).
		EnumVar(&format, shellfmt.BashFormat, shellfmt.BatchFormat, shellfmt.PowershellFormat)

	return func(manager *gvm.Manager) error {
		if version == "" {
			return fmt.Errorf("no version specified")
		}
		ver, err := gvm.ParseVersion(version)
		if err != nil {
			return err
		}
		log.Debugf("Using Go version %v", ver)

		shellFmt, err := shellfmt.New(format)
		if err != nil {
			return err
		}

		goroot, err := manager.Install(ver)
		if err != nil {
			return err
		}

		shellFmt.Set("GOROOT", goroot)
		shellFmt.Prepend("PATH", filepath.Join(goroot, "bin"))
		if strings.HasPrefix(version, "1.5") {
			shellFmt.Set("GO15VENDOREXPERIMENT", "1")
		}

		return nil
	}
}
