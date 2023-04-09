package main

import (
	"fmt"
	"io"
	"os"

	git "github.com/libgit2/git2go/v28"
)

type runner struct {
	repo *git.Repository
	out  io.Writer
}

var commands = map[string]func(*runner, ...string) error{
	"diff": diff,
}

func Main() {
	// TODO: real CLI parser
	command := commands[os.Args[1]]
	if command == nil {
		fmt.Println("unknown command", os.Args[1])
		os.Exit(1)
	}

	r, err := git.OpenRepositoryExtended(".", 0, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer r.Free()

	err = command(&runner{r, os.Stdout}, os.Args[2:]...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Main()
}
