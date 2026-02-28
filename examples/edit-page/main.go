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

	pageName := "my-classpages/aaa"
	newContent := "* 見出し（更新）\n\n更新後の本文です。\n"

	if err := client.EditPage(pageName, newContent); err != nil {
		if errors.Is(err, libpuki.ErrPageNotFound) {
			log.Fatalf("Page %q not found", pageName)
		}
		if errors.Is(err, libpuki.ErrOutOfScope) {
			log.Fatalf("Page %q is outside the configured scope", pageName)
		}
		log.Fatalf("Failed to edit page: %v", err)
	}

	fmt.Printf("Page %q updated successfully.\n", pageName)
}
