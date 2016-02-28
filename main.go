package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"git.cafebazaar.ir/alaee/kupak/kubectl"
	"git.cafebazaar.ir/alaee/kupak/manager"
	"git.cafebazaar.ir/alaee/kupak/pak"
	"git.cafebazaar.ir/alaee/kupak/util"
	"github.com/codegangsta/cli"
	"github.com/ghodss/yaml"
)

var pakManager *manager.Manager

func main() {
	kc, err := kubectl.NewKubectlRunner()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	pakManager, err = manager.NewManager(kc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	app := cli.NewApp()
	app.Name = "kupak"
	app.Usage = "Kubernetes Package Manager"
	app.Version = "0.1"
	app.Commands = []cli.Command{
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
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "repo, r",
			Value:  "src/kupak/example/index.yaml",
			Usage:  "specify repo url",
			EnvVar: "KUPAK_REPO",
		},
		cli.StringFlag{
			Name:   "namespace",
			Value:  "default",
			Usage:  "namespace",
			EnvVar: "KUPAK_NAMESPACE",
		},
	}
	app.Run(os.Args)
}

func install(c *cli.Context) {
	pakURL := c.Args().First()
	valuesFile := c.Args().Get(1)
	if pakURL == "" {
		fmt.Fprintln(os.Stderr, "please specify the pak")
		os.Exit(-1)
	}

	p, err := pak.FromURL(pakURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	values := make(map[string]interface{})

	if c.Bool("interactive") {
		// interactive
		for i := range p.Properties {
			var prompt string
			if p.Properties[i].Default == nil {
				prompt = fmt.Sprintf("Field \"%s\" [type: %s]? ", p.Properties[i].Name, p.Properties[i].Type)
			} else {
				prompt = fmt.Sprintf("Field \"%s\" [type: %s, default: %v] (return for default)? ", p.Properties[i].Name, p.Properties[i].Type, p.Properties[i].Default)
			}
			value, err := scanValue(prompt, p.Properties[i].Type, p.Properties[i].Default == nil)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(-1)
			}
			if value != nil {
				values[p.Properties[i].Name] = value
			}
		}
	} else {
		// read value file
		var valuesData []byte
		if valuesFile == "" {
			valuesData, err = ioutil.ReadAll(os.Stdin)
		} else {
			valuesData, err = ioutil.ReadFile(valuesFile)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		err = yaml.Unmarshal(valuesData, &values)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	}

	_, err = pakManager.Install(p, c.GlobalString("namespace"), values)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

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

func list(c *cli.Context) {
	paks, err := pakManager.List(c.GlobalString("namespace"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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

func scanValue(prompt string, valueType string, required bool) (interface{}, error) {
	bio := bufio.NewReader(os.Stdin)
	var value []byte
	var err error
	for {
		fmt.Printf(prompt)
		value, _, err = bio.ReadLine()
		if err != nil {
			return nil, err
		}
		if len(value) == 0 && !required {
			return nil, nil
		} else if len(value) == 0 && required {
			continue
		}

		switch valueType {
		case "string":
			return string(value), nil
		case "int":
			i, err := strconv.Atoi(string(value))
			if err != nil {
				fmt.Println("Bad value, try again")
				continue
			}
			return i, nil
		case "bool":
			b, err := util.StringToBool(string(value))
			if err != nil {
				fmt.Println(err)
				continue
			}
			return b, nil
		default:
			return value, nil
		}
	}
}
