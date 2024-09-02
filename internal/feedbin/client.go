package feedbin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	username   string
	password   string
}

type Entry struct {
	ID          int       `json:"id"`
	FeedID      int       `json:"feed_id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewClient(username, password string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		username:   username,
		password:   password,
	}
}

func (c *Client) GetLatestFeeds() ([]Entry, error) {
	// Create a new request for the entries API
	req, err := http.NewRequest("GET", "https://api.feedbin.com/v2/entries.json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters for pagination and sorting
	q := req.URL.Query()
	q.Add("per_page", "15") // Adjust the number of entries as needed
	q.Add("order", "desc")  // Get the latest entries first
	req.URL.RawQuery = q.Encode()

	// Add Basic Authentication
	req.SetBasicAuth(c.username, c.password)

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feeds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var entries []Entry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return entries, nil
}
