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

func (c *Client) fetchAllRepos(cfg config.Config) ([]*gh.Repository, error) {
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
			if !cfg.IncludeForks && repo.GetFork() {
				continue
			}

			if helpers.IsRepoExcluded(cfg.ExcludeRepos, repo.GetName()) {
				continue
			}

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

func (c *Client) fetchRepos(repos []string) ([]*gh.Repository, error) {
	ctx := context.Background()
	var allRepos []*gh.Repository

	for _, repo := range repos {
		repo, _, err := c.Client.Repositories.Get(ctx, c.Username, repo)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repo)
	}

	return allRepos, nil
}

func GetGitHubRepos(cfg config.Config) ([]*gh.Repository, error) {
	client := NewClient(cfg.Username, cfg.Token)

	var repos []*gh.Repository
	// If include_repos is set, only fetch those repos and that's it
	// no checks for forks or exclude_repos
	if len(cfg.IncludeRepos) > 0 {
		r, err := client.fetchRepos(cfg.IncludeRepos)
		if err != nil {
			return nil, err
		}

		repos = r
	} else {
		// for any other case, fetch all repos and then filter them
		r, err := client.fetchAllRepos(cfg)
		if err != nil {
			return nil, err
		}

		repos = r
	}

	return repos, nil
}
