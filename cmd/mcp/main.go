package main

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/moriT958/pukiwiki-mcp/internal/auth"
	"github.com/moriT958/pukiwiki-mcp/internal/tools"
)

var version = "dev" // GoReleaser の ldflags で上書き

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "setup":
			runSetup()
			return
		case "--version", "-v", "version":
			fmt.Printf("pukiwiki-mcp %s\n", version)
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", os.Args[1])
			os.Exit(1)
		}
	}

	runServer()
}

// runSetup は認証情報ウィザードを起動して Keychain に保存する。
// パスワード変更時など手動で再設定したい場合に使う。
func runSetup() {
	cfg, err := auth.(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "setup wizard failed: %v\n", err)
		os.Exit(1)
	}
	if err := auth.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save credentials: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Credentials updated.")
}

// runServer は MCP サーバーをスタートする。
// 認証情報は最初のツール呼び出し時に auth.Provider 経由で取得する。
func runServer() {
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
