package main

import (
	"fmt"
	"math/rand"
	"strings"

	git "github.com/libgit2/git2go/v28"
)

func referenceNameToID(repo *git.Repository, fullName string) (*git.Oid, error) {
	ref, err := repo.References.Lookup(fullName)
	if err != nil {
		return nil, err
	}
	for ref.Type() != git.ReferenceOid {
		ref, err = ref.Resolve()
		if err != nil {
			return nil, err
		}
	}
	return ref.Target(), nil
}

func randomID() string {
	n := 16
	var buf strings.Builder
	for i := 0; i < n; i++ {
		d := byte(rand.Intn(36))
		if d < 10 {
			buf.WriteByte('0' + d)
		} else {
			buf.WriteByte('a' + d - 10)
		}
	}
	return buf.String()
}

func parseRemote(url string) (owner, name string, err error) {
	rest, ok := strings.CutPrefix(url, "git@github.com:")
	if !ok {
		return "", "", fmt.Errorf("unknown GitHub URL prefix: %s", url)
	}
	rest, ok = strings.CutSuffix(rest, ".git")
	if !ok {
		return "", "", fmt.Errorf("unknown GitHub URL suffix: %s", url)
	}
	parts := strings.Split(rest, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("GitHub URL has too many parts: %s", url)
	}
	return parts[0], parts[1], nil
}
