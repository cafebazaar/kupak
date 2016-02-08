package main

import (
	"kupak"
	"os"
)

func main() {
	repo, err := kupak.RepoFromUrl(os.Args[1])
	if err != nil {
		panic(err)
	}

	pak, err := repo.Pak("redis", "1.0")
	if err != nil {
		panic(err)
	}

	manager, err := kupak.NewManager()
	if err != nil {
		panic(err)
	}

	err = manager.Install(pak, "default", map[string]interface{}{})
	if err != nil {
		panic(err)
	}
}
