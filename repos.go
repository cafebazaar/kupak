package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
)

const (
	REPOS_FILE_NAME = "repos"
)

func reposFileEntries() []string {
	f, err := ioutil.ReadFile(reposFilePath())
	if err != nil {
		return make([]string, 0)
	}
	return strings.Split(string(f), "\n")
}

func reposFileNonEmptyEntries() []string {
	ret := make([]string, 0)
	for _, line := range reposFileEntries() {
		if len(line) > 0 && strings.TrimSpace(line)[0] != '#' {
			ret = append(ret, line)
		}
	}
	return ret
}

func reposFileDir() string {
	homeDir := os.Getenv("HOME")
	if len(homeDir) == 0 {
		homeDir = "/"
	}
	return path.Join(homeDir, ".kupak")
}

func reposFilePath() string {
	return path.Join(reposFileDir(), REPOS_FILE_NAME)
}

func reposFileExists() bool {
	_, err := os.Stat(reposFilePath())
	return err == nil
}

func reposFileCreateIfNotExist() error {
	if reposFileExists() {
		return nil
	}
	if err := os.MkdirAll(reposFileDir(), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(reposFilePath(), []byte("#official repo\ngithub.com/cafebazaar/paks\n"), 0644); err != nil {
		return err
	}
	return nil
}

func reposAdd(c *cli.Context) error {
	repo := c.Args().Get(0)
	if repo == "" {
		return cli.NewExitError("Please specify the repository", -1)
	}

	reposFileCreateIfNotExist()

	description := strings.Join(c.Args().Tail(), " ")

	toWrite := "\n"
	if description != "" {
		toWrite += "#" + description + "\n"
	}
	toWrite += repo
	toWrite += "\n"

	reposFile, err := os.OpenFile(reposFilePath(), os.O_APPEND|os.O_WRONLY, 0600)
	defer reposFile.Close()
	if err != nil {
		return cli.NewExitError("Cannot open repos File for appending", -1)
	}

	_, err = reposFile.WriteString(toWrite)
	if err != nil {
		println(err.Error())
		return cli.NewExitError("Cannot write to repos file", -1)
	}

	return nil
}

func reposList(c *cli.Context) error {
	reposFileCreateIfNotExist()
	println("Repositories:")
	for _, repo := range reposFileNonEmptyEntries() {
		println("  -", repo)
	}
	println("\n")
	return nil
}
