package grinex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const getDepthEndpoint = "/api/v2/depth"

// Client represents a Grinex API client.
// It contains the base URL for the API and an HTTP client for making requests.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// NewClient creates a new Grinex API client.
func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{}
	return &Client{
		baseURL:    parsedURL,
		httpClient: httpClient,
	}, nil
}

// GetDepth retrieves the depth of the specified market from the Grinex API.
// It returns a DepthResponse containing the timestamp, asks, and bids.
// The market parameter should be a string representing the market, e.g., "usdtrub".
func (c *Client) GetDepth(ctx context.Context, market string) (*DepthResponse, error) {
	depthURL, err := c.baseURL.Parse(getDepthEndpoint)
	if err != nil {
		return nil, err
	}
	query := depthURL.Query()
	query.Set("market", market)
	depthURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, depthURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	var depthResponse DepthResponse
	if err = json.NewDecoder(resp.Body).Decode(&depthResponse); err != nil {
		return nil, err
	}

	return &depthResponse, nil
}
