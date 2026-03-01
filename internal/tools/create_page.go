package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
)

type CreatePageInput struct {
	PageName string `json:"page_name" jsonschema:"Name of the wiki page to create"`
	Content  string `json:"content"   jsonschema:"Wiki source content for the new page"`
}

func RegisterCreatePage(s *mcp.Server, p *auth.Provider) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_page",
		Description: "PukiWiki に新規ページを作成する",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CreatePageInput) (*mcp.CallToolResult, any, error) {
		if input.PageName == "" {
			return errResult("page_name is required")
		}
		if input.Content == "" {
			return errResult("content is required")
		}

		c, err := p.Get(ctx)
		if err != nil {
			return errResult(fmt.Sprintf("auth error: %v", err))
		}

		err = c.CreatePage(input.PageName, input.Content)
		if err != nil {
			return handlePukiwikiErr(p, err, input.PageName, "create_page")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q created successfully", input.PageName)}},
		}, nil, nil
	})
}
