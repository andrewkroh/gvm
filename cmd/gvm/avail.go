package main

import (
	"fmt"

	"github.com/andrewkroh/gvm"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func availCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	return func(manager *gvm.Manager) error {
		list, hasBin, err := manager.Available()
		if err != nil {
			return err
		}

		for i, v := range list {
			if hasBin[i] {
				fmt.Printf("%v\t(source, binary)\n", v)
			} else {
				fmt.Printf("%v\t(source)\n", v)
			}
		}
		return nil
	}
}
