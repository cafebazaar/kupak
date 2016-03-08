package main

import "github.com/codegangsta/cli"

func remove(c *cli.Context) {
	group := c.Args().First()
	pakManager.Remove(c.GlobalString("namespace"), group)
}
