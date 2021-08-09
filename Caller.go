package api

import (
	"context"
	"fmt"
	"github.com/gojekfarm/albatross-client-go/flags"
	ApiHelper "github.com/gojekfarm/albatross-client-go/httpclient"
	"net/url"
)

func main() {
	apiclient := new(ApiHelper.Client)
	//apiresponse, err := json.Marshal(&uninstallResponse{
	//	Status: "uninstalled",
	//})
	//if err != nil {
	//	t.Error("Unable to encode uninstall response")
	//}
	//httpresponse := &http.Response{
	//	Status:     "200 OK",
	//	StatusCode: 200,
	//	Body:       ioutil.NopCloser(bytes.NewReader(apiresponse)),
	//}
	//apiclient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(httpresponse, apiresponse, nil)

	baseUrl, _ := url.ParseRequestURI("http://localhost:8080")

	httpclient := &HttpClient{
		baseUrl: baseUrl,
		client:  apiclient,
	}

	values := Values{
		"test": "test",
	}

	fl := flags.UninstallFlags{
		CommonFlags: flags.CommonFlags{
			Namespace: "quota-test",
		},
	}
	result, err := httpclient.Uninstall(context.Background(), "arpit-test-release", "workload", values, fl)
	if err != nil {
		fmt.Println("Error in the code")
	}
	fmt.Println("Result is  ",  result)

}
