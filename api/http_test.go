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
	if args.Get(1) == nil {
		return args.Get(0).(*http.Response), nil, args.Error(2)
	}
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
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases", cluster, namespace)

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
		Name:   releaseName,
	})
	assert.NoError(t, err)
	apiclient.On("Send", expectedURL, http.MethodPost, bytes.NewBuffer(expectedReq)).Return(httpresponse, apiresponse, nil).Once()
	result, err := httpclient.Install(context.Background(), releaseName, "testchart", values, fl)
	assert.NoError(t, err)
	assert.Equal(t, result, "deployed")
	apiclient.AssertExpectations(t)
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
		Name:   "testrelease",
	})
	require.NoError(t, err)
	apiclient.On("Send", mock.Anything, http.MethodPost, bytes.NewBuffer(jsonRequest)).Return(httpresponse, apiresponse, nil)
	result, err := httpclient.Install(context.Background(), "testrelease", "", values, fl)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.EqualError(t, err, "Install API returned an error: Invalid Request")
	apiclient.AssertExpectations(t)
}

func TestHttpClientInstallAPIReturnsErrorWhenNameIsEmptyString(t *testing.T) {
	apiclient := new(mockAPIClient)

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
	result, err := httpclient.Install(context.Background(), "", "", values, fl)
	assert.Error(t, err)
	assert.Empty(t, result)
	assert.EqualError(t, err, "name cannot be empty")
	apiclient.AssertExpectations(t)
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
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, namespace, releaseName)
	apiclient.On("Send", expectedURL, http.MethodPut, bytes.NewBuffer(req)).Return(httpresponse, apiresponse, nil)

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
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, namespace, releaseName)
	apiclient.On("Send", expectedURL, http.MethodPut, bytes.NewBuffer(req)).Return(httpresponse, apiresponse, nil)

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
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/releases?failed=true", cluster)
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
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases?failed=true", cluster, namespace)
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

func TestHttpClientStatusAPIOnSuccess(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&statusResponse{
		Release: release.Release{
			Name:       "test",
			Namespace:  "test",
			Version:    1,
			Status:     "deployed",
			Chart:      "testchart",
			AppVersion: "v1",
		},
	})
	httpresponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}
	if err != nil {
		t.Error("Unable to encode status response")
	}

	cluster := "integration"
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, "test", "test")
	apiclient.On("Send", expectedURL, http.MethodGet, nil).Return(httpresponse, apiresponse, nil)

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.StatusFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   "test",
		},
	}
	release, err := httpclient.Status(context.Background(), "test", fl)
	assert.NoError(t, err)
	assert.Equal(t, release.Name, "test")
	assert.Equal(t, release.Version, 1)
	assert.Equal(t, release.AppVersion, "v1")
}

func TestHttpClientStatusAPIOnNotFoundFailure(t *testing.T) {
	apiclient := new(mockAPIClient)
	httpresponse := &http.Response{
		Status:     "404 Not Found",
		StatusCode: 404,
		Body:       http.NoBody,
	}

	cluster := "integration"
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, "test", "test")
	apiclient.On("Send", expectedURL, http.MethodGet, nil).Return(httpresponse, nil, nil).Once()

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.StatusFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   "test",
		},
	}
	_, err := httpclient.Status(context.Background(), "test", fl)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("no release found: %s", "test"), err.Error())
}

func TestHttpClientStatusAPIOnServerFailure(t *testing.T) {
	apiclient := new(mockAPIClient)
	apiresponse, err := json.Marshal(&statusResponse{
		Error: "server error",
	})
	if err != nil {
		t.Error("Unable to encode status response")
	}
	httpresponse := &http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}

	cluster := "integration"
	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, "test", "test")
	apiclient.On("Send", expectedURL, http.MethodGet, nil).Return(httpresponse, apiresponse, nil).Once()

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.StatusFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   "test",
		},
	}
	_, err = httpclient.Status(context.Background(), "test", fl)
	assert.Error(t, err)
	assert.Equal(t, "Status API returned an error: server error", err.Error())
}

func TestHttpClientUninstallApiOnSuccess(t *testing.T) {
	apiclient := new(mockAPIClient)
	cluster := "integration"
	expectedRelease := release.Release{
		Name:       "test",
		Namespace:  "test",
		Version:    1,
		Status:     "deployed",
		Chart:      "testchart",
		AppVersion: "v1",
	}
	apiresponse, err := json.Marshal(&unintstallResponse{
		Release: expectedRelease,
	})
	if err != nil {
		t.Error(err)
	}
	httpresponse := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}

	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s?keep_history=true", cluster, "test", "test")
	apiclient.On("Send", expectedURL, http.MethodDelete, nil).Return(httpresponse, apiresponse, nil).Once()

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.UninstallFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   "test",
		},
		KeepHistory: true,
	}
	release, err := httpclient.Uninstall(context.Background(), "test", fl)
	assert.NoError(t, err)
	assert.Equal(t, expectedRelease, release)
}

func TestHttpClientUninstallApiOnFailure(t *testing.T) {
	apiclient := new(mockAPIClient)
	cluster := "integration"
	expectedError := "Uninstall API returned an error: Something went wrong on server end"
	apiresponse, err := json.Marshal(&unintstallResponse{
		Error: "Something went wrong on server end",
	})
	if err != nil {
		t.Error("Unable to encode uninstall response")
	}
	httpresponse := &http.Response{
		Status:     "500 Internal Server Error",
		StatusCode: 500,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}

	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, "test", "test")
	apiclient.On("Send", expectedURL, http.MethodDelete, nil).Return(httpresponse, apiresponse, nil).Once()

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.UninstallFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   "test",
		},
	}
	release, err := httpclient.Uninstall(context.Background(), "test", fl)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	assert.NotNil(t, release)
	assert.Empty(t, release.Name)
}

func TestHttpClientUninstallApiOnNotFound(t *testing.T) {
	apiclient := new(mockAPIClient)
	cluster := "integration"
	expectedError := "no release found: test"
	apiresponse, err := json.Marshal(&unintstallResponse{
		Error: "Something went wrong on server end",
	})
	if err != nil {
		t.Error("Unable to encode uninstall response")
	}
	httpresponse := &http.Response{
		Status:     "404 Not Found",
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	}

	expectedURL := fmt.Sprintf("http://localhost:8080/clusters/%s/namespaces/%s/releases/%s", cluster, "test", "test")
	apiclient.On("Send", expectedURL, http.MethodDelete, nil).Return(httpresponse, apiresponse, nil).Once()

	baseURL, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseURL,
		client:  apiclient,
	}

	fl := flags.UninstallFlags{
		CommonFlags: flags.CommonFlags{
			KubeContext: cluster,
			Namespace:   "test",
		},
	}
	release, err := httpclient.Uninstall(context.Background(), "test", fl)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err.Error())
	assert.NotNil(t, release)
	assert.Empty(t, release.Name)
}
