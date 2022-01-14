package main

import (
	"io/ioutil"
	"os"

	"github.com/andrewkroh/gvm"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const usage = `gvm is a Go version manager. gvm installs a Go version and prints
the commands to configure your environment to use it. gvm can only install
binary versions of Go from https://storage.googleapis.com/golang. Below are
examples for common shells.

  bash:
    eval "$(gvm 1.17.7)"

  batch (windows cmd.exe):
    FOR /f "tokens=*" %i IN ('"gvm.exe" 1.17.7') DO %i

  powershell:
    gvm --format=powershell 1.17.7 | Invoke-Expression

gvm flags can be set via environment variables by setting GVM_<flag>. For
example --http-timeout can be set via GVM_HTTP_TIMEOUT=10m.
`

var (
	version = "SNAPSHOT"

	log = logrus.WithField("package", "main")
)

type commandFactory func(*kingpin.CmdClause) func(*gvm.Manager) error

func main() {
	app := kingpin.New("gvm", usage)
	debug := app.Flag("debug", "Enable debug logging to stderr.").Short('d').Bool()

	manager := &gvm.Manager{}
	commands := map[string]func(*gvm.Manager) error{}
	command := func(factory commandFactory, name, doc string) *kingpin.CmdClause {
		cmd := app.Command(name, doc)
		act := factory(cmd)
		if act != nil {
			commands[name] = act
		}
		return cmd
	}

	app.Flag("os", "Go binaries target os.").StringVar(&manager.GOOS)
	app.Flag("arch", "Go binaries target architecture.").StringVar(&manager.GOOS)
	app.Flag("home", "GVM home directory.").StringVar(&manager.Home)
	app.Flag("url", "Go binaries repository base URL.").StringVar(&manager.GoStorageHome)
	app.Flag("repository", "Go upstream git repository.").StringVar(&manager.GoSourceURL)
	app.Flag("http-timeout", "Timeout for HTTP requests.").Default("3m").DurationVar(&manager.HTTPTimeout)

	command(useCommand, "use", "prepare go version and print environment variables").
		Default()
	command(initCommand, "init", "init .gvm and update source cache")
	command(installCommand, "install", "install go version if not already installed")
	command(availCommand, "available", "list all installable go versions")
	command(listCommand, "list", "list installed versions")
	command(removeCommand, "remove", "remove a go version")
	command(purgeCommand, "purge", "remove all but the newest go version")

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
