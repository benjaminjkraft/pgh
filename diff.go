package main

import (
	"fmt"
	"os"
	"strings"
)

func diff(runner *runner, args ...string) error {
	branch := "main"
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

	mainID, err := referenceNameToID(runner.repo, "refs/heads/"+branch)
	if err != nil {
		return err
	}
	if *headID != *mainID {
		return fmt.Errorf("TODO: must be main for now, but head was %v and main is %v", headID, mainID)
	}

	config, err := runner.repo.Config()
	if err != nil {
		return err
	}
	remoteName, err := config.LookupString(fmt.Sprintf("branch.%s.remote", branch))
	if err != nil {
		return err
	}
	if remoteName == "" {
		return fmt.Errorf("branch doesn't have remote?")
	}

	upstreamID, err := referenceNameToID(runner.repo, "refs/remotes/"+remoteName+"/"+branch)
	if err != nil {
		return err
	}
	if *parent.Id() != *upstreamID {
		return fmt.Errorf("TODO: must be one commit ahead of upstream for now, but parent was %v and main is %v", parent.Id(), mainID)
	}

	remote, err := runner.repo.Remotes.Lookup(remoteName)
	if err != nil {
		return err
	}
	pushUrl := remote.PushUrl()
	if pushUrl == "" {
		pushUrl = remote.Url()
	}
	if pushUrl == "" {
		return fmt.Errorf("branch doesn't have remote URL?")
	}

	remoteRepoOwner, remoteRepoName, err := parseRemote(pushUrl)
	if err != nil {
		return err
	}

	remoteRepoDetails, err := getRepoID(runner.ctx, runner.client, remoteRepoOwner, remoteRepoName)
	if err != nil {
		return err
	}
	remoteRepoID := remoteRepoDetails.Repository.Id

	// TODO: text wrapping? ugh
	msgParts := strings.SplitN(headCommit.Message(), "\n", 2)
	title := msgParts[0]
	body := strings.TrimLeft(msgParts[1], "\n")

	// TODO: customize, literally anything else
	branchName := fmt.Sprintf("%s.pgh.%s", os.Getenv("USER"), randomID())
	refspec := fmt.Sprintf("%s:refs/heads/%s", headID.String(), branchName)
	err = remote.Push([]string{refspec}, nil)
	if err != nil {
		return err
	}

	resp, err := createPR(runner.ctx, runner.client, CreatePullRequestInput{
		BaseRefName:         branch,
		Body:                body,
		Draft:               false,
		HeadRefName:         branchName,
		HeadRepositoryId:    remoteRepoID, // TODO: could be a fork
		MaintainerCanModify: true,
		RepositoryId:        remoteRepoID,
		Title:               title,
	})
	if err != nil {
		return err
	}

	prURL := resp.CreatePullRequest.PullRequest.Url

	fmt.Fprintln(runner.out, "created", prURL)

	return nil
}
