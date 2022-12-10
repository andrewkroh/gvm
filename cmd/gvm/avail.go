package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/andrewkroh/gvm"
	"github.com/sirupsen/logrus"
)

func availCommand(_ *kingpin.CmdClause) func(*gvm.Manager) error {
	return func(manager *gvm.Manager) error {
		list, err := manager.Available()
		if err != nil {
			logrus.Errorf("failed to list available source, err=%s", err.Error())
			return err
		}

		for _, v := range list {
			fmt.Println(v)
		}
		return nil
	}
}
