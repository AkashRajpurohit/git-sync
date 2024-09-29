package forgejo

import (
	"fmt"
	"sync"

	fg "codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
)

type ForgejoClient struct {
	Client *fg.Client
}

func NewForgejoClient(serverConfig config.Server, token string) *ForgejoClient {
	logger.Debug("Creating new Forgejo client ⏳")

	client, err := fg.NewClient(
		fmt.Sprintf("%s://%s", serverConfig.Protocol, serverConfig.Domain),
		fg.SetToken(token))
	if err != nil {
		logger.Error("Error creating Forgejo client: ", err)
		return nil
	}

	logger.Debug("Forgejo client created ✅")

	return &ForgejoClient{
		Client: client,
	}
}

func (c ForgejoClient) Sync(cfg config.Config) error {
	repos, err := c.getUserRepos(cfg)
	if err != nil {
		return err
	}

	logger.Info("Total repositories: ", len(repos))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, repo := range repos {
		wg.Add(1)
		go func(repo *fg.Repository) {
			defer wg.Done()
			sem <- struct{}{}
			gitSync.CloneOrUpdateRepo(repo.Owner.UserName, repo.Name, cfg)
			if cfg.IncludeWiki && repo.HasWiki {
				gitSync.SyncWiki(repo.Owner.UserName, repo.Name, cfg)
			}
			<-sem
		}(repo)
	}

	wg.Wait()
	logger.Info("All repositories synced ✅")

	return nil
}

func (c ForgejoClient) getUserRepos(cfg config.Config) ([]*fg.Repository, error) {
	logger.Debug("Fetching list of repositories ⏳")
	var allRepos []*fg.Repository

	pageOpt := fg.ListOptions{
		PageSize: 100,
	}

	for {
		repos, resp, err := c.Client.ListMyRepos(fg.ListReposOptions{ListOptions: pageOpt})
		if err != nil {
			return nil, err
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
				if helpers.IsExcludedInList(cfg.ExcludeRepos, repo.FullName) {
					logger.Debug("[exclude_repos] Repo excluded: ", repo.Name)
					continue
				}
			}

			if !cfg.IncludeForks && repo.Fork {
				logger.Debug("[include_forks] Repo excluded: ", repo.Name)
				continue
			}

			// If none of the above conditions are met, include the repo
			// This usually means that you don't have include_repos or the current repo is not in exclude_repos
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
