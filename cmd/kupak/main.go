package main

import (
	"fmt"
	"kupak"
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kupak"
	app.Usage = "Kubernetes Package Manager"
	app.Version = "0.1"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list packages of specified repo",
			Action:  list,
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install the specified pak (full url or relative to --repo)",
			Action:  install,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo, r",
			Value:  "src/kupak/example/index.yaml",
			Usage:  "specify repo url",
			EnvVar: "KUPAK_REPO",
		},
	}
	app.Run(os.Args)
}

func install(c *cli.Context) {
	manager, err := kupak.NewManager()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	pak, err := kupak.PakFromURL("src/kupak/example/paks/redis-1.0/pak.yaml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	err = manager.Install(pak, "default", map[string]interface{}{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func list(c *cli.Context) {
	repo, err := kupak.RepoFromURL(c.GlobalString("repo"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	fmt.Println("List of available Paks:")
	for i := range repo.Paks {
		fmt.Println("- Name:", repo.Paks[i].Name)
		fmt.Println("  Version:", repo.Paks[i].Version)
		fmt.Println("  Tags:", "["+strings.Join(repo.Paks[i].Tags, ", ")+"]")
		fmt.Println()
	}
}
