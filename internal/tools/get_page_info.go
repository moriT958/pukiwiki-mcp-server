package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/libpuki"
)

type GetPageInfoInput struct {
	PageName string `json:"page_name" jsonschema:"Name of the wiki page"`
}

func RegisterGetPageInfo(s *mcp.Server, c *libpuki.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_page_info",
		Description: "PukiWiki ページの存在確認と最終更新日時を返す",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetPageInfoInput) (*mcp.CallToolResult, any, error) {
		if input.PageName == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "page_name is required"}},
				IsError: true,
			}, nil, nil
		}

		info, err := c.GetPageInfo(input.PageName)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("get_page_info failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		result, err := json.Marshal(map[string]any{
			"page_name":     info.Name,
			"exists":        info.Exists,
			"last_modified": info.LastModified,
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
