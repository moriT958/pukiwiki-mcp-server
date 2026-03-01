package pukiwiki

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// PukiWiki 内ののページ一覧を返す
// WithScope が設定されている場合、そのプレフィックス配下のページのみを返す
func (c *Client) ListPages() ([]string, error) {
	reqURL := fmt.Sprintf("%s/?cmd=list", c.baseURL)

	doc, err := c.fetchDocument(reqURL)
	if err != nil {
		return nil, err
	}

	var pages []string
	doc.Find("#body ul li ul li a").Each(func(_ int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		if name == "" {
			return
		}
		if c.scope != "" {
			if name != c.scope && !strings.HasPrefix(name, c.scope+"/") {
				return
			}
		}
		pages = append(pages, name)
	})

	return pages, nil
}
