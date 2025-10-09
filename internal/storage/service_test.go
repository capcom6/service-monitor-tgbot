package storage_test

import (
	"reflect"
	"testing"

	"github.com/capcom6/service-monitor-tgbot/internal/storage"
)

var (
	emptyService = storage.MonitoredService{}
	tcpService   = storage.MonitoredService{
		ID:   "tcpService",
		Name: "tcpService",
		TCPSocket: storage.TCPSocket{
			Host: "tcpService",
			Port: 9999,
		},
	}
	tcpServiceResult = storage.MonitoredService{
		ID:                     "tcpService",
		Name:                   "tcpService",
		InitialDelaySecondsRaw: 0,
		PeriodSeconds:          10,
		TimeoutSeconds:         1,
		SuccessThreshold:       1,
		FailureThreshold:       3,
		TCPSocket: storage.TCPSocket{
			Host: "tcpService",
			Port: 9999,
		},
	}
	httpService = storage.MonitoredService{
		ID:   "httpService",
		Name: "httpService",
		HTTPGet: storage.HTTPGet{
			TCPSocket: storage.TCPSocket{
				Host: "google.com",
			},
		},
	}
	httpServiceResult = storage.MonitoredService{
		ID:                     "httpService",
		Name:                   "httpService",
		InitialDelaySecondsRaw: 0,
		PeriodSeconds:          10,
		TimeoutSeconds:         1,
		SuccessThreshold:       1,
		FailureThreshold:       3,
		HTTPGet:                storage.HTTPGet{TCPSocket: storage.TCPSocket{Host: "google.com", Port: 80}, Scheme: "http", Path: "/"},
		TCPSocket:              storage.TCPSocket{},
	}
)

type mockStorage struct {
	Services []storage.MonitoredService
	Error    error
}

func (m *mockStorage) Load() ([]storage.MonitoredService, error) {
	return m.Services, m.Error
}

func TestStorageService_Load(t *testing.T) {
	type fields struct {
		storage storage.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    []storage.MonitoredService
		wantErr bool
	}{
		{
			name: "empty service",
			fields: fields{
				storage: &mockStorage{
					Services: []storage.MonitoredService{emptyService, tcpService, httpService},
					Error:    nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "two services",
			fields: fields{
				storage: &mockStorage{
					Services: []storage.MonitoredService{tcpService, httpService},
					Error:    nil,
				},
			},
			want:    []storage.MonitoredService{tcpServiceResult, httpServiceResult},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage.NewService(tt.fields.storage)
			got, err := s.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("StorageService.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StorageService.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
