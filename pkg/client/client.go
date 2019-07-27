package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// XKeyHeader header
	XKeyHeader = "X-AUTH-KEY"
	// XSecretHeader header
	XSecretHeader = "X-AUTH-SECRET"
)

// Client struct
type Client struct {
	host, key, secret string
	httpClient        *http.Client
}

// NotificationInput struct
type NotificationInput struct {
}

// NotificationOutput struct
type NotificationOutput struct {
	MessageID string `json:"messageId"`
	Error     error  `json:"error,omitempty"`
}

// VersionOutput struct
type VersionOutput struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

// Option struct
type Option func(*Client)

// SetHTTPClient func
func SetHTTPClient(httpClient *http.Client) Option {
	return func(cli *Client) {
		cli.httpClient = httpClient
	}
}

// NewClient creates client
func NewClient(host, key, secret string, options ...Option) *Client {
	cli := Client{
		host:   host,
		key:    key,
		secret: secret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for i := range options {
		options[i](&cli)
	}

	return &cli

}

// Version returns health check status
func (a *Client) Version() (*VersionOutput, error) {
	var output VersionOutput
	byts, err := a.do("GET", "/version", nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(byts, &output)
	return &output, err
}

// Send notification message
func (a *Client) Send(input *NotificationInput) (*NotificationOutput, error) {
	byts, _ := json.Marshal(input)
	var output NotificationOutput
	byts, err := a.do("POST", "/v1/notify", byts)
	json.Unmarshal(byts, &output)
	return &output, err
}

func (a *Client) do(method, url string, byts []byte) ([]byte, error) {
	req, err := http.NewRequest(method, a.host+url, bytes.NewBuffer(byts))
	if err != nil {
		return nil, err
	}
	req.Header.Set(XKeyHeader, a.key)
	req.Header.Set(XSecretHeader, a.secret)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
