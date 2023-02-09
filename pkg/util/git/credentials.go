package git

import git "github.com/libgit2/git2go/v34"

func getUserpassPlaintextCredentialsCallback(username string, password string) git.CredentialsCallback {
	return func(url string, username_from_url string, allowed_types git.CredentialType) (*git.Credential, error) {
		cred, err := git.NewCredentialUserpassPlaintext(username, password)
		return cred, err
	}
}
