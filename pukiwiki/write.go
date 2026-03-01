package pukiwiki

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ページ名が管理可能なスコープ内かどうかを確認する
func (c *Client) checkWriteScope(pageName string) error {
	if c.scope == "" {
		return nil
	}
	if pageName == c.scope || strings.HasPrefix(pageName, c.scope+"/") {
		return nil
	}
	return ErrOutOfScope
}

// 編集フォームを取得し、digest とページ存在フラグを返す
// textarea[name="msg"] が空でなければページが既存と判断する
func (c *Client) getEditForm(pageName string) (digest string, exists bool, err error) {
	reqURL := fmt.Sprintf("%s/?cmd=edit&page=%s", c.baseURL, url.QueryEscape(pageName))

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch edit form: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse edit form: %w", err)
	}

	if isLoginPage(doc) {
		return "", false, ErrSessionExpired
	}

	digest, ok := doc.Find(`input[name="digest"]`).Attr("value")
	if !ok || digest == "" {
		return "", false, fmt.Errorf("digest field not found in edit form")
	}

	existingContent := strings.TrimSpace(doc.Find(`textarea[name="msg"]`).Text())
	return digest, existingContent != "", nil
}

func (c *Client) postWrite(pageName, digest, content string) error {
	c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() { c.httpClient.CheckRedirect = nil }()

	form := url.Values{}
	form.Set("cmd", "edit")
	form.Set("page", pageName)
	form.Set("digest", digest)
	form.Set("msg", content)
	form.Set("write", "true")

	resp, err := c.httpClient.PostForm(c.baseURL+"/", form)
	if err != nil {
		return fmt.Errorf("failed to post write: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return nil
	}
	return fmt.Errorf("write failed: unexpected status code %d", resp.StatusCode)
}

// 既存ページを上書き編集する
func (c *Client) EditPage(pageName, content string) error {
	if err := c.checkWriteScope(pageName); err != nil {
		return err
	}

	digest, exists, err := c.getEditForm(pageName)
	if err != nil {
		return err
	}

	// ページが存在しない場合は ErrPageNotFound を返す
	if !exists {
		return ErrPageNotFound
	}

	return c.postWrite(pageName, digest, content)
}

// 新規ページを作成する
func (c *Client) CreatePage(pageName, content string) error {
	if err := c.checkWriteScope(pageName); err != nil {
		return err
	}

	digest, exists, err := c.getEditForm(pageName)
	if err != nil {
		return err
	}

	// ページが既に存在する場合は ErrPageAlreadyExists を返す
	if exists {
		return ErrPageAlreadyExists
	}

	return c.postWrite(pageName, digest, content)
}
