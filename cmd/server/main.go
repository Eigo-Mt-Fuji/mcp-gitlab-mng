package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mcp-gitlab-mng/internal/repository"
	"mcp-gitlab-mng/internal/usecase"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var gitlabService *usecase.GitLabService

func main() {
	// Initialize GitLab repository
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		log.Fatal("GITLAB_TOKEN environment variable is required")
	}

	baseURL := os.Getenv("GITLAB_BASE_URL")
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	gitlabRepo, err := repository.NewGitLabRepository(token, baseURL)
	if err != nil {
		log.Fatalf("Failed to initialize GitLab client: %v", err)
	}

	gitlabService = usecase.NewGitLabService(gitlabRepo)

	mcpServer := server.NewMCPServer(
		"mcp-gitlab-mng",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Register list_repository tool
	listRepoTool := mcp.NewTool("list_repository",
		mcp.WithDescription("List GitLab repositories"),
		mcp.WithString("group",
			mcp.Description("GitLab group name (optional)"),
		),
	)

	mcpServer.AddTool(listRepoTool, handleListRepository)

	// Register list_terraform_versions tool
	listTerraformVersionsTool := mcp.NewTool("list_terraform_versions",
		mcp.WithDescription("List Terraform versions across GitLab repositories"),
		mcp.WithString("group",
			mcp.Description("GitLab group name (optional)"),
		),
	)

	mcpServer.AddTool(listTerraformVersionsTool, handleListTerraformVersions)

	// Register list_code_by_keyword tool
	listCodeByKeywordTool := mcp.NewTool("list_code_by_keyword",
		mcp.WithDescription("Search for code by keyword across GitLab repositories"),
		mcp.WithString("group",
			mcp.Description("GitLab group name (optional)"),
		),
		mcp.WithString("keyword",
			mcp.Description("Keyword to search for in code"),
		),
		mcp.WithObject("regex_patterns",
			mcp.Description("Optional regex patterns to extract specific data from matched files"),
		),
	)

	mcpServer.AddTool(listCodeByKeywordTool, handleListCodeByKeyword)

	// Start stdio transport
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatal(err)
	}
}

func handleListRepository(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	var groupPath string
	if group, ok := args["group"].(string); ok {
		groupPath = group
	}

	repositories, err := gitlabService.ListRepositories(ctx, groupPath, true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list repositories: %v", err)), nil
	}

	// Convert to JSON for better formatting
	repoJSON, err := json.MarshalIndent(repositories, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format repositories: %v", err)), nil
	}

	return mcp.NewToolResultText(string(repoJSON)), nil
}

func handleListTerraformVersions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	var groupPath string
	if group, ok := args["group"].(string); ok {
		groupPath = group
	}

	// First get repositories
	repositories, err := gitlabService.ListRepositories(ctx, groupPath, true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list repositories: %v", err)), nil
	}

	// Then get Terraform versions for those repositories
	terraformVersions, err := gitlabService.ListTerraformVersions(ctx, repositories)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list Terraform versions: %v", err)), nil
	}

	// Convert to JSON for better formatting
	terraformJSON, err := json.MarshalIndent(terraformVersions, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format Terraform versions: %v", err)), nil
	}

	return mcp.NewToolResultText(string(terraformJSON)), nil
}

func handleListCodeByKeyword(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()

	var groupPath string
	if group, ok := args["group"].(string); ok {
		groupPath = group
	}

	keyword, ok := args["keyword"].(string)
	if !ok || keyword == "" {
		return mcp.NewToolResultError("keyword parameter is required"), nil
	}

	// Parse regex patterns if provided
	regexPatterns := make(map[string]string)
	if patterns, ok := args["regex_patterns"].(map[string]interface{}); ok {
		for name, pattern := range patterns {
			if patternStr, ok := pattern.(string); ok {
				regexPatterns[name] = patternStr
			}
		}
	}

	// First get repositories
	repositories, err := gitlabService.ListRepositories(ctx, groupPath, true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to list repositories: %v", err)), nil
	}

	// Then search for code by keyword
	codeSearchResults, err := gitlabService.ListCodeByKeyword(ctx, repositories, keyword, regexPatterns)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search code by keyword: %v", err)), nil
	}

	// Convert to JSON for better formatting
	codeSearchJSON, err := json.MarshalIndent(codeSearchResults, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format code search results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(codeSearchJSON)), nil
}
