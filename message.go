package pushie

import (
	"time"

	firebase "firebase.google.com/go/messaging"
	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
)

// Message is the notification data
type Message struct {
	Apple  *Destination `json:"apple,omitempty"`
	Google *Destination `json:"google,omitempty"`

	Priority   Priority      `json:"priority,omitempty"`
	TTL        time.Duration `json:"ttl,omitempty"`
	Title      string        `json:"title,omitempty"`
	Body       string        `json:"body,omitempty"`
	CollapseID string        `json:"collapse_id,omitempty"`

	Data map[string]string `json:"data,omitempty"`
}

// ToFirebase converts the Message to a the Firebase version
func (m Message) ToFirebase() *firebase.Message {
	var ttl *time.Duration
	if m.TTL != 0 {
		ttl = &m.TTL
	}
	return &firebase.Message{
		Topic: m.Google.Topic,
		Token: m.Google.Device,
		Android: &firebase.AndroidConfig{
			Data: m.Data,
			Notification: &firebase.AndroidNotification{
				Title: m.Title,
				Body:  m.Body,
			},
			TTL:         ttl,
			CollapseKey: m.CollapseID,
			Priority:    m.Priority.String(),
		},
	}
}

// ToApns converts the Message to a the Apns version
func (m Message) ToApns() *apns.Notification {
	payload := payload.NewPayload().AlertTitle(m.Title).AlertBody(m.Body)
	for k, v := range m.Data {
		payload = payload.Custom(k, v)
	}
	var exp time.Time
	if m.TTL > 0 {
		exp = time.Now().Add(m.TTL)
	}
	return &apns.Notification{
		Payload:     payload,
		Topic:       m.Apple.Topic,
		DeviceToken: m.Apple.Device,
		CollapseID:  m.CollapseID,
		Expiration:  exp,
		Priority:    int(m.Priority),
	}
}

// Destination is message recepient or topic
type Destination struct {
	Topic  string `json:"topic,omitempty"`
	Device string `json:"device,omitempty"`
}

// Priority is the message priority
type Priority int

// Priorities
const (
	PriorityNormal Priority = apns.PriorityLow
	PriorityHigh   Priority = apns.PriorityHigh
)

func (p Priority) String() string {
	switch {
	case p > 0 && p <= PriorityNormal:
		return "normal"
	case p > PriorityNormal && p <= PriorityHigh:
		return "high"
	default:
		return ""
	}
}
