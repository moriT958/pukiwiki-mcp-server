package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zalando/go-keyring"
)

func init() {
	keyring.MockInit()
}

// ログイン成功をモックするテストサーバー
func newLoginServer(t *testing.T) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Query().Get("plugin") == "loginform" {
			if err := r.ParseForm(); err != nil {
				t.Fatalf("failed to init mock server: %v", err)
			}

			// usename="testuser", password="testpass" で認証成功
			if r.FormValue("username") == "testuser" && r.FormValue("password") == "testpass" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
}

// Provider が libpuki.Client をキャッシュするか検証
func TestProvider_Get_WithStoredCredentials(t *testing.T) {
	keyring.MockInit()

	srv := newLoginServer(t)
	defer srv.Close()

	cfg := &config{
		URL:      srv.URL,
		Username: "testuser",
		Password: "testpass",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	p := &Provider{}
	client, err := p.Get(context.Background())
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if client == nil {
		t.Fatal("Get() returned nil client")
	}
	client2, err := p.Get(context.Background())
	if err != nil {
		t.Fatalf("Get() 2nd call error = %v", err)
	}

	if client != client2 {
		t.Error("Get() 2nd call should return cached client")
	}
}

// Provider.Reset() で Keychain が削除されるか検証
func TestProvider_Reset(t *testing.T) {
	keyring.MockInit()

	srv := newLoginServer(t)
	defer srv.Close()

	cfg := &config{URL: srv.URL, Username: "testuser", Password: "testpass"}
	_ = Save(cfg)

	p := &Provider{}
	_, _ = p.Get(context.Background())

	if err := p.Reset(); err != nil {
		t.Errorf("failed to run Provider.Reset(): %v", err)
	}

	_, err := load()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("after Reset(), Load() error = %v, want ErrNotFound", err)
	}
}
