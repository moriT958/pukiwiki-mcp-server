package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
			return errResult("page_name is required")
		}

		c, err := p.Get(ctx)
		if err != nil {
			return errResult(fmt.Sprintf("auth error: %v", err))
		}

		source, err := c.GetPageSource(input.PageName)
		if err != nil {
			return handlePukiwikiErr(p, err, input.PageName, "get_page")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: source}},
		}, nil, nil
	})
}
