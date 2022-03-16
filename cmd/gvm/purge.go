package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/andrewkroh/gvm"
)

func purgeCommand(_ *kingpin.CmdClause) func(*gvm.Manager) error {
	return func(manager *gvm.Manager) error {
		versions, err := manager.Installed()
		if err != nil {
			return err
		}

		// find installed highest stable release
		stable := -1
		for i := len(versions) - 1; i != -1; i-- {
			if versions[i].Stable() {
				stable = i
				break
			}
		}

		if stable <= 0 {
			fmt.Println("No versions to remove")
		} else {
			removeVersions(manager, versions[:stable])

			// unstable versions > last stable version
			versions = versions[stable+1:]
		}

		if len(versions) <= 1 {
			return nil
		}

		// remove all but highest unstable version
		removeVersions(manager, versions[:len(versions)-1])

		return nil
	}
}
