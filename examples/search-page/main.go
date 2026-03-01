package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moriT958/pukiwiki-mcp/pukiwiki"
)

func main() {
	pukiwikiURL := os.Getenv("PUKIWIKI_URL")
	user := os.Getenv("PUKIWIKI_USER")
	pass := os.Getenv("PUKIWIKI_PASS")
	scope := os.Getenv("PUKIWIKI_SCOPE")

	client, err := pukiwiki.New(pukiwikiURL,
		pukiwiki.WithAuth(user, pass),
		pukiwiki.WithScope(scope),
	)
	if err != nil {
		log.Fatalf("Failed to init pukiwiki client: %v", err)
	}

	if err := client.Login(); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	// AND 検索
	results, err := client.SearchPages("word1 word2", pukiwiki.MatchAll)
	if err != nil {
		log.Fatalf("Failed to search pages: %v", err)
	}

	fmt.Printf("AND search: Found %d page(s):\n", len(results))
	for _, r := range results {
		fmt.Printf(" - %s (updated: %s)\n", r.Name, r.UpdatedAt)
	}

	// OR 検索
	results, err = client.SearchPages("word1 word2", pukiwiki.MatchAny)
	if err != nil {
		log.Fatalf("Failed to search pages: %v", err)
	}

	fmt.Printf("OR search: Found %d page(s):\n", len(results))
	for _, r := range results {
		fmt.Printf(" - %s (updated: %s)\n", r.Name, r.UpdatedAt)
	}
}
