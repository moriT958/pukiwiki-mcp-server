package libpuki

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
	scope      string
}

type Option func(*Client)

func New(baseURL string, opts ...Option) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL must not be empty")
	}

	jar, _ := cookiejar.New(nil)
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}

	for _, opt := range opts {
		opt(c)
	}

	// opts で追加した httpClient で CookieJar が nil で上書きされた場合の対策
	if c.httpClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		c.httpClient.Jar = jar
	}

	return c, nil
}

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = d
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

func WithAuth(username, password string) Option {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

func WithScope(prefix string) Option {
	return func(c *Client) {
		c.scope = strings.TrimRight(prefix, "/")
	}
}
