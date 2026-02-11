package libpuki

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

// 指定パスの HTML フィクスチャを読み込みレスポンスとして書き込む
func serveFixture(t *testing.T, w http.ResponseWriter, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", path, err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

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
