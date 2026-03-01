package tools

import (
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
	"github.com/moriT958/pukiwiki-mcp/pukiwiki"
)

// MCP Tool のエラー結果を生成する
func errResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
		IsError: true,
	}, nil, nil
}

func handlePukiwikiErr(p *auth.Provider, err error, pageName, toolName string) (*mcp.CallToolResult, any, error) {
	switch {
	case errors.Is(err, pukiwiki.ErrSessionExpired):
		if resetErr := p.Reset(); resetErr != nil {
			return errResult(fmt.Sprintf("session expired but failed to clear credentials: %v. please retry.", resetErr))
		}
		return errResult("session expired. setup wizard launched. please retry after login.")
	case errors.Is(err, pukiwiki.ErrPageNotFound):
		return errResult(fmt.Sprintf("page %q not found", pageName))
	case errors.Is(err, pukiwiki.ErrPageAlreadyExists):
		return errResult(fmt.Sprintf("page %q already exists", pageName))
	case errors.Is(err, pukiwiki.ErrOutOfScope):
		return errResult(fmt.Sprintf("page %q is outside the configured scope", pageName))
	default:
		return errResult(fmt.Sprintf("%s failed: %v", toolName, err))
	}
}
