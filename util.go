package main

import git "github.com/libgit2/git2go/v28"

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
