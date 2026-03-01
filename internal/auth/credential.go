package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

// MacOS Keychain に保存するエントリの識別情報
const (
	serviceName = "pukiwiki-mcp"
	accountName = "config"
)

var (
	// MacOS Keychain に認証情報が存在しない場合に返す
	ErrNotFound = errors.New("credentials not found in keychain")
)

// libpuki.Client の初期化に必要な認証情報を保持する
type config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Scope    string `json:"scope,omitempty"`
}

// macOS Keychain から認証情報を読み込む
func load() (*config, error) {
	raw, err := keyring.Get(serviceName, accountName)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("keychain read failed: %w", err)
	}

	var cfg config
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, fmt.Errorf("credential parse failed: %w", err)
	}

	return &cfg, nil
}

// 認証情報を JSON エンコードして macOS Keychain に保存する
func Save(cfg *config) error {
	raw, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("credential marshal failed: %w", err)
	}

	if err := keyring.Set(serviceName, accountName, string(raw)); err != nil {
		return fmt.Errorf("keychain write failed: %w", err)
	}

	return nil
}

// macOS Keychain から認証情報を削除する
func delete() error {
	err := keyring.Delete(serviceName, accountName)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return fmt.Errorf("keychain delete failed: %w", err)
	}

	return nil
}
