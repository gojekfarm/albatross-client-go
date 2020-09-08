package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientWithNoConfFuncs(t *testing.T) {
	host := "http://localhost:8080"
	_, err := NewClient(host)
	assert.NoError(t, err)
}

func TestNewClientShouldFailForInvalidUrl(t *testing.T) {
	host := "http/localhost:888"
	_, err := NewClient(host)
	assert.Error(t, err)
}
