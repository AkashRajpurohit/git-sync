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
	projects, err := c.GetProjects(cfg)
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
			<-sem
		}(project)
	}

	wg.Wait()

	return nil
}

func (c GitlabClient) GetProjects(cfg config.Config) ([]*gl.Project, error) {
	projects, _, err := c.Client.Projects.ListProjects(&gl.ListProjectsOptions{Owned: &[]bool{true}[0]})
	if err != nil {
		return nil, err
	}

	var projectsToInclude []*gl.Project
	for _, project := range projects {
		projectName := project.Path
		IsGroupProject := project.Namespace.Kind == "group"
		groupName := project.Namespace.FullPath

		if len(cfg.IncludeOrgs) > 0 {
			if IsGroupProject && helpers.IsOrgIncluded(cfg.IncludeOrgs, groupName) {
				logger.Debug("[include_orgs] Project included: ", projectName)
				projectsToInclude = append(projectsToInclude, project)
			}

			continue
		}

		// If exclude orgs are set, exclude those and move to next checks if any
		if len(cfg.ExcludeOrgs) > 0 {
			if IsGroupProject && helpers.IsOrgExcluded(cfg.ExcludeOrgs, groupName) {
				logger.Debug("[exclude_orgs] Repo excluded: ", projectName)
				continue
			}
		}

		if len(cfg.IncludeRepos) > 0 {
			if helpers.IsRepoIncluded(cfg.IncludeRepos, projectName) {
				logger.Debug("[include_repos] Project included: ", projectName)
				projectsToInclude = append(projectsToInclude, project)
			}

			continue
		}

		if len(cfg.ExcludeRepos) > 0 {
			if helpers.IsRepoExcluded(cfg.ExcludeRepos, projectName) {
				logger.Debug("[exclude_repos] Project excluded: ", projectName)
				continue
			}
		}

		logger.Debug("Project included: ", projectName)
		projectsToInclude = append(projectsToInclude, project)
	}

	return projectsToInclude, nil
}
