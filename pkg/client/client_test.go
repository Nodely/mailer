package main

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sendOkResponse = `{
		"messageId": "1-COOL-MESSAGE-ID"
	}`
	versionOkResponse = `{
		"version": "1.0",
		"build": "some-build-number"
	}`
)

func TestVersionCheck(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "key", r.Header.Get(XKeyHeader))
		assert.Equal(t, "secret", r.Header.Get(XSecretHeader))
		w.Write([]byte(versionOkResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("https://localhost", "key", "secret", SetHTTPClient(httpClient))

	output, err := cli.Version()

	assert.Nil(t, err)
	assert.Equal(t, "1.0", output.Version)
	assert.Equal(t, "some-build-number", output.Build)
}

func TestSend(t *testing.T) {

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "key", r.Header.Get(XKeyHeader))
		assert.Equal(t, "secret", r.Header.Get(XSecretHeader))
		w.Write([]byte(sendOkResponse))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("https://localhost", "key", "secret", SetHTTPClient(httpClient))

	output, err := cli.Send(&NotificationInput{})

	assert.Nil(t, err)
	assert.Nil(t, output.Error)
	assert.Equal(t, "1-COOL-MESSAGE-ID", output.MessageID)
}

func TestInvalidURL(t *testing.T) {
	cli := NewClient("localhost/%%2", "key", "secret")
	_, err := cli.Version()
	assert.NotNil(t, err)
}

func TestInvalidHost(t *testing.T) {
	cli := NewClient("http://localhost,22", "key", "secret")
	_, err := cli.Version()
	assert.NotNil(t, err)
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}
