package httpclient

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gojekfarm/albatross-client-go/config"
	"github.com/gojekfarm/albatross-client-go/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestHttpClientSendOnSuccess(t *testing.T) {
	mc := new(mockClient)
	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
	}

	mc.On("Do", mock.Anything).Return(response, nil)
	client := &Client{
		client: mc,
		retry: &config.Retry{
			RetryCount: 3,
			Backoff:    2 * time.Second,
		},
		logger: &logger.DefaultLogger{},
	}

	resp, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.NoError(t, err)
	assert.Equal(t, data, []byte("abcde"))
	assert.Equal(t, resp.StatusCode, 200)
}

func TestHttpClientSendWithRetry(t *testing.T) {

	t.Run("When Body is nil and call is Successful", func(t *testing.T) {
		mc := new(mockClient)
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("ab"))),
		}
		mc.On("Do", mock.Anything).Return(response, nil)
		httpClient := &Client{
			client: mc,
			retry: &config.Retry{
				RetryCount: 3,
				Backoff:    2 * time.Second,
			},
			logger: &logger.DefaultLogger{},
		}

		resp, err := httpClient.sendWithRetry("http://localhost:444", "GET", nil)

		assert.NoError(t, err)
		assert.Equal(t, resp, response)
		assert.Equal(t, resp.StatusCode, 200)
	})

	t.Run("When Body is nil and call fails with 500", func(t *testing.T) {
		mc := new(mockClient)
		response := &http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
		}

		mc.On("Do", mock.Anything).Return(response, nil)
		client := &Client{
			client: mc,
			retry: &config.Retry{
				RetryCount: 3,
				Backoff:    2 * time.Second,
			},
			logger: &logger.DefaultLogger{},
		}

		resp, err := client.sendWithRetry("http://localhost:444", "GET", nil)

		assert.Nil(t, err)
		assert.Equal(t, resp, response)
		assert.Equal(t, resp.StatusCode, 500)
	})

	t.Run("When Body is nil and call fails with 400", func(t *testing.T) {
		mc := new(mockClient)
		response := &http.Response{
			Status:     "400 Bad Request",
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
		}

		mc.On("Do", mock.Anything).Return(response, nil)
		client := &Client{
			client: mc,
			retry: &config.Retry{
				RetryCount: 3,
				Backoff:    2 * time.Second,
			},
			logger: &logger.DefaultLogger{},
		}

		resp, err := client.sendWithRetry("http://localhost:444", "GET", nil)

		assert.Nil(t, err)
		assert.Equal(t, resp, response)
		assert.Equal(t, resp.StatusCode, 400)
	})

	t.Run("When Body is nil and call fails with Network Error On Retries", func(t *testing.T) {
		mc := new(mockClient)

		mc.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Network Error"))
		client := &Client{
			client: mc,
			retry: &config.Retry{
				RetryCount: 3,
				Backoff:    500 * time.Millisecond,
			},
			logger: &logger.DefaultLogger{},
		}

		resp, err := client.sendWithRetry("http://localhost:444", "GET", nil)

		assert.Error(t, err)
		assert.EqualError(t, err, "Max retries exceeded: Network Error")
		assert.Nil(t, resp)
	})

	t.Run("When Body is nil and call fails with Network Error With zero Retries", func(t *testing.T) {
		mc := new(mockClient)

		mc.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Network Error"))
		client := &Client{
			client: mc,
			retry: &config.Retry{
				RetryCount: 0,
			},
			logger: &logger.DefaultLogger{},
		}

		resp, err := client.sendWithRetry("http://localhost:444", "GET", nil)

		assert.Error(t, err)
		assert.EqualError(t, err, "Max retries exceeded: Network Error")
		assert.Nil(t, resp)
	})

	t.Run("When Body is nil and call fails with Network Error With Recovery", func(t *testing.T) {
		mc := new(mockClient)
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
		}

		mc.On("Do", mock.Anything).Return(&http.Response{}, &url.Error{}).Once()
		mc.On("Do", mock.Anything).Return(&http.Response{}, &url.Error{}).Once()
		mc.On("Do", mock.Anything).Return(response, nil).Once()
		client := &Client{
			client: mc,
			retry: &config.Retry{
				RetryCount: 3,
				Backoff:    500 * time.Millisecond,
			},
			logger: &logger.DefaultLogger{},
		}

		resp, err := client.sendWithRetry("http://localhost:444", "GET", nil)

		assert.NoError(t, err)
		assert.Equal(t, resp, response)
		assert.Equal(t, resp.StatusCode, 200)
	})
}

func TestHttpClientSendOnServerError(t *testing.T) {
	mc := new(mockClient)
	response := &http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
	}

	mc.On("Do", mock.Anything).Return(response, nil)
	client := &Client{
		client: mc,
		retry: &config.Retry{
			RetryCount: 3,
			Backoff:    2 * time.Second,
		},
		logger: &logger.DefaultLogger{},
	}

	resp, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.Nil(t, err)
	assert.Equal(t, data, []byte("abcde"))
	assert.Equal(t, resp.StatusCode, 500)
}

func TestHttpClientSendOnBadRequest(t *testing.T) {
	mc := new(mockClient)
	response := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
	}

	mc.On("Do", mock.Anything).Return(response, nil)
	client := &Client{
		client: mc,
		retry: &config.Retry{
			RetryCount: 3,
			Backoff:    2 * time.Second,
		},
		logger: &logger.DefaultLogger{},
	}

	resp, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, data, []byte("abcde"))
	assert.Equal(t, resp.StatusCode, 400)
}

func TestHttpClientSendOnNetworkErrorWithRetries(t *testing.T) {
	mc := new(mockClient)

	mc.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Network Error"))
	client := &Client{
		client: mc,
		retry: &config.Retry{
			RetryCount: 3,
			Backoff:    500 * time.Millisecond,
		},
		logger: &logger.DefaultLogger{},
	}

	_, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.Error(t, err)
	assert.EqualError(t, err, "Max retries exceeded: Network Error")
	assert.Nil(t, data)
}

func TestHttpClientSendOnNetworkErrorWithoutRetries(t *testing.T) {
	mc := new(mockClient)

	mc.On("Do", mock.Anything).Return(&http.Response{}, errors.New("Network Error"))
	client := &Client{
		client: mc,
		retry:  nil,
		logger: &logger.DefaultLogger{},
	}

	_, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.Error(t, err)
	assert.EqualError(t, err, "Network Error")
	assert.Nil(t, data)
}

func TestHttpClientSendOnNetworkErrorWithRecovery(t *testing.T) {
	mc := new(mockClient)
	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("abcde"))),
	}

	mc.On("Do", mock.Anything).Return(&http.Response{}, &url.Error{}).Once()
	mc.On("Do", mock.Anything).Return(&http.Response{}, &url.Error{}).Once()
	mc.On("Do", mock.Anything).Return(response, nil).Once()
	client := &Client{
		client: mc,
		retry: &config.Retry{
			RetryCount: 3,
			Backoff:    500 * time.Millisecond,
		},
		logger: &logger.DefaultLogger{},
	}

	resp, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.NoError(t, err)
	assert.Equal(t, data, []byte("abcde"))
	assert.Equal(t, resp.StatusCode, 200)
}
