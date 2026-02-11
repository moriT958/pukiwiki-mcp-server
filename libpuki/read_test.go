package libpuki

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

// testdata/ の HTML フィクスチャに基づいて PukiWiki のレスポンスをモックする httptest.Server を返す
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.RawQuery

		// cmd=source&page=<pagename> でソース取得
		if values, err := url.ParseQuery(q); err == nil && values.Get("cmd") == "source" {
			page := values.Get("page")
			switch page {
			case "seminar-personal/morita2023":
				serveFixture(t, w, "testdata/source_page.html")
			default:
				serveFixture(t, w, "testdata/source_page_empty.html")
			}
			return
		}

		// ?<pagename> でページ表示
		pageName, err := url.QueryUnescape(q)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch pageName {
		case "seminar-personal/morita2023":
			serveFixture(t, w, "testdata/page_exists.html")
		default:
			serveFixture(t, w, "testdata/page_not_found.html")
		}
	}))
}

func serveFixture(t *testing.T, w http.ResponseWriter, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

func TestGetPageSource(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	c, err := New(ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Run("存在するページのソースを取得できる", func(t *testing.T) {
		source, err := c.GetPageSource("seminar-personal/morita2023")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := "* 見出し\n- リスト1\n- リスト2\n\nテスト本文です。"
		if source != want {
			t.Errorf("source = %q, want %q", source, want)
		}
	})

	t.Run("存在しないページはErrPageNotFoundを返す", func(t *testing.T) {
		_, err := c.GetPageSource("nonexistent")
		if !errors.Is(err, ErrPageNotFound) {
			t.Errorf("err = %v, want ErrPageNotFound", err)
		}
	})

	t.Run("HTTPエラー時はエラーを返す", func(t *testing.T) {
		errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errServer.Close()

		ec, err := New(errServer.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = ec.GetPageSource("any-page")
		if err == nil {
			t.Fatal("expected error for HTTP 500, got nil")
		}
	})
}

func TestGetPageInfo(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	c, err := New(ts.URL)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	t.Run("存在するページの情報を取得できる", func(t *testing.T) {
		info, err := c.GetPageInfo("seminar-personal/morita2023")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Name != "seminar-personal/morita2023" {
			t.Errorf("Name = %q, want %q", info.Name, "seminar-personal/morita2023")
		}
		if !info.Exists {
			t.Error("Exists = false, want true")
		}
		if info.LastModified != "2024-01-15 (Mon) 10:30:00" {
			t.Errorf("LastModified = %q, want %q", info.LastModified, "2024-01-15 (Mon) 10:30:00")
		}
	})

	t.Run("存在しないページはExists=falseで返る", func(t *testing.T) {
		info, err := c.GetPageInfo("nonexistent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.Name != "nonexistent" {
			t.Errorf("Name = %q, want %q", info.Name, "nonexistent")
		}
		if info.Exists {
			t.Error("Exists = true, want false")
		}
		if info.LastModified != "" {
			t.Errorf("LastModified = %q, want empty", info.LastModified)
		}
	})

	t.Run("HTTPエラー時はエラーを返す", func(t *testing.T) {
		errServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer errServer.Close()

		ec, err := New(errServer.URL)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		_, err = ec.GetPageInfo("any-page")
		if err == nil {
			t.Fatal("expected error for HTTP 500, got nil")
		}
	})
}
