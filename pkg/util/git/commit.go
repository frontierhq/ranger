package git

import (
	"fmt"

	git "github.com/libgit2/git2go/v34"
)

func Commit(repository *git.Repository, branch *git.Branch, userEmail string, userName string, message string) (*git.Oid, error) {
	signature := &git.Signature{
		Email: userEmail,
		Name:  userName,
	}

	idx, err := repository.Index()
	if err != nil {
		return nil, err
	}

	idx.AddAll([]string{}, git.IndexAddDefault, nil)
	if err != nil {
		return nil, err
	}

	treeId, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}
	err = idx.Write()
	if err != nil {
		return nil, err
	}
	tree, err := repository.LookupTree(treeId)
	if err != nil {
		return nil, err
	}
	commitTarget, err := repository.LookupCommit(branch.Target())
	if err != nil {
		return nil, err
	}

	branchName, err := branch.Name()
	if err != nil {
		return nil, err
	}

	commit, err := repository.CreateCommit(fmt.Sprintf("refs/heads/%s", branchName), signature, signature, message, tree, commitTarget)
	if err != nil {
		return nil, err
	}

	return commit, nil
}
