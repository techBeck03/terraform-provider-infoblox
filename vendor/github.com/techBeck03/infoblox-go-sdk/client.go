package infoblox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// Config - Configuration details for connecting to infoblox
type Config struct {
	Host                   string
	Port                   string
	Version                string
	Username               string
	Password               string
	DisableTLSVerification bool
}

// Client - base client for infoblox interactions
type Client struct {
	client          *http.Client
	config          Config
	baseURL         string
	cookies         []*http.Cookie
	eaDefinitions   []EADefinition
	OrchestratorEAs *ExtensibleAttribute
	SequentialLock  sync.Mutex
}

// New - creates a new infoblox client
func New(config Config) Client {
	var client *http.Client
	if config.DisableTLSVerification {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: transport}
	} else {
		client = http.DefaultClient
	}
	return Client{
		client:  client,
		config:  config,
		baseURL: fmt.Sprintf("https://%s:%s/wapi/v%s", config.Host, config.Port, config.Version),
	}
}

// BuildQuery creates query string
func (c *Client) BuildQuery(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		q.Add(k, v)
	}
	return q.Encode()
}

// CreateJSONRequest - helper function for creating json based http requests
func (c *Client) CreateJSONRequest(method string, path string, params interface{}) (*http.Request, error) {
	var request *http.Request
	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(&params)
	if err != nil {
		return request, err
	}
	combinedPath := fmt.Sprintf("%s/%s", c.baseURL, path)
	request, err = http.NewRequest(method, combinedPath, &buf)
	if err != nil {
		return request, err
	}
	if buf.Len() == 0 {
		request.Body = http.NoBody
	}
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

// Call - function for handling http requests
func (c *Client) Call(request *http.Request, result interface{}) *ResponseError {
	request.SetBasicAuth(c.config.Username, c.config.Password)

	// Use cookies for auth if set
	if len(c.cookies) > 0 {
		for i := range c.cookies {
			request.AddCookie(c.cookies[i])
		}
	}
	response, err := c.client.Do(request)
	if err != nil {
		return &ResponseError{
			StatusCode:   0,
			Request:      fmt.Sprintf("%+v", request),
			ResponseBody: "",
			ErrorMessage: fmt.Sprint(err),
		}
	}
	defer response.Body.Close()
	if !(response.StatusCode >= 200 && response.StatusCode <= 299) {
		var rawBodyBuffer bytes.Buffer
		// Decode raw response, usually contains
		// additional error details
		body := io.TeeReader(response.Body, &rawBodyBuffer)
		var responseBody interface{}
		json.NewDecoder(body).Decode(&responseBody)
		return &ResponseError{
			StatusCode:   response.StatusCode,
			Request:      fmt.Sprintf("%+v", request),
			ResponseBody: fmt.Sprintf("%+v", responseBody),
			ErrorMessage: fmt.Sprintf("Request %+v\n failed with status code %d\n response %+v", request,
				response.StatusCode, responseBody),
		}
	}

	// Add cookies if none exist
	if len(c.cookies) == 0 {
		c.cookies = response.Request.Cookies()
	}
	// If no result is expected, don't attempt to decode a potentially
	// empty response stream and avoid incurring EOF errors
	if result == nil {
		return nil
	}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return &ResponseError{
			StatusCode:   0,
			Request:      fmt.Sprintf("%+v", request),
			ResponseBody: "",
			ErrorMessage: fmt.Sprint(err),
		}
	}
	return nil
}

// Logout clears auth cookie
func (c *Client) Logout() error {
	request, err := c.CreateJSONRequest(http.MethodPost, "logout", nil)
	if err != nil {
		return err
	}
	response := c.Call(request, nil)
	if response != nil {
		return fmt.Errorf(response.ErrorMessage)
	}
	return nil
}
