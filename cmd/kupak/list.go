package main

import (
	"fmt"
	"os"

	"git.cafebazaar.ir/alaee/kupak/logging"
	"github.com/codegangsta/cli"
)

func list(c *cli.Context) {
	paks, err := pakManager.List(c.GlobalString("namespace"))
	if err != nil {
		logging.Error(fmt.Sprint(err))
		os.Exit(-1)
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
}
