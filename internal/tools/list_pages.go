package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	libpuki "github.com/moriT958/pukiwiki-mcp"
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
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("auth error: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		pages, err := c.ListPages()
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
