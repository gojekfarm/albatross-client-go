package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/gojekfarm/albatross-client-go/config"
	"github.com/gojekfarm/albatross-client-go/logger"
)

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client acts as a wrapper around the net/http.Client to take care of exponential retries.
// TODO: Discuss if we should use go-retryablehttp
type Client struct {
	client client
	retry  *config.Retry
	logger logger.Logger
}

// Request executes a given request with the provided retry policy
// The send implementation follows the same behaviour as the default implementation
// The error is returned only in case of network errors, non 2xx codes do not cause errors
// The assumption is that since it's a json api, all responses should result in a valid
// json body unless under exceptional circumstances. The users can check the response status code
// and parse the bytestream accordingly.
// The client can be extended to handle authentication failures
func (c *Client) Send(url string, method string, body io.Reader) (*http.Response, []byte, error) {
	resp, err := c.send(url, method, body)
	if err != nil {
		c.logger.Errorf("Error sending request: %s", err)
		return nil, nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	// Log all 5xx errors, and returns the resp for additional parsing
	if resp.StatusCode >= 500 {
		// We just log that we recieved a 5xx and pass the data to the
		// the caller to handle the 5xx data
		c.logger.Errorf("server error for albatross api: %s - %s", url, data)
	}

	return resp, data, nil
}

func (c *Client) send(url string, method string, body io.Reader) (*http.Response, error) {
	if c.retry == nil {
		return c.sendOnce(url, method, body)
	}

	return c.sendWithRetry(url, method, body)
}

func (c *Client) sendOnce(url string, method string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		c.logger.Errorf("Unable to create a new request: %s", err)
		return nil, err
	}

	return c.client.Do(request)
}

func (c *Client) getBackoffForRetry(count int) time.Duration {
	if count <= 0 {
		return 0
	}
	return time.Duration(math.Exp2(float64(count))) * c.retry.Backoff
}

func (c *Client) sendWithRetry(url string, method string, body io.Reader) (*http.Response, error) {
	// reqBytes is used to populate the body for the request for each retry,
	var reqBytes []byte = nil


	if body != nil {
		var err error = nil

		if reqBytes, err = ioutil.ReadAll(body); err != nil {
			return nil, fmt.Errorf("Error reading the request body: %s", err)
		}
	}



	var retryError error
	for count := 0; count <= c.retry.RetryCount; count++ {
		timeout := c.getBackoffForRetry(count)

		select {
		case <-time.After(timeout):
			// We are creating a new request for every retry, which is not ideal,
			// but the Request struct does not provide convenient methods to reset seek offset of
			// the request body for subsequent retries. To do it without creating a new request object
			// everytime, the body needs to be recreated for every retry, and the response body
			// needs to be drained as well to prevent corruption of response object.
			// For now, adopting NewRequest on each retry. We can easily adopt
			// hashicorp/retryablehttp here, it satifies the default http client(and our) interface.
			request, err := http.NewRequest(method, url, bytes.NewBuffer(reqBytes))
			if err != nil {
				c.logger.Errorf("Unable to create a new request: %s", err)
				return nil, err
			}

			resp, err := c.client.Do(request)
			if err != nil {
				c.logger.Errorf("Error connecting to albatross API: %s - retrying", err)
				retryError = err
				break
			}

			return resp, nil
		}
	}

	return nil, fmt.Errorf("Max retries exceeded: %s", retryError)
}

// NewClient returns a new http client
// It sets the client timeout using the timeout specified in config
// and sets retry policy
func NewClient(config *config.Config) *Client {
	return &Client{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		retry:  config.Retry,
		logger: config.Logger,
	}
}
