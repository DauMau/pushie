package apple

import (
	"os"
	"testing"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
)

func TestClient(t *testing.T) {
	t.Log(os.Getenv("IOS_TOKEN"))
	cert, err := CertFile(os.Getenv("IOS_CERT"), os.Getenv("IOS_PASS"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := New(cert, os.Getenv("IOS_ENV") == "prod").Send(&apns2.Notification{
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
