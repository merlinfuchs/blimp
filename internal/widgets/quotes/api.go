package quotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type QuoteData struct {
	Content string   `json:"content"`
	Author  string   `json:"author"`
	Tags    []string `json:"tags"`
}

func getRandomQuote(tags []string, tagsRequireAll bool) (QuoteData, error) {
	tagsJoinChar := "|"
	if tagsRequireAll {
		tagsJoinChar = ","
	}

	url := fmt.Sprintf(
		"https://api.quotable.io/quotes/random?tags=%s",
		strings.Join(tags, tagsJoinChar),
	)

	resp, err := http.Get(url)
	if err != nil {
		return QuoteData{}, fmt.Errorf("failed to get next departures data: %w", err)
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return QuoteData{}, fmt.Errorf("Failed to get next departures aata, status code %d", resp.StatusCode)
	}

	var data []QuoteData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return QuoteData{}, fmt.Errorf("Failed to decode next departures data, %w", err)
	}

	if len(data) == 0 {
		return QuoteData{}, fmt.Errorf("No quotes found")
	}

	return data[0], nil
}
