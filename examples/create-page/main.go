package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	libpuki "github.com/moriT958/pukiwiki-mcp"
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

	pageName := "my-class-pages/aaa"
	content := "* 見出し\n\nページの本文です。\n"

	if err := client.CreatePage(pageName, content); err != nil {
		if errors.Is(err, libpuki.ErrPageAlreadyExists) {
			log.Fatalf("Page %q already exists", pageName)
		}
		if errors.Is(err, libpuki.ErrOutOfScope) {
			log.Fatalf("Page %q is outside the configured scope", pageName)
		}
		log.Fatalf("Failed to create page: %v", err)
	}

	fmt.Printf("Page %q created successfully.\n", pageName)
}
