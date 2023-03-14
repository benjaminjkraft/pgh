package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func head(args ...string) error {
	r, err := git.PlainOpen(".")
	if err != nil {
		return err
	}
	ref, err := r.Head()
	if err != nil {
		return err
	}
	fmt.Println(ref)
	return nil
}

func diff(args ...string) error {
	return nil
}

var commands = map[string]func(...string) error{
	"diff": diff,
	"head": head,
}

func Main() {
	// TODO: real CLI parser
	command := commands[os.Args[1]]
	if command == nil {
		fmt.Println("unknown command", os.Args[1])
		os.Exit(1)
	}

	err := command(os.Args[2:]...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Main()
}
