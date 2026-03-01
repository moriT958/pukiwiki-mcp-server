package pukiwiki

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestEditPage(t *testing.T) {
	t.Run("既存ページを編集できる", func(t *testing.T) {
		var receivedForm url.Values
		ts := newWriteTestServer(t, true, func(form url.Values) {
			receivedForm = form
		})
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := c.EditPage("existing-page", "* 更新内容\nコンテンツ"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if receivedForm.Get("cmd") != "edit" {
			t.Errorf("cmd = %q, want %q", receivedForm.Get("cmd"), "edit")
		}
		if receivedForm.Get("write") != "true" {
			t.Errorf("write = %q, want %q", receivedForm.Get("write"), "true")
		}
		if receivedForm.Get("msg") != "* 更新内容\nコンテンツ" {
			t.Errorf("msg = %q, want %q", receivedForm.Get("msg"), "* 更新内容\nコンテンツ")
		}
	})

	t.Run("存在しないページは ErrPageNotFound を返す", func(t *testing.T) {
		ts := newWriteTestServer(t, false, func(form url.Values) {
			t.Error("POST should not be called when page does not exist")
		})
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = c.EditPage("new-page", "内容")
		if !errors.Is(err, ErrPageNotFound) {
			t.Errorf("err = %v, want ErrPageNotFound", err)
		}
	})

	t.Run("スコープ外は ErrOutOfScope を返す", func(t *testing.T) {
		c, err := New("https://example.com", WithScope("scoped"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = c.EditPage("other/page", "内容")
		if !errors.Is(err, ErrOutOfScope) {
			t.Errorf("err = %v, want ErrOutOfScope", err)
		}
	})

	t.Run("編集フォームの GET が HTTP エラーの場合はエラーを返す", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := c.EditPage("existing-page", "内容"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestCreatePage(t *testing.T) {
	t.Run("新規ページを作成できる", func(t *testing.T) {
		var receivedForm url.Values
		ts := newWriteTestServer(t, false, func(form url.Values) {
			receivedForm = form
		})
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := c.CreatePage("new-page", "* 新しいページ\nコンテンツ"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if receivedForm.Get("cmd") != "edit" {
			t.Errorf("cmd = %q, want %q", receivedForm.Get("cmd"), "edit")
		}
		if receivedForm.Get("write") != "true" {
			t.Errorf("write = %q, want %q", receivedForm.Get("write"), "true")
		}
		if receivedForm.Get("digest") != "d41d8cd98f00b204e9800998ecf8427e" {
			t.Errorf("digest = %q, want %q", receivedForm.Get("digest"), "d41d8cd98f00b204e9800998ecf8427e")
		}
		if receivedForm.Get("msg") != "* 新しいページ\nコンテンツ" {
			t.Errorf("msg = %q, want %q", receivedForm.Get("msg"), "* 新しいページ\nコンテンツ")
		}
	})

	t.Run("既存ページはErrPageAlreadyExistsを返す", func(t *testing.T) {
		ts := newWriteTestServer(t, true, func(form url.Values) {
			t.Error("POST should not be called when page already exists")
		})
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = c.CreatePage("existing-page", "内容")
		if !errors.Is(err, ErrPageAlreadyExists) {
			t.Errorf("err = %v, want ErrPageAlreadyExists", err)
		}
	})

	t.Run("スコープ外は ErrOutOfScope を返す", func(t *testing.T) {
		c, err := New("https://example.com", WithScope("scoped"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = c.CreatePage("other/page", "内容")
		if !errors.Is(err, ErrOutOfScope) {
			t.Errorf("err = %v, want ErrOutOfScope", err)
		}
	})

	t.Run("編集フォームの GET が HTTP エラーの場合はエラーを返す", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := c.CreatePage("new-page", "内容"); err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("POST が 200 を返した場合はエラーを返す", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				serveFixture(t, w, "testdata/edit_form_new_page.html")
				return
			}
			// 200 = 書き込み失敗
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		c, err := New(ts.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := c.CreatePage("new-page", "内容"); err == nil {
			t.Error("expected error, got nil")
		}
	})
}
