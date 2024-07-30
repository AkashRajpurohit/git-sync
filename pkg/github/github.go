package github

import (
	"context"
	"os"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	ghSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	gh "github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	Client *gh.Client
}

func NewGitHubClient(token string) *GitHubClient {
	logger.Debug("Creating new GitHub client ⏳")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)

	logger.Debug("GitHub client created ✅")

	return &GitHubClient{
		Client: client,
	}
}

func (c GitHubClient) Sync(cfg config.Config) error {
	backupDir := cfg.BackupDir
	os.MkdirAll(backupDir, os.ModePerm)

	repos, err := c.getRepos(cfg)
	if err != nil {
		return err
	}

	logger.Info("Total repositories: ", len(repos))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, repo := range repos {
		wg.Add(1)
		go func(repo *gh.Repository) {
			defer wg.Done()
			sem <- struct{}{}
			ghSync.CloneOrUpdateRepo(repo.GetOwner().GetLogin(), repo.GetName(), cfg)
			<-sem
		}(repo)
	}

	wg.Wait()

	return nil
}

func (c GitHubClient) getRepos(cfg config.Config) ([]*gh.Repository, error) {
	logger.Debug("Fetching list of repositories ⏳")
	ctx := context.Background()
	opt := &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: gh.ListOptions{PerPage: 100},
	}

	var allRepos []*gh.Repository
	for {
		repos, resp, err := c.Client.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			return nil, err
		}

		var reposToInclude []*gh.Repository
		for _, repo := range repos {
			repoName := repo.GetName()
			isOrganizationRepo := repo.Owner.GetType() == "Organization"
			orgName := repo.Owner.GetLogin()

			// If include orgs are set, only include those and skip the rest
			if len(cfg.IncludeOrgs) > 0 {
				if isOrganizationRepo && helpers.IsOrgIncluded(cfg.IncludeOrgs, orgName) {
					logger.Debug("[include_orgs] Repo included: ", repoName)
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			// If exclude orgs are set, exclude those and move to next checks if any
			if len(cfg.ExcludeOrgs) > 0 {
				if isOrganizationRepo && helpers.IsOrgExcluded(cfg.ExcludeOrgs, orgName) {
					logger.Debug("[exclude_orgs] Repo excluded: ", repoName)
					continue
				}
			}

			// If include repos are set, only include those and skip the rest
			if len(cfg.IncludeRepos) > 0 {
				if helpers.IsRepoIncluded(cfg.IncludeRepos, repoName) {
					logger.Debug("[include_repos] Repo included: ", repoName)
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			// If exclude repos are set, exclude those and move to next checks if any
			if len(cfg.ExcludeRepos) > 0 {
				if helpers.IsRepoExcluded(cfg.ExcludeRepos, repoName) {
					logger.Debug("[exclude_repos] Repo excluded: ", repoName)
					continue
				}
			}

			// If include forks is not set, skip forks
			if !cfg.IncludeForks && repo.GetFork() {
				logger.Debug("[include_forks] Repo excluded: ", repoName)
				continue
			}

			// If none of the above conditions are met, include the repo
			// This usually means that you don't have include_repos or the current repo is not in exclude_repos
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
