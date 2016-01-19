package main

import (
	"fmt"
	"kupak"
	"os"
)

func main() {
	repo, err := kupak.RepoFromUrl(os.Args[1])
	if err != nil {
		panic(err)
	}
	for i := range repo.Index {
		fmt.Println(repo.Index[i].String())
	}

	pak, err := repo.Pak("redis", "1.0")
	if err != nil {
		panic(err)
	}
	fmt.Println(pak.Name)
}
