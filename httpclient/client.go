package httpclient

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gojekfarm/albatross-client-go/config"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client acts as a wrapper around the net/http.Client to take care of exponential retries.
// TODO: Discuss if we should use go-retryablehttp
type Client struct {
	client client
	retry  *config.Retry
}

// Request executes a given request with the provided retry policy
// TODO: Handle timeouts and implement retries
// The send implementation follows the same behaviour as the default implementation
// The error is returned only in case of network errors, non 2xx codes do not cause errors
// The assumption is that since it's a json api, all responses should result in a valid
// json body unless under exceptional circumstances. The users can check the response status code
// and parse the bytestream accordingly.
// The client can be extended to handle authentication failures
func (c *Client) Send(url string, method string, body io.Reader) (*http.Response, []byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.client.Do(request)
	if err != nil {
		// TODO: log, this most likely a network error
		return nil, nil, err
	}

	defer resp.Body.Close()

	// Log all 5xx errors, and returns the resp for additional parsing
	if resp.StatusCode >= 500 {
		// TODO: log server errors here
		// With a 5xx, the response most likely won't be parsable
		return resp, nil, nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, data, nil
}

// NewClient returns a new http client
// It sets the client timeout using the timeout specified in config
// and sets retry policy
func NewClient(config *config.Config) *Client {
	return &Client{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		retry: config.Retry,
	}
}
