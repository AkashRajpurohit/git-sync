package forgejo

import (
	"fmt"

	fg "codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
)

type ForgejoClient struct {
	tokenManager *token.Manager
	serverConfig config.Server
}

func NewForgejoClient(serverConfig config.Server, tokens []string) *ForgejoClient {
	return &ForgejoClient{
		tokenManager: token.NewManager(tokens),
		serverConfig: serverConfig,
	}
}

func (c *ForgejoClient) GetTokenManager() *token.Manager {
	return c.tokenManager
}

func (c *ForgejoClient) createClient() (*fg.Client, error) {
	client, err := fg.NewClient(
		fmt.Sprintf("%s://%s", c.serverConfig.Protocol, c.serverConfig.Domain),
		fg.SetToken(c.tokenManager.GetNextToken()))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ForgejoClient) Sync(cfg config.Config) error {
	repos, err := c.getUserRepos(cfg)
	if err != nil {
		return err
	}

	gitSync.LogRepoCount(len(repos), cfg.Platform)

	gitSync.SyncWithConcurrency(cfg, repos, func(repo *fg.Repository) {
		gitSync.CloneOrUpdateRepo(repo.Owner.UserName, repo.Name, cfg)
		if cfg.IncludeWiki && repo.HasWiki {
			gitSync.SyncWiki(repo.Owner.UserName, repo.Name, cfg)
		}
	})

	gitSync.LogSyncSummary(&cfg)
	return nil
}

func (c *ForgejoClient) getUserRepos(cfg config.Config) ([]*fg.Repository, error) {
	logger.Debug("Fetching list of repositories ⏳")
	client, err := c.createClient()
	if err != nil {
		return nil, err
	}

	var allRepos []*fg.Repository
	pageOpt := fg.ListOptions{
		PageSize: 100,
	}

	for {
		repos, resp, err := client.ListMyRepos(fg.ListReposOptions{ListOptions: pageOpt})
		if err != nil {
			logger.Debugf("Error with current token, trying next token: %v", err)
			client, err = c.createClient()
			if err != nil {
				return nil, err
			}
			continue
		}

		var reposToInclude []*fg.Repository
		for _, repo := range repos {
			if len(cfg.IncludeRepos) > 0 {
				if helpers.IsIncludedInList(cfg.IncludeRepos, repo.FullName) {
					logger.Debug("[include_repos] Repo included: ", repo.Name)
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			if len(cfg.ExcludeRepos) > 0 {
				if helpers.IsIncludedInList(cfg.ExcludeRepos, repo.FullName) {
					logger.Debug("[exclude_repos] Repo excluded: ", repo.Name)
					continue
				}
			}

			if !cfg.IncludeForks && repo.Fork {
				logger.Debug("[include_forks] Repo excluded: ", repo.Name)
				continue
			}

			logger.Debug("Repo included: ", repo.Name)
			reposToInclude = append(reposToInclude, repo)
		}

		allRepos = append(allRepos, reposToInclude...)
		if resp.NextPage == 0 {
			break
		}

		logger.Debug("Fetching next page: ", resp.NextPage)
		pageOpt.Page = resp.NextPage
	}

	return allRepos, nil
}
