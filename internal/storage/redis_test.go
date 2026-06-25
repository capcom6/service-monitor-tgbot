package storage_test

import (
	"encoding/json"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
)

func newRedisTestServer(t *testing.T) *miniredis.Miniredis {
	t.Helper()

	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	t.Cleanup(s.Close)

	return s
}

func TestRedisStorage_Load(t *testing.T) {
	srv := newRedisTestServer(t)

	services := []storage.MonitoredService{
		{
			ID:   "service1",
			Name: "Service 1",
			TCPSocket: storage.TCPSocket{
				Host: "example.com",
				Port: 80,
			},
		},
	}

	data, err := json.Marshal(services)
	if err != nil {
		t.Fatalf("failed to marshal services: %v", err)
	}

	srv.Set("service-monitor:services", string(data))

	store, err := storage.NewFromDSN("redis://" + srv.Addr() + "/0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = store.Close()
	}()

	got, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 service, got %d", len(got))
	}

	if got[0].ID != "service1" {
		t.Fatalf("expected ID 'service1', got '%s'", got[0].ID)
	}
}

func TestRedisStorage_Load_MissingKey(t *testing.T) {
	srv := newRedisTestServer(t)

	store, err := storage.NewFromDSN("redis://" + srv.Addr() + "/0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = store.Close()
	}()

	_, err = store.Load()
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRedisStorage_Load_InvalidJSON(t *testing.T) {
	srv := newRedisTestServer(t)

	srv.Set("service-monitor:services", "invalid json")

	store, err := storage.NewFromDSN("redis://" + srv.Addr() + "/0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = store.Close()
	}()

	_, err = store.Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRedisStorage_Load_CustomKey(t *testing.T) {
	srv := newRedisTestServer(t)

	services := []storage.MonitoredService{
		{ID: "svc1", Name: "Svc 1"},
	}
	data, _ := json.Marshal(services)
	srv.Set("custom:key", string(data))

	store, err := storage.NewFromDSN("redis://" + srv.Addr() + "/0?key=custom:key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = store.Close()
	}()

	got, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 service, got %d", len(got))
	}
}

func TestRedisStorage_Close(t *testing.T) {
	srv := newRedisTestServer(t)

	store, err := storage.NewFromDSN("redis://" + srv.Addr() + "/0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if closeErr := store.Close(); closeErr != nil {
		t.Fatalf("unexpected error on close: %v", closeErr)
	}
}
