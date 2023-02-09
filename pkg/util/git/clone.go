package git

import (
	"os"

	git "github.com/libgit2/git2go/v34"
)

func CloneOverHttp(url string, username string, password string) (*git.Repository, string, error) {
	repositoryPath, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, "", err
	}

	err = os.RemoveAll(repositoryPath)
	if err != nil {
		return nil, "", err
	}

	cloneOptions := &git.CloneOptions{
		FetchOptions: git.FetchOptions{
			RemoteCallbacks: git.RemoteCallbacks{
				CredentialsCallback: getUserpassPlaintextCredentialsCallback(username, password),
			},
		},
	}
	repository, err := git.Clone(url, repositoryPath, cloneOptions)
	if err != nil {
		return nil, "", err
	}

	return repository, repositoryPath, nil
}
