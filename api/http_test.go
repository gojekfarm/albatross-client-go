package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gojekfarm/albatross-client-go/flags"
	"github.com/gojekfarm/albatross-client-go/release"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type mockAPIClient struct {
	mock.Mock
}

func (m *mockAPIClient) Send(url string, method string, body io.Reader) (*http.Response, []byte, error) {
	args := m.Called(url, method, body)
	return args.Get(0).(*http.Response), args.Get(1).([]byte), args.Error(2)
}

func TestHttpClientInstallAPIOnSuccess(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&installResponse{
		Status: "deployed",
	})
	if err != nil {
		t.Error("Unable to encode install response")
	}
	httpresponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.InstallFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "testnamespace",
		},
	}
	result, err := httpclient.Install(context.Background(), "testrelease", "testchart", values, fl)
	assert.NoError(t, err)
	assert.Equal(t, result, "deployed")
}

func TestHttpClientInstallAPIOnFailure(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&installResponse{
		Error: "Invalid Request",
	})
	if err != nil {
		t.Error("Unable to encode install response")
	}
	httpresponse := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.InstallFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "testnamespace",
		},
	}
	result, err := httpclient.Install(context.Background(), "testrelease", "testchart", values, fl)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.EqualError(t, err, "Install API returned an error: Invalid Request")
}

func TestHttpClientUpgradeAPIOnSuccess(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&upgradeResponse{
		Status: "deployed",
	})
	if err != nil {
		t.Error("Unable to encode upgrade response")
	}
	httpresponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.UpgradeFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "testnamespace",
		},
	}
	result, err := httpclient.Upgrade(context.Background(), "testrelease", "testchart", values, fl)
	assert.NoError(t, err)
	assert.Equal(t, result, "deployed")
}

func TestHttpClientUpgradeAPIOnFailure(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&upgradeResponse{
		Error: "Invalid Request",
	})
	if err != nil {
		t.Error("Unable to encode upgrade response")
	}
	httpresponse := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.UpgradeFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "testnamespace",
		},
	}
	result, err := httpclient.Upgrade(context.Background(), "testrelease", "testchart", values, fl)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.EqualError(t, err, "Upgrade API returned an error: Invalid Request")
}

func TestHttpClientListAPIOnSuccess(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&listResponse{
		Releases: []release.Release{
			{
				Name:       "test",
				Namespace:  "test",
				Version:    1,
				Status:     "deployed",
				Chart:      "testchart",
				AppVersion: "v1",
			},
		},
	})
	httpresponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	if err != nil {
		t.Error("Unable to encode list response")
	}
	apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	fl := flags.ListFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "testnamespace",
		},
	}
	releases, err := httpclient.List(context.Background(), fl)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, releases[0].Name, "test")
	assert.Equal(t, releases[0].Version, 1)
	assert.Equal(t, releases[0].AppVersion, "v1")
}

func TestHttpClientListAPIOnFailure(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&installResponse{
		Error: "cluster unavailable",
	})
	if err != nil {
		t.Error("Unable to encode install response")
	}

	httpresponse := &http.Response{
		Status:     "400 Bad Request",
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	fl := flags.ListFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "testnamespace",
		},
	}
	result, err := httpclient.List(context.Background(), fl)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.EqualError(t, err, "List API returned an error: cluster unavailable")
}

func TestHttpClientUninstallAPIOnSuccess(t *testing.T) {
	apiClient := new(mockAPIClient)
	apiResponse, err := json.Marshal(&uninstallResponse{
		Status: "uninstalled",
		Release: release.Release{
			Name:"testrelease",
			Namespace: "testnamespace",
			Version: 1,
			Updated: time.Now(),
			Status: "uninstalled",
			Chart: "testchart",
			AppVersion: "1.16.0",
		},
	})
	if err != nil {
		t.Error("Unable to encode install response")
	}
	httpResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiResponse)),
	}
	apiClient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpResponse, apiResponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiClient,
	}

	fl := flags.UninstallFlags{
		DryRun: false,
	}
	result, release, err := httpclient.Uninstall(context.Background(), "testrelease", "testchart", "testnamespace", fl)
	assert.NoError(t, err)
	assert.Equal(t, result, "uninstalled")
	assert.Equal(t, release.Name, "testrelease")
	assert.Equal(t, release.Namespace, "testnamespace")
	assert.Equal(t, release.Version, 1)
	assert.Equal(t, release.AppVersion, "1.16.0")
}

func TestHttpClientUninstallAPIWithDryRunFlag(t *testing.T) {
	apiClient := new(mockAPIClient)
	apiResponse, err := json.Marshal(&uninstallResponse{
		Status: "deployed",
		Release: release.Release{
			Name:"testrelease",
			Namespace: "testnamespace",
			Version: 1,
			Updated: time.Now(),
			Status: "deployed",
			Chart: "testchart",
			AppVersion: "1.16.0",
		},
	})
	if err != nil {
		t.Error("Unable to encode install response")
	}
	httpResponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiResponse)),
	}
	apiClient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpResponse, apiResponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiClient,
	}


	fl := flags.UninstallFlags{
		DryRun: true,
	}
	result, release, err := httpclient.Uninstall(context.Background(), "testrelease", "testchart", "testnamespace", fl)
	assert.NoError(t, err)
	assert.Equal(t, result, "deployed")
	assert.Equal(t, release.Name, "testrelease")
	assert.Equal(t, release.Version, 1)
	assert.Equal(t, release.Namespace, "testnamespace")
	assert.Equal(t, release.AppVersion, "1.16.0")
}

func TestHttpClientUninstallAPIOnFailure(t *testing.T) {
	apiClient := new(mockAPIClient)
	apiResponse, err := json.Marshal(&uninstallResponse{
		Error: "uninstall: Release not loaded: testrelease: release: not found",
	})

	if err != nil {
		t.Error("Unable to encode install response")
	}

	httpResponse := &http.Response{
		Status:     "404 Not Found",
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader(apiResponse)),
	}
	apiClient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpResponse, apiResponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiClient,
	}

	fl := flags.UninstallFlags{
		DryRun: false,
	}
	result, release, err := httpclient.Uninstall(context.Background(), "testrelease", "testchart", "testnamespace", fl)
	assert.Empty(t, result)
	assert.Empty(t, release)
	assert.EqualError(t, err, "Uninstall API returned an error: uninstall: Release not loaded: testrelease: release: not found")
}
