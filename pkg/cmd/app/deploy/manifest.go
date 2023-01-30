package deploy

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/frontierdigital/ranger/pkg/util/output"
	git "github.com/libgit2/git2go/v34"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	adoGit "github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/segmentio/ksuid"
	"golang.org/x/exp/slices"

	cp "github.com/otiai10/copy"
)

func DeployManifest(config *app.Config, projectName string, organisationName string) error {
	organisationUrl := fmt.Sprintf("https://dev.azure.com/%s", organisationName)
	connection := azuredevops.NewPatConnection(organisationUrl, config.ADO.PAT)

	ctx := context.Background()

	buildClient, err := build.NewClient(ctx, connection)
	if err != nil {
		return err
	}

	gitClient, err := adoGit.NewClient(ctx, connection)
	if err != nil {
		return err
	}

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := manifest.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifestName := fmt.Sprintf("%s-%s", manifest.Environment, manifest.Layer)

	manifest.PrintHeader(manifestName, manifest.Layer, manifest.Environment, manifest.Version)

	for _, workload := range manifest.Workloads {
		sourceProjectName, sourceRepositoryName := workload.GetSourceProjectAndRepositoryNames()

		pipelineName := fmt.Sprintf("%s (deploy)", sourceRepositoryName)
		getDefinitionsArgs := build.GetDefinitionsArgs{
			Name:    &pipelineName,
			Project: &sourceProjectName,
		}
		pipelines, err := buildClient.GetDefinitions(ctx, getDefinitionsArgs)
		if err != nil {
			return err
		}
		if len(pipelines.Value) == 0 {
			return fmt.Errorf("pipeline with name '%s' not found", pipelineName)
		}
		if len(pipelines.Value) > 1 {
			return fmt.Errorf("multiple pipeline with name '%s' found", pipelineName)
		}
		pipeline := pipelines.Value[0]

		output.PrintlnfInfo("Found deploy pipeline definition with Id '%d' for workload '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
			*pipeline.Id, workload.Source, organisationName, projectName, *pipeline.Id)

		getRepositoriesArgs := adoGit.GetRepositoriesArgs{
			Project: &sourceProjectName,
		}
		repositories, err := gitClient.GetRepositories(ctx, getRepositoriesArgs)
		if err != nil {
			return err
		}
		findRepositoryFunc := func(r adoGit.GitRepository) bool { return *r.Name == sourceRepositoryName }
		repositoryIdx := slices.IndexFunc(*repositories, findRepositoryFunc)

		repository := (*repositories)[repositoryIdx]

		output.PrintlnfInfo("Found repository with Id '%s' for workload '%s' (%s)", repository.Id, workload.Source, *repository.WebUrl)

		workloadConfigPath := path.Join("config", "workloads", workload.Name)
		workloadConfigExists := true
		_, err = os.Stat(workloadConfigPath)
		if err != nil {
			workloadConfigExists = !os.IsNotExist(err)
		}

		workloadSecretsPath := path.Join("secrets", "workloads", workload.Name)
		workloadSecretsExists := true
		_, err = os.Stat(workloadSecretsPath)
		if err != nil {
			workloadSecretsExists = !os.IsNotExist(err)
		}

		if workloadConfigExists || workloadSecretsExists {
			configRepoPath, err := ioutil.TempDir("", "rangerconfig")
			if err != nil {
				return err
			}

			err = os.RemoveAll(configRepoPath)
			if err != nil {
				return err
			}

			configRepoName := "generated-manifest-config"
			configRepoUrl := fmt.Sprintf("https://frontierdigital@dev.azure.com/%s/%s/_git/%s", organisationName, projectName, configRepoName)

			cloneOptions := &git.CloneOptions{
				FetchOptions: git.FetchOptions{
					RemoteCallbacks: git.RemoteCallbacks{
						CredentialsCallback: makeCredentialsCallback(config.ADO.PAT, "x-oauth-basic"),
					},
				},
			}
			configRepo, err := git.Clone(configRepoUrl, configRepoPath, cloneOptions)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if err != nil {
				return err
			}

			configBranchName := fmt.Sprintf("%s_%s_%s_%s", workload.Name, manifest.Environment, manifest.Layer, ksuid.New().String())
			head, err := configRepo.Head()
			if err != nil {
				return err
			}
			headCommit, err := configRepo.LookupCommit(head.Target())
			if err != nil {
				return err
			}
			branch, err := configRepo.CreateBranch(configBranchName, headCommit, false)
			if err != nil {
				return err
			}
			_, err = configRepo.References.CreateSymbolic("HEAD", fmt.Sprintf("refs/heads/%s", configBranchName), true, "headOne")
			if err != nil {
				return err
			}
			checkoutOptions := &git.CheckoutOptions{
				Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
			}
			if err := configRepo.CheckoutHead(checkoutOptions); err != nil {
				return err
			}

			signature := &git.Signature{
				Email: config.Git.UserEmail,
				Name:  config.Git.UserName,
			}

			if workloadConfigExists {
				err = cp.Copy(workloadConfigPath, path.Join(configRepoPath, ".config"))
				if err != nil {
					return err
				}
			}

			if workloadSecretsExists {
				err = cp.Copy(workloadSecretsPath, path.Join(configRepoPath, ".secrets"))
				if err != nil {
					return err
				}
			}

			idx, err := configRepo.Index()
			if err != nil {
				return err
			}

			idx.AddAll([]string{}, git.IndexAddDefault, nil)
			if err != nil {
				return err
			}

			treeId, err := idx.WriteTree()
			if err != nil {
				return err
			}
			err = idx.Write()
			if err != nil {
				return err
			}
			tree, err := configRepo.LookupTree(treeId)
			if err != nil {
				return err
			}
			commitTarget, err := configRepo.LookupCommit(branch.Target())
			if err != nil {
				return err
			}
			commitMessage := fmt.Sprintf("Generate config for workload '%s', environment '%s' and layer '%s'", workload.Name, manifest.Environment, manifest.Layer)
			commitId, err := configRepo.CreateCommit(fmt.Sprintf("refs/heads/%s", configBranchName), signature, signature, commitMessage, tree, commitTarget)
			if err != nil {
				return err
			}
			_ = commitId

			remote, err := configRepo.Remotes.Lookup("origin")
			if err != nil {
				return err
			}

			pushOptions := &git.PushOptions{
				RemoteCallbacks: git.RemoteCallbacks{
					CredentialsCallback: makeCredentialsCallback(config.ADO.PAT, "x-oauth-basic"),
				},
			}
			if err := remote.Push([]string{fmt.Sprintf("refs/heads/%s", configBranchName)}, pushOptions); err != nil {
				return err
			}

			output.PrintlnfInfo("Pushed config (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", organisationName, projectName, configRepoName, configBranchName)

			defer os.RemoveAll(configRepoPath)
		}
	}

	return nil
}

func makeCredentialsCallback(username, password string) git.CredentialsCallback {
	// If we're trying it means the credentials are invalid
	called := false
	return func(url string, username_from_url string, allowed_types git.CredentialType) (*git.Credential, error) {
		if called {
			return nil, git.MakeGitError2(2)
		}
		called = true
		cred, err := git.NewCredentialUserpassPlaintext(username, password)
		return cred, err
	}
}
