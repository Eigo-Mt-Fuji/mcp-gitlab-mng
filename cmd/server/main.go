package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mcp-gitlab-mng/internal/repository"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var gitlabRepo *repository.GitLabRepository

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

	var err error
	gitlabRepo, err = repository.NewGitLabRepository(token, baseURL)
	if err != nil {
		log.Fatalf("Failed to initialize GitLab client: %v", err)
	}

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

	repositories, err := gitlabRepo.ListRepositories(groupPath)
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
