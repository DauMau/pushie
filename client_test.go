package pushie

import (
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
	google, err := google.New("")
	if err != nil {
		t.Fatal(err)
	}
	var client = Client{Google: google}

	id, err := client.SendGoogle(&m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("message_id", id)
}

func TestApple(t *testing.T) {
	apple, err := apple.New(os.Getenv("IOS_CERT"))
	if err != nil {
		t.Fatal(err)
	}
	var client = Client{Apple: apple}

	id, err := client.SendApple(&m)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("message_id", id)
}
