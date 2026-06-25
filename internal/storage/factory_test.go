package storage_test

import (
	"testing"

	"github.com/capcom6/service-monitor-tgbot/internal/storage"
)

func TestNewFromDSN_Empty(t *testing.T) {
	_, err := storage.NewFromDSN("")
	if err == nil {
		t.Fatal("expected error for empty DSN")
	}
}

func TestNewFromDSN_MissingScheme(t *testing.T) {
	_, err := storage.NewFromDSN("./configs/services.yml")
	if err == nil {
		t.Fatal("expected error for DSN missing scheme")
	}
}

func TestNewFromDSN_UnsupportedScheme(t *testing.T) {
	_, err := storage.NewFromDSN("unsupported://path")
	if err == nil {
		t.Fatal("expected error for unsupported scheme")
	}
}

func TestNewFromDSN_FileScheme_Relative(t *testing.T) {
	s, err := storage.NewFromDSN("file://./configs/services.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ys, ok := s.(*storage.YamlStorage)
	if !ok {
		t.Fatal("expected yamlStorage")
	}

	if ys.Path != "./configs/services.yml" {
		t.Fatalf("expected path './configs/services.yml', got '%s'", ys.Path)
	}
}

func TestNewFromDSN_FileScheme_Absolute(t *testing.T) {
	s, err := storage.NewFromDSN("file:///etc/service-monitor/services.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ys, ok := s.(*storage.YamlStorage)
	if !ok {
		t.Fatal("expected yamlStorage")
	}

	if ys.Path != "/etc/service-monitor/services.yml" {
		t.Fatalf("expected path '/etc/service-monitor/services.yml', got '%s'", ys.Path)
	}
}

func TestNewFromDSN_RedisScheme(t *testing.T) {
	s, err := storage.NewFromDSN("redis://localhost:6379/0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = s.Close()
	}()

	rs, ok := s.(*storage.RedisStorage)
	if !ok {
		t.Fatal("expected RedisStorage")
	}

	if rs.Key() != "service-monitor:services" {
		t.Fatalf("expected default key, got '%s'", rs.Key())
	}
}

func TestNewFromDSN_RedisScheme_WithCustomKeyAndChannel(t *testing.T) {
	s, err := storage.NewFromDSN("redis://localhost:6379/0?key=mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = s.Close()
	}()

	rs, ok := s.(*storage.RedisStorage)
	if !ok {
		t.Fatal("expected RedisStorage")
	}

	if rs.Key() != "mykey" {
		t.Fatalf("expected key 'mykey', got '%s'", rs.Key())
	}
}

func TestNewFromDSN_FileScheme_Opaque(t *testing.T) {
	s, err := storage.NewFromDSN("file:./configs/services.yml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ys, ok := s.(*storage.YamlStorage)
	if !ok {
		t.Fatal("expected yamlStorage")
	}

	if ys.Path != "./configs/services.yml" {
		t.Fatalf("expected path './configs/services.yml', got '%s'", ys.Path)
	}
}
