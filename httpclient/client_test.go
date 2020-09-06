package httpclient

import (
	"bytes"
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
	assert.Nil(t, data)
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

func TestHttpClientSendOnNetworkError(t *testing.T) {
	mc := new(mockClient)

	mc.On("Do", mock.Anything).Return(&http.Response{}, &url.Error{})
	client := &Client{
		client: mc,
		retry: &config.Retry{
			RetryCount: 3,
			Backoff:    2 * time.Second,
		},
		logger: &logger.DefaultLogger{},
	}

	_, data, err := client.Send("http://localhost:444", "GET", bytes.NewReader([]byte("abcde")))

	assert.Error(t, err)
	assert.Nil(t, data)
}
