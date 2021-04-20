package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/gojekfarm/albatross-client-go/flags"
	"github.com/gojekfarm/albatross-client-go/release"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	cluster, namespace, releaseName := "integration", "testnamespace", "testrelease"
	expectedURL := fmt.Sprintf("http://localhost:8080/releases/%s/%s/%s", cluster, namespace, releaseName)

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.InstallFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   namespace,
		},
	}
	expectedReq, err := json.Marshal(&installRequest{
		Chart:  "testchart",
		Values: values,
		Flags:  fl,
	})
	assert.NoError(t, err)
	apiclient.On("Send", expectedURL, http.MethodPut, bytes.NewBuffer(expectedReq)).Return(httpresponse, apiresponse, nil).Once()
	result, err := httpclient.Install(context.Background(), releaseName, "testchart", values, fl)
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

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.InstallFlags{
		CommonFlags: flags.CommonFlags{
			Namespace:   "testnamespace",
			KubeContext: "staging",
		},
	}
	jsonRequest, err := json.Marshal(&installRequest{
		Chart:  "",
		Values: values,
		Flags:  fl,
	})
	require.NoError(t, err)
	apiclient.On("Send", mock.Anything, http.MethodPut, bytes.NewBuffer(jsonRequest)).Return(httpresponse, apiresponse, nil)
	result, err := httpclient.Install(context.Background(), "testrelease", "", values, fl)
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

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	cluster, namespace, releaseName := "integration", "testnamespace", "testrelease"
	fl := flags.UpgradeFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   namespace,
		},
	}
	req, err := json.Marshal(&upgradeRequest{
		Chart:  "testchart",
		Values: values,
		Flags:  fl,
	})
	assert.NoError(t, err)
	expectedURL := fmt.Sprintf("http://localhost:8080/releases/%s/%s/%s", cluster, namespace, releaseName)
	apiclient.On("Send", expectedURL, http.MethodPost, bytes.NewBuffer(req)).Return(httpresponse, apiresponse, nil)

	result, err := httpclient.Upgrade(context.Background(), releaseName, "testchart", values, fl)

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
	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}
	cluster, namespace, releaseName := "integration", "testnamespace", "testrelease"
	fl := flags.UpgradeFlags{
		CommonFlags: flags.CommonFlags{
			Namespace:   "testnamespace",
			KubeContext: cluster,
		},
	}
	req, err := json.Marshal(&upgradeRequest{
		Chart:  "testchart",
		Values: values,
		Flags:  fl,
	})
	assert.NoError(t, err)
	expectedURL := fmt.Sprintf("http://localhost:8080/releases/%s/%s/%s", cluster, namespace, releaseName)
	apiclient.On("Send", expectedURL, http.MethodPost, bytes.NewBuffer(req)).Return(httpresponse, apiresponse, nil)

	result, err := httpclient.Upgrade(context.Background(), releaseName, "testchart", values, fl)
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

	cluster := "integration"
	expectedURL := fmt.Sprintf("http://localhost:8080/releases/%s?failed=true", cluster)
	apiclient.On("Send", expectedURL, http.MethodGet, nil).Return(httpresponse, apiresponse, nil)

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.ListFlags{
		Failed:        true,
		AllNamespaces: true,
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
		},
	}
	releases, err := httpclient.List(context.Background(), fl)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, releases[0].Name, "test")
	assert.Equal(t, releases[0].Version, 1)
	assert.Equal(t, releases[0].AppVersion, "v1")
}

func TestHttpClientListWithNamespaceAPIOnSuccess(t *testing.T) {
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

	cluster, namespace := "integration", "testnamespace"
	expectedURL := fmt.Sprintf("http://localhost:8080/releases/%s/%s?failed=true", cluster, namespace)
	apiclient.On("Send", expectedURL, http.MethodGet, nil).Return(httpresponse, apiresponse, nil)

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.ListFlags{
		Failed: true,
		CommonFlags: flags.CommonFlags{
			Namespace:   namespace,
			KubeContext: cluster,
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

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.ListFlags{
		CommonFlags: flags.CommonFlags{
			Namespace:   "testnamespace",
			KubeContext: "unavailable_cluster",
		},
	}
	result, err := httpclient.List(context.Background(), fl)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.EqualError(t, err, "List API returned an error: cluster unavailable")
}
