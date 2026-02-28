package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	libpuki "github.com/moriT958/pukiwiki-mcp"
)

type EditPageInput struct {
	PageName   string `json:"page_name"   jsonschema:"Name of the wiki page to edit"`
	NewContent string `json:"new_content" jsonschema:"New wiki source content (full replacement)"`
}

func RegisterEditPage(s *mcp.Server, c *libpuki.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "edit_page",
		Description: "PukiWiki の既存ページを上書き編集する",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input EditPageInput) (*mcp.CallToolResult, any, error) {
		if input.PageName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "page_name is required"}},
				IsError: true,
			}, nil, nil
		}
		if input.NewContent == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "new_content is required"}},
				IsError: true,
			}, nil, nil
		}

		if err := c.EditPage(input.PageName, input.NewContent); err != nil {
			if errors.Is(err, libpuki.ErrPageNotFound) {
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q not found", input.PageName)}},
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
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("edit_page failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q updated successfully", input.PageName)}},
		}, nil, nil
	})
}
