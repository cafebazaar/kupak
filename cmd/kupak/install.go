package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"git.cafebazaar.ir/alaee/kupak/pkg/pak"
	"git.cafebazaar.ir/alaee/kupak/util"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
)

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
		values = readValuesInteractively(p)
	} else {
		values = readValuesFromFile(p, valuesFile)
	}

	_, err = pakManager.Install(p, c.GlobalString("namespace"), values)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func readValuesFromFile(p *pak.Pak, path string) map[string]interface{} {
	values := make(map[string]interface{})
	var valuesData []byte
	var err error
	if path == "" {
		valuesData, err = ioutil.ReadAll(os.Stdin)
	} else {
		valuesData, err = ioutil.ReadFile(path)
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
	return values
}

func readValuesInteractively(p *pak.Pak) map[string]interface{} {
	values := make(map[string]interface{})

	// ask for group
	groupValue, err := scanValue("Group Name (return for random): ", "string", false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
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
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		if value != nil {
			values[p.Properties[i].Name] = value
		}
	}
	return values
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
