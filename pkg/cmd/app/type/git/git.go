package git

type Git struct {
	repositoryPath string
}

// NewGit creates a new Git
func NewGit(repositoryPath string) *Git {
	return &Git{
		repositoryPath: repositoryPath,
	}
}
