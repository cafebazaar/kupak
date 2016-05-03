package main

import (
	"golang.org/x/crypto/ssh/terminal"

	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/cafebazaar/kupak/pkg/pak"
	"github.com/cafebazaar/kupak/pkg/util"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

func install(c *cli.Context) error {
	pakURL := c.Args().First()
	valuesFile := c.Args().Get(1)
	if pakURL == "" {
		return cli.NewExitError("please specify the pak", -1)
	}

	if strings.Index(pakURL, "/") == -1 &&
		!strings.HasSuffix(pakURL, ".json") &&
		!strings.HasSuffix(pakURL, ".yaml") {

		nameOfPakToInstall := pakURL
		repoAddr := c.GlobalString("repo")

		if len(repoAddr) > 0 {
			// TODO: change JoinURL
			repoPaks, err := pak.RepoFromURL(repoAddr)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("can't fetch the repo: %v", err.Error()), -1)
			}
			for _, pak := range repoPaks.Paks {
				if pak.Name == nameOfPakToInstall {
					pakURL = pak.URL
					if util.Relative(pakURL) {
						pakURL = util.JoinURL(repoAddr, pakURL)
					}
				}
			}
		}
	}

	p, err := pak.FromURL(pakURL)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("can't fetch the pak: %v", err.Error()), -1)
	}

	values := make(map[string]interface{})
	if c.Bool("interactive") || terminal.IsTerminal(int(os.Stdin.Fd())) {
		values, err = readValuesInteractively(p)
	} else {
		values, err = readValuesFromFile(p, valuesFile)
	}
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("error in reading value: %v", err.Error()), -1)
	}

	_, err = pakManager.Install(p, c.GlobalString("namespace"), values)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("installation error: %v", err.Error()), -1)
	}
	return nil
}

func readValuesFromFile(p *pak.Pak, path string) (map[string]interface{}, error) {
	values := make(map[string]interface{})
	var valuesData []byte
	var err error
	if path == "" {
		valuesData, err = ioutil.ReadAll(os.Stdin)
	} else {
		valuesData, err = ioutil.ReadFile(path)
	}
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(valuesData, &values)
	return values, err
}

func readValuesInteractively(p *pak.Pak) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	// ask for group
	groupValue, err := scanValue("Group Name (return for random): ", "string", false)
	if err != nil {
		return nil, err
	}
	if groupValue != nil {
		values["group"] = groupValue.(string)
	}

	// ask for all properties
	for i := range p.Properties {
		var prompt string
		if p.Properties[i].Default == nil {
			prompt = fmt.Sprintf("Field \"%s\" [type: %s]? ", p.Properties[i].Name, p.Properties[i].Type)
		} else {
			prompt = fmt.Sprintf("Field \"%s\" [type: %s, default: %v] (return for default)? ", p.Properties[i].Name, p.Properties[i].Type, p.Properties[i].Default)
		}
		value, err := scanValue(prompt, p.Properties[i].Type, p.Properties[i].Default == nil)
		if err != nil {
			return nil, err
		}
		if value != nil {
			values[p.Properties[i].Name] = value
		}
	}
	return values, nil
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
		case "int":
			i, err := strconv.Atoi(string(value))
			if err != nil {
				fmt.Println("given value is not an int, try again")
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
		case "string":
			fallthrough
		default:
			return string(value), nil
		}
	}
}
