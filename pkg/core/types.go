package core

type GitRepository struct {
	LocalPath string
	UserName  string
	UserEmail string
}

type AzureDevOps struct {
	OrganisationName string
	ProjectName      string
	PAT              string
	WorkloadFeedName string
	WikiRepoName     string
	WikiRepo         *GitRepository
}

type Workload struct {
	Name      string
	Version   string
	Build     string
	Readme    string
	Instances int
}
