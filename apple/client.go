package apple

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

// CertFile retuns a new certificate using a file
func CertFile(path, password string) (tls.Certificate, error) {
	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".p12":
		return certificate.FromP12File(path, password)
	case ".pem":
		return certificate.FromPemFile(path, password)
	default:
		return tls.Certificate{}, fmt.Errorf("Unknown extension: %s", ext)
	}
}

// CertBytes retuns a new certificate using a bytes
func CertBytes(bytes []byte, password string) (tls.Certificate, error) {
	cert, err := certificate.FromP12Bytes(bytes, password)
	if err != nil {
		cert, err = certificate.FromPemBytes(bytes, password)
	}
	return cert, err
}

// New creates a new APNS2 Client
func New(cert tls.Certificate, production bool) *Client {
	c := apns2.NewClient(cert)
	if production {
		c = c.Production()
	} else {
		c = c.Development()
	}
	return &Client{client: c}
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
