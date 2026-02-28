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

		// cmd=list でページ一覧取得
		if values, err := url.ParseQuery(q); err == nil && values.Get("cmd") == "list" {
			serveFixture(t, w, "testdata/list_pages.html")
			return
		}

		// cmd=search2&action=query で JSON 検索結果を返す
		if values, err := url.ParseQuery(q); err == nil && values.Get("cmd") == "search2" && values.Get("action") == "query" {
			switch values.Get("q") {
			case "keyword", "keyword OR other":
				serveFixture(t, w, "testdata/search_results.json")
			default:
				serveFixture(t, w, "testdata/search_no_results.json")
			}
			return
		}

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

// PukiWiki の編集フォームと書き込みをモックする httptest.Server を返す
func newWriteTestServer(t *testing.T, pageExists bool, onWrite func(url.Values)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if pageExists {
				serveFixture(t, w, "testdata/edit_form_existing_page.html")
			} else {
				serveFixture(t, w, "testdata/edit_form_new_page.html")
			}
			return
		}
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if onWrite != nil {
				onWrite(r.Form)
			}
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
}

// PukiWiki のログインフォームをモックする httptest.Server を返す
func newAuthTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Query().Get("plugin") != "loginform" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "testuser" && password == "testpass" {
			http.SetCookie(w, &http.Cookie{
				Name:  "PHPSESSID",
				Value: "testsession",
				Path:  "/",
			})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		serveFixture(t, w, "testdata/login_failure.html")
	}))
}
