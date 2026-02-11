package libpuki

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

// PukiWiki のソース表示ページから取得する
func (c *Client) GetPageSource(pageName string) (string, error) {
	reqURL := fmt.Sprintf("%s/?cmd=source&page=%s", c.baseURL, url.QueryEscape(pageName))

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed to request page source: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
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

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to request page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
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
