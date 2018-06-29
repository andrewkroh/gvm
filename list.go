package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/andrewkroh/gvm/golang"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func listCommand(cmd *kingpin.CmdClause) func() error {
	return func() error {
		dir, err := golang.VersionsDir()
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}

		commonPrefix := "go"
		versionSuffix := fmt.Sprintf(".%v.%v", runtime.GOOS, golang.GOARCH())

		for _, fi := range files {
			name := fi.Name()
			if strings.HasSuffix(name, versionSuffix) {
				name = name[:len(name)-len(versionSuffix)]
			}
			if strings.HasPrefix(name, commonPrefix) {
				name = name[len(commonPrefix):]
			}

			fmt.Println(name)
		}

		return nil
	}
}
