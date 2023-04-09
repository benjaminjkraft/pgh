package main

import (
	"fmt"
)

func diff(runner *runner, args ...string) error {
	headID, err := referenceNameToID(runner.repo, "HEAD")
	if err != nil {
		return err
	}
	headCommit, err := runner.repo.LookupCommit(headID)
	if err != nil {
		return err
	}

	switch headCommit.ParentCount() {
	case 0:
		return fmt.Errorf("can't diff on initial commit")
	case 1:
	default:
		return fmt.Errorf("can't diff on merge commit")
	}
	parent := headCommit.Parent(0)

	mainID, err := referenceNameToID(runner.repo, "refs/heads/main")
	if err != nil {
		return err
	}
	if *parent.Id() != *mainID {
		return fmt.Errorf("TODO: must be one commit ahead of main for now, but parent was %v and main is %v", parent.Id(), mainID)
	}
	fmt.Fprintln(runner.out, "parent", parent.Id(), parent.Message())

	return nil
}
