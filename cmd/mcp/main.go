package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/tools"
)

func main() {
	baseURL := os.Getenv("PUKIWIKI_URL")
	if baseURL == "" {
		fmt.Fprintln(os.Stderr, "PUKIWIKI_URL is required")
		os.Exit(1)
	}

	opts := []libpuki.Option{
		libpuki.WithAuth(os.Getenv("PUKIWIKI_USER"), os.Getenv("PUKIWIKI_PASS")),
	}
	if scope := os.Getenv("PUKIWIKI_SCOPE"); scope != "" {
		opts = append(opts, libpuki.WithScope(scope))
	}

	client, err := libpuki.New(baseURL, opts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create client: %v\n", err)
		os.Exit(1)
	}

	if err := client.Login(); err != nil {
		if errors.Is(err, libpuki.ErrAuthFailed) {
			fmt.Fprintln(os.Stderr, "authentication failed")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "login error: %v\n", err)
		os.Exit(1)
	}

	s := mcp.NewServer(
		&mcp.Implementation{Name: "pukiwiki-mcp", Version: "1.0.0"},
		nil,
	)

	tools.RegisterListPages(s, client)
	tools.RegisterGetPage(s, client)
	tools.RegisterGetPageInfo(s, client)
	tools.RegisterSearchPages(s, client)
	tools.RegisterCreatePage(s, client)
	tools.RegisterEditPage(s, client)

	if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
