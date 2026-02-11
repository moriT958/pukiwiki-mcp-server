package libpuki

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// PukiWiki のログインフォームをモックする httptest.Server を返す
// username=testuser&password=testpass なら 302 + Set-Cookie を返す
// それ以外は 200 + login_failure.html を返す。
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

func TestLogin(t *testing.T) {
	t.Run("正しい認証情報でログイン成功", func(t *testing.T) {
		ts := newAuthTestServer(t)
		defer ts.Close()

		c, err := New(ts.URL, WithHTTPClient(&http.Client{}), WithAuth("testuser", "testpass"))
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		if err := c.Login(); err != nil {
			t.Fatalf("Login() returned unexpected error: %v", err)
		}

		// Cookie が保存されたことを検証
		tsURL, _ := url.Parse(ts.URL)
		cookies := c.httpClient.Jar.Cookies(tsURL)
		var found bool
		for _, cookie := range cookies {
			if cookie.Name == "PHPSESSID" {
				found = true
				break
			}
		}
		if !found {
			t.Error("PHPSESSID cookie not found after successful Login()")
		}
	})

	t.Run("誤った認証情報でErrAuthFailedを返す", func(t *testing.T) {
		ts := newAuthTestServer(t)
		defer ts.Close()

		c, err := New(ts.URL, WithHTTPClient(&http.Client{}), WithAuth("wrong", "creds"))
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		err = c.Login()
		if !errors.Is(err, ErrAuthFailed) {
			t.Errorf("Login() err = %v, want ErrAuthFailed", err)
		}
	})

	t.Run("認証情報なしの場合は何もせずnilを返す", func(t *testing.T) {
		c, err := New("https://example.com")
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}

		if err := c.Login(); err != nil {
			t.Fatalf("Login() returned unexpected error: %v", err)
		}
	})
}
