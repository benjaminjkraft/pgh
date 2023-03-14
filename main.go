package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func diff(r *git.Repository, args ...string) error {
	head, err := r.Head()
	if err != nil {
		return err
	}
	headCommit, err := r.CommitObject(head.Hash())
	if err != nil {
		return err
	}
	switch headCommit.NumParents() {
	case 0:
		return fmt.Errorf("can't diff on initial commit")
	case 1:
	default:
		return fmt.Errorf("can't diff on merge commit")
	}

	parent, err := headCommit.Parent(0)
	if err != nil {
		return err
	}

	mainBranch, err := r.Branch("main")
	if err != nil {
		return err
	}
	main, err := r.Reference(mainBranch.Merge, true)
	if err != nil {
		return err
	}
	if parent.Hash != main.Hash() {
		return fmt.Errorf("TODO: must be one commit ahead of main for now, but parent was %v and main is %v", parent.Hash, main.Hash())
	}
	fmt.Println(parent, main)

	return nil
}

var commands = map[string]func(*git.Repository, ...string) error{
	"diff": diff,
}

func Main() {
	// TODO: real CLI parser
	command := commands[os.Args[1]]
	if command == nil {
		fmt.Println("unknown command", os.Args[1])
		os.Exit(1)
	}

	r, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = command(r, os.Args[2:]...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Main()
}
