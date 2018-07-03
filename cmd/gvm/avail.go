package main

import (
	"fmt"

	"github.com/andrewkroh/gvm"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func availCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	return func(manager *gvm.Manager) error {
		list, err := manager.Available()
		if err != nil {
			return err
		}

		for _, v := range list {
			fmt.Println(v)
		}
		return nil
	}
}
