package repository

import (
	"context"
	"mcp-gitlab-mng/internal/domain"
	"regexp"
	"strings"

	"gitlab.com/gitlab-org/api/client-go"
)

type GitLabRepository struct {
	client *gitlab.Client
}

func NewGitLabRepository(token, baseURL string) (*GitLabRepository, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, err
	}

	return &GitLabRepository{
		client: client,
	}, nil
}

func (r *GitLabRepository) ListRepositories(ctx context.Context, groupPath string, onlyPrivate bool) ([]*domain.Repository, error) {
	var projects []*gitlab.Project
	var err error
	var visibility gitlab.VisibilityValue

	if onlyPrivate {
		visibility = gitlab.PrivateVisibility
	}

	if groupPath != "" {
		// List projects for specific group
		opts := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
			Visibility: &visibility,
		}
		projects, _, err = r.client.Groups.ListGroupProjects(groupPath, opts)
	} else {
		// List all projects accessible to user
		opts := &gitlab.ListProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
			},
			Visibility: &visibility,
		}
		projects, _, err = r.client.Projects.ListProjects(opts)
	}

	if err != nil {
		return nil, err
	}

	repositories := make([]*domain.Repository, len(projects))
	for i, project := range projects {
		repositories[i] = convertProjectToRepository(project)
	}

	return repositories, nil
}

func (r *GitLabRepository) GetRepository(projectID int) (*domain.Repository, error) {
	project, _, err := r.client.Projects.GetProject(projectID, nil)
	if err != nil {
		return nil, err
	}

	return convertProjectToRepository(project), nil
}

func convertProjectToRepository(project *gitlab.Project) *domain.Repository {
	repo := &domain.Repository{
		ID:                project.ID,
		Name:              project.Name,
		NameWithNamespace: project.NameWithNamespace,
		Path:              project.Path,
		PathWithNamespace: project.PathWithNamespace,
		HTTPURLToRepo:     project.HTTPURLToRepo,
		SSHURLToRepo:      project.SSHURLToRepo,
		WebURL:            project.WebURL,
		Description:       project.Description,
		DefaultBranch:     project.DefaultBranch,
		Visibility:        string(project.Visibility),
	}

	if project.CreatedAt != nil {
		repo.CreatedAt = *project.CreatedAt
	}
	if project.LastActivityAt != nil {
		repo.LastActivityAt = *project.LastActivityAt
	}

	return repo
}

func (r *GitLabRepository) ListTerraformVersions(ctx context.Context, repositories []*domain.Repository) ([]*domain.RepositoryTerraformVersions, error) {
	var result []*domain.RepositoryTerraformVersions

	for _, repo := range repositories {
		components, err := r.searchTerraformVersionsInRepository(ctx, repo.ID)
		if err != nil {
			return nil, err
		}

		if len(components) > 0 {
			result = append(result, &domain.RepositoryTerraformVersions{
				Repository: repo,
				Components: components,
			})
		}
	}

	return result, nil
}

func (r *GitLabRepository) searchTerraformVersionsInRepository(ctx context.Context, projectID int) ([]*domain.TerraformComponent, error) {
	opts := &gitlab.SearchOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	searchResults, _, err := r.client.Search.BlobsByProject(projectID, "required_version", opts)
	if err != nil {
		return nil, err
	}

	var components []*domain.TerraformComponent
	requiredVersionRegex := regexp.MustCompile(`required_version\s*=\s*"([^"]+)"`)

	for _, result := range searchResults {
		if result.Data != "" {
			content := result.Data
			matches := requiredVersionRegex.FindAllStringSubmatch(content, -1)

			for _, match := range matches {
				if len(match) > 1 {
					// Extract directory path from file path
					filePath := result.Path

					dirPath := strings.TrimSuffix(filePath, "/"+getFileName(filePath))
					if dirPath == "" {
						dirPath = "/"
					}

					component := &domain.TerraformComponent{
						Path:            dirPath,
						RequiredVersion: match[1],
						FilePath:        filePath,
					}
					components = append(components, component)
				}
			}
		}
	}

	return components, nil
}

func getFileName(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return filePath
}

func (r *GitLabRepository) ListCodeByKeyword(ctx context.Context, repositories []*domain.Repository, keyword string, regexPatterns map[string]string) ([]*domain.RepositoryCodeSearch, error) {
	var result []*domain.RepositoryCodeSearch

	for _, repo := range repositories {
		codeMatches, err := r.searchCodeByKeywordInRepository(ctx, repo.ID, keyword, regexPatterns)
		if err != nil {
			return nil, err
		}

		if len(codeMatches) > 0 {
			result = append(result, &domain.RepositoryCodeSearch{
				Repository: repo,
				Results:    codeMatches,
			})
		}
	}

	return result, nil
}

func (r *GitLabRepository) searchCodeByKeywordInRepository(ctx context.Context, projectID int, keyword string, regexPatterns map[string]string) ([]*domain.CodeMatch, error) {
	opts := &gitlab.SearchOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	searchResults, _, err := r.client.Search.BlobsByProject(projectID, keyword, opts)
	if err != nil {
		return nil, err
	}

	var codeMatches []*domain.CodeMatch
	compiledRegexes := make(map[string]*regexp.Regexp)

	// Compile all regex patterns once
	for name, pattern := range regexPatterns {
		compiledRegex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		compiledRegexes[name] = compiledRegex
	}

	for _, result := range searchResults {
		if result.Data != "" {
			content := result.Data
			matches := make(map[string]string)

			// Apply all regex patterns to the content
			for name, compiledRegex := range compiledRegexes {
				regexMatches := compiledRegex.FindAllStringSubmatch(content, -1)
				for _, match := range regexMatches {
					if len(match) > 1 {
						// Use the first capturing group, or full match if no groups
						if len(match) > 1 {
							matches[name] = match[1]
						} else {
							matches[name] = match[0]
						}
					}
				}
			}

			// Only create CodeMatch if we found matches
			if len(matches) > 0 {
				filePath := result.Path
				dirPath := strings.TrimSuffix(filePath, "/"+getFileName(filePath))
				if dirPath == "" {
					dirPath = "/"
				}

				codeMatch := &domain.CodeMatch{
					Path:     dirPath,
					FilePath: filePath,
					Matches:  matches,
				}
				codeMatches = append(codeMatches, codeMatch)
			}
		}
	}

	return codeMatches, nil
}
