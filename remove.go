package main

import (
	"fmt"
	"os"

	"github.com/cafebazaar/kupak/logging"
	"github.com/codegangsta/cli"
)

func remove(c *cli.Context) {
	group := c.Args().First()
	err := pakManager.Remove(c.GlobalString("namespace"), group)
	if err != nil {
		logging.Error(fmt.Sprint(err))
		os.Exit(-1)
	}
}
