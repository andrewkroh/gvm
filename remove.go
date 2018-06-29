package main

import (
	"fmt"
	"os"

	"github.com/andrewkroh/gvm/golang"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func removeCommand(cmd *kingpin.CmdClause) func() error {
	var versions []string
	cmd.Arg("versions", "Go versions to remove").StringsVar(&versions)
	return func() error {
		for _, version := range versions {
			dir, err := golang.VersionDir(version)
			if err != nil {
				return err
			}

			fi, err := os.Stat(dir)
			if os.IsNotExist(err) {
				fmt.Printf("Version %v not installed\n", version)
				continue
			}
			if err != nil {
				return err
			}

			if !fi.IsDir() {
				fmt.Println("Path %v is no directory", dir)
				continue
			}

			if err := os.RemoveAll(dir); err != nil {
				return err
			}

		}
		return nil
	}
}
