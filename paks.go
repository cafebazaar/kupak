package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/cafebazaar/kupak/pkg/pak"
	"github.com/cafebazaar/kupak/pkg/util"
)

func paks(c *cli.Context) error {
	repoAddresses := make([]string, 0)
	if c.GlobalIsSet("repo") {
		repoAddresses = append(repoAddresses, c.GlobalString("repo"))
	} else {
		reposFileCreateIfNotExist()
		for _, repoEntry := range reposFileNonEmptyEntries() {
			trimmed := strings.TrimSpace(repoEntry)
			if len(trimmed) > 0 {
				repoAddresses = append(repoAddresses, trimmed)
			}
		}
		if len(repoAddresses) == 0 {
			repoAddresses = append(repoAddresses, "github.com/cafebazaar/paks")
		}
	}

	for _, repoAddress := range repoAddresses {
		repo, err := pak.RepoFromURL(repoAddress)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("can't fetch paks list: %v", err.Error()), -1)
		}
		fmt.Println("- Repository:", repoAddress)
		for i := range repo.Paks {
			fmt.Println("  - Name:", repo.Paks[i].Name)
			fmt.Println("    Version:", repo.Paks[i].Version)
			fmt.Println("    URL:", repo.Paks[i].URL)
			if len(repo.Paks[i].Tags) > 0 {
				fmt.Println("    Tags:", "["+strings.Join(repo.Paks[i].Tags, ", ")+"]")
			}
			fmt.Println("   ", strings.Trim(repo.Paks[i].Description, "\n"))
		}

	}
	return nil
}

func spec(c *cli.Context) error {
	pakURL := c.Args().First()
	if pakURL == "" {
		return cli.NewExitError("please specify the pak name", -1)
	}
	if strings.Index(pakURL, "/") == -1 &&
		!strings.HasSuffix(pakURL, ".json") &&
		!strings.HasSuffix(pakURL, ".yaml") {

		nameOfPakToCheck := pakURL

		repoAddresses := make([]string, 0)
		if c.GlobalIsSet("repo") {
			repoAddresses = append(repoAddresses, c.GlobalString("repo"))
		} else {
			reposFileCreateIfNotExist()
			for _, repoEntry := range reposFileNonEmptyEntries() {
				trimmed := strings.TrimSpace(repoEntry)
				if len(trimmed) > 0 {
					repoAddresses = append(repoAddresses, trimmed)
				}
			}
			if len(repoAddresses) == 0 {
				repoAddresses = append(repoAddresses, "github.com/cafebazaar/paks")
			}
		}

		for _, repoAddr := range repoAddresses {
			if len(repoAddr) > 0 {
				// TODO: change JoinURL
				repoPaks, err := pak.RepoFromURL(repoAddr)
				if err != nil {
					return cli.NewExitError(fmt.Sprintf("can't fetch the repo: %v", err.Error()), -1)
				}
				for _, pak := range repoPaks.Paks {
					if pak.Name == nameOfPakToCheck {
						pakURL = pak.URL
						if util.Relative(pakURL) {
							pakURL = util.JoinURL(repoAddr, pakURL)
						}
						specOfPakWithAbsoluteURLFromRepo(pakURL, repoAddr)
					}
				}
			}
		}
	} else {
		specOfPakWithAbsoluteURL(pakURL)
	}
	return nil
}

func specOfPakWithAbsoluteURLFromRepo(pakURL, repo string) error {
	p, err := pak.FromURL(pakURL)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("can't fetch the specified pak: %v", err.Error()), -1)
	}
	fmt.Println("-Repository:", repo)
	fmt.Println("  Name:", p.Name)
	fmt.Println("  Version:", p.Version)
	if len(p.Tags) > 0 {
		fmt.Println("  Tags:", "["+strings.Join(p.Tags, ", ")+"]")
	}
	fmt.Println(" ", strings.Trim(p.Description, "\n"))

	fmt.Println("\n  Properties:")
	for i := range p.Properties {
		property := p.Properties[i]
		fmt.Println("   - Name:", property.Name)
		fmt.Println("     Description:", strings.Trim(property.Description, "\n"))
		fmt.Println("     Type:", property.Type)
		if property.Default != nil {
			fmt.Println("     Default:", property.Default)
		}

	}
	return nil

}

func specOfPakWithAbsoluteURL(pakURL string) error {
	p, err := pak.FromURL(pakURL)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("can't fetch the specified pak: %v", err.Error()), -1)
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
	return nil

}
