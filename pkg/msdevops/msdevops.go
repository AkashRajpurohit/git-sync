package msdevops

import (
	"context"
	"fmt"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
)

type MSDevOpsClient struct {
	tokenManager *token.Manager
	serverConfig config.Server
}

func NewMSDevOpsClient(serverConfig config.Server, tokens []string) *MSDevOpsClient {
	return &MSDevOpsClient{
		tokenManager: token.NewManager(tokens),
		serverConfig: serverConfig,
	}
}

func (c *MSDevOpsClient) GetTokenManager() *token.Manager {
	return c.tokenManager
}

// createClient initializes and returns a new Azure DevOps Git client using the provided token manager and server configuration.
func (c *MSDevOpsClient) createClient() (git.Client, error) {
	//logger.InitLogger("debug")

	organizationURL := fmt.Sprintf("%s://%s", c.serverConfig.Protocol, c.serverConfig.Domain)
	logger.Debugf("organizationURL: %s", organizationURL)
	ctx := context.Background()
	token := c.tokenManager.GetNextToken()
	if token == "" {
		return nil, fmt.Errorf("a valid token was not available")
	}
	connection := azuredevops.NewPatConnection(organizationURL, token)

	client, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// derefString returns the value of a string pointer or an empty string if the pointer is nil.
func derefString(ref *string) string {
	if ref == nil {
		return ""
	}
	return *ref
}

func (c *MSDevOpsClient) Sync(cfg config.Config) error {
	repos, err := c.getUserRepos(cfg)
	if err != nil {
		return err
	}

	gitSync.LogRepoCount(len(repos), cfg.Platform)

	gitSync.SyncReposWithConcurrency(cfg, repos, func(repo git.GitRepository) {
		repoOwner := derefString(repo.Project.Name)
		repoName := derefString(repo.Name)
		repoURL := derefString(repo.WebUrl)
		protoLen := len(cfg.Server.Protocol + "://")

		// Need to manually construct the repo URL by inserting the user token into the URL
		repoAuthURL := repoURL[:protoLen] + c.tokenManager.GetNextToken() + "@" + repoURL[protoLen:]
		logger.Debugf("repoAuthURL: %s", repoAuthURL)

		if *repo.IsDisabled {
			logger.Warnf("Skipping repo %s as it is disabled", repoName)
		} else {
			gitSync.CloneOrUpdateRawRepo(repoOwner, repoName, repoAuthURL, cfg)
		}

	})
	gitSync.LogSyncSummary()
	return nil
}

func (c *MSDevOpsClient) getUserRepos(cfg config.Config) ([]git.GitRepository, error) {
	logger.Debug("Fetching list of repositories ⏳")
	client, err := c.createClient()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	allRepos, err := client.GetRepositories(ctx, git.GetRepositoriesArgs{
		Project: &cfg.Workspace,
	})
	if err != nil {
		return nil, err
	}

	for _, repo := range *allRepos {
		logger.Debugf("Found repo: %s", derefString(repo.Name))
		logger.Debugf("Repo WebURL: %s", derefString(repo.WebUrl))
	}

	return *allRepos, nil
}
