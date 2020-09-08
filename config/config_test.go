package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigWithTimeout(t *testing.T) {
	timeout := 100 * time.Second
	config := DefaultConfig()
	WithTimeout(timeout)(config)
	assert.Equal(t, config.Timeout, timeout)
}

func TestConfigWithRetry(t *testing.T) {
	retry := &Retry{
		RetryCount: 10,
		Backoff:    100 * time.Second,
	}
	config := DefaultConfig()
	WithRetry(retry)(config)

	assert.Equal(t, config.Retry.RetryCount, retry.RetryCount)
	assert.Equal(t, config.Retry.Backoff, retry.Backoff)
}
