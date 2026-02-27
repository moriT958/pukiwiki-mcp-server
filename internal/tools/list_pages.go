package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp"
)

type ListPagesInput struct{}

func RegisterListPages(s *mcp.Server, c *libpuki.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_pages",
		Description: "PukiWiki 内のページ一覧を返す。WithScope が設定されている場合はそのプレフィックス配下のみ。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListPagesInput) (*mcp.CallToolResult, any, error) {
		pages, err := c.ListPages()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("list_pages failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		result, err := json.Marshal(map[string]any{
			"pages": pages,
			"count": len(pages),
		})
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("marshal failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(result)}},
		}, nil, nil
	})
}
