package api

import (
	"context"
	"net/url"

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

	// Upgrade installs a release, specified by the params, and returns a status.
	// UpgradeFlags govern the actions of upgrade action, i.e whether it should be installed if not present
	Upgrade(ctx context.Context, name string, chart string, values Values, fl flags.UpgradeFlags) (string, error)

	// Status returns the status of a release with the specific release and revision
	Status(ctx context.Context, name string, fl flags.StatusFlags) (release.Release, error)
}

// NewClient returns a new http client for the corresponding host
// and config options.
// In case of invalid host, it returns an error
func NewClient(host string, opts ...config.Option) (Client, error) {
	baseUrl, err := url.ParseRequestURI(host)
	if err != nil {
		return nil, err
	}

	cfg := config.DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	return &HttpClient{
		baseUrl: baseUrl,
		client:  httpclient.NewClient(cfg),
	}, nil
}
