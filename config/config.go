package config

import (
	"time"

	"github.com/gojekfarm/albatross-client-go/logger"
)

// Option represents the contract of a config modifier function
type Option func(config *Config)

// Retry keeps the retry policy for api calls
type Retry struct {
	// Max number of retries, the implementation follows exponotial retries
	RetryCount int

	// Backoff determines time between retries, the backoff for successive retries will
	// be exponential
	Backoff time.Duration
}

// Config defines settings for a new client
type Config struct {
	// Timeout for API calls
	Timeout time.Duration

	// The retry configuration
	Retry *Retry

	// The logger instance for the client
	Logger logger.Logger
}

// DefaultConfig returns a default Config struct with sensible defaults set
func DefaultConfig() *Config {
	return &Config{
		Timeout: 5 * time.Second,
		Logger:  &logger.DefaultLogger{},
	}
}

// WithRetry allows the user to set a custom timeout for api calls
func WithTimeout(timeout time.Duration) Option {
	return func(config *Config) {
		config.Timeout = timeout
	}
}

// WithRetry sets the retry policy
func WithRetry(retryConfig *Retry) Option {
	return func(config *Config) {
		config.Retry = retryConfig
	}
}

// WithLogger sets the logger for the client
func WithLogger(logger logger.Logger) Option {
	return func(config *Config) {
		config.Logger = logger
	}
}
