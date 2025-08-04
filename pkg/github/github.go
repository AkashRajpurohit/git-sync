package github

import (
	"context"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
	gh "github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	tokenManager *token.Manager
}

func NewGitHubClient(tokens []string) *GitHubClient {
	return &GitHubClient{
		tokenManager: token.NewManager(tokens),
	}
}

func (c *GitHubClient) GetTokenManager() *token.Manager {
	return c.tokenManager
}

func (c *GitHubClient) createClient() *gh.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.tokenManager.GetNextToken()},
	)
	tc := oauth2.NewClient(ctx, ts)
	return gh.NewClient(tc)
}

func (c *GitHubClient) Sync(cfg config.Config) error {
	repos, err := c.getRepos(cfg)
	if err != nil {
		return err
	}

	gitSync.LogRepoCount(len(repos), cfg.Platform)

	gitSync.SyncWithConcurrency(cfg, repos, func(repo *gh.Repository) {
		gitSync.CloneOrUpdateRepo(repo.GetOwner().GetLogin(), repo.GetName(), cfg)
		if cfg.IncludeWiki && repo.GetHasWiki() {
			gitSync.SyncWiki(repo.GetOwner().GetLogin(), repo.GetName(), cfg)
		}
	})

	gitSync.LogSyncSummary(&cfg)
	return nil
}

func (c *GitHubClient) getRepos(cfg config.Config) ([]*gh.Repository, error) {
	logger.Debug("Fetching list of repositories ⏳")
	ctx := context.Background()
	client := c.createClient()
	opt := &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}

	var allRepos []*gh.Repository
	for {
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			logger.Debugf("Error with current token, trying next token: %v", err)
			client = c.createClient()
			continue
		}

		var reposToInclude []*gh.Repository
		for _, repo := range repos {
			repoName := repo.GetName()
			isOrganizationRepo := repo.Owner.GetType() == "Organization"
			orgName := repo.Owner.GetLogin()

			if len(cfg.IncludeOrgs) > 0 {
				if isOrganizationRepo && helpers.IsIncludedInList(cfg.IncludeOrgs, orgName) {
					logger.Debug("[include_orgs] Repo included: ", repoName)
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			if len(cfg.ExcludeOrgs) > 0 {
				if isOrganizationRepo && helpers.IsIncludedInList(cfg.ExcludeOrgs, orgName) {
					logger.Debug("[exclude_orgs] Repo excluded: ", repoName)
					continue
				}
			}

			if len(cfg.IncludeRepos) > 0 {
				if helpers.IsIncludedInList(cfg.IncludeRepos, repoName) {
					logger.Debug("[include_repos] Repo included: ", repoName)
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			if len(cfg.ExcludeRepos) > 0 {
				if helpers.IsIncludedInList(cfg.ExcludeRepos, repoName) {
					logger.Debug("[exclude_repos] Repo excluded: ", repoName)
					continue
				}
			}

			if !cfg.IncludeForks && repo.GetFork() {
				logger.Debug("[include_forks] Repo excluded: ", repoName)
				continue
			}

			logger.Debug("Repo included: ", repoName)
			reposToInclude = append(reposToInclude, repo)
		}

		allRepos = append(allRepos, reposToInclude...)
		if resp.NextPage == 0 {
			break
		}

		logger.Debug("Fetching next page: ", resp.NextPage)
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}
