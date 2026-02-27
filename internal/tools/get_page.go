package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/libpuki"
)

type GetPageInput struct {
	PageName string `json:"page_name" jsonschema:"Name of the wiki page to retrieve"`
}

func RegisterGetPage(s *mcp.Server, c *libpuki.Client) {
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

		source, err := c.GetPageSource(input.PageName)
		if err != nil {
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
