package main

import (
	"fmt"
	"os"
	"strings"

	"git.cafebazaar.ir/alaee/kupak/pkg/pak"
	"github.com/codegangsta/cli"
)

func paks(c *cli.Context) {
	repo, err := pak.RepoFromURL(c.GlobalString("repo"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	for i := range repo.Paks {
		fmt.Println("- Name:", repo.Paks[i].Name)
		fmt.Println("  Version:", repo.Paks[i].Version)
		fmt.Println("  URL:", repo.Paks[i].URL)
		if len(repo.Paks[i].Tags) > 0 {
			fmt.Println("  Tags:", "["+strings.Join(repo.Paks[i].Tags, ", ")+"]")
		}
		fmt.Println(" ", strings.Trim(repo.Paks[i].Description, "\n"))
	}
}

func spec(c *cli.Context) {
	pakURL := c.Args().First()
	if pakURL == "" {
		fmt.Fprintln(os.Stderr, "please specify the pak")
		os.Exit(-1)
	}
	p, err := pak.FromURL(pakURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(-1)
	}
	fmt.Println("Name:", p.Name)
	fmt.Println("Version:", p.Version)
	if len(p.Tags) > 0 {
		fmt.Println("Tags:", "["+strings.Join(p.Tags, ", ")+"]")
	}
	fmt.Println(strings.Trim(p.Description, "\n"))

	fmt.Println("\nProperties:")
	for i := range p.Properties {
		property := p.Properties[i]
		fmt.Println(" - Name:", property.Name)
		fmt.Println("   Description:", strings.Trim(property.Description, "\n"))
		fmt.Println("   Type:", property.Type)
		if property.Default != nil {
			fmt.Println("   Default:", property.Default)
		}

	}
}
