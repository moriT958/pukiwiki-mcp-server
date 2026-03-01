package pukiwiki

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

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
