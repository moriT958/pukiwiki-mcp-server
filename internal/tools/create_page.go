package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	libpuki "github.com/moriT958/pukiwiki-mcp"
)

type CreatePageInput struct {
	PageName string `json:"page_name" jsonschema:"Name of the wiki page to create"`
	Content  string `json:"content"   jsonschema:"Wiki source content for the new page"`
}

func RegisterCreatePage(s *mcp.Server, c *libpuki.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_page",
		Description: "PukiWiki に新規ページを作成する",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreatePageInput) (*mcp.CallToolResult, any, error) {
		if input.PageName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "page_name is required"}},
				IsError: true,
			}, nil, nil
		}
		if input.Content == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "content is required"}},
				IsError: true,
			}, nil, nil
		}

		err := c.CreatePage(input.PageName, input.Content)
		if err != nil {
			if errors.Is(err, libpuki.ErrPageAlreadyExists) {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q already exists", input.PageName)}},
					IsError: true,
				}, nil, nil
			}
			if errors.Is(err, libpuki.ErrOutOfScope) {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q is outside the configured scope", input.PageName)}},
					IsError: true,
				}, nil, nil
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("create_page failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q created successfully", input.PageName)}},
		}, nil, nil
	})
}
