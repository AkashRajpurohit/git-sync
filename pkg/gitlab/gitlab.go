package gitlab

import (
	"fmt"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	gl "github.com/xanzy/go-gitlab"
)

type GitlabClient struct {
	Client *gl.Client
}

func NewGitlabClient(serverConfig config.Server, token string) *GitlabClient {
	baseURL := fmt.Sprintf("%s://%s/api/v4", serverConfig.Protocol, serverConfig.Domain)
	client, err := gl.NewClient(token, gl.WithBaseURL(baseURL))
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

	gitSync.LogRepoCount(len(projects), "GitLab")

	gitSync.SyncReposWithConcurrency(cfg, projects, func(project *gl.Project) {
		gitSync.CloneOrUpdateRepo(project.Namespace.FullPath, project.Path, cfg)
		if cfg.IncludeWiki && project.WikiEnabled {
			gitSync.SyncWiki(project.Namespace.FullPath, project.Path, cfg)
		}
	})

	gitSync.LogSyncComplete("GitLab")
	return nil
}

func (c GitlabClient) getProjects(cfg config.Config) ([]*gl.Project, error) {
	requestOpts := &gl.ListProjectsOptions{
		ListOptions: gl.ListOptions{
			OrderBy:    "id",
			Pagination: "keyset",
			PerPage:    100,
			Sort:       "asc",
		},
		Owned: &[]bool{true}[0],
	}

	options := []gl.RequestOptionFunc{}
	var projects []*gl.Project
	for {
		pageResults, response, err := c.Client.Projects.ListProjects(requestOpts, options...)

		if err != nil {
			return nil, err
		}

		projects = append(projects, pageResults...)

		if response.NextLink == "" {
			break
		}

		options = []gl.RequestOptionFunc{
			gl.WithKeysetPaginationParameters(response.NextLink),
		}
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
			if isGroupProject && helpers.IsIncludedInList(cfg.ExcludeOrgs, groupName) {
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
			if helpers.IsIncludedInList(cfg.ExcludeRepos, projectName) {
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
