package bitbucket

import (
	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
	bb "github.com/ktrysmt/go-bitbucket"
)

type BitbucketClient struct {
	tokenManager *token.Manager
	username     string
}

func NewBitbucketClient(username string, tokens []string) *BitbucketClient {
	return &BitbucketClient{
		tokenManager: token.NewManager(tokens),
		username:     username,
	}
}

func (c *BitbucketClient) GetTokenManager() *token.Manager {
	return c.tokenManager
}

func (c *BitbucketClient) createClient() *bb.Client {
	return bb.NewBasicAuth(c.username, c.tokenManager.GetNextToken())
}

func (c *BitbucketClient) Sync(cfg config.Config) error {
	repos, err := c.getRepos(cfg)
	if err != nil {
		return err
	}

	gitSync.LogRepoCount(len(repos), cfg.Platform)

	gitSync.SyncWithConcurrency(cfg, repos, func(repo *bb.Repository) {
		gitSync.CloneOrUpdateRepo(cfg.Workspace, repo.Name, cfg)
		if cfg.IncludeWiki && repo.Has_wiki {
			gitSync.SyncWiki(cfg.Workspace, repo.Name, cfg)
		}
	})

	gitSync.LogSyncSummary(&cfg)
	return nil
}

func (c *BitbucketClient) getRepos(cfg config.Config) ([]*bb.Repository, error) {
	client := c.createClient()
	opt := &bb.RepositoriesOptions{
		Owner: cfg.Workspace,
		Page:  &[]int{1}[0],
	}

	var allRepos []*bb.Repository
	for {
		repos, err := client.Repositories.ListForAccount(opt)
		if err != nil {
			logger.Debugf("Error with current token, trying next token: %v", err)
			client = c.createClient()
			continue
		}

		var reposToInclude []*bb.Repository

		for _, repo := range repos.Items {
			repoName := repo.Name

			if len(cfg.IncludeRepos) > 0 {
				if helpers.IsIncludedInList(cfg.IncludeRepos, repoName) {
					logger.Debug("[include_repos] Repo included: ", repoName)
					reposToInclude = append(reposToInclude, &repo)
				}

				continue
			}

			if len(cfg.ExcludeRepos) > 0 {
				if helpers.IsIncludedInList(cfg.ExcludeRepos, repoName) {
					logger.Debug("[exclude_repos] Repo excluded: ", repoName)
					continue
				}
			}

			if !cfg.IncludeForks && repo.Parent != nil {
				logger.Debug("[include_forks] Repo excluded: ", repoName)
				continue
			}

			logger.Debug("Repo included: ", repoName)
			reposToInclude = append(reposToInclude, &repo)
		}

		allRepos = append(allRepos, reposToInclude...)

		if repos.Size < repos.Pagelen {
			break
		}

		opt.Page = &[]int{*opt.Page + 1}[0]
	}

	return allRepos, nil
}
