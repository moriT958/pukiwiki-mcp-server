package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

func init() {
	keyring.MockInit()
}

func newTestMux(t *testing.T) http.Handler {
	t.Helper()

	setupTmpl, doneTmpl, err := parseTemplates()
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	done := make(chan struct{}, 1)
	result := make(chan *config, 1)
	return newServeMux(setupTmpl, doneTmpl, done, result)
}

// セットアップウィザードの HTML が正しいか検証
func TestWizardHandler_GetSetup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	newTestMux(t).ServeHTTP(w, req)

	// assert: status 200
	if w.Code != http.StatusOK {
		t.Errorf("GET / status = %d, want %d", w.Code, http.StatusOK)
	}
	// assert: has form tag
	if !strings.Contains(w.Body.String(), "<form") {
		t.Error("GET / response does not contain form")
	}
}

// submit が正しいレスポンスを返すか検証
func TestWizardHandler_ValidSubmit(t *testing.T) {
	keyring.MockInit()

	form := url.Values{
		"url":      {"https://wiki.example.com"},
		"username": {"testuser"},
		"password": {"testpass"},
		"scope":    {"user/testuser"},
	}

	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	newTestMux(t).ServeHTTP(w, req)

	// assert: status 303
	if w.Code != http.StatusSeeOther {
		t.Errorf("POST /submit status = %d, want %d", w.Code, http.StatusSeeOther)
	}
	// assert: redirect /done
	if loc := w.Header().Get("Location"); loc != "/done" {
		t.Errorf("Location = %q, want /done", loc)
	}
}

// url が空の場合のレスポンスを検証
func TestWizardHandler_MissingURL(t *testing.T) {
	form := url.Values{
		"url":      {""},
		"username": {"testuser"},
		"password": {"pass"},
	}

	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	newTestMux(t).ServeHTTP(w, req)

	// assert: status 200 (リダイレクトせずに再表示されるか)
	if w.Code != http.StatusOK {
		t.Errorf("POST /submit (missing url) status = %d, want %d", w.Code, http.StatusOK)
	}
	// assert: has validation message in response body
	if !strings.Contains(w.Body.String(), "URL とユーザー名は必須です。") {
		t.Error("response does not contain validation error message")
	}
}

// /done へのリクエストを検証
func TestWizardHandler_GetDone(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/done", nil)
	w := httptest.NewRecorder()
	newTestMux(t).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /done status = %d, want %d", w.Code, http.StatusOK)
	}
}
