package storage

import (
	"reflect"
	"testing"
)

var (
	emptyService = Service{}
	tcpService   = Service{
		Id:   "tcpService",
		Name: "tcpService",
		TCPSocket: TCPSocket{
			Host: "tcpService",
			Port: 9999,
		},
	}
	tcpServiceResult = Service{
		Id:                  "tcpService",
		Name:                "tcpService",
		InitialDelaySeconds: 0,
		PeriodSeconds:       10,
		TimeoutSeconds:      1,
		SuccessThreshold:    1,
		FailureThreshold:    3,
		TCPSocket: TCPSocket{
			Host: "tcpService",
			Port: 9999,
		},
	}
	httpService = Service{
		Id:   "httpService",
		Name: "httpService",
		HTTPGet: HTTPGet{
			TCPSocket: TCPSocket{
				Host: "google.com",
			},
		},
	}
	httpServiceResult = Service{
		Id:                  "httpService",
		Name:                "httpService",
		InitialDelaySeconds: 0,
		PeriodSeconds:       10,
		TimeoutSeconds:      1,
		SuccessThreshold:    1,
		FailureThreshold:    3,
		HTTPGet:             HTTPGet{TCPSocket: TCPSocket{Host: "google.com", Port: 80}, Scheme: "http", Path: "/"},
		TCPSocket:           TCPSocket{},
	}
)

type mockStorage struct {
	Services []Service
	Error    error
}

func (m *mockStorage) Load() ([]Service, error) {
	return m.Services, m.Error
}

func TestStorageService_Load(t *testing.T) {
	type fields struct {
		storage Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Service
		wantErr bool
	}{
		{
			name: "empty service",
			fields: fields{
				storage: &mockStorage{
					Services: []Service{emptyService, tcpService, httpService},
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
					Services: []Service{tcpService, httpService},
					Error:    nil,
				},
			},
			want:    []Service{tcpServiceResult, httpServiceResult},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StorageService{
				storage: tt.fields.storage,
			}
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
