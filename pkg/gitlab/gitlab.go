package gitlab

import (
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	gl "github.com/xanzy/go-gitlab"
)

type GitlabClient struct {
	Client *gl.Client
}

func NewGitlabClient(token string) *GitlabClient {
	client, err := gl.NewClient(token)
	if err != nil {
		return nil
	}

	return &GitlabClient{
		Client: client,
	}
}

func (c GitlabClient) Sync(cfg config.Config) error {
	projects, err := c.getProjects(cfg)
	if err != nil {
		return err
	}

	logger.Info("Total projects: ", len(projects))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, project := range projects {
		wg.Add(1)
		go func(project *gl.Project) {
			defer wg.Done()
			sem <- struct{}{}
			gitSync.CloneOrUpdateRepo(project.Namespace.FullPath, project.Path, cfg)
			if cfg.IncludeWiki && project.WikiEnabled {
				gitSync.SyncWiki(project.Namespace.FullPath, project.Path, cfg)
			}
			<-sem
		}(project)
	}

	wg.Wait()
	logger.Info("All projects synced ✅")

	return nil
}

func (c GitlabClient) getProjects(cfg config.Config) ([]*gl.Project, error) {
	projects, _, err := c.Client.Projects.ListProjects(&gl.ListProjectsOptions{Owned: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}

	var projectsToInclude []*gl.Project
	for _, project := range projects {
		projectName := project.Path
		isFork := project.ForkedFromProject != nil
		isGroupProject := project.Namespace.Kind == "group"
		groupName := project.Namespace.FullPath

		if len(cfg.IncludeOrgs) > 0 {
			if isGroupProject && helpers.IsIncludedInList(cfg.IncludeOrgs, groupName) {
				logger.Debug("[include_groups] Project included: ", projectName)
				projectsToInclude = append(projectsToInclude, project)
			}

			continue
		}

		// If exclude orgs are set, exclude those and move to next checks if any
		if len(cfg.ExcludeOrgs) > 0 {
			if isGroupProject && helpers.IsExcludedInList(cfg.ExcludeOrgs, groupName) {
				logger.Debug("[exclude_groups] Project excluded: ", projectName)
				continue
			}
		}

		if len(cfg.IncludeRepos) > 0 {
			if helpers.IsIncludedInList(cfg.IncludeRepos, projectName) {
				logger.Debug("[include_projects] Project included: ", projectName)
				projectsToInclude = append(projectsToInclude, project)
			}

			continue
		}

		if len(cfg.ExcludeRepos) > 0 {
			if helpers.IsExcludedInList(cfg.ExcludeRepos, projectName) {
				logger.Debug("[exclude_projects] Project excluded: ", projectName)
				continue
			}
		}

		if !cfg.IncludeForks && isFork {
			logger.Debug("[include_forks] Fork excluded: ", projectName)
			continue
		}

		logger.Debug("Project included: ", projectName)
		projectsToInclude = append(projectsToInclude, project)
	}

	return projectsToInclude, nil
}
