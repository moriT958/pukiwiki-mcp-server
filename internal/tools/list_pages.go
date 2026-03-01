package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
)

type ListPagesInput struct{}

func RegisterListPages(s *mcp.Server, p *auth.Provider) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pages",
		Description: "PukiWiki 内のページ一覧を返す。WithScope が設定されている場合はそのプレフィックス配下のみ。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListPagesInput) (*mcp.CallToolResult, any, error) {
		c, err := p.Get(ctx)
		if err != nil {
			return errResult(fmt.Sprintf("auth error: %v", err))
		}

		pages, err := c.ListPages()
		if err != nil {
			return handlePukiwikiErr(p, err, "", "list_pages")
		}

		result, err := json.Marshal(map[string]any{
			"pages": pages,
			"count": len(pages),
		})
		if err != nil {
			return errResult(fmt.Sprintf("marshal failed: %v", err))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(result)}},
		}, nil, nil
	})
}
