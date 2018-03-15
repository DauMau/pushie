package google

import (
	"os"
	"testing"
	"time"

	"firebase.google.com/go/messaging"
)

func TestClient(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	id, err := c.Send(&messaging.Message{
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
	t.Log("message_id", id)
}
