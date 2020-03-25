package harbor

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

const (
	apiURL = "/api"
)

func NewClient(baseURL string, username string, password string, tlsInsecureSkipVerify bool) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := &Client{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}

	return client, nil
}

func (client *Client) sendRequest(request *http.Request) ([]byte, error) {
	request.SetBasicAuth(client.username, client.password)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, errors.New("Bad Request")
	}

	return body, nil
}

func (client *Client) get(path string, resource interface{}, params map[string]string) error {
	resourceURL := client.baseURL + apiURL + path

	request, err := http.NewRequest(http.MethodGet, resourceURL, nil)
	if err != nil {
		return err
	}

	if params != nil {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, v)
		}
		request.URL.RawQuery = query.Encode()
	}

	body, err := client.sendRequest(request)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, resource)
}

func (client *Client) post(path string, requestBody interface{}) ([]byte, error) {
	resourceURL := client.baseURL + apiURL + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, resourceURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	body, err := client.sendRequest(request)

	return body, err
}

func (client *Client) put(path string, requestBody interface{}) error {
	resourceURL := client.baseURL + apiURL + path

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPut, resourceURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	_, err = client.sendRequest(request)

	return err
}

func (client *Client) delete(path string, requestBody interface{}) error {
	resourceURL := client.baseURL + apiURL + path

	var body io.Reader

	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(payload)
	}

	request, err := http.NewRequest(http.MethodDelete, resourceURL, body)
	if err != nil {
		return err
	}

	_, err = client.sendRequest(request)

	return err
}
