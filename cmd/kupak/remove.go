package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func remove(c *cli.Context) {
	group := c.Args().First()
	err := pakManager.Remove(c.GlobalString("namespace"), group)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(-1)
	}
}
