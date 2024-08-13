package bitbucket

import (
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"

	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	bb "github.com/ktrysmt/go-bitbucket"
)

type BitbucketClient struct {
	Client *bb.Client
}

func NewBitbucketClient(username, password string) *BitbucketClient {
	client := bb.NewBasicAuth(username, password)

	return &BitbucketClient{
		Client: client,
	}
}

func (c BitbucketClient) Sync(cfg config.Config) error {
	repos, err := c.getRepos(cfg)
	if err != nil {
		return err
	}

	logger.Info("Total repositories: ", len(repos))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, repo := range repos {
		wg.Add(1)
		go func(repo *bb.Repository) {
			defer wg.Done()
			sem <- struct{}{}
			gitSync.CloneOrUpdateRepo(cfg.Workspace, repo.Name, cfg)
			if repo.Has_wiki {
				gitSync.SyncWiki(cfg.Workspace, repo.Name, cfg)
			}
			<-sem
		}(repo)
	}

	wg.Wait()
	logger.Info("All repositories synced ✅")

	return nil
}

func (c BitbucketClient) getRepos(cfg config.Config) ([]*bb.Repository, error) {
	opt := &bb.RepositoriesOptions{
		Owner: cfg.Workspace,
		Page:  &[]int{1}[0],
	}

	var allRepos []*bb.Repository
	for {
		repos, err := c.Client.Repositories.ListForAccount(opt)
		if err != nil {
			return nil, err
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
				if helpers.IsExcludedInList(cfg.ExcludeRepos, repoName) {
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
