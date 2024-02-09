package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

// TODO: put repo root in runner, simplify tmpdr plumbing + recomputing etc
type runner struct {
	repo *git.Repository
	out  io.Writer
}

func help(_ *runner, _ ...string) error {
	fmt.Println("commands:")
	for name := range commands {
		if name[0] == '-' {
			continue
		}
		fmt.Println("\t" + name)
	}
	return nil
}

var commands map[string]func(*runner, ...string) error

func init() {
	commands = map[string]func(*runner, ...string) error{
		"fake-merge": fakeMerge,
		"merge-up":   mergeUp,
		"-h":         help,
		"--help":     help,
	}
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

	err = command(&runner{r, os.Stdout}, os.Args[2:]...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Main()
}
