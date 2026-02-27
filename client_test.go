package libpuki

import (
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("有効なbaseURLでClientが生成される", func(t *testing.T) {
		c, err := New("https://example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.baseURL != "https://example.com" {
			t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
		}
	})

	t.Run("末尾スラッシュが除去される", func(t *testing.T) {
		c, err := New("https://example.com/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.baseURL != "https://example.com" {
			t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
		}
	})

	t.Run("空のbaseURLはエラーを返す", func(t *testing.T) {
		_, err := New("")
		if err == nil {
			t.Fatal("expected error for empty baseURL, got nil")
		}
	})
}

func TestOption(t *testing.T) {
	t.Run("WithTimeoutでタイムアウトが設定される", func(t *testing.T) {
		c, err := New("https://example.com", WithTimeout(5*time.Second))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient.Timeout != 5*time.Second {
			t.Errorf("timeout = %v, want %v", c.httpClient.Timeout, 5*time.Second)
		}
	})

	t.Run("WithHTTPClientでカスタムクライアントが注入される", func(t *testing.T) {
		custom := &http.Client{Timeout: 99 * time.Second}
		c, err := New("https://example.com", WithHTTPClient(custom))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient != custom {
			t.Error("httpClient was not set to the custom client")
		}
	})

	t.Run("WithAuthでusernameとpasswordがセットされる", func(t *testing.T) {
		c, err := New("https://example.com", WithAuth("user", "pass"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.username != "user" {
			t.Errorf("username = %q, want %q", c.username, "user")
		}
		if c.password != "pass" {
			t.Errorf("password = %q, want %q", c.password, "pass")
		}
	})

	t.Run("NewでデフォルトCookieJarが設定される", func(t *testing.T) {
		c, err := New("https://example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient.Jar == nil {
			t.Error("CookieJar is nil, want non-nil")
		}
	})
}
