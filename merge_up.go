package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type graph[T comparable] map[T][]T

func dfsInner[T comparable](t graph[T], start T, b []T) []T {
	// TODO: check for cycles (lol)
	b = append(b, start)
	for _, child := range t[start] {
		b = dfsInner(t, child, b)
	}
	return b
}

func dfs[T comparable](t graph[T], start T) []T {
	return dfsInner(t, start, nil)
}

func makeBranchGraph(repo *git.Repository) (graph[plumbing.ReferenceName], error) {
	g := graph[plumbing.ReferenceName]{}
	bs, err := repo.Branches()
	if err != nil {
		return nil, err
	}
	err = bs.ForEach(func(b *plumbing.Reference) error {
		u := upstream(repo, b.Name(), false)
		if u != "" {
			g[u] = append(g[u], b.Name())
		}
		return nil
	})
	return g, err
}

func mergeUp(runner *runner, _ ...string) error {
	head, err := runner.repo.Head()
	if err != nil {
		return err
	}
	branchGraph, err := makeBranchGraph(runner.repo)
	if err != nil {
		return err
	}
	// slice bc no need to merge head itself.
	branches := dfs(branchGraph, head.Name())[1:]

	for _, refName := range branches {
		// short bc checkout interprets a full ref as implying --detach
		// TODO: does Worktree.Checkout have the reset bug too?
		err = callGit(runner, "checkout", refName.Short())
		if err != nil {
			return err
		}
		err = callGit(runner, "merge") // (defaults to the upstream)
		if err != nil {
			// TODO: nice message + tell you what's left on merge
			return err
		}
	}

	return nil
}
