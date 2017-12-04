package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/andrewkroh/gvm/golang"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

const usage = `gvm is a Go version manager. gvm installs a Go version and prints
the commands to configure your environment to use it. gvm can only install
binary versions of Go from https://golang.org/dl/. Below are examples for
common shells.

  bash:
    eval "$(gvm 1.9.2)"

  batch (windows cmd.exe):
    FOR /f "tokens=*" %i IN ('"gvm.exe" 1.9.2') DO %i

  powershell:
    gvm --format=powershell 1.9.2 | Invoke-Expression
`

// Output formats.
const (
	BashFormat       = "bash"
	BatchFormat      = "batch"
	PowershellFormat = "powershell"
)

var (
	version = "SNAPSHOT"

	log = logrus.WithField("package", "main")
)

type GVM struct {
	Version      string
	UseProjectGo bool
	Format       string

	out io.Writer // Stdout writer (used for capturing output in tests).
}

func (g *GVM) Run(_ *kingpin.ParseContext) error {
	version := g.Version
	if g.UseProjectGo {
		ver, err := getProjectGoVersion()
		if err != nil {
			return err
		}
		version = ver
	}

	if version == "" {
		return fmt.Errorf("no version specified")
	}
	log.Debugf("Using Go version %v", version)

	goroot, err := golang.SetupGolang(version)
	if err != nil {
		return err
	}

	switch g.Format {
	case BashFormat:
		fmt.Fprintf(g.out, `export GOROOT="%v"`+"\n", goroot)
		fmt.Fprintf(g.out, `export PATH="$GOROOT/bin:$PATH"`+"\n")
		if strings.HasPrefix(version, "1.5") {
			fmt.Fprintln(g.out, `export GO15VENDOREXPERIMENT=1`)
		}
	case BatchFormat:
		fmt.Fprintf(g.out, `set GOROOT=%v`+"\n", goroot)
		fmt.Fprintf(g.out, `set PATH=%s\bin;%s`+"\n", goroot, os.Getenv("PATH"))
		if strings.HasPrefix(version, "1.5") {
			fmt.Fprintln(g.out, `set GO15VENDOREXPERIMENT=1`)
		}
	case PowershellFormat:
		fmt.Fprintf(g.out, `$env:GOROOT = "%v"`+"\n", goroot)
		fmt.Fprintf(g.out, `$env:PATH = "$env:GOROOT\bin;$env:PATH"`+"\n")
		if strings.HasPrefix(version, "1.5") {
			fmt.Fprintln(g.out, `$env:GO15VENDOREXPERIMENT=1`)
		}
	default:
		return errors.Errorf("invalid format option '%v'", g.Format)
	}

	return nil
}

func getProjectGoVersion() (string, error) {
	ver, err := parseTravisYml(".travis.yml")
	if err != nil {
		return "", fmt.Errorf("failed to detect the project's golang version: %v", err)
	}

	return ver, nil
}

func parseTravisYml(name string) (string, error) {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	var re = regexp.MustCompile(`(?mi)^go:\s*\r?\n\s*-\s+(\S+)\s*$`)
	matches := re.FindAllStringSubmatch(string(file), 1)
	if len(matches) == 0 {
		return "", fmt.Errorf("go not found in %v", name)
	}

	goVersion := matches[0][1]
	return goVersion, nil
}

func defaultFormat() string {
	if runtime.GOOS == "windows" {
		return BatchFormat
	}
	return BashFormat
}

func main() {
	app := kingpin.New("gvm", usage)
	debug := app.Flag("debug", "Enable debug logging to stderr.").Short('d').Bool()

	g := &GVM{out: os.Stdout}
	app.Flag("project-go", "Use the project's Go version.").BoolVar(&g.UseProjectGo)
	app.Flag("format", "Format to use for the shell commands. Options: bash, batch, powershell").Short('f').Default(defaultFormat()).EnumVar(&g.Format, BashFormat, BatchFormat, PowershellFormat)
	app.Arg("version", "Go version to install (e.g. 1.9.2).").StringVar(&g.Version)
	app.Action(g.Run)

	app.Version(version)
	app.HelpFlag.Short('h')
	app.DefaultEnvars()
	app.UsageTemplate(kingpin.SeparateOptionalFlagsUsageTemplate)

	// Enable debug.
	app.PreAction(func(ctx *kingpin.ParseContext) error {
		logrus.SetLevel(logrus.DebugLevel)
		if *debug {
			logrus.SetOutput(os.Stderr)
		} else {
			logrus.SetOutput(ioutil.Discard)
		}
		return nil
	})

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Errorf("%v", err)
		os.Exit(1)
	}
}
