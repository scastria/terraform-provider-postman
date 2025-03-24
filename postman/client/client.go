package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	APIKeyHeader    = "X-API-Key"
	ApplicationJson = "application/json"
	TextJS          = "text/javascript"
	Bearer          = "Bearer"
	PostmanDomain   = "api.getpostman.com"
	IdSeparator     = ":"
)

type EntityId struct {
	Id string `json:"id"`
}

type Client struct {
	apiKey     string
	numRetries int
	retryDelay int
	httpClient *http.Client
}

func NewClient(apiKey string, numRetries int, retryDelay int) (*Client, error) {
	c := &Client{
		apiKey:     apiKey,
		numRetries: numRetries,
		retryDelay: retryDelay,
		httpClient: &http.Client{},
	}
	return c, nil
}

func (c *Client) HttpRequest(ctx context.Context, method string, path string, query url.Values, headerMap http.Header, body *bytes.Buffer) (*bytes.Buffer, error) {
	req, err := http.NewRequest(method, c.RequestPath(path), body)
	if err != nil {
		return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
	}
	//Set default headers
	req.Header.Set("User-Agent", "Postman-Terraform-Provider")
	//Handle query values
	if query != nil {
		requestQuery := req.URL.Query()
		for key, values := range query {
			for _, value := range values {
				requestQuery.Add(key, value)
			}
		}
		req.URL.RawQuery = requestQuery.Encode()
	}
	//Handle header values
	if headerMap != nil {
		for key, values := range headerMap {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	//Handle authentication
	if c.apiKey != "" {
		req.Header.Set(APIKeyHeader, c.apiKey)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		tflog.Info(ctx, "Postman API:", map[string]any{"error": err})
	} else {
		tflog.Info(ctx, "Postman API: ", map[string]any{"request": string(requestDump)})
	}
	try := 0
	var resp *http.Response
	for {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, &RequestError{StatusCode: http.StatusInternalServerError, Err: err}
		}
		if (resp.StatusCode == http.StatusBadRequest) || (resp.StatusCode == http.StatusTooManyRequests) || (resp.StatusCode >= http.StatusInternalServerError) {
			try++
			if try >= c.numRetries {
				break
			}
			time.Sleep(time.Duration(c.retryDelay) * time.Second)
			continue
		}
		break
	}
	defer resp.Body.Close()
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: err}
	}
	if (resp.StatusCode < http.StatusOK) || (resp.StatusCode >= http.StatusMultipleChoices) {
		return nil, &RequestError{StatusCode: resp.StatusCode, Err: fmt.Errorf("%s", respBody.String())}
	}
	return respBody, nil
}

func (c *Client) RequestPath(path string) string {
	return fmt.Sprintf("https://%s/%s", PostmanDomain, path)
}
