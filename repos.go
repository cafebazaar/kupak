package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
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
	if err := ioutil.WriteFile(reposFilePath(), []byte("github.com/cafebazaar/paks"), 0644); err != nil {
		return err
	}
	return nil
}
