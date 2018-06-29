package main

import (
	"fmt"

	"github.com/andrewkroh/gvm/golang"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type gvmInstall struct {
	version      string
	useProjectGo bool
}

func installCommand(cmd *kingpin.CmdClause) func() error {
	g := &gvmInstall{}
	cmd.Flag("project-go", "Use the project's Go version.").BoolVar(&g.useProjectGo)
	cmd.Arg("version", "Go version to install (e.g. 1.10).").StringVar(&g.version)
	return g.run
}

func (g *gvmInstall) run() error {
	version, err := useVersion(g.version, g.useProjectGo)
	if err != nil {
		return err
	}

	has, err := golang.HasVersion(version)
	if err != nil {
		return err
	}

	if has {
		fmt.Printf("Version %v already installed\n", version)
		return nil
	}

	fmt.Printf("Installing go-%v. Please wait...\n", version)
	dir, err := golang.SetupGolang(version)
	if err != nil {
		fmt.Println("Installation failed with:\n", err)
		return err
	}

	fmt.Printf("Sucessfully installed go-%v to %v", version, dir)
	return nil
}
