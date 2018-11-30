package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"

	"github.com/valyala/fasthttp"

	"firebase.google.com/go/messaging"
	"google.golang.org/api/googleapi"
)

var (
	projectExtractor = regexp.MustCompile(`"project_id":\s*"([^"]+)",`)
	errProjectID     = errors.New("project_id is missing from json")
)

// ConfigFromFile read configuration from a file
func ConfigFromFile(path string) (*jwt.Config, string, error) {
	if path == "" {
		path = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if path == "" {
			return nil, "", fmt.Errorf("No credentials nor $GOOGLE_APPLICATION_CREDENTIALS specified")
		}
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return ConfigFromBytes(bytes)
}

// ConfigFromBytes read configuration from bytes
func ConfigFromBytes(bytes []byte) (*jwt.Config, string, error) {
	matches := projectExtractor.FindSubmatch(bytes)
	if len(matches) == 0 {
		return nil, "", errProjectID
	}
	conf, err := google.JWTConfigFromJSON(bytes, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return nil, "", errProjectID
	}
	return conf, string(matches[1]), nil
}

// New creates a new FCM Client
func New(conf *jwt.Config, projectID string) *Client {
	return &Client{
		conf:      conf,
		projectID: projectID,
		client:    new(fasthttp.Client),
	}
}

// Client is a FCM client
type Client struct {
	conf      *jwt.Config
	token     *oauth2.Token
	projectID string
	client    *fasthttp.Client
}

func (c *Client) checkToken() error {
	if c.token != nil && c.token.Valid() {
		return nil
	}
	t, err := c.conf.TokenSource(context.Background()).Token()
	if err != nil {
		return err
	}
	c.token = t
	return nil
}

// Send makes a request to FCM and returns the message ID
func (c *Client) Send(m *messaging.Message) (string, int, error) {
	if err := c.checkToken(); err != nil {
		return "", http.StatusInternalServerError, err
	}
	var (
		req   = fasthttp.AcquireRequest()
		uri   = fasthttp.AcquireURI()
		query = fasthttp.AcquireArgs()
		buff  = fasthttp.AcquireByteBuffer()
		resp  = fasthttp.AcquireResponse()
	)
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseURI(uri)
		fasthttp.ReleaseArgs(query)
		fasthttp.ReleaseByteBuffer(buff)
		fasthttp.ReleaseResponse(resp)
	}()
	req.Header.SetMethod(http.MethodPost)
	req.SetRequestURI(fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", c.projectID))
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", c.token.TokenType, c.token.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	if err := json.NewEncoder(buff).Encode(fcmRequest{Message: m}); err != nil {
		return "", http.StatusInternalServerError, err
	}
	req.SetBody(buff.B)
	if err := c.client.Do(req, resp); err != nil {
		return "", http.StatusInternalServerError, err
	}
	if status := resp.StatusCode(); status != http.StatusOK {
		var v struct {
			Error googleapi.Error
		}
		if err := json.Unmarshal(resp.Body(), &v); err != nil {
			return "", http.StatusInternalServerError, err
		}
		return "", status, &v.Error
	}
	var v struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(resp.Body(), &v); err != nil {
		return "", http.StatusInternalServerError, err
	}
	return v.Name, http.StatusOK, nil
}

type fcmRequest struct {
	ValidateOnly bool               `json:"validate_only,omitempty"`
	Message      *messaging.Message `json:"message,omitempty"`
}
