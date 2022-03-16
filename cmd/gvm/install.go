package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/andrewkroh/gvm"
)

func installCommand(cmd *kingpin.CmdClause) func(*gvm.Manager) error {
	var version string
	var build bool
	cmd.Flag("build", "Build go version from source").Short('b').BoolVar(&build)
	cmd.Arg("version", "Go version to install (e.g. 1.10.3).").StringVar(&version)

	return func(manager *gvm.Manager) error {
		if version == "" {
			return fmt.Errorf("no version specified")
		}
		ver, err := gvm.ParseVersion(version)
		if err != nil {
			return err
		}

		has, err := manager.HasVersion(ver)
		if err != nil {
			return err
		}
		if has {
			fmt.Printf("Version %v already installed\n", version)
			return nil
		}

		var dir string
		if build {
			fmt.Printf("Building go-%v. Please wait...\n", version)
			dir, err = manager.Build(ver)
		} else {
			fmt.Printf("Installing go-%v. Please wait...\n", version)
			dir, err = manager.Install(ver)
		}
		if err != nil {
			fmt.Println("Installation failed with:\n", err)
			return err
		}

		fmt.Printf("Sucessfully installed go-%v to %v\n", version, dir)
		return nil
	}
}
