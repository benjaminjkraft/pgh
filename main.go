package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func Main() {
	r, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}
	ref, err := r.Head()
	if err != nil {
		panic(err)
	}
	fmt.Println(ref)
}

func main() {
	Main()
}
