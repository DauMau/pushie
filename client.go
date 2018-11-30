package pushie

import (
	"errors"
	"net/http"

	"github.com/DauMau/pushie/apple"
	"github.com/DauMau/pushie/google"
)

// Error List
var (
	ErrDestination = errors.New("No destination: please specifify a device or a topic")
)

// Client is a common interface for both APNS and FCM
type Client struct {
	Google *google.Client
	Apple  *apple.Client
}

// SendGoogle calls to Google service
func (c *Client) SendGoogle(m *Message) (string, int, error) {
	if m.Google == nil {
		return "", http.StatusInternalServerError, ErrDestination
	}
	return c.Google.Send(m.ToFirebase())
}

// SendApple calls to Apple service
func (c *Client) SendApple(m *Message) (string, int, error) {
	if m.Apple == nil {
		return "", http.StatusInternalServerError, ErrDestination
	}
	return c.Apple.Send(m.ToApns())
}
