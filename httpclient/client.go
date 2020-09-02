package httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gojekfarm/albatross-client-go/config"
)

// Client acts as a wrapper around the net/http.Client to take care of exponential retries.
// TODO: Discuss if we should use go-retryablehttp
type Client struct {
	client http.Client
	retry  *config.Retry
}

// Request executes a given request with the provided retry policy
// TODO: Handle timeouts and implement retries
func (c *Client) Send(url string, method string, body io.Reader) (io.Reader, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}

// NewClient returns a new http client
// It sets the client timeout using the timeout specified in config
// and sets retry policy
func NewClient(config *config.Config) *Client {
	return &Client{
		client: http.Client{
			Timeout: config.Timeout,
		},
		retry: config.Retry,
	}
}
