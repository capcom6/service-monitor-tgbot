package messages_test

import (
	"testing"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"go.uber.org/zap"
)

func TestHeartbeat_AllOnline(t *testing.T) {
	svc := messages.NewService(messages.Config{}, zap.NewNop())

	result, err := svc.Heartbeat(messages.HeartbeatContext{
		TotalServices:   5,
		OnlineServices:  5,
		OfflineServices: 0,
		CheckedAt:       time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "💓 Heartbeat: 5/5 services online" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestHeartbeat_SomeOffline(t *testing.T) {
	svc := messages.NewService(messages.Config{}, zap.NewNop())

	result, err := svc.Heartbeat(messages.HeartbeatContext{
		TotalServices:   5,
		OnlineServices:  3,
		OfflineServices: 2,
		CheckedAt:       time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "💓 Heartbeat: 3/5 services online (2 offline)"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestHeartbeat_AllOffline(t *testing.T) {
	svc := messages.NewService(messages.Config{}, zap.NewNop())

	result, err := svc.Heartbeat(messages.HeartbeatContext{
		TotalServices:   3,
		OnlineServices:  0,
		OfflineServices: 3,
		CheckedAt:       time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "💓 Heartbeat: 0/3 services online (3 offline)"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestHeartbeat_CustomTemplate(t *testing.T) {
	svc := messages.NewService(messages.Config{
		Templates: map[string]string{
			messages.TemplateHeartbeat: "Status: {{.OnlineServices}} up, {{.OfflineServices}} down",
		},
	}, zap.NewNop())

	result, err := svc.Heartbeat(messages.HeartbeatContext{
		TotalServices:   10,
		OnlineServices:  7,
		OfflineServices: 3,
		CheckedAt:       time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Status: 7 up, 3 down"
	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestHeartbeat_TemplateConstant(t *testing.T) {
	if messages.TemplateHeartbeat != "heartbeat" {
		t.Errorf("TemplateHeartbeat = %q, want %q", messages.TemplateHeartbeat, "heartbeat")
	}
}
