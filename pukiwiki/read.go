package pukiwiki

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageInfo struct {
	Name         string
	Exists       bool
	LastModified string
}

// goquery でページを取得する
func (c *Client) fetchDocument(rawURL string) (*goquery.Document, error) {
	resp, err := c.httpClient.Get(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	if isLoginPage(doc) {
		return nil, ErrSessionExpired
	}

	return doc, nil
}

// isLoginPage はレスポンス HTML にログインフォームが含まれるか確認する
func isLoginPage(doc *goquery.Document) bool {
	return doc.Find(`input[name="username"]`).Length() > 0
}

// PukiWiki のソース表示ページから取得する
func (c *Client) GetPageSource(pageName string) (string, error) {
	reqURL := fmt.Sprintf("%s/?cmd=source&page=%s", c.baseURL, url.QueryEscape(pageName))

	doc, err := c.fetchDocument(reqURL)
	if err != nil {
		return "", err
	}

	source := strings.TrimSpace(doc.Find("#source").Text())
	if source == "" {
		return "", ErrPageNotFound
	}

	return source, nil
}

// ページの存在確認と更新日時を取得して返す
func (c *Client) GetPageInfo(pageName string) (*PageInfo, error) {
	reqURL := fmt.Sprintf("%s/?%s", c.baseURL, url.QueryEscape(pageName))

	doc, err := c.fetchDocument(reqURL)
	if err != nil {
		return nil, err
	}

	info := &PageInfo{Name: pageName}

	// .message_box に "not found" が含まれていればページが存在しない
	msgBox := doc.Find("#body .message_box")
	if msgBox.Length() > 0 {
		text := msgBox.Text()
		if strings.Contains(text, "not found") || strings.Contains(text, "見つかりません") {
			return info, nil
		}
	}

	info.Exists = true

	// #lastmodified からタイムスタンプを取得
	lastMod := strings.TrimSpace(doc.Find("#lastmodified").Text())
	lastMod = strings.TrimPrefix(lastMod, "Last-modified:")
	lastMod = strings.TrimSpace(lastMod)
	info.LastModified = lastMod

	return info, nil
}
