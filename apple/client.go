package apple

import (
	"crypto/tls"
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

// New creates a new APNS2 Client
func New(path string) (c *Client, err error) {
	var cert tls.Certificate
	switch strings.ToLower(filepath.Ext(path)) {
	case ".p12":
		cert, err = certificate.FromP12File(path, "")
		if err != nil {
			return
		}
	case ".pem":
		cert, err = certificate.FromPemFile(path, "")
		if err != nil {
			return
		}
	}
	return &Client{client: apns2.NewClient(cert).Production()}, nil
}

// Client is a FCM client
type Client struct {
	client *apns2.Client
}

// Send makes a request to APNS and returns the message ID
func (c *Client) Send(m *apns2.Notification) (string, error) {
	resp, err := c.client.Push(m)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Reason)
	}
	return resp.ApnsID, nil
}
