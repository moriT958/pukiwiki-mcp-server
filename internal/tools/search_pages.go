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

type SearchPagesInput struct {
	Query     string `json:"query"                jsonschema:"Search query string"`
	MatchType string `json:"match_type,omitempty" jsonschema:"Match type: AND or OR (default: AND)"`
}

type searchPagesOutput struct {
	Query     string                 `json:"query"`
	MatchType string                 `json:"match_type"`
	Count     int                    `json:"count"`
	Results   []libpuki.SearchResult `json:"results"`
}

func RegisterSearchPages(s *mcp.Server, p *auth.Provider) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_pages",
		Description: "キーワードを含むページを検索し、ページ名・更新日時・本文を返す。match_type で AND/OR 検索を選択できる (Default: AND)。",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchPagesInput) (*mcp.CallToolResult, any, error) {
		if input.Query == "" {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: "query is required"}},
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

		matchType := libpuki.MatchAll
		if input.MatchType == string(libpuki.MatchAny) {
			matchType = libpuki.MatchAny
		}

		results, err := c.SearchPages(input.Query, matchType)
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
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("search_pages failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		out, err := json.Marshal(searchPagesOutput{
			Query:     input.Query,
			MatchType: string(matchType),
			Count:     len(results),
			Results:   results,
		})
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("marshal failed: %v", err)}},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(out)}},
		}, nil, nil
	})
}
