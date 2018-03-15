package apple

import (
	"os"
	"testing"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
)

func TestClient(t *testing.T) {
	t.Log(os.Getenv("IOS_TOKEN"))
	c, err := New(os.Getenv("IOS_CERT"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := c.Send(&apns2.Notification{
		Payload:     payload.NewPayload().AlertTitle("Hello").AlertBody("Hello from the other side"),
		DeviceToken: os.Getenv("IOS_TOKEN"),
		Priority:    apns2.PriorityHigh,
		Topic:       os.Getenv("IOS_TOPIC"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("message_id", id)
}
