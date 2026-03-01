package tools

import (
	"context"
	"encoding/json"
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
			return errResult("query is required")
		}

		c, err := p.Get(ctx)
		if err != nil {
			return errResult(fmt.Sprintf("auth error: %v", err))
		}

		matchType := libpuki.MatchAll
		if input.MatchType == string(libpuki.MatchAny) {
			matchType = libpuki.MatchAny
		}

		results, err := c.SearchPages(input.Query, matchType)
		if err != nil {
			return handlePukiwikiErr(p, err, "", "search_pages")
		}

		out, err := json.Marshal(searchPagesOutput{
			Query:     input.Query,
			MatchType: string(matchType),
			Count:     len(results),
			Results:   results,
		})
		if err != nil {
			return errResult(fmt.Sprintf("marshal failed: %v", err))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(out)}},
		}, nil, nil
	})
}
