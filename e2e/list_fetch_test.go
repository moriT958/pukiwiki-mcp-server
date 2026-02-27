package e2e

import (
	"os"
	"testing"

	"github.com/moriT958/pukiwiki-mcp"
)

func TestListThenGetSource(t *testing.T) {
	baseURL := os.Getenv("PUKIWIKI_URL")
	user := os.Getenv("PUKIWIKI_USER")
	pass := os.Getenv("PUKIWIKI_PASS")
	scope := os.Getenv("PUKIWIKI_SCOPE")

	if baseURL == "" {
		t.Skip("PUKIWIKI_URL not set")
	}

	c, err := libpuki.New(baseURL,
		libpuki.WithAuth(user, pass),
		libpuki.WithScope(scope),
	)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if err := c.Login(); err != nil {
		t.Fatalf("Login: %v", err)
	}

	pages, err := c.ListPages()
	if err != nil {
		t.Fatalf("ListPages: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("ListPages returned no pages")
	}
	t.Logf("Found %d page(s): %v", len(pages), pages)

	source, err := c.GetPageSource(pages[3])
	if err != nil {
		t.Fatalf("GetPageSource(%q): %v", pages[0], err)
	}
	if source == "" {
		t.Errorf("GetPageSource(%q) returned empty source", pages[0])
	}
	t.Logf("Source of %q:\n%s", pages[0], source)
}
