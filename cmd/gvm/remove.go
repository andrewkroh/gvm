package main

import (
	"fmt"

	"github.com/andrewkroh/gvm"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func removeCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	var versions []string
	cmd.Arg("versions", "Go versions to remove").StringsVar(&versions)

	return func(manager *gvm.Manager) error {
		if len(versions) == 0 {
			return fmt.Errorf("no versions specified")
		}

		for _, version := range versions {
			ver, err := gvm.ParseVersion(version)
			if err != nil {
				fmt.Printf("Invalid version '%v': %v\n", version, err)
				continue
			}

			fmt.Printf("Removing version %v...\n", version)
			if err := manager.Remove(ver); err != nil {
				fmt.Printf("Can not remove verions %v:\n%v\n", version, err)
			} else {
				fmt.Println("Removed version", version)
			}
		}
		return nil
	}
}
