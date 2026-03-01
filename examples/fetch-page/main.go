package main

import (
	"fmt"
	"log"

	"github.com/moriT958/pukiwiki-mcp/pukiwiki"
)

func main() {
	client, err := pukiwiki.New("https://base-url.jp",
		pukiwiki.WithAuth("username", "password"),
	)
	if err != nil {
		log.Fatalf("Failed to init pukiwiki client: %v", err)
	}

	if err := client.Login(); err != nil {
		log.Fatalf("Failed to login my pukiwiki site: %v", err)
	}

	// base-url/?<pagename>
	pageSrc, err := client.GetPageSource("pagename")
	if err != nil {
		log.Printf("Failed to get page source: %v", err)
	}

	pageInfo, err := client.GetPageInfo("pagename")
	if err != nil {
		log.Printf("Failed to get page info: %v", err)
	}

	fmt.Println("===== Page Source =====")
	fmt.Println(pageSrc)

	fmt.Println("===== Page Info =====")
	fmt.Println(pageInfo)
}
