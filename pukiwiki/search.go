package pukiwiki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MatchType string

const (
	MatchAll MatchType = "AND"
	MatchAny MatchType = "OR"
)

// 検索結果の 1 件分
type SearchResult struct {
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
	Body      string `json:"body"`
}

type search2Response struct {
	Results      []SearchResult `json:"results"`
	NextStartIdx int            `json:"next_start_index"`
	SearchDone   bool           `json:"search_done"`
}

// PukiWiki の search2 JSON API を使ってページを検索する
func (c *Client) SearchPages(query string, matchType MatchType) ([]SearchResult, error) {
	formattedQuery := formatQuery(query, matchType)
	searchStartTime := strconv.FormatInt(time.Now().Unix(), 10)

	var allResults []SearchResult
	start := 0

	for {
		reqURL := fmt.Sprintf("%s/?cmd=search2&action=query&encode_hint=%s&q=%s&start=%d&search_start_time=%s",
			c.baseURL,
			url.QueryEscape("ぷ"),
			url.QueryEscape(formattedQuery),
			start,
			searchStartTime,
		)

		resp, err := c.httpClient.Get(reqURL)
		if err != nil {
			return nil, fmt.Errorf("failed to request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		var result search2Response
		decodeErr := json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		if decodeErr != nil {
			return nil, fmt.Errorf("failed to parse response: %w", decodeErr)
		}

		for _, r := range result.Results {
			if c.scope != "" {
				if r.Name != c.scope && !strings.HasPrefix(r.Name, c.scope+"/") {
					continue
				}
			}
			allResults = append(allResults, r)
		}

		if result.SearchDone {
			break
		}
		start = result.NextStartIdx
	}

	return allResults, nil
}

// MachType が AND の場合クエリをそのまま使う
// OR の場合、単語を " OR " で結合する
func formatQuery(query string, matchType MatchType) string {
	if matchType != MatchAny {
		return query
	}
	words := strings.Fields(query)
	return strings.Join(words, " OR ")
}
