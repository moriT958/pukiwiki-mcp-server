package libpuki

import (
	"fmt"
	"net/http"
	"net/url"
)

// PukiWiki のフォーム認証 (AUTH_TYPE_FORM) を実行する
func (c *Client) Login() error {
	if c.username == "" && c.password == "" {
		return nil
	}

	loginURL := fmt.Sprintf("%s/?plugin=loginform", c.baseURL)

	form := url.Values{}
	form.Set("username", c.username)
	form.Set("password", c.password)

	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		c.httpClient.CheckRedirect = nil
	}()

	resp, err := c.httpClient.PostForm(loginURL, form)
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return nil
	}

	return ErrAuthFailed
}
