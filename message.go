package pushie

import (
	"fmt"
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

	Data map[string]interface{} `json:"data,omitempty"`
}

// ToFirebase converts the Message to a the Firebase version
func (m Message) ToFirebase() *firebase.Message {
	var ttl *time.Duration
	if m.TTL != 0 {
		ttl = &m.TTL
	}
	data := map[string]string{}
	for k, v := range m.Data {
		if v != nil {
			data[k] = fmt.Sprintf("%v", v)
		}
	}
	return &firebase.Message{
		Topic: m.Google.Topic,
		Token: m.Google.Device,
		Android: &firebase.AndroidConfig{
			Data: data,
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
	payload := payload.NewPayload()
	dataAps := map[string]interface{}{}
	dataHer := map[string]interface{}{}
	for k, v := range m.Data {
		if v != nil {
			switch k {
			case "alert", "badge", "category", "mutable-content", "sound":
				dataAps[k] = v
			default:
				dataAps[k] = v
				dataHer[k] = v
			}
		}
	}
	payload = payload.Custom("aps", dataAps)
	payload = payload.Custom("her", dataHer)
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
