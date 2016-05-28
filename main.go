package main

import (
	"fmt"
	"os"

	"github.com/cafebazaar/kupak/logging"
	"github.com/cafebazaar/kupak/pkg/kubectl"
	"github.com/cafebazaar/kupak/pkg/manager"
	"github.com/cafebazaar/kupak/pkg/version"
	"github.com/codegangsta/cli"
)

var pakManager *manager.Manager

func main() {
	kc, err := kubectl.NewKubectlRunner()
	if err != nil {
		logging.Error(fmt.Sprintln(err))
		os.Exit(-1)
	}
	pakManager, err = manager.NewManager(kc)
	if err != nil {
		logging.Error(fmt.Sprintln(err))
		os.Exit(-1)
	}

	app := cli.NewApp()
	app.Name = "kupak"
	app.Usage = "Kubernetes Package Manager"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:   "version",
			Usage:  "print the current version of Kupak",
			Action: printVersion,
		},
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
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "remove the pak specified by group name",
			Action:  remove,
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
		{
			Name:    "repos",
			Aliases: []string{},
			Usage:   "Managing kupak repositories",
			Subcommands: []cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage: `Adds a pak repo to kupak repos file
   kupak repos add github.com/foo/repo [description]`,
					Action: reposAdd,
				},
				{
					Name:    "list",
					Aliases: []string{"l", "ls"},
					Usage:   `Lists the repos in kupak repos file`,
					Action:  reposList,
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo, r",
			Value:  "github.com/cafebazaar/paks",
			Usage:  "specify repo url",
			EnvVar: "KUPAK_REPO",
		},
		cli.StringFlag{
			Name:   "namespace",
			Value:  "default",
			Usage:  "namespace",
			EnvVar: "KUPAK_NAMESPACE",
		},
		cli.BoolFlag{
			Name:        "verbose, V",
			Usage:       "be verbose",
			EnvVar:      "KUPAK_VERBOSE",
			Destination: &logging.Verbose,
		},
	}
	app.Run(os.Args)
}

func printVersion(c *cli.Context) error {
	fmt.Printf("Kupak %v\n", version.KupakVersion)
	return nil
}
