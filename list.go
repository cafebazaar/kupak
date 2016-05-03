package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func list(c *cli.Context) error {
	paks, err := pakManager.List(c.GlobalString("namespace"))
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error in fetching installed paks: %v", err.Error()), -1)
	}
	for i := range paks {
		fmt.Printf("Pak URL:  %s\n", paks[i].PakURL)
		fmt.Printf("Group: %s\n", paks[i].Group)
		fmt.Printf("Objects:\n")
		for j := range paks[i].Objects {
			obj := paks[i].Objects[j]
			md, _ := obj.Metadata()
			fmt.Printf("\t(%s) %s\n", md.Kind, md.Name)
			if md.Kind == "Pod" {
				status, _ := obj.Status()
				fmt.Printf("\t  State:     %s\n", status.Phase)
				fmt.Printf("\t  Pod IP:    %s\n", status.PodIP)
				if status.Reason != "" {
					fmt.Printf("\t  Reason:  %s\n", status.Reason)
				}
				if status.Message != "" {
					fmt.Printf("\t  Message: %s\n", status.Message)
				}
			}
		}
		fmt.Println()
	}
	return nil
}
