package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
			return errResult("page_name is required")
		}
		if input.NewContent == "" {
			return errResult("new_content is required")
		}

		c, err := p.Get(ctx)
		if err != nil {
			return errResult(fmt.Sprintf("auth error: %v", err))
		}

		if err := c.EditPage(input.PageName, input.NewContent); err != nil {
			return handlePukiwikiErr(p, err, input.PageName, "edit_page")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("page %q updated successfully", input.PageName)}},
		}, nil, nil
	})
}
