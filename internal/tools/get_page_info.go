package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
)

type GetPageInfoInput struct {
	PageName string `json:"page_name" jsonschema:"Name of the wiki page"`
}

func RegisterGetPageInfo(s *mcp.Server, p *auth.Provider) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_page_info",
		Description: "PukiWiki ページの存在確認と最終更新日時を返す",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetPageInfoInput) (*mcp.CallToolResult, any, error) {
		if input.PageName == "" {
			return errResult("page_name is required")
		}

		c, err := p.Get(ctx)
		if err != nil {
			return errResult(fmt.Sprintf("auth error: %v", err))
		}

		info, err := c.GetPageInfo(input.PageName)
		if err != nil {
			return handlePukiwikiErr(p, err, input.PageName, "get_page_info")
		}

		result, err := json.Marshal(map[string]any{
			"page_name":     info.Name,
			"exists":        info.Exists,
			"last_modified": info.LastModified,
		})
		if err != nil {
			return errResult(fmt.Sprintf("marshal failed: %v", err))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(result)}},
		}, nil, nil
	})
}
