package main

import (
	"fmt"
	"path/filepath"

	"github.com/andrewkroh/gvm"
	"github.com/andrewkroh/gvm/cmd/gvm/internal/shellfmt"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type useCmd struct {
	version string
	build   bool
	format  string
}

func useCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	ctx := &useCmd{}

	cmd.Arg("version", "Go version to install (e.g. 1.10).").StringVar(&ctx.version)
	cmd.Flag("build", "Build go version from source").Short('b').BoolVar(&ctx.build)
	cmd.Flag("format", "Format to use for the shell commands. Options: bash, batch, powershell").
		Short('f').
		Default(shellfmt.DefaultFormat()).
		EnumVar(&ctx.format, shellfmt.BashFormat, shellfmt.BatchFormat, shellfmt.PowershellFormat)

	return ctx.Run
}

func (cmd *useCmd) Run(manager *gvm.Manager) error {
	if cmd.version == "" {
		return fmt.Errorf("no version specified")
	}
	ver, err := gvm.ParseVersion(cmd.version)
	if err != nil {
		return err
	}
	log.Debugf("Using Go version %v", ver)

	shellFmt, err := shellfmt.New(cmd.format)
	if err != nil {
		return err
	}

	var goroot string
	if cmd.build {
		goroot, err = manager.Build(ver)
	} else {
		goroot, err = manager.Install(ver)
	}
	if err != nil {
		return err
	}

	shellFmt.Set("GOROOT", goroot)
	shellFmt.Prepend("PATH", filepath.Join(goroot, "bin"))
	if _, experimental := ver.VendorSupport(); experimental {
		shellFmt.Set("GO15VENDOREXPERIMENT", "1")
	}

	return nil
}
