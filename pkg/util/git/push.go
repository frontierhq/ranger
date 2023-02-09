package git

import (
	"fmt"

	git "github.com/libgit2/git2go/v34"
)

func Push(repository *git.Repository, branch *git.Branch, username string, password string) error {
	remote, err := repository.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	branchName, err := branch.Name()
	if err != nil {
		return err
	}

	pushOptions := &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: getUserpassPlaintextCredentialsCallback(username, password),
		},
	}
	if err := remote.Push([]string{fmt.Sprintf("refs/heads/%s", branchName)}, pushOptions); err != nil {
		return err
	}

	return nil
}
