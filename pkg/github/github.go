package github

import (
	"context"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	gh "github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
)

type Client struct {
	Username string
	Token    string
	Client   *gh.Client
}

func NewClient(username, token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := gh.NewClient(tc)

	return &Client{
		Username: username,
		Token:    token,
		Client:   client,
	}
}

func (c *Client) fetchListOfRepos(cfg config.Config) ([]*gh.Repository, error) {
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

			// If include repos are set, only include those and skip the rest
			if len(cfg.IncludeRepos) > 0 {
				if helpers.IsRepoIncluded(cfg.IncludeRepos, repo.GetName()) {
					reposToInclude = append(reposToInclude, repo)
				}

				continue
			}

			// If exclude repos are set, exclude those and move to next checks if any
			if len(cfg.ExcludeRepos) > 0 {
				if helpers.IsRepoExcluded(cfg.ExcludeRepos, repo.GetName()) {
					continue
				}
			}

			// If include forks is not set, skip forks
			if !cfg.IncludeForks && repo.GetFork() {
				continue
			}

			// If none of the above conditions are met, include the repo
			// This usually means that you don't have include_repos or the current repo is not in exclude_repos
			reposToInclude = append(reposToInclude, repo)
		}

		allRepos = append(allRepos, reposToInclude...)
		if resp.NextPage == 0 {
			break
		}

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
