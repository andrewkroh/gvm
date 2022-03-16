package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/andrewkroh/gvm"
)

func listCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	return func(manager *gvm.Manager) error {
		versions, err := manager.Installed()
		if err != nil {
			return err
		}

		for _, version := range versions {
			fmt.Println(version)
		}
		return nil
	}
}
