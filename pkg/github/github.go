package github

import (
	"context"

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

func (c *Client) FetchAllRepos() ([]*gh.Repository, error) {
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

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

func (c *Client) FetchRepos(repos []string) ([]*gh.Repository, error) {
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
