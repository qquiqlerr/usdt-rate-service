//nolint:testpackage // This file needs to be in the same package because we are testing the package internals.
package grinex

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(handler http.Handler) *Client {
	server := httptest.NewServer(handler)
	parsedURL, _ := url.Parse(server.URL)
	return &Client{
		baseURL:    parsedURL,
		httpClient: server.Client(),
	}
}

func TestGetDepth_TableDriven(t *testing.T) {
	expected := DepthResponse{
		Timestamp: time.Now().Unix(),
		Asks:      []Order{{"1.1", "2.2", "3.3", "4.4", "limit"}, {"5.5", "6.6", "7.7", "8.8", "limit"}},
		Bids:      []Order{{"5.5", "6.6", "7.7", "8.8", "limit"}, {"9.9", "10.10", "11.11", "12.12", "limit"}},
	}

	tests := []struct {
		name    string
		client  *Client
		handler http.HandlerFunc
		want    *DepthResponse
		wantErr bool
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, getDepthEndpoint, r.URL.Path)
				assert.Equal(t, "usdtrub", r.URL.Query().Get("market"))
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(expected)
			},
			want:    &expected,
			wantErr: false,
		},
		{
			name: "non-200 status",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "bad request", http.StatusBadRequest)
			},
			wantErr: true,
		},
		{
			name: "invalid json",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, "{invalid json}")
			},
			wantErr: true,
		},
		{
			name: "http client error",
			client: &Client{
				baseURL:    &url.URL{Scheme: "http", Host: "invalid"},
				httpClient: &http.Client{Transport: nil},
			},
			wantErr: true,
		},
		{
			name: "bad base url",
			client: &Client{
				baseURL:    &url.URL{Scheme: "http", Host: "%%invalid%%"},
				httpClient: http.DefaultClient,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *Client
			if tt.client != nil {
				client = tt.client
			} else {
				client = newTestClient(tt.handler)
			}
			ctx := context.Background()
			resp, err := client.GetDepth(ctx, "usdtrub")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, tt.want.Timestamp, resp.Timestamp)
				assert.Equal(t, tt.want.Asks, resp.Asks)
				assert.Equal(t, tt.want.Bids, resp.Bids)
			}
		})
	}
}
