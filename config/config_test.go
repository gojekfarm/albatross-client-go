package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigWithValidHost(t *testing.T) {
	host := "http://localhost:8080"
	config := DefaultConfig()
	err := WithHost(host)(config)

	assert.NoError(t, err)
	assert.Equal(t, config.Host, host)
}

func TestConfigWithValidInvalidHost(t *testing.T) {
	host := "http/localhost:8080"
	config := DefaultConfig()
	err := WithHost(host)(config)

	assert.Error(t, err)
}

func TestConfigWithTimeout(t *testing.T) {
	timeout := 100 * time.Second
	config := DefaultConfig()
	err := WithTimeout(timeout)(config)

	assert.NoError(t, err)
	assert.Equal(t, config.Timeout, timeout)
}

func TestConfigWithRetry(t *testing.T) {
	retry := &Retry{
		RetryCount: 10,
		Backoff:    100 * time.Second,
	}
	config := DefaultConfig()
	err := WithRetry(retry)(config)

	assert.NoError(t, err)
	assert.Equal(t, config.Retry.RetryCount, retry.RetryCount)
	assert.Equal(t, config.Retry.Backoff, retry.Backoff)
}
