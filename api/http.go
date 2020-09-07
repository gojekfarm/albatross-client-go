package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gojekfarm/albatross-client-go/flags"
	"github.com/gojekfarm/albatross-client-go/release"
)

// APIClient defines the contract for the http client implementation to send requests to
// the albatross api server
type APIClient interface {
	Send(url string, method string, body io.Reader) (*http.Response, []byte, error)
}

// HttpClient is responsible to sending api requests and parsing their responses
// It embeds the base url of the albatross service and an underlying http apiclient
// that handles sending requests to the albatross api server
type HttpClient struct {
	baseUrl *url.URL
	client  APIClient
}

// installRequest is the json schema for the install api
type installRequest struct {
	Name   string
	Chart  string
	Values Values
	Flags  flags.InstallFlags
}

// installResponse is the json schema to parse the install api response
type installResponse struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
	Data   string `json:"data,omitempty"`
}

// upgradeRequest is the json schema for the upgrade api
type upgradeRequest struct {
	Name   string
	Chart  string
	Values Values
	Flags  flags.UpgradeFlags
}

// upgradeResponse is the json schema to parse the upgrade api response
type upgradeResponse struct {
	Error  string `json:"error,omitempty"`
	Status string `json:"status,omitempty"`
	Data   string `json:"data,omitempty"`
}

// listRequest is the json schema for the list api
type listRequest struct {
	flags.ListFlags
}

// listResponse is the json schema to parse the list api response
type listResponse struct {
	Error    string            `json:"error,omitempty"`
	Releases []release.Release `json:"releases,omitempty"`
}

// request is a helper function to append the path to baseUrl and send the request to the APIClient
func (c *HttpClient) request(ctx context.Context, reqPath string, method string, body io.Reader) (*http.Response, []byte, error) {
	u := *c.baseUrl
	u.Path = path.Join(strings.TrimRight(u.Path, "/"), reqPath)
	return c.client.Send(u.String(), method, body)
}

// List sends the list api request to the APIClient and returns a list of releases if successfull.
func (c *HttpClient) List(ctx context.Context, fl flags.ListFlags) ([]release.Release, error) {
	reqBody, err := json.Marshal(&listRequest{
		ListFlags: fl,
	})
	if err != nil {
		return nil, err
	}

	_, data, err := c.request(ctx, "list", http.MethodGet, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	var result listResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf("List API returned an error: %s", result.Error)

	}

	return result.Releases, nil
}

// Install calls the install api and returns the status
// TODO: Make install api return an installed release rather than just the status
func (c *HttpClient) Install(ctx context.Context, name string, chart string, values Values, fl flags.InstallFlags) (string, error) {
	reqBody, err := json.Marshal(&installRequest{
		Name:   name,
		Chart:  chart,
		Values: values,
		Flags:  fl,
	})
	if err != nil {
		return "", err
	}

	_, data, err := c.request(ctx, "install", http.MethodPut, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	var result installResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("Install API returned an error: %s", result.Error)

	}

	return result.Status, nil
}

// Upgrade calls the upgrade api and returns the status
func (c *HttpClient) Upgrade(ctx context.Context, name string, chart string, values Values, fl flags.UpgradeFlags) (string, error) {
	reqBody, err := json.Marshal(&upgradeRequest{
		Name:   name,
		Chart:  chart,
		Values: values,
		Flags:  fl,
	})
	if err != nil {
		return "", err
	}

	_, data, err := c.request(ctx, "upgrade", http.MethodPost, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	var result upgradeResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("Upgrade API returned an error: %s", result.Error)

	}

	return result.Status, nil
}
