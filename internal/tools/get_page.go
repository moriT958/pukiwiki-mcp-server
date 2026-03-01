package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	libpuki "github.com/moriT958/pukiwiki-mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
)

type GetPageInput struct {
	PageName string `json:"page_name" jsonschema:"Name of the wiki page to retrieve"`
}

func RegisterGetPage(s *mcp.Server, p *auth.Provider) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_page",
		Description: "PukiWiki ページのソースを返す",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetPageInput) (*mcp.CallToolResult, any, error) {
		if input.PageName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "page_name is required"}},
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

		source, err := c.GetPageSource(input.PageName)
		if err != nil {
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
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("get_page failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: source}},
		}, nil, nil
	})
}
