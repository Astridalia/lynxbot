package mediawiki

import (
	"context"
	"net/http"
	"time"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

type HttpClient struct {
	Client *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		Client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: cloudflarebp.AddCloudFlareByPass(&http.Transport{}),
		},
	}
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}
