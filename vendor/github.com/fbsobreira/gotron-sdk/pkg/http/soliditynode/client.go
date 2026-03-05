package soliditynode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
}

func NewClient(baseURL string, client *http.Client) *Client {
	return &Client{
		BaseURL:    baseURL,
		HttpClient: client,
	}
}

func (tc *Client) request(method string, path string, reqBody interface{}, responseBody interface{}) error {
	endpoint := tc.BaseURL + path

	var req *http.Request

	if reqBody != nil {
		var jsonbytes []byte
		var err error
		jsonbytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON request body (%s %s): %w", method, endpoint, err)
		}

		req, err = http.NewRequest(method, endpoint, bytes.NewBuffer(jsonbytes))
		if err != nil {
			return fmt.Errorf("failed to create new HTTP request with body (%s %s): %w", method, endpoint, err)
		}
	} else {
		var err error
		req, err = http.NewRequest(method, endpoint, nil)
		if err != nil {
			return fmt.Errorf("failed to create new HTTP request (%s %s): %w", method, endpoint, err)
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := tc.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request (%s %s): %w", method, endpoint, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response body (%s %s): %w", method, endpoint, err)
	}

	// this is fine because TRON node only returns 200 for success.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid http status (%s %s): %d", method, endpoint, resp.StatusCode)
	}

	// Check for possible Error response in response body
	errResponse := make(map[string]interface{})

	if err = json.Unmarshal(body, &errResponse); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response for error check (%s %s): %w", method, endpoint, err)
	}

	if responseError, exists := errResponse["Error"]; exists {
		responseErrorStr, ok := responseError.(string)
		if !ok {
			return fmt.Errorf("failed to read JSON error field as string (%s %s): %+v", method, endpoint, responseError)
		}
		return fmt.Errorf("RPC returned error (%s %s): %s", method, endpoint, responseErrorStr)
	}

	// TODO: consider using mapstructure instead of Unmarshaling twice from JSON
	if err = json.Unmarshal(body, responseBody); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response (%s %s): %w", method, endpoint, err)
	}

	return nil

}

func (tc *Client) Post(endpoint string, reqBody, responseBody interface{}) error {
	return tc.request("POST", endpoint, reqBody, responseBody)
}

func (tc *Client) Get(endpoint string, responseBody interface{}) error {
	return tc.request("GET", endpoint, nil, responseBody)
}
