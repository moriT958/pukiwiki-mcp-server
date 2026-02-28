package libpuki

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestSearchPages(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	t.Run("マッチした結果が返る", func(t *testing.T) {
		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		results, err := c.SearchPages("keyword", MatchAll)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := map[string]bool{"test-pages/page1": false, "test-pages/page2": false}
		if len(results) != len(want) {
			t.Fatalf("len(results) = %d, want %d: got %v", len(results), len(want), results)
		}
		for _, r := range results {
			if _, ok := want[r.Name]; !ok {
				t.Errorf("unexpected page %q in results", r.Name)
			}
			if r.Body == "" {
				t.Errorf("expected non-empty body for %q", r.Name)
			}
		}
	})

	t.Run("0 件の場合は空スライスを返す", func(t *testing.T) {
		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		results, err := c.SearchPages("nomatch", MatchAll)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 0 {
			t.Errorf("expected empty slice, got %v", results)
		}
	})

	t.Run("MatchAny で OR 検索のクエリが送られる", func(t *testing.T) {
		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		results, err := c.SearchPages("keyword other", MatchAny)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("len(results) = %d, want 2: got %v", len(results), results)
		}
	})

	t.Run("WithScope でスコープ外のページが除外される", func(t *testing.T) {
		c, err := New(ts.URL, WithScope("test-pages/page1"))
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		results, err := c.SearchPages("keyword", MatchAll)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 1 || results[0].Name != "test-pages/page1" {
			t.Errorf("expected [test-pages/page1], got %v", results)
		}
	})

	t.Run("HTTP エラー時はエラーを返す", func(t *testing.T) {
		errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errServer.Close()

		c, err := New(errServer.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = c.SearchPages("keyword", MatchAll)
		if err == nil {
			t.Fatal("expected error for HTTP 500, got nil")
		}
	})

	t.Run("ページネーションで全結果が取得できる", func(t *testing.T) {
		paginationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			values := r.URL.Query()
			start, _ := strconv.Atoi(values.Get("start"))

			var resp map[string]any
			if start == 0 {
				resp = map[string]any{
					"search_done":      false,
					"next_start_index": 1,
					"results": []map[string]any{
						{"name": "test-pages/page1", "updated_at": "2025-01-01T00:00:00+09:00", "body": "content1"},
					},
				}
			} else {
				resp = map[string]any{
					"search_done":      true,
					"next_start_index": 2,
					"results": []map[string]any{
						{"name": "test-pages/page2", "updated_at": "2025-01-02T00:00:00+09:00", "body": "content2"},
					},
				}
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer paginationServer.Close()

		c, err := New(paginationServer.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		results, err := c.SearchPages("keyword", MatchAll)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(results) != 2 {
			t.Fatalf("expected 2 results from pagination, got %d: %v", len(results), results)
		}
	})
}
