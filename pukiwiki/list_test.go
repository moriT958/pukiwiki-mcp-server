package pukiwiki

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestListPages(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	t.Run("スコープなしで全ページが返る", func(t *testing.T) {
		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		pages, err := c.ListPages()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{"FrontPage", "SandBox", "test-pages", "test-pages/notes"}
		if len(pages) != len(want) {
			t.Fatalf("len(pages) = %d, want %d: got %v", len(pages), len(want), pages)
		}
		for _, w := range want {
			if !slices.Contains(pages, w) {
				t.Errorf("pages does not contain %q", w)
			}
		}
	})

	t.Run("WithScopeでプレフィックス配下のみ返る", func(t *testing.T) {
		c, err := New(ts.URL, WithScope("test-pages"))
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		pages, err := c.ListPages()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := []string{"test-pages", "test-pages/notes"}
		if len(pages) != len(want) {
			t.Fatalf("len(pages) = %d, want %d: got %v", len(pages), len(want), pages)
		}
		for _, w := range want {
			if !slices.Contains(pages, w) {
				t.Errorf("pages does not contain %q", w)
			}
		}
	})

	t.Run("HTTPエラー時はエラーを返す", func(t *testing.T) {
		errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errServer.Close()

		c, err := New(errServer.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = c.ListPages()
		if err == nil {
			t.Fatal("expected error for HTTP 500, got nil")
		}
	})
}
