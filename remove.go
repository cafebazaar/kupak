package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func remove(c *cli.Context) error {
	group := c.Args().First()
	err := pakManager.Remove(c.GlobalString("namespace"), group)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error while removing the pak: %v", err.Error()), -1)
	}
	return nil
}
