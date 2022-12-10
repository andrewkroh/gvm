package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/andrewkroh/gvm"
)

const usage = `gvm is a Go version manager. gvm installs a Go version and prints
the commands to configure your environment to use it. gvm can only install
binary versions of Go from https://storage.googleapis.com/golang. Below are
examples for common shells.

  bash:
    eval "$(gvm 1.18.5)"

  batch (windows cmd.exe):
    FOR /f "tokens=*" %i IN ('"gvm.exe" 1.18.5') DO %i

  powershell:
    gvm --format=powershell 1.18.5 | Invoke-Expression

gvm flags can be set via environment variables by setting GVM_<flag>. For
example --http-timeout can be set via GVM_HTTP_TIMEOUT=10m.
`

var log = logrus.WithField("package", "main")

// Build info.
var (
	version string
	commit  string
)

func init() {
	if version == "" && commit == "" {
		// Fall back to Go module data when not built with goreleaser.
		if info, ok := debug.ReadBuildInfo(); ok {
			if info.Main.Sum == "" {
				info.Main.Sum = "unknown"
			}
			version = info.Main.Version
			commit = info.Main.Sum
		}
	}
}

type commandFactory func(*kingpin.CmdClause) func(*gvm.Manager) error

var (
	commands = map[string]func(*gvm.Manager) error{}
)

func registerCommand(app *kingpin.Application, factory commandFactory, name, doc string) *kingpin.CmdClause {
	cmd := app.Command(name, doc)
	act := factory(cmd)
	if act != nil {
		commands[name] = act
	}
	return cmd
}

func main() {
	app := kingpin.New("gvm", usage)
	debug := app.Flag("debug", "Enable debug logging to stderr.").Short('d').Bool()

	// manager := &gvm.Manager{}
	manager, err := gvm.NewDefaultWithProxy()
	if err != nil {
		panic(fmt.Sprintf("failed to build default gvm, err=%s", err.Error()))
	}

	app.Flag("os", "Go binaries target os.").StringVar(&manager.GOOS)
	app.Flag("arch", "Go binaries target architecture.").StringVar(&manager.GOARCH)
	app.Flag("home", "GVM home directory.").StringVar(&manager.WorkDir)
	app.Flag("url", "Go binaries repository base URL.").StringVar(&manager.GoStorageHome)
	app.Flag("repository", "Go upstream git repository.").StringVar(&manager.GoSourceURL)
	app.Flag("http-timeout", "Timeout for HTTP requests.").Default("3m").DurationVar(&manager.HTTPTimeout)

	// register subcommand
	registerCommand(app, useCommand, "use", "prepare go version and print environment variables").Default()
	registerCommand(app, initCommand, "init", "init .gvm and update source cache")
	registerCommand(app, installCommand, "install", "install go version if not already installed")
	registerCommand(app, availCommand, "available", "list all installable go versions")
	registerCommand(app, listCommand, "list", "list installed versions")
	registerCommand(app, removeCommand, "remove", "remove a go version")
	registerCommand(app, purgeCommand, "purge", "remove all but the newest go version")
	registerCommand(app, configCommand, "config", "setting gvm configs")

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

	selCommand, err := app.Parse(os.Args[1:])
	if err != nil {
		app.Errorf("%v", err)
		os.Exit(1)
	}

	logrus.Debug("GVM version: ", version)
	logrus.Debug("GVM commit: ", commit)
	logrus.Debug("GVM arch: ", runtime.GOARCH)

	action, exists := commands[selCommand]
	if !exists {
		app.Errorf("unknown command: %v", selCommand)
		app.Usage(os.Args[1:])
		os.Exit(2)
	}

	if err := manager.Init(); err != nil {
		app.Errorf("%v", err)
		os.Exit(1)
	}

	if err := action(manager); err != nil {
		app.Errorf("%v", err)
		os.Exit(1)
	}
}
