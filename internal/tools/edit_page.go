package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	libpuki "github.com/moriT958/pukiwiki-mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
)

type EditPageInput struct {
	PageName   string `json:"page_name"   jsonschema:"Name of the wiki page to edit"`
	NewContent string `json:"new_content" jsonschema:"New wiki source content (full replacement)"`
}

func RegisterEditPage(s *mcp.Server, p *auth.Provider) {
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

		c, err := p.Get(ctx)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("auth error: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		if err := c.EditPage(input.PageName, input.NewContent); err != nil {
			if errors.Is(err, libpuki.ErrSessionExpired) {
				if resetErr := p.Reset(); resetErr != nil {
					return &mcp.CallToolResult{
						Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("session expired but failed to clear credentials: %v. please retry.", resetErr)}},
						IsError: true,
					}, nil, nil
				}
				return &mcp.CallToolResult{
					Content: []mcp.Content{&mcp.TextContent{Text: "session expired. setup wizard launched. please retry after login."}},
					IsError: true,
				}, nil, nil
			}
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
