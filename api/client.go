package api

import (
	"context"

	"github.com/gojekfarm/albatross-client-go/config"
	"github.com/gojekfarm/albatross-client-go/flags"
	"github.com/gojekfarm/albatross-client-go/httpclient"
	"github.com/gojekfarm/albatross-client-go/release"
)

// Values represents the chart values that need to be overriden
type Values map[string]interface{}

// Client represents a contract that a concrete client types(http/grpc) must implement
type Client interface {
	// List returns a list of release corresponding to the provided list flags
	List(ctx context.Context, fl flags.ListFlags) ([]release.Release, error)

	// Install installs a release, specified by the params, and returns a status
	// TODO: We should have the api return the release object, instead of just the status
	Install(ctx context.Context, name string, chart string, values Values, fl flags.InstallFlags) (string, error)
}

// NewClient returns a new based on the passed Config
// It looks at the connection type in the config struct and hands out appropriate clients
func NewClient(host string, confFuncs ...config.ConfigFunc) Client {
	config := config.DefaultConfig()
	config.Host = host

	for _, confFunc := range confFuncs {
		confFunc(config)
	}

	return &HttpClient{
		baseUrl: config.Host,
		client:  httpclient.NewClient(config),
	}
}
