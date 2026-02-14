package gitlab

import (
	"fmt"
	"time"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/helpers"
	"github.com/AkashRajpurohit/git-sync/pkg/issues"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
	gl "github.com/xanzy/go-gitlab"
)

type GitlabClient struct {
	tokenManager *token.Manager
	serverConfig config.Server
}

func NewGitlabClient(serverConfig config.Server, tokens []string) *GitlabClient {
	return &GitlabClient{
		tokenManager: token.NewManager(tokens),
		serverConfig: serverConfig,
	}
}

func (c *GitlabClient) GetTokenManager() *token.Manager {
	return c.tokenManager
}

func (c *GitlabClient) createClient() (*gl.Client, error) {
	baseURL := fmt.Sprintf("%s://%s/api/v4", c.serverConfig.Protocol, c.serverConfig.Domain)
	client, err := gl.NewClient(c.tokenManager.GetNextToken(), gl.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *GitlabClient) Sync(cfg config.Config) error {
	projects, err := c.getProjects(cfg)
	if err != nil {
		return err
	}

	gitSync.LogRepoCount(len(projects), cfg.Platform)

	gitSync.SyncWithConcurrency(cfg, projects, func(project *gl.Project) {
		gitSync.CloneOrUpdateRepo(project.Namespace.FullPath, project.Path, cfg)
		if cfg.IncludeWiki && project.WikiEnabled {
			gitSync.SyncWiki(project.Namespace.FullPath, project.Path, cfg)
		}
		if cfg.IncludeIssues && project.IssuesEnabled {
			since, hasPrevSync := issues.ReadLastSyncTime(cfg.BackupDir, project.Namespace.FullPath, project.Path)
			allIssues, err := c.fetchIssues(project.ID, since, hasPrevSync)
			if err != nil {
				logger.Errorf("Failed to fetch issues for %s/%s: %v", project.Namespace.FullPath, project.Path, err)
			} else {
				gitSync.SyncIssues(project.Namespace.FullPath, project.Path, allIssues, cfg)
			}
		}
	})

	gitSync.LogSyncSummary(&cfg)
	return nil
}

func (c *GitlabClient) getProjects(cfg config.Config) ([]*gl.Project, error) {
	client, err := c.createClient()
	if err != nil {
		return nil, err
	}

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
		pageResults, response, err := client.Projects.ListProjects(requestOpts, options...)

		if err != nil {
			logger.Debugf("Error with current token, trying next token: %v", err)
			client, err = c.createClient()
			if err != nil {
				return nil, err
			}
			continue
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

func (c *GitlabClient) fetchIssues(projectID int, since time.Time, incremental bool) ([]issues.Issue, error) {
	client, err := c.createClient()
	if err != nil {
		return nil, err
	}

	stateAll := "all"
	opt := &gl.ListProjectIssuesOptions{
		State: &stateAll,
		ListOptions: gl.ListOptions{
			PerPage: 100,
		},
	}

	if incremental {
		opt.UpdatedAfter = &since
		logger.Debugf("Incremental fetch for project %d (issues updated since %s)", projectID, since.Format(time.RFC3339))
	} else {
		logger.Debugf("Full fetch for project %d ⏳", projectID)
	}

	var allIssues []issues.Issue
	for {
		glIssues, resp, err := client.Issues.ListProjectIssues(projectID, opt)
		if err != nil {
			return nil, err
		}

		for _, glIssue := range glIssues {
			issue := convertGitLabIssue(glIssue)

			logger.Debugf("Fetching notes for issue #%d in project %d", glIssue.IID, projectID)
			notes, err := fetchIssueNotes(client, projectID, glIssue.IID)
			if err != nil {
				logger.Warnf("Failed to fetch notes for issue #%d: %v", glIssue.IID, err)
			} else {
				issue.Comments = notes
			}

			allIssues = append(allIssues, issue)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	logger.Debugf("Fetched %d issues for project %d", len(allIssues), projectID)
	return allIssues, nil
}

func fetchIssueNotes(client *gl.Client, projectID, issueIID int) ([]issues.Comment, error) {
	opt := &gl.ListIssueNotesOptions{
		ListOptions: gl.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []issues.Comment
	for {
		notes, resp, err := client.Notes.ListIssueNotes(projectID, issueIID, opt)
		if err != nil {
			return nil, err
		}

		for _, note := range notes {
			if note.System {
				continue
			}
			allComments = append(allComments, issues.Comment{
				ID:        int64(note.ID),
				Body:      note.Body,
				Author:    issues.User{Login: note.Author.Username, URL: note.Author.WebURL},
				URL:       "",
				CreatedAt: *note.CreatedAt,
				UpdatedAt: *note.UpdatedAt,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allComments, nil
}

func convertGitLabIssue(glIssue *gl.Issue) issues.Issue {
	issue := issues.Issue{
		Number:    glIssue.IID,
		Title:     glIssue.Title,
		Body:      glIssue.Description,
		State:     glIssue.State,
		Author:    issues.User{Login: glIssue.Author.Username, URL: glIssue.Author.WebURL},
		URL:       glIssue.WebURL,
		CreatedAt: *glIssue.CreatedAt,
		UpdatedAt: *glIssue.UpdatedAt,
	}

	if glIssue.ClosedAt != nil {
		issue.ClosedAt = glIssue.ClosedAt
	}

	if glIssue.Milestone != nil {
		issue.Milestone = glIssue.Milestone.Title
	}

	labels := make([]string, len(glIssue.Labels))
	copy(labels, glIssue.Labels)
	issue.Labels = labels

	assignees := make([]issues.User, 0, len(glIssue.Assignees))
	for _, a := range glIssue.Assignees {
		assignees = append(assignees, issues.User{Login: a.Username, URL: a.WebURL})
	}
	issue.Assignees = assignees

	return issue
}
