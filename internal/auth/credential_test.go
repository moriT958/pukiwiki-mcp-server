package auth

import (
	"errors"
	"testing"
)

// 認証情報が正しく保存されるか検証
func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	want := &config{
		URL:      "https://wiki.example.com",
		Username: "testuser",
		Password: "testpass",
		Scope:    "user/testuser",
	}

	if err := Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got.URL != want.URL || got.Username != want.Username ||
		got.Password != want.Password || got.Scope != want.Scope {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

// 認証情報が未登録の状態でロードすると ErrNotFound が返るか検証
func TestLoad_NotFound(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	_, err := load()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Load() error = %v, want ErrNotFound", err)
	}
}

// 認証情報の削除後は ErrNotFound が返るか検証
func TestDelete(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	cfg := &config{URL: "https://wiki.example.com", Username: "u", Password: "p"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := delete(); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := load()
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Load() after Delete() error = %v, want ErrNotFound", err)
	}
}

// 存在しないエントリを削除してもエラーとして扱わない
func TestDelete_NotFound(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	if err := delete(); err != nil {
		t.Errorf("Delete() on empty config error = %v, want nil", err)
	}
}

// TODO: これは仕様的にどうなの？ (Scope は必須でもいいかも)
// Scope 未指定の場合は空で保存されるか検証
func TestSave_NoScope(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("XDG_CONFIG_HOME", tmp)

	want := &config{
		URL:      "https://wiki.example.com",
		Username: "u",
		Password: "p",
		// Scope は空
	}
	if err := Save(want); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got.Scope != "" {
		t.Errorf("Load().Scope = %q, want empty", got.Scope)
	}
}
