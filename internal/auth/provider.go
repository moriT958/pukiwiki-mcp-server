package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	libpuki "github.com/moriT958/pukiwiki-mcp"
)

// 認証済みの libpuki.Client を保持する
type Provider struct {
	mu     sync.Mutex
	client *libpuki.Client
}

// 認証済みの libpuki.Client を返す
func (p *Provider) Get(ctx context.Context) (*libpuki.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		return p.client, nil
	}

	cfg, err := load()
	if errors.Is(err, ErrNotFound) {
		cfg, err = runcWizard(ctx)
		if err != nil {
			return nil, fmt.Errorf("setup wizard failed: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	client, err := p.buildClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	p.client = client
	return client, nil
}

// キャッシュ済みクライアントを破棄し Keychain の認証情報を削除する
func (p *Provider) Reset() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.client = nil
	if err := delete(); err != nil {
		return fmt.Errorf("failed to delete credentials from keychain: %w", err)
	}
	return nil
}

// Config から libpuki.Client を生成してログインする
func (p *Provider) buildClient(ctx context.Context, cfg *config) (*libpuki.Client, error) {
	opts := []libpuki.Option{
		libpuki.WithAuth(cfg.Username, cfg.Password),
	}

	// TODO: Scope 設定は必須にするか検討
	if cfg.Scope != "" {
		opts = append(opts, libpuki.WithScope(cfg.Scope))
	}

	client, err := libpuki.New(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Login(); err != nil {
		if errors.Is(err, libpuki.ErrAuthFailed) {
			if delErr := delete(); delErr != nil {
				fmt.Fprintf(os.Stderr, "pukiwiki-mcp: failed to delete credentials from keychain: %v\n", delErr)
			}

			newCfg, wizardErr := runcWizard(ctx)
			if wizardErr != nil {
				return nil, fmt.Errorf("re-authentication wizard failed: %w", wizardErr)
			}

			return p.buildClient(ctx, newCfg)
		}
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return client, nil
}
