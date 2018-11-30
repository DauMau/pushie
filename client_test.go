package pushie

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/DauMau/pushie/apple"
	"github.com/DauMau/pushie/google"
)

var m = Message{
	Google: &Destination{
		Device: os.Getenv("FCM_REGID"),
	},
	Apple: &Destination{
		Device: os.Getenv("IOS_TOKEN"),
		Topic:  os.Getenv("IOS_TOPIC"),
	},
	Title:    "Hello",
	Body:     "Hello from the other side",
	Priority: PriorityHigh,
}

func TestGoogle(t *testing.T) {
	conf, projectID, err := google.ConfigFromFile("")
	if err != nil {
		t.Fatal(err)
	}
	var client = Client{Google: google.New(conf, projectID)}

	id, code, err := client.SendGoogle(&m)
	if err != nil {
		t.Fatal(err)
	}
	if code != http.StatusOK {
		t.Fatal(errors.New("Bad http status"))
	}
	t.Log("message_id", id)
}

func TestApple(t *testing.T) {
	cert, err := apple.CertFile(os.Getenv("IOS_CERT"), os.Getenv("IOS_PASS"))
	if err != nil {
		t.Fatal(err)
	}
	var client = Client{Apple: apple.New(cert, os.Getenv("IOS_ENV") == "prod")}

	id, code, err := client.SendApple(&m)
	if err != nil {
		t.Fatal(err)
	}
	if code != http.StatusOK {
		t.Fatal(errors.New("Bad http status"))
	}
	t.Log("message_id", id)
}
