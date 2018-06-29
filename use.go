package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewkroh/gvm/golang"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type gvmUse struct {
	version      string
	useProjectGo bool
	format       string

	out io.Writer // Stdout writer (used for capturing output in tests).
}

func useCommand(cmd *kingpin.CmdClause) func() error {
	g := &gvmUse{out: os.Stdout}
	cmd.Flag("project-go", "Use the project's Go version.").BoolVar(&g.useProjectGo)
	cmd.Flag("format", "Format to use for the shell commands. Options: bash, batch, powershell").
		Short('f').
		Default(defaultFormat()).
		EnumVar(&g.format, BashFormat, BatchFormat, PowershellFormat)
	cmd.Arg("version", "Go version to install (e.g. 1.10).").StringVar(&g.version)
	return g.run
}

func (g *gvmUse) run() error {
	version, err := useVersion(g.version, g.useProjectGo)
	if err != nil {
		return err
	}

	format, err := getEnvFormatter(g.format)
	if err != nil {
		return err
	}

	if version == "" {
		return fmt.Errorf("no version specified")
	}
	log.Debugf("Using Go version %v", version)

	goroot, err := golang.SetupGolang(version)
	if err != nil {
		return err
	}

	binDir := filepath.Join(goroot, "bin")
	fmt.Fprintln(g.out, format.Set("GOROOT", goroot))
	fmt.Fprintln(g.out, format.Prepend("PATH", binDir))
	if strings.HasPrefix(version, "1.5") {
		fmt.Fprintln(g.out, format.Set("GO15VENDOREXPERIMENT", "1"))
	}

	return nil
}
