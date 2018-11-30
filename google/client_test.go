package google

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"firebase.google.com/go/messaging"
)

func TestClient(t *testing.T) {
	conf, projectID, err := ConfigFromFile("")
	if err != nil {
		t.Fatal(err)
	}
	id, code, err := New(conf, projectID).Send(&messaging.Message{
		Data:  map[string]string{"time": time.Now().String(), "sender": "Test Sender"},
		Token: os.Getenv("FCM_REGID"),
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{Title: "Hello", Body: "Hello from the other side"},
			CollapseKey:  "hello",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if code != http.StatusOK {
		t.Fatal(errors.New("Bad http status"))
	}
	t.Log("message_id", id)
}
