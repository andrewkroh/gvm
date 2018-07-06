package main

import (
	"github.com/andrewkroh/gvm"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func initCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	return func(manager *gvm.Manager) error {
		return manager.UpdateCache()
	}
}
