package main

import (
	"fmt"
	"log"
	"os"

	"github.com/moriT958/pukiwiki-mcp"
)

func main() {
	pukiwiki := os.Getenv("PUKIWIKI_URL")
	user := os.Getenv("PUKIWIKI_USER")
	pass := os.Getenv("PUKIWIKI_PASS")
	scope := os.Getenv("PUKIWIKI_SCOPE")

	client, err := libpuki.New(pukiwiki,
		libpuki.WithAuth(user, pass),
		libpuki.WithScope(scope),
	)
	if err != nil {
		log.Fatalf("Failed to init pukiwiki client: %v", err)
	}

	if err := client.Login(); err != nil {
		log.Fatalf("Failed to login: %v", err)
	}

	pages, err := client.ListPages()
	if err != nil {
		log.Fatalf("Failed to list pages: %v", err)
	}

	fmt.Printf("Found %d page(s):\n", len(pages))
	for _, p := range pages {
		fmt.Println(" -", p)
	}
}
