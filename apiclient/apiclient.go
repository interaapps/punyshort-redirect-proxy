package apiclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type PunyshortAPI struct {
	baseUrl    string
	apiToken   string
	httpClient http.Client
}

func NewClient(baseUrl string, key string) PunyshortAPI {
	api := PunyshortAPI{
		baseUrl:    "https://api.punyshort.ga",
		httpClient: http.Client{},
	}
	api.SetBaseURL(baseUrl)
	api.SetApiToken(key)
	return api
}

func (apiClient *PunyshortAPI) SetApiToken(token string) {
	apiClient.apiToken = token
}

func (apiClient *PunyshortAPI) SetBaseURL(baseURL string) {
	apiClient.baseUrl = baseURL
}

func (apiClient PunyshortAPI) Request(method string, url string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader = nil

	if body != nil {
		bodyJson, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyJson)
	}

	req, err := http.NewRequest(method, apiClient.baseUrl+url, bodyReader)

	if err != nil {
		return nil, err
	}
	if apiClient.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+apiClient.apiToken)
	}
	res, err := apiClient.httpClient.Do(req)

	return res, err
}

func (apiClient PunyshortAPI) RequestMap(method string, url string, body interface{}, ma interface{}) (*http.Response, error) {
	response, err := apiClient.Request(method, url, body)
	if err != nil {
		return response, err
	}
	all, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	err2 := json.Unmarshal(all, &ma)
	if err2 != nil {
		return nil, err2
	}
	return response, err
}
