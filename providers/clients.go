package providers

import (
	"encoding/json"
	"fmt"
	"github.com/gojektech/heimdall"
	"io"
	"io/ioutil"
	"time"
)

type ApiClient interface {
	Get(url string, str interface{}) error
	Post(url string, body io.Reader, str interface{}) error
}

func NewHttpApiClient(timeout time.Duration) ApiClient {
	client := heimdall.NewHTTPClient(timeout)
	return HttpApiClient{http: client}
}

type HttpApiClient struct {
	http heimdall.Client
}

func (client HttpApiClient) Post(url string, body io.Reader, str interface{}) error {

	headers := make(map[string][]string)

	headers["Content-Type"] = []string{"application/json"}

	res, err := client.http.Post(url, body, headers)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d != 200", res.StatusCode)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		return fmt.Errorf("empty response from: %s", url)
	}

	err = json.Unmarshal(content, &str)

	if err != nil {
		return err
	}

	return nil
}

func (client HttpApiClient) Get(url string, str interface{}) error {

	res, err := client.http.Get(url, nil)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &str)

	if err != nil {
		return err
	}

	return nil
}
