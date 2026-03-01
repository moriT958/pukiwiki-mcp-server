package main

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
	"github.com/moriT958/pukiwiki-mcp/internal/tools"
)

var version = "dev"

func main() {
	provider := &auth.Provider{}

	s := mcp.NewServer(
		&mcp.Implementation{Name: "pukiwiki-mcp", Version: version},
		nil,
	)

	tools.RegisterListPages(s, provider)
	tools.RegisterGetPage(s, provider)
	tools.RegisterGetPageInfo(s, provider)
	tools.RegisterSearchPages(s, provider)
	tools.RegisterCreatePage(s, provider)
	tools.RegisterEditPage(s, provider)

	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
