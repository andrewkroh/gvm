package main

import (
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const usage = `gvm is a Go version manager. gvm installs a Go version and prints
the commands to configure your environment to use it. gvm can only install
binary versions of Go from https://golang.org/dl/. Below are examples for
common shells.

  bash:
    eval "$(gvm 1.10)"

  batch (windows cmd.exe):
    FOR /f "tokens=*" %i IN ('"gvm.exe" 1.10') DO %i

  powershell:
    gvm --format=powershell 1.10 | Invoke-Expression
`

var (
	version = "SNAPSHOT"

	log = logrus.WithField("package", "main")
)

type commandFactory func(*kingpin.CmdClause) func() error

func main() {

	app := kingpin.New("gvm", usage)
	debug := app.Flag("debug", "Enable debug logging to stderr.").Short('d').Bool()

	addCommand(app, useCommand, "use", "prepare go version and print environment variables").
		Default()
	addCommand(app, installCommand, "install", "install go version if not already installed")
	// addCommand(app, availCommand, "available", "list all installable go versions")
	addCommand(app, listCommand, "list", "list installed versions")
	addCommand(app, removeCommand, "remove", "remove a go version")
	// addCommand(app, purgeCommand, "purge", "remove all but the newest go version")

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

func addCommand(app *kingpin.Application, factory commandFactory, name, doc string) *kingpin.CmdClause {
	cmd := app.Command(name, doc)
	act := factory(cmd)
	if act != nil {
		cmd.Action(func(_ *kingpin.ParseContext) error {
			return act()
		})
	}
	return cmd
}
