package git

import (
	errors "errors"
	"fmt"

	git "github.com/libgit2/git2go/v34"
)

func CheckoutBranch(repository *git.Repository, name string, create bool) (*git.Branch, error) {
	if !create {
		return nil, errors.New("branch checkout without create not implemented yet")
	}

	head, err := repository.Head()
	if err != nil {
		return nil, err
	}
	headCommit, err := repository.LookupCommit(head.Target())
	if err != nil {
		return nil, err
	}
	branch, err := repository.CreateBranch(name, headCommit, false)
	if err != nil {
		return nil, err
	}
	_, err = repository.References.CreateSymbolic("HEAD", fmt.Sprintf("refs/heads/%s", name), true, "headOne")
	if err != nil {
		return nil, err
	}
	checkoutOptions := &git.CheckoutOptions{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
	}
	if err := repository.CheckoutHead(checkoutOptions); err != nil {
		return nil, err
	}

	return branch, nil
}
