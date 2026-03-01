package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// pukiwiki.Client の初期化に必要な認証情報を保持する
type config struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Scope    string `json:"scope,omitempty"`
}

const appDirName = "PukiwikiMCP"
const configFileName = "config.json"

var (
	// 認証情報ファイルが存在しない場合に返す
	ErrNotFound = errors.New("credentials not found")
)

// os.UserConfigDir() に準拠した設定ファイルパスを返す
func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config directory: %w", err)
	}
	return filepath.Join(dir, appDirName, configFileName), nil
}

// 設定ファイルから認証情報を読み込む
func load() (*config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("config file read failed: %w", err)
	}
	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("credential parse failed: %w", err)
	}
	return &cfg, nil
}

// 認証情報を JSON エンコードして設定ファイルに保存する
func Save(cfg *config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("credential marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("config file write failed: %w", err)
	}
	return nil
}

// 設定ファイルを削除する
func delete() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("config file delete failed: %w", err)
	}
	return nil
}
