package github

import (
	"context"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gh "github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

type Client struct {
	Username string
	Token    string
	Client   *gh.Client
}

func NewClient(username, token string) *Client {
	logger.Debug("Creating new GitHub client ⏳")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)

	logger.Debug("GitHub client created ✅")

	return &Client{
		Username: username,
		Token:    token,
		Client:   client,
	}
}

func (c *Client) fetchListOfRepos(cfg config.Config) ([]*gh.Repository, error) {
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
			orgName := "personal"

			if repo.Owner.GetType() == "Organization" {
				orgName = repo.GetOrganization().GetName()
			}

			// If include orgs are set, only include those and skip the rest
			if len(cfg.IncludeOrgs) > 0 {
				if helpers.IsOrgIncluded(cfg.IncludeOrgs, orgName) {
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			// If exclude orgs are set, exclude those and move to next checks if any
			if len(cfg.ExcludeOrgs) > 0 {
				if helpers.IsOrgExcluded(cfg.ExcludeOrgs, repoName) {
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

func GetGitHubRepos(cfg config.Config) ([]*gh.Repository, error) {
	client := NewClient(cfg.Username, cfg.Token)

	repos, err := client.fetchListOfRepos(cfg)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
