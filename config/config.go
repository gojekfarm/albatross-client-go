package config

import "time"

// ConfiConfigFunc represents the contract of a config modifier function
type ConfigFunc func(config *Config)


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
	// The host to connect to albatross API
	Host string

	// Timeout for API calls
	Timeout time.Duration

	// The retry configuration
	Retry *Retry
}

// DefaultConfig returns a default Config struct with sensible defaults set
func DefaultConfig() *Config {
	return &Config{
		Host:    "http://localhost:8080",
		Timeout: 5 * time.Second,
	}
}

// WithRetry allows the user to set a custom timeout for api calls
func WithTimeout(timeout time.Duration) ConfigFunc {
	return func(config *Config) {
		config.Timeout = timeout
	}
}

// WithRetry sets the retry policy
func WithRetry(retryConfig *Retry) ConfigFunc {
	return func(config *Config) {
		config.Retry = retryConfig
	}
}
