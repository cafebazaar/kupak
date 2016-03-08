package main

import (
	"fmt"
	"os"

	"git.cafebazaar.ir/alaee/kupak/pkg/kubectl"
	"git.cafebazaar.ir/alaee/kupak/pkg/manager"
	"github.com/codegangsta/cli"
)

var pakManager *manager.Manager

func main() {
	kc, err := kubectl.NewKubectlRunner()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	pakManager, err = manager.NewManager(kc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	app := cli.NewApp()
	app.Name = "kupak"
	app.Usage = "Kubernetes Package Manager"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:    "paks",
			Aliases: []string{"p"},
			Usage:   "list all available paks of specified repo",
			Action:  paks,
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install the specified pak (full url or a plain name that exists in specified repo)",
			Action:  install,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "interactive, i",
					Usage: "interactive installation",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all installed packages",
			Action:  list,
		},
		{
			Name:    "spec",
			Aliases: []string{"s"},
			Usage:   "print specification of a pak",
			Action:  spec,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo, r",
			Value:  "src/kupak/example/index.yaml",
			Usage:  "specify repo url",
			EnvVar: "KUPAK_REPO",
		},
		cli.StringFlag{
			Name:   "namespace",
			Value:  "default",
			Usage:  "namespace",
			EnvVar: "KUPAK_NAMESPACE",
		},
	}
	app.Run(os.Args)
}
